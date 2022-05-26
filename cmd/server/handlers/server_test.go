package handlers

import (
	"fmt"
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
			url:  "/update/gauge/test_metric/303",
			name: "test usual gauge metric",
			want: want{
				code: http.StatusOK,
			},
		},
			{
				url:  "/update/counter/test_counter/3",
				name: "test usual counter metric",
				want: want{
					code: http.StatusOK,
				},
			},
			{
			url:  "/update/gauge/test_metric",
			name: "test gauge without a value",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			url:  "/update/counter/test_counter",
			name: "test counter without a value",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			url:  "/update/gauge/test_metric/some",
			name: "test gauge with an incorrect value",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			url:  "/update/counter/test_counter/some",
			name: "test counter with an incorrect value",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			url:  "/update/rand/test_counter/some",
			name: "test case with an incorrect metric type",
			want: want{
				code: http.StatusNotImplemented,
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

func TestGetMetric(t *testing.T) {
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
			url:  "/value/gauge/test_metric",
			name: "test usual gauge metric",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			url:  "/value/counter/test_counter",
			name: "test usual counter metric",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			url:  "/value/gauge/test_metric3",
			name: "test gauge with a non-existent value",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			url:  "/value/counter/test_counter3",
			name: "test counter with a non-existent value",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			url:  "/value/rand/test_counter",
			name: "test case with an incorrect metric type",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			url:  "/value/gauge",
			name: "test gauge without a metric name",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			url:  "/value/counter",
			name: "test counter without a metric name",
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := CreateRouter()
			ts := httptest.NewServer(r)
			defer ts.Close()
			url := "/update/gauge/test_metric/303"
			testRequest(t, ts, http.MethodPost, url, nil)
			url = "/update/counter/test_counter/3"
			testRequest(t, ts, http.MethodPost, url, nil)
			req, _ := testRequest(t, ts, http.MethodGet, tt.url, nil)
			assert.Equal(t, tt.want.code, req.StatusCode, "Возвращаемый код не равен ожидаемому")
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	fmt.Println(req.URL.Path)

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