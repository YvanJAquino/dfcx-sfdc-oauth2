package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/YvanJAquino/dfcx-sfdc-oauth2/dfcx"
	"github.com/go-redis/redis/v8"
)

// Argolis Address: https://sfdc-oauth2-a67fdjzmma-uc.a.run.app/generate-login

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

func RichHyperLink(s string) *dfcx.RichContents {
	url := "https://cloudcolosseum-dev-ed.my.salesforce.com/services/oauth2/authorize?client_id=3MVG9p1Q1BCe9GmCjkWZeRQ0vWIHjTtfeOjOsKl0XUWOVuSNQdZ4QzNogu25T_GNO3G3BmaNz.dbQIlYlctCV&redirect_uri=https://sfdc-oauth2-63ietzwyxq-uk.a.run.app/callback&response_type=code&state=" + s
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

var parent = context.Background()

// Refactor to use Secret Manager
var opts = &redis.Options{
	Addr:     "10.62.49.107:6378",
	Password: "", // no password set
	DB:       0,  // use default DB
}
var rdb = redis.NewClient(opts)

type Session struct {
	Session string
	Token   string
}

func (d *Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Session) ToRedis(ctx context.Context, rdb *redis.Client, key string) error {
	err := rdb.Set(ctx, key, d, 12*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (d *Session) FromRedis(ctx context.Context, rdb *redis.Client, key string) error {
	s, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(s), d)
	return nil
}

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8081"
	}

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			msg := &struct {
				Message string
			}{
				Message: "Hello World!",
			}
			json.NewEncoder(w).Encode(msg)
		})

	http.HandleFunc("/generate-login",
		func(w http.ResponseWriter, r *http.Request) {
			wr, err := dfcx.FromRequest(r)
			if err != nil {
				log.Fatal(err)
			}
			session, err := wr.SessionInfo.ExtractSession()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(session)
			var s Session
			s.Session = wr.SessionInfo.Session
			err = s.ToRedis(parent, rdb, session)
			if err != nil {
				log.Fatal(err)
			}

			msg := &dfcx.RichContentsMessage{
				Payload: RichHyperLink(wr.SessionInfo.Session),
			}
			resp := dfcx.NewTextResponse("Click on the link below to login. Once you've logged in, you can say 'done' to move forward. ")
			resp.AddMessage(msg)
			resp.Respond(w)
		})

	http.HandleFunc("/callback",
		func(w http.ResponseWriter, r *http.Request) {
			code, err := url.QueryUnescape(r.URL.Query().Get("code"))
			if err != nil {
				log.Fatal("Error during QueryUnescape: ", err)
			}
			state, err := url.QueryUnescape(r.URL.Query().Get("state"))
			fmt.Println("Callback state: ", state)
			if err != nil {
				log.Fatal("Error during QueryUnescape: ", err)
			}
			var s Session
			err = s.FromRedis(parent, rdb, state)
			if err != nil {
				log.Fatal(err)
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
			io.Copy(os.Stdout, resp.Body)
			htmlFile, err := os.Open("static/callback-response.html")
			if err != nil {
				log.Fatal(err)
			}
			defer htmlFile.Close()
			w.Header().Set("Content-Type", "text/html")
			io.Copy(w, htmlFile)
			// io.Copy(w, resp.Body)
			// https://cloud.google.com/go/docs/reference/cloud.google.com/go/dialogflow/latest/cx/apiv3beta1#cloud_google_com_go_dialogflow_cx_apiv3beta1_SessionsClient_DetectIntent
		},
	)
	fmt.Printf("Starting HTTP Server on :%s\n", PORT)
	fmt.Println("Please visit: https://cloudcolosseum-dev-ed.my.salesforce.com/services/oauth2/authorize?client_id=3MVG9p1Q1BCe9GmCjkWZeRQ0vWIHjTtfeOjOsKl0XUWOVuSNQdZ4QzNogu25T_GNO3G3BmaNz.dbQIlYlctCV&redirect_uri=https://sfdc-oauth2-63ietzwyxq-uk.a.run.app/callback&response_type=code")
	http.ListenAndServe(":"+PORT, nil)
}
