package sfdc

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	LOGIN_URL          = "https://login.salesforce.com/services/oauth2/token"
	INSTANCE_LOGIN_URL = "https://cloudcolosseum-dev-ed.my.salesforce.com/services/oauth2/token"
	GRANT_TYPE         = "authorization_code"
	CLIENT_ID          = "3MVG9p1Q1BCe9GmCjkWZeRQ0vWIHjTtfeOjOsKl0XUWOVuSNQdZ4QzNogu25T_GNO3G3BmaNz.dbQIlYlctCV"
	CLIENT_SECRET      = "52B4F8AAD37457E84E5D22029FBA3E9DD1C593E4FB66622DE65ECC0B580A6B65"
	REDIRECT_URI       = "https://sfdc-oauth2-63ietzwyxq-uk.a.run.app/callback"
)

type OAuth2Request struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
}

func NewLoginRequest(code string) *http.Request {
	var or OAuth2Request
	req := or.PrepareRequest(code)
	return req
}

func (r *OAuth2Request) PrepareRequest(code string) *http.Request {

	r.GrantType = GRANT_TYPE
	r.Code = code
	r.ClientId = CLIENT_ID
	r.ClientSecret = CLIENT_SECRET
	r.RedirectUri = REDIRECT_URI

	data := url.Values{}
	data.Set("grant_type", r.GrantType)
	data.Set("code", r.Code)
	data.Set("client_id", r.ClientId)
	data.Set("client_secret", r.ClientSecret)
	data.Set("redirect_uri", r.RedirectUri)
	encoded := data.Encode()
	req, err := http.NewRequest(http.MethodPost, INSTANCE_LOGIN_URL, strings.NewReader(encoded))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}
