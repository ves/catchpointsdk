package catchpointsdk

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/kelseyhightower/envconfig"
	"github.com/parnurzeal/gorequest"
	"strconv"
	"time"
	"encoding/json"
	"github.com/jehiah/go-strftime"
)

type ProductPayload struct {
	ChangeDate     string                `json:"change_date"`
	DivisionID     int                   `json:"division_id"`
	Id             int                   `json:"id"`
	Name           string                `json:"name"`
	Status         ProductPayloadStatus     `json:"status"`
}

type ProductPayloadStatus struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func init() {
	if division_id == "" {
		division_id = strconv.Itoa(GetDefaultDivisionId())
	}
}

func AddProduct(product *ProductPayload) string {
	token := Authenticate()
	if product.DivisionID == 0 {
		divid,_ := strconv.Atoi(division_id)
		product.DivisionID = divid
	}
	product.ChangeDate = strftime.Format("%m/%d/%Y %I:%M:%S %p", time.Unix(time.Now().Unix(), 0))
	product.Status.Id = 0
	product.Status.Name = "Active"
	b, _ := json.Marshal(product)
	gr := gorequest.New()
	_, body, _ := gr.Post(fmt.Sprintf("%s/ui/api/v1/products/0", c.Endpoint)).
		Set("Accept", "*/*").
		Set("Authorization", fmt.Sprintf("Bearer %s", token)).
		Send(string(b)).
		End()
	return body
}

func GetDefaultProductId() int {
	product_id, _ := jsonparser.GetInt(getProducts(), "id")
	return int(product_id)
}

func GetDefaultDivisionId() int {
	division_id, _ := jsonparser.GetInt(getProducts(), "division_id")
	//fmt.Println(division_id)
	return int(division_id)
}

func GetProducts() map[int]string {
	m := make(map[int]string)
	products, _, _, _ := jsonparser.Get(getProducts(), "items")
	jsonparser.ArrayEach(products, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		product_name, _, _, _ := jsonparser.Get(value, "name")
		product_id, _ := jsonparser.GetInt(value, "id")
		m[int(product_id)] = string(product_name)
	})
	return m
}

func GetProductIdByName(product_name string) int {
	var pid float64
	products, _, _, _ := jsonparser.Get(getProducts(), "items")
	jsonparser.ArrayEach(products, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		pname, _, _, _ := jsonparser.Get(value, "name")
		if string(pname) == product_name {
			pid, _ = jsonparser.GetFloat(value, "id")
		}
	})
	return int(pid)
}

func getProducts() []byte {
	envconfig.Process("catchpointsdk", &c)
	token := Authenticate()
	_, body, errs := gorequest.New().Get(fmt.Sprintf("%s/ui/api/v1/products?divisionId=%s", c.Endpoint, division_id)).
		Set("Accept", "*/*").
		Set("Authorization", fmt.Sprintf("Bearer %s", token)).
		EndBytes()
	if errs != nil {
		panic(errs)
	}
	return body
}
