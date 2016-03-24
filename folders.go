package catchpointsdk

import (
  "github.com/parnurzeal/gorequest"
  "fmt"
  "github.com/kelseyhightower/envconfig"
  "github.com/buger/jsonparser"
  "strconv"
)

func init() {
  if division_id == "" {
    division_id = strconv.Itoa(GetDefaultDivisionId())
  }
}

func GetFoldersJson() string {
  return string(getFolders())
}

func GetFolders() (map[int]string) {
  m := make(map[int]string)
  folders, _, _, _ := jsonparser.Get(getFolders(), "items")
  jsonparser.ArrayEach(folders, func(value []byte, dataType int, offset int, err error) {
    folder_name,_,_,_ := jsonparser.Get(value, "name")
    folder_id,_,_ := jsonparser.GetNumber(value, "id")
    m[int(folder_id)] = string(folder_name)
  })
  return m
}

func GetFolderIdByName(folder_name string) int {
  var folder_id float64
  folders, _, _, _ := jsonparser.Get(getFolders(), "items")
  jsonparser.ArrayEach(folders, func(value []byte, dataType int, offset int, err error) {
    fname,_,_,_ := jsonparser.Get(value, "name")
    if string(fname) == folder_name {
      folder_id,_,_ = jsonparser.GetNumber(value, "id")
    }
  })
  return int(folder_id)
}

func getFolders() []byte {
  envconfig.Process("catchpointsdk", &c)
  token := Authenticate()
  resp, body, errs := gorequest.New().Get(fmt.Sprintf("%s/ui/api/v1/folders?division_id=%s", c.Endpoint, division_id)).
    Set("Accept", "*/*").
    Set("Authorization", fmt.Sprintf("Bearer %s", token)).
    EndBytes()
  if errs != nil { panic(errs) }

  if resp.Status == "200 OK" {
    return body
  } else {
    return nil
  }
}
