package storage

import (
	"context"
	"database/sql"
	"encoding/json"
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
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS gauge (ID text, Value double precision);`)
	if err != nil {
		return DBStorage{}, err
	}
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS counter (ID text, Delta BIGINT);`)
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
	return nil
}

func (s DBStorage) GetGuageByID(id string) float64 {
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	db, err := sql.Open("postgres", s.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var fl float64
	row := db.QueryRow(`SELECT Value FROM gauge WHERE id = $1;`, id)
	log.Println("row is", row)
	err = row.Scan(&fl)
	if err != nil {
		log.Fatal(err)
	}
	return fl
}

func (s DBStorage) GetCounterByID(id string) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", s.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var val int64
	row := db.QueryRowContext(ctx, `SELECT Delta FROM counter WHERE ID = $1;`, id)
	err = row.Scan(&val)
	if err != nil {
		log.Fatal(err)
	}
	return val
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
		`INSERT INTO gauge (ID, Value) VALUES ($1, $2) ON CONFLICT(ID) DO UPDATE SET Value = $3`,
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
		`INSERT INTO counter (ID, Delta) VALUES ($1, $2) ON CONFLICT(ID) DO UPDATE SET Value = counter.Value = $3`,
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
		`INSERT INTO gauge (ID, Value) VALUES ($1, $2) ON CONFLICT(ID) DO UPDATE SET Value = $3`,
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
		`INSERT INTO counter (ID, Delta) VALUES ($1, $2) ON CONFLICT(ID) DO UPDATE SET Value = counter.Value = $3`,
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
	rows, err := db.QueryContext(ctx, `SELECT (*) FROM counter;`)
	if err != nil {
		log.Fatal(err)
	}
	var name string
	var value float64
	for rows.Next() {
		err = rows.Scan(&name, &value)
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
	rows, err := db.QueryContext(ctx, `SELECT (*) FROM counter;`)
	if err != nil {
		log.Fatal(err)
	}
	var name string
	var value int64
	for rows.Next() {
		err = rows.Scan(&name, &value)
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
	_, err = db.QueryContext(ctx, `SELECT Value FROM gauge WHERE ID = $1;`, name)
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
	_, err = db.QueryContext(ctx, `SELECT Delta FROM counter WHERE ID = $1;`, name)
	if err != nil {
		return false
	}
	return true
}
