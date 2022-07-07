package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/dayterr/go_agent_metrics/internal/metric"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func NewDB(dsn string) (DBStorage, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return DBStorage{}, err
	}
	defer db.Close()
	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS gauge (id serial PRIMARY KEY, name text UNIQUE NOT NULL, Value double precision NOT NULL);`)
	if err != nil {
		return DBStorage{}, err
	}
	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS counter (id serial PRIMARY KEY, name text UNIQUE, Delta BIGINT);`)
	if err != nil {
		return DBStorage{}, err
	}
	return DBStorage{
		GaugeField:   make(map[string]Gauge),
		CounterField: make(map[string]Counter),
		DSN: dsn,
	}, nil
}

func (s DBStorage) LoadMetricsFromFile(filename string) error {
	if _, err := os.Stat(filename); err != nil {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &s)
	if err != nil {
		return err
	}
	for key, value := range s.GaugeField {
		s.SetGaugeFromMemStats(key, value.ToFloat())
	}
	for key, value := range s.CounterField {
		s.SetCounterFromMemStats(key, value.ToInt64())
	}
	return nil
}

func (s DBStorage) GetGuageByID(id string) (float64, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	db, err := sql.Open("postgres", s.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var fl float64
	row := db.QueryRow(`SELECT Value FROM gauge WHERE name = $1;`, id)
	err = row.Scan(&fl)
	if err != nil {
		return 0, err
	}
	return fl, nil
}

func (s DBStorage) GetCounterByID(id string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", s.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var val int64
	row := db.QueryRowContext(ctx, `SELECT Delta FROM counter WHERE name = $1;`, id)
	err = row.Scan(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (s DBStorage) SetGuage(id string, v *float64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", s.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.ExecContext(ctx,
		`INSERT INTO gauge (name, Value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET Value = $3`,
		id, Gauge(*v), Gauge(*v))
	if err != nil {
		log.Fatal(err)
	}
}

func (s DBStorage) SetCounter(id string, v *int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", s.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.ExecContext(ctx,
		`INSERT INTO counter (name, Delta) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET Delta = counter.Delta + $3`,
		id, Gauge(*v), Gauge(*v))

	if err != nil {
		log.Fatal(err)
	}
}

func (s DBStorage) SetGaugeFromMemStats(id string, value float64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", s.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.ExecContext(ctx,
		`INSERT INTO gauge (name, Value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET Value = $3`,
		id, Gauge(value), Gauge(value))
	if err != nil {
		log.Fatal(err)
	}
}

func (s DBStorage) SetCounterFromMemStats(id string, value int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", s.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.ExecContext(ctx,
		`INSERT INTO counter (name, Delta) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET Delta = counter.Delta + $3`,
		id, Counter(value), Counter(value))

	if err != nil {
		log.Fatal(err)
	}
}

func (s DBStorage) ReadMetrics() {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	s.SetGaugeFromMemStats("Alloc", float64(m.Alloc))
	s.SetGaugeFromMemStats("BuckHashSys", float64(m.BuckHashSys))
	s.SetGaugeFromMemStats("Frees", float64(m.Frees))
	s.SetGaugeFromMemStats("GCCPUFraction", m.GCCPUFraction)
	s.SetGaugeFromMemStats("GCSys", float64(m.GCSys))
	s.SetGaugeFromMemStats("HeapAlloc", float64(m.HeapAlloc))
	s.SetGaugeFromMemStats("HeapIdle", float64(m.HeapIdle))
	s.SetGaugeFromMemStats("HeapInuse", float64(m.HeapInuse))
	s.SetGaugeFromMemStats("HeapObjects", float64(m.HeapObjects))
	s.SetGaugeFromMemStats("HeapReleased", float64(m.HeapReleased))
	s.SetGaugeFromMemStats("HeapSys", float64(m.HeapSys))
	s.SetGaugeFromMemStats("LastGC", float64(m.HeapAlloc))
	s.SetGaugeFromMemStats("Lookups", float64(m.Lookups))
	s.SetGaugeFromMemStats("MCacheInuse", float64(m.MCacheInuse))
	s.SetGaugeFromMemStats("MCacheSys", float64(m.MCacheSys))
	s.SetGaugeFromMemStats("MSpanInuse", float64(m.MSpanInuse))
	s.SetGaugeFromMemStats("MSpanSys", float64(m.MSpanSys))
	s.SetGaugeFromMemStats("Mallocs", float64(m.Mallocs))
	s.SetGaugeFromMemStats("NextGC", float64(m.NextGC))
	s.SetGaugeFromMemStats("NumForcedGC", float64(m.NumForcedGC))
	s.SetGaugeFromMemStats("NumGC", float64(m.NumGC))
	s.SetGaugeFromMemStats("OtherSys", float64(m.OtherSys))
	s.SetGaugeFromMemStats("PauseTotalNs", float64(m.PauseTotalNs))
	s.SetGaugeFromMemStats("StackInuse", float64(m.StackInuse))
	s.SetGaugeFromMemStats("StackSys", float64(m.StackSys))
	s.SetGaugeFromMemStats("Sys", float64(m.Sys))
	s.SetGaugeFromMemStats("TotalAlloc", float64(m.TotalAlloc))
	s.SetGaugeFromMemStats("RandomValue", rand.Float64())
	s.SetCounterFromMemStats("PollCount", 1)
}

func (s DBStorage) GetGauges() map[string]Gauge {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", s.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.QueryContext(ctx, `SELECT * FROM gauge;`)
	if err != nil {
		log.Fatal(err)
	}
	var name string
	var value float64
	var id int
	for rows.Next() {
		err = rows.Scan(&id, &name, &value)
		if err != nil {
			log.Fatal(err)
		}
		s.GaugeField[name] = Gauge(value)
	}
	return s.GaugeField
}

func (s DBStorage) GetCounters() map[string]Counter {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", s.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.QueryContext(ctx, `SELECT * FROM counter;`)
	if err != nil {
		log.Fatal(err)
	}
	var name string
	var value int64
	var id int
	for rows.Next() {
		err = rows.Scan(&id, &name, &value)
		if err != nil {
			log.Fatal(err)
		}
		s.CounterField[name] = Counter(value)
	}
	return s.CounterField
}

func (s DBStorage) CheckGaugeByName(name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", s.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.QueryContext(ctx, `SELECT Value FROM gauge WHERE name = $1;`, name)
	if err != nil {
		return false
	}
	return true
}

func (s DBStorage) CheckCounterByName(name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", s.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.QueryContext(ctx, `SELECT Delta FROM counter WHERE name = $1;`, name)
	if err != nil {
		return false
	}
	return true
}

func (s DBStorage) SaveMany(metricsList []metric.Metrics) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", s.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	//
	stmtGauge, err := tx.PrepareContext(ctx,
		`INSERT INTO gauge (name, Value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET Value = $3`)
	if err != nil {
		return err
	}
	defer stmtGauge.Close()

	stmtCounter, err := tx.PrepareContext(ctx,
		`INSERT INTO counter (name, Delta) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET Delta = counter.Delta + $3`)
	if err != nil {
		return err
	}
	defer stmtCounter.Close()

	for _, metric := range metricsList {
		if metric.MType == "gauge" {
			_, err := stmtGauge.ExecContext(ctx, metric.ID, metric.Value, metric.Value)
			if err != nil {
				return err
			}
		} else {
			_, err := stmtCounter.ExecContext(ctx, metric.ID, metric.Delta, metric.Delta)
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}