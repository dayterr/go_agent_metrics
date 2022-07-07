package handlers

import (
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"github.com/dayterr/go_agent_metrics/internal/hash"
	"github.com/dayterr/go_agent_metrics/internal/metric"
	"github.com/dayterr/go_agent_metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/dayterr/go_agent_metrics/internal/agent"
)

type AsyncHandler struct {
	storage storage.Storager
	key     string
	dsn string
}

type SyncHandler struct {
	storage storage.Storager
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

func (ah AsyncHandler) MarshallMetrics() ([]byte, error) {
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

	switch m.MType {
	case agent.GaugeType:
		v, err := ah.storage.GetGuageByID(m.ID)
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
	case agent.CounterType:
		d, err := ah.storage.GetCounterByID(m.ID)
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

func (as AsyncHandler) PostJSON(w http.ResponseWriter, r *http.Request) {
	var m metric.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	hashGot := r.Header.Get("Hash")
	if hashGot != "" && as.key != "" {
		hashCheck := hash.EncryptMetric(m, as.key)
		if !hash.CheckHash(m, hashCheck) {
			w.WriteHeader(http.StatusBadRequest)
		}

	}
	switch m.MType {
	case agent.GaugeType:
		as.storage.SetGuage(m.ID, m.Value)
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
		as.storage.SetCounter(m.ID, m.Delta)
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
	switch metricType {
	case agent.GaugeType:
		valFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ah.storage.SetGuage(metricName, &valFloat)
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
		valInt, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ah.storage.SetCounterFromMemStats(metricName, int64(valInt))
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
	switch metricType {
	case agent.GaugeType:
		if ah.storage.CheckGaugeByName(metricName) {
			v, err := ah.storage.GetGuageByID(metricName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			value := strconv.FormatFloat(v, 'f', -1, 64)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(value))
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case agent.CounterType:
		if ah.storage.CheckCounterByName(metricName) {
			v, err := ah.storage.GetCounterByID(metricName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			c := strconv.Itoa(int(v))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(c))
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (ah AsyncHandler) GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "index.html", ah.storage.GetGauges())
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
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
	return
}

func (ah AsyncHandler) PostMany(w http.ResponseWriter, r *http.Request) {
	var metricList []metric.Metrics

	err := json.NewDecoder(r.Body).Decode(&metricList)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = ah.storage.SaveMany(metricList)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}