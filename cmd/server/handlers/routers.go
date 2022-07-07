package handlers

import (
	"github.com/dayterr/go_agent_metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
)

func NewAsyncHandler(key, dsn string, isDB bool) AsyncHandler {
	var s storage.Storager
	var err error
	if isDB {
		s, err = storage.NewDB(dsn)
		if err != nil {
			log.Println(err)
		}
	} else {
		s = storage.NewIMS()
	}
	h := AsyncHandler{storage: s, key: key, dsn: dsn}
	return h
}

func CreateRouterWithAsyncHandler(filename string, isRestored bool, h AsyncHandler) chi.Router {
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
	r.Get("/ping", h.Ping)
	r.Post("/updates/", h.PostMany)
	return r
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
