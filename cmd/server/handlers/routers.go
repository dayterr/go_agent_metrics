package handlers

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dayterr/go_agent_metrics/internal/server"
	"github.com/dayterr/go_agent_metrics/internal/storage"
)

func NewAsyncHandler(key, dsn string, isDB bool) (AsyncHandler, error) {
	var s storage.Storager
	var err error
	if isDB {
		s, err = storage.NewDB(dsn)
		if err != nil {
			log.Println(err)
			return AsyncHandler{}, err
		}
	} else {
		s = storage.NewIMS()
	}
	h := AsyncHandler{storage: s, key: key, dsn: dsn}
	return h, nil
}

func CreateRouterWithAsyncHandler(filename string, isRestored bool, h AsyncHandler) (chi.Router, error) {
	// Функция для создания нового роутера
	if isRestored {
		var err error
		h.storage, err = server.LoadMetricsFromFile(filename)
		log.Println("uploaded", h.storage)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	r := chi.NewRouter()
	r.Use(gzipHandle)
	r.Mount("/debug", middleware.Profiler())

	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.PostJSON)
		r.Post("/{metricType}/{metricName}/{value}", h.PostMetric)
	})
	r.Post("/value/", h.GetValue)
	r.Get("/value/{metricType}/{metricName}", h.GetMetric)
	r.Get("/", h.GetIndex)
	r.Get("/ping", h.Ping)
	r.Post("/updates/", h.PostMany)
	return r, nil
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
