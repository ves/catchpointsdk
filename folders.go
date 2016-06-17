package catchpointsdk

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/kelseyhightower/envconfig"
	"github.com/parnurzeal/gorequest"
	"strconv"
	"encoding/json"
	"github.com/jehiah/go-strftime"
	"time"
)

type FolderPayload struct {
	ChangeDate     string                `json:"change_date"`
	DivisionID     int                   `json:"division_id"`
	Id             int                   `json:"id"`
	Name           string                `json:"name"`
	ProductID      int                   `json:"product_id"`
	ParentFolderID int                   `json:"parent_folder_id"`
}

func init() {
	if division_id == "" {
		division_id = strconv.Itoa(GetDefaultDivisionId())
	}
}

func AddFolder(folder *FolderPayload) string {
	token := Authenticate()
	if folder.DivisionID == 0 {
		divid,_ := strconv.Atoi(division_id)
		folder.DivisionID = divid
	}
	if folder.ProductID == 0 {
		folder.ProductID = GetDefaultProductId()
	}
	folder.ChangeDate = strftime.Format("%m/%d/%Y %I:%M:%S %p", time.Unix(time.Now().Unix(), 0))
	b, _ := json.Marshal(folder)
	gr := gorequest.New()
	_, body, _ := gr.Post(fmt.Sprintf("%s/ui/api/v1/folders/0", c.Endpoint)).
		Set("Accept", "*/*").
		Set("Authorization", fmt.Sprintf("Bearer %s", token)).
		Send(string(b)).
		End()
	return body
}

func GetFoldersJson() string {
	return string(getFolders())
}

func GetFolders() map[int]string {
	m := make(map[int]string)
	data := getFolders()
	if data != nil {
		folders, _, _, _ := jsonparser.Get(getFolders(), "items")
		jsonparser.ArrayEach(folders, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			folder_name, _, _, _ := jsonparser.Get(value, "name")
			folder_id, _ := jsonparser.GetInt(value, "id")
			m[int(folder_id)] = string(folder_name)
		})
	}
	return m
}

func GetFolderIdByName(folder_name string) int {
	var folder_id float64
	folders, _, _, _ := jsonparser.Get(getFolders(), "items")
	jsonparser.ArrayEach(folders, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		fname, _, _, _ := jsonparser.Get(value, "name")
		if string(fname) == folder_name {
			folder_id, _ = jsonparser.GetFloat(value, "id")
		}
	})
	return int(folder_id)
}

func getFolders() []byte {
	envconfig.Process("catchpointsdk", &c)
	token := Authenticate()
	resp, body, errs := gorequest.New().Get(fmt.Sprintf("%s/ui/api/v1/folders?divisionId=%s", c.Endpoint, division_id)).
		Set("Accept", "*/*").
		Set("Authorization", fmt.Sprintf("Bearer %s", token)).
		EndBytes()
	if errs != nil {
		panic(errs)
	}

	if resp.Status == "200 OK" {
		return body
	} else {
		return nil
	}
}
