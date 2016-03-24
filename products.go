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

func GetDefaultProductId() int {
  fmt.Println(string(getProducts()))
  product_id,_,_ := jsonparser.GetNumber(getProducts(), "id")
  return int(product_id)
}

func GetDefaultDivisionId() int {
  division_id,_,_ := jsonparser.GetNumber(getProducts(), "division_id")
  return int(division_id)
}

func GetProducts() (map[int]string) {
  m := make(map[int]string)
  products, _, _, _ := jsonparser.Get(getProducts(), "items")
  jsonparser.ArrayEach(products, func(value []byte, dataType int, offset int, err error) {
    product_name,_,_,_ := jsonparser.Get(value, "name")
    product_id,_,_ := jsonparser.GetNumber(value, "id")
    m[int(product_id)] = string(product_name)
  })
  return m
}

func GetProductIdByName(product_name string) int {
  var pid float64
  products, _, _, _ := jsonparser.Get(getProducts(), "items")
  jsonparser.ArrayEach(products, func(value []byte, dataType int, offset int, err error) {
    pname,_,_,_ := jsonparser.Get(value, "name")
    if string(pname) == product_name {
      pid,_,_ = jsonparser.GetNumber(value, "id")
    }
  })
  return int(pid)
}


func getProducts() []byte {
  envconfig.Process("catchpointsdk", &c)
  token := Authenticate()
  _, body, errs := gorequest.New().Get(fmt.Sprintf("%s/ui/api/v1/products?division_id=%s", c.Endpoint, division_id)).
    Set("Accept", "*/*").
    Set("Authorization", fmt.Sprintf("Bearer %s", token)).
    EndBytes()
  if errs != nil { panic(errs) }
  return body
}
