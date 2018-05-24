package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type TokenHandler struct {
	httpUrl    string
	AppId      int64  `json:"appId"`
	AppKey     string `json:"appKey"`
	AppSecret  string `json:"appSecret"`
	AppAccount string `json:"appAccount"`
}

func NewTokenHandler(httpUrl, appKey, appSecret, appAccount *string, appId *int64) *TokenHandler {
	tokenHandler := new(TokenHandler)
	tokenHandler.httpUrl = *httpUrl
	tokenHandler.AppId = *appId
	tokenHandler.AppKey = *appKey
	tokenHandler.AppSecret = *appSecret
	tokenHandler.AppAccount = *appAccount
	return tokenHandler
}

func (this *TokenHandler) FetchToken() *string {
	jsonBytes, err := json.Marshal(*this)
	if err != nil {
		return nil
	}
	requestJsonBody := bytes.NewBuffer(jsonBytes).String()
	request, err := http.Post(this.httpUrl, "application/json", strings.NewReader(requestJsonBody))
	if err != nil {
		return nil
	}
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil
	}
	token := string(body)
	return &token
}
