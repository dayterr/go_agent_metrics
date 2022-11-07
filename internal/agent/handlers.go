package agent

import (
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dayterr/go_agent_metrics/internal/hash"
	"github.com/dayterr/go_agent_metrics/internal/metric"
	"github.com/dayterr/go_agent_metrics/internal/server"
	"github.com/dayterr/go_agent_metrics/internal/storage"
)

type AsyncHandler struct {
	storage storage.Storager
	key     string
	dsn     string
}

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

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func CreateRouterWithAsyncHandler(filename string, isRestored bool, h AsyncHandler) chi.Router {
	// Функция для создания нового роутера
	if isRestored {
		var err error
		h.storage, err = server.LoadMetricsFromFile(filename)
		log.Println("uploaded", h.storage)
		if err != nil {
			log.Fatal(err)
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
	//r.Get("/", h.GetIndex)
	r.Get("/ping", h.Ping)
	r.Post("/updates/", h.PostMany)
	return r
}

func (ah AsyncHandler) MarshallMetrics() ([]byte, error) {
	// Метод возвращает данные из хранилища в json-формате
	jsn, err := json.Marshal(ah.storage)
	if err != nil {
		return nil, err
	}
	return jsn, nil
}

func (ah AsyncHandler) GetValue(w http.ResponseWriter, r *http.Request) {
	var m metric.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	ctx := r.Context()
	switch m.MType {
	case GaugeType:
		v, err := ah.storage.GetGuageByID(ctx, m.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		m.Value = &v
		if ah.key != "" {
			m.Hash = hash.EncryptMetric(m, ah.key)
		}
		mJSON, err := json.Marshal(m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("content-type", "application/json")
		w.Write(mJSON)
	case CounterType:
		d, err := ah.storage.GetCounterByID(ctx, m.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		m.Delta = &d
		if ah.key != "" {
			m.Hash = hash.EncryptMetric(m, ah.key)
		}
		mJSON, err := json.Marshal(m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("content-type", "application/json")
		w.Write(mJSON)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (ah AsyncHandler) PostJSON(w http.ResponseWriter, r *http.Request) {
	var m metric.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	hashGot := r.Header.Get("Hash")
	if hashGot != "" && ah.key != "" {
		hashCheck := hash.EncryptMetric(m, ah.key)
		if !hash.CheckHash(m, hashCheck) {
			w.WriteHeader(http.StatusBadRequest)
		}

	}
	ctx := r.Context()
	switch m.MType {
	case GaugeType:
		ah.storage.SetGuage(ctx, m.ID, m.Value)
		w.WriteHeader(http.StatusOK)
	case CounterType:
		ah.storage.SetCounter(ctx, m.ID, m.Delta)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (ah AsyncHandler) PostMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	value := strings.Trim(chi.URLParam(r, "value"), "\n")
	if value == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	ctx := r.Context()
	switch metricType {
	case GaugeType:
		valFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ah.storage.SetGuage(ctx, metricName, &valFloat)
		w.WriteHeader(http.StatusOK)
	case CounterType:
		valInt, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ah.storage.SetCounterFromMemStats(ctx, metricName, int64(valInt))
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (ah AsyncHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := r.Context()
	switch metricType {
	case GaugeType:
		if ah.storage.CheckGaugeByName(ctx, metricName) {
			v, err := ah.storage.GetGuageByID(ctx, metricName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			value := strconv.FormatFloat(v, 'f', -1, 64)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(value))
			return
		}
	case CounterType:
		if ah.storage.CheckCounterByName(ctx, metricName) {
			v, err := ah.storage.GetCounterByID(ctx, metricName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			c := strconv.Itoa(int(v))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(c))
			return
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !ah.storage.CheckCounterByName(ctx, metricName) && !ah.storage.CheckGaugeByName(ctx, metricName) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (ah AsyncHandler) Ping(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", ah.dsn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ah AsyncHandler) PostMany(w http.ResponseWriter, r *http.Request) {
	var metricList []metric.Metrics

	err := json.NewDecoder(r.Body).Decode(&metricList)
	if err != nil {
		log.Println("err decoding", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = ah.storage.SaveMany(ctx, metricList)
	if err != nil {
		log.Println("err saving", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
