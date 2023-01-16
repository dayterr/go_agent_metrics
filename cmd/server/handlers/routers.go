package handlers

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dayterr/go_agent_metrics/internal/encryption"
	"github.com/dayterr/go_agent_metrics/internal/server"
	"github.com/dayterr/go_agent_metrics/internal/storage"
)

type NextHandler func(next http.Handler) http.Handler
type Salt struct{}

func PassSalt(salt []byte) NextHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(writer, req.WithContext(context.WithValue(req.Context(), Salt{}, salt)))
		})
	}
}

func NewAsyncHandler(key, dsn string) (AsyncHandler, error) {
	var s storage.Storager
	var err error
	if dsn != "" {
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

func CreateRouterWithAsyncHandler(filename string, isRestored bool, h AsyncHandler, e encryption.Encryptor,
	salt []byte, ts string) (chi.Router, error) {
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

	_, cidr, err := net.ParseCIDR(ts)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	//r.Use(gzipHandle)
	r.Use(DecryptingMiddleware(e), gzipHandle, PassSalt(salt))
	r.Use(DecryptingMiddleware(e), CheckIPMiddleware(cidr), PassSalt(salt))
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
