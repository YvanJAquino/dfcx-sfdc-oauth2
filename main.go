package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/YvanJAquino/dfcx-sfdc-oauth2/dfcx"
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

func (r *OAuth2Request) Encode(code string) *http.Request {

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
	return req
}

func RichHyperLink(url string) *dfcx.RichContents {
	contents := &dfcx.RichContent{
		Type: "button",
		Icon: &dfcx.Icon{
			Type: "account_circle",
		},
		Text: "Salesforce Login",
		Link: url,
	}
	return &dfcx.RichContents{
		RichContent: [][]*dfcx.RichContent{{contents}},
	}
}

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8081"
	}

	http.HandleFunc("/generate-login",
		func(w http.ResponseWriter, r *http.Request) {
			wr, err := dfcx.FromRequest(r)
			if err != nil {
				log.Fatal(err)
			}
			msg := &dfcx.RichContentsMessage{
				Payload: RichHyperLink(wr.SessionInfo.Session),
			}
			var wrsp dfcx.WebhookResponse
			wrsp.TextResponse(w, "Login here!") // not gonna work...!
			wrsp.FulfillmentResponse.Messages = append(wrsp.FulfillmentResponse.Messages, msg)
		})

	http.HandleFunc("/callback",
		func(w http.ResponseWriter, r *http.Request) {
			code, err := url.QueryUnescape(r.URL.Query().Get("code"))
			if err != nil {
				log.Fatal("Error during QueryUnescape: ", err)
			}
			state, err := url.QueryUnescape(r.URL.Query().Get("state"))
			fmt.Println(state)
			if err != nil {
				log.Fatal("Error during QueryUnescape: ", err)
			}
			var oauth OAuth2Request
			req := oauth.Encode(code)
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			io.Copy(w, resp.Body)
			// https://cloud.google.com/go/docs/reference/cloud.google.com/go/dialogflow/latest/cx/apiv3beta1#cloud_google_com_go_dialogflow_cx_apiv3beta1_SessionsClient_DetectIntent
		},
	)
	fmt.Printf("Starting HTTP Server on :%s\n", PORT)
	fmt.Println("Please visit: https://cloudcolosseum-dev-ed.my.salesforce.com/services/oauth2/authorize?client_id=3MVG9p1Q1BCe9GmCjkWZeRQ0vWIHjTtfeOjOsKl0XUWOVuSNQdZ4QzNogu25T_GNO3G3BmaNz.dbQIlYlctCV&redirect_uri=https://sfdc-oauth2-63ietzwyxq-uk.a.run.app/callback&response_type=code")
	http.ListenAndServe(":"+PORT, nil)
}
