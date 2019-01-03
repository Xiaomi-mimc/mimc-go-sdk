package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetAddressFromResolver(url, list string) *string {
	params := map[string]string{}
	params["ver"] = "4.0"
	params["type"] = "wifi"
	params["list"] = list

	url = url + "?"
	for key, value := range params {
		url = url + key + "=" + value + "&"
	}
	response, err := http.Get(url)
	if err != nil {
		//fmt.Print("error:%s", err)
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//fmt.Print("error1:%s", err)
		return nil
	}
	token := string(body)
	return &token
}

func GetFEAddress(url, list string) []string {
	response := GetAddressFromResolver(url, list)
	var jsonData = bytes.NewBufferString(*response).Bytes()
	domain := make(map[string]interface{})
	err := json.Unmarshal(jsonData, &domain)
	if err != nil {
		return nil
	}
	if strings.Compare(domain["S"].(string), "Ok") != 0 {
		fmt.Printf("\nstatus:%v\n", domain["S"].(string))
		return nil
	}
	feInfo := domain["R"].(map[string]interface{})["wifi"].(map[string]interface{})[list].([]interface{})
	var feAddrs []string
	for _, value := range feInfo {
		feAddrs = append(feAddrs, value.(string))
	}
	return feAddrs
}

//
