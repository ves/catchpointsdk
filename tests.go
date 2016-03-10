package catchpointsdk

import (
  "net/http"
  "fmt"
  "io/ioutil"
  "github.com/kelseyhightower/envconfig"
  "strings"
  "github.com/jmoiron/jsonq"
  "encoding/json"
  "github.com/fatih/color"
)


func ListTestsCli() {
  jsonstring := ListTestsJson()
  data := map[string]interface{}{}
  dec := json.NewDecoder(strings.NewReader(jsonstring))
  dec.Decode(&data)
  jq := jsonq.NewQuery(data)
  tests, err := jq.Array("items")
  if err != nil { panic(err) }
  color.Red(fmt.Sprintf("%d tests found.\n\n", len(tests)))
  for _, value := range tests {
    jq := jsonq.NewQuery(value)
    tname, _ := jq.String("name")
    ttype, _ := jq.String("monitors", "0", "name")
    tstatus, _ := jq.String("status", "name")
    turl, _ := jq.String("test_url")
    color.Green(fmt.Sprintf("%s \t <type: %s> \t <status: %s> \t <url: %s>", tname, ttype, tstatus, turl))
  }
}

func ListTestsJson() string {
  envconfig.Process("catchpointsdk", &c)
  token := Authenticate()
  client := &http.Client{}
  resp, err := client.Get(c.Endpoint)
  if err != nil { panic(err) }
  uri := fmt.Sprintf("%s/ui/api/v1/tests", c.Endpoint)
  req, err := http.NewRequest("GET", uri, nil)
  if err != nil { panic(err) }
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
  req.Header.Add("Accept", "*/*")
  resp, err = client.Do(req)
  if err != nil { panic(err) }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil { panic(err) }

  if resp.Status == "200 OK" {
    return string(body)
  } else {
    return ""
  }
}
