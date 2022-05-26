package handlers

import (
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostMetric(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		url string
		name string
		want want
	}{
		{
			url: "/update/gauge/test_metric/303",
			name: "test usual gauge metric",
			want: want{
				code:        http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := CreateRouter()
			ts := httptest.NewServer(r)
			defer ts.Close()
			req, _ := testRequest(t, ts, http.MethodPost, tt.url, nil)
			assert.Equal(t, tt.want.code, req.StatusCode, "Возвращаемый код не равен ожидаемому")
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	return resp, string(respBody)
}