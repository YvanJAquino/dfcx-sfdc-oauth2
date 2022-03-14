package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/YvanJAquino/dfcx-sfdc-oauth2/dfcx"
	"github.com/YvanJAquino/dfcx-sfdc-oauth2/helpers"
)

var HttpMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodHead,
	http.MethodPut,
	http.MethodDelete,
}

func TestHealthCheck(t *testing.T) {
	for _, method := range HttpMethods {
		req := httptest.NewRequest(method, "/health-check", nil)
		w := httptest.NewRecorder()
		HealthCheckHandler(w, req)
		res := w.Result()
		got := res.StatusCode
		var want int
		if method == http.MethodGet {
			want = http.StatusOK
		} else {
			want = http.StatusMethodNotAllowed
		}
		if want != got {
			t.Error("Method: ", method, " ~ Status: ", got)
		}
	}
}

func TestGenerateLoginHandler(t *testing.T) {
	var bodies dfcx.WebhookRequests
	file, err := os.Open("generate_login_tests.json")
	helpers.HandleError(err, "os.Open", false)
	readers, err := bodies.UnmarshalJSONToReaders(file)
	helpers.HandleError(err, "bodies.UmarshalJSONToReaders", false)
	for index := range readers {
		r := readers[index]
		req := httptest.NewRequest(http.MethodPost, "/generate-login", r)
		w := httptest.NewRecorder()
		GenerateLoginHandler(w, req)
		res := w.Result()
		defer res.Body.Close()
		webhook, err := dfcx.FromReader(res.Body)
		helpers.HandleError(err, "dfcx.FromRequest", false)
		fmt.Println(webhook)
	}
}
