package catchpointsdk

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	"fmt"
	"github.com/jehiah/go-strftime"
	"github.com/kelseyhightower/envconfig"
	"github.com/parnurzeal/gorequest"
	"strconv"
	"time"
)

type TestProperties struct {
	AdditionalProperties bool        `json:"additionalProperties"`
	Properties           TestPayload `json:"properties"`
}

type TestPayload struct {
	ChangeDate     string                `json:"change_date"`
	DivisionID     int                   `json:"division_id"`
	Id             int                   `json:"id"`
	Monitor        TestPayloadMonitor    `json:"monitor"`
	Advanced       TestAdvancedOnFailure `json:"advanced_settings"`
	Name           string                `json:"name"`
	ProductID      int                   `json:"product_id"`
	ParentFolderID int                   `json:"parent_folder_id"`
	StartDate      string                `json:"start"`
	Status         TestPayloadStatus     `json:"status"`
	TestType       TestPayloadType       `json:"test_type"`
	TestURL        string                `json:"test_url"`
}

type TestPayloadMonitor struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TestPayloadStatus struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TestPayloadType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TestAdvancedOnFailure struct {
	OnFailure TestAdvancedOnFailureSettings `json:"on_failure"`
}

type TestAdvancedOnFailureSettings struct {
	VerifyTest       bool `json:"verify_test"`
	DebugPrimaryHost bool `json:"debug_primary_host"`
}

type TestAdvancedCaptureHttpHeaders struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TestAdvancedCaptureResponseContent struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func init() {
	if division_id == "" {
		division_id = strconv.Itoa(GetDefaultDivisionId())
	}
}

func AddTest(folder_name string, product_name string, test *TestPayload) string {
	token := Authenticate()
	test.TestType.Id = getTestTypeId(test.TestType.Name)
	test.Monitor.Id = getMonitorId(test.Monitor.Name)
	test.ChangeDate = strftime.Format("%m/%d/%Y %I:%M:%S %p", time.Unix(time.Now().Unix(), 0))
	test.StartDate = test.ChangeDate
	test.Status.Id = 0
	test.Status.Name = "Active"
	if folder_name != "" {
		test.ParentFolderID = GetFolderIdByName(folder_name)
	}
	if product_name == "" {
		test.ProductID = GetDefaultProductId()
	} else {
		test.ProductID = GetProductIdByName(product_name)
	}
	b, _ := json.Marshal(test)
	gr := gorequest.New()
	_, body, _ := gr.Post(fmt.Sprintf("%s/ui/api/v1/tests/0", c.Endpoint)).
		Set("Accept", "*/*").
		Set("Authorization", fmt.Sprintf("Bearer %s", token)).
		Send(string(b)).
		End()
	return body
}

func ListTests() []TestPayload {
	testitems, _, _, _ := jsonparser.Get(listTests(), "items")
	keys := make([]TestPayload, 0)
	json.Unmarshal(testitems, &keys)
	return keys
}

func ListTestsJson() string {
	return string(listTests())
}

func listTests() []byte {
	envconfig.Process("catchpointsdk", &c)
	token := Authenticate()
	gr := gorequest.New()
	resp, body, errs := gr.Get(fmt.Sprintf("%s/ui/api/v1/tests?divisionId=%s", c.Endpoint, division_id)).
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

func getMonitorId(monitor_name string) int {
	m := map[string]int{
		"Object":         2,
		"Emulated":       3,
		"IEBrowser":      0,
		"Tcp":            15,
		"Ftp":            16,
		"PingIcmp":       8,
		"PingTcp":        15,
		"PingUdp":        23,
		"DnsDirect":      13,
		"DnsExperience":  12,
		"ChromeBrowser":  18,
		"Mobile":         26,
		"Playback":       19,
		"MobilePlayback": 20,
		"SMTP":           21,
		"Api":            25,
		"Streaming":      24,
		"Ssh":            28,
		"TraceRouteIcmp": 9,
		"TraceRouteUdp":  14,
		"TraceRouteTcp":  29,
	}
	return m[monitor_name]
}

func getTestTypeId(test_type_name string) int {
	m := map[string]int{
		"Web":         0,
		"Transaction": 1,
		"HtmlCode":    2,
		"Ftp":         3,
		"Tcp":         4,
		"Dns":         5,
		"Ping":        6,
		"Smtp":        7,
		"Api":         9,
		"Streaming":   10,
		"Ssh":         11,
		"TraceRoute":  12,
	}
	return m[test_type_name]
}
