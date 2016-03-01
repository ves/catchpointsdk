package catchpointsdk

import (
  "net/http"
  "github.com/kelseyhightower/envconfig"
  //"encoding/json"
  "fmt"
  "log"
  "bytes"
  "io/ioutil"
)

type Authentication struct {
  ClientID string
  ClientSecret  string
  Endpoint  string
}

/*
Authenticate into the Catchpoint API and return a Bearer token
*/
func Authenticate() {
  var a Authentication
  err := envconfig.Process("catchpointsdk", &a)
  if err != nil { log.Fatal(err.Error()) }
  fmt.Println(a.Endpoint)
  b := []byte(`{"grant_type":"client_credentials","client_credentials":"%s","client_secret":"%s"}, a.ClientID, a.ClientSecret`)
  req, err := http.NewRequest("POST", a.Endpoint, bytes.NewBuffer(b))
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()

  fmt.Println("response Status:", resp.Status)
  fmt.Println("response Headers:", resp.Header)
  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println("response Body:", string(body))
  //return body
}
