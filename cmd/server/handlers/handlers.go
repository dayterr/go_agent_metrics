package handlers

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/encryption"
	"github.com/dayterr/go_agent_metrics/internal/hash"
	"github.com/dayterr/go_agent_metrics/internal/metric"
	"github.com/dayterr/go_agent_metrics/internal/storage"
)

type AsyncHandler struct {
	storage storage.Storager
	key     string
	dsn     string
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

func DecryptingMiddleware(e encryption.Encryptor) func(http.Handler) http.Handler {
	var b bytes.Buffer

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !(reflect.DeepEqual(e, encryption.Encryptor{})) {
				switch r.Method {
				case http.MethodPost:
					if _, err := b.ReadFrom(r.Body); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					buf, err := e.DecryptMessage(b.Bytes())
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					r.Body = io.NopCloser(bytes.NewReader(buf))
					r.ContentLength = int64(len(buf))

				default:
					next.ServeHTTP(w, r)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CheckIPMiddleware(ts *net.IPNet) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ts == nil {
				next.ServeHTTP(w, r)
			}

			if !ts.Contains(net.ParseIP(r.RemoteAddr)) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
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
		return
	}

	ctx := r.Context()
	switch m.MType {
	case agent.GaugeType:
		v, err := ah.storage.GetGuageByID(ctx, m.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		m.Value = &v
		if ah.key != "" {
			m.Hash = hash.EncryptMetric(m, ah.key)
		}
		mJSON, err := json.Marshal(m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(mJSON)
	case agent.CounterType:
		d, err := ah.storage.GetCounterByID(ctx, m.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		m.Delta = &d
		if ah.key != "" {
			m.Hash = hash.EncryptMetric(m, ah.key)
		}
		mJSON, err := json.Marshal(m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
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
		return
	}

	hashGot := r.Header.Get("Hash")
	if hashGot != "" && ah.key != "" {
		hashCheck := hash.EncryptMetric(m, ah.key)
		if !hash.CheckHash(m, hashCheck) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
	ctx := r.Context()
	switch m.MType {
	case agent.GaugeType:
		ah.storage.SetGuage(ctx, m.ID, m.Value)
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
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
	case agent.GaugeType:
		valFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ah.storage.SetGuage(ctx, metricName, &valFloat)
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
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
	case agent.GaugeType:
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
	case agent.CounterType:
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

func (ah AsyncHandler) GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Println("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	gs, err := ah.storage.GetGauges(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "index.html", gs)
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
