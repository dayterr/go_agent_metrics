package handlers

import (
	"github.com/dayterr/go_agent_metrics/internal/storage"
	"log"
	"github.com/go-chi/chi/v5"
)

func CreateRouterWithAsyncHandler(filename string, isRestored bool) (chi.Router, AsyncHandler) {
	h := AsyncHandler{storage: storage.InMemoryStorage{}}
	if isRestored {
		err := h.storage.LoadMetricsFromFile(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	r := chi.NewRouter()
	r.Use(gzipHandle)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.PostJSON)
		r.Post("/{metricType}/{metricName}/{value}", h.PostMetric)
	})
	r.Post("/value/", h.GetValue)
	r.Get("/value/{metricType}/{metricName}", h.GetMetric)
	r.Get("/", h.GetIndex)
	return r, h
}

/*func CreateRouterWithSyncHandler(filename string, isRestored bool) chi.Router {
	h := SyncHandler{}
	if isRestored {
		err := allMetrics.LoadMetricsFromFile(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	r := chi.NewRouter()
	r.Use(gzipHandle)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.PostJSON)
		r.Post("/{metricType}/{metricName}/{value}", h.PostMetric)
	})
	r.Post("/value/", h.GetValue)
	r.Get("/value/{metricType}/{metricName}", h.GetMetric)
	r.Get("/", h.GetIndex)
	return r
}*/
