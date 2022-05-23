package main

import (
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
			name: "positive test #1",
			want: want{
				code:        200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostMetric)
			h.ServeHTTP(w, request)
			res := w.Result()
			/*if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}*/
			defer res.Body.Close()
		})
	}
}
