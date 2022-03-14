package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/YvanJAquino/dfcx-sfdc-oauth2/dfcx"
	"github.com/YvanJAquino/dfcx-sfdc-oauth2/helpers"
	"github.com/YvanJAquino/dfcx-sfdc-oauth2/sfdc"
	"github.com/go-redis/redis"
)

var client = &http.Client{}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Health checks must be HTTP/S GET. ")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GenerateLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Generate Login requests must be HTTP/S POST")
		return
	}
	webhook, err := dfcx.FromRequest(r)
	helpers.HandleError(err, "dfcx.FromRequest", false)
	webhookSessionId, err := webhook.SessionInfo.ExtractSession()
	helpers.HandleError(err, "webhook.SessionInfo.ExtractSession", false)
	resp := dfcx.NewTextResponse("Login with the link below.  When your done, just hit enter. ")
	resp.AddMessage(&dfcx.RichContentsMessage{Payload: RichHyperLink(webhookSessionId)})
	resp.Respond(w)

}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code, err := url.QueryUnescape(r.URL.Query().Get("code"))
	helpers.HandleError(err, "url.QueryUnescape", false)
	state, err := url.QueryUnescape(r.URL.Query().Get("state"))
	helpers.HandleError(err, "url.QueryUnescape", false)
	fmt.Println("Code: ", code, " ~ State: ", state)
	req := sfdc.NewLoginRequest(code)
	resp, err := client.Do(req)
	helpers.HandleError(err, "client.Do", false)
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
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

type RedisStateManager struct {
	Client *redis.Client
	State  *State
}

func NewRedisStateManager(c *redis.Client) *RedisStateManager {
	var rsm RedisStateManager
	rsm.Client = c
	rsm.State = &State{}
	return &rsm
}

func (rsm *RedisStateManager) SetKey(ctx context.Context) error {
	err := rsm.Client.Set(ctx, rsm.State.SessionId, rsm.State, 12*time.Hour).Err()
	return err
}

func (rsm *RedisStateManager) GetKey(ctx context.Context) error {
	s, err := rsm.Get(ctx, rsm.SessionId).Result()
	if err != nil {
		return err
	}
	err = rsm.UnmarshalString(s)
	return err
}

func (rsm *RedisStateManager) FromRequest(r *http.Request) (*dfcx.WebhookRequest, error) {
	webhook, err := dfcx.FromRequest(r)
	if err != nil {
		return nil, err
	}
	rsm.SetSessionId(webhook)
	rsm.SetKey(r.Context())
	return webhook, nil
}

type State struct {
	Session     string `json:"session"`
	SessionId   string `json:"sessionId,omitempty"`
	AccessToken string `json:"accessToken,omitempty"`
}

func (s *State) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *State) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *State) UnmarshalString(data string) error {
	return json.Unmarshal([]byte(data), s)
}

func (s *State) SetSessionId(r *dfcx.WebhookRequest) error {
	sessionId, err := r.SessionInfo.ExtractSession()
	if err != nil {
		return err
	}
	s.SessionId = sessionId
	return nil
}

type GenerateLoginHandle struct {
	*RedisStateManager
}

func NewGenerateLoginHandle(c *redis.Client) *GenerateLoginHandle {
	var h GenerateLoginHandle
	h.Init(c)
	return &h
}

func (h *GenerateLoginHandle) Init(c *redis.Client) {
	h.RedisStateManager = NewRedisStateManager(c)
}

func (h *GenerateLoginHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := h.FromRequest(r)
	if err != nil {
		fmt.Println("Error during h.FromRequest")
		return
	}
	resp := dfcx.NewTextResponse("Login with the link below.  When your done, just hit enter. ")
	resp.AddMessage(&dfcx.RichContentsMessage{Payload: RichHyperLink(h.SessionId)})
	resp.Respond(w)
}
