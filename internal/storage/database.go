package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/dayterr/go_agent_metrics/internal/metric"
)

func NewDB(dsn string) (DBStorage, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return DBStorage{}, err
	}
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
		DB:           db,
		DSN:          dsn,
		GaugeField:   make(map[string]Gauge),
		CounterField: make(map[string]Counter),
	}, nil
}

func (s DBStorage) GetGuageByID(ctx context.Context, id string) (float64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		var fl float64
		row := s.DB.QueryRowContext(ctx, `SELECT Value FROM gauge WHERE name = $1;`, id)
		err := row.Scan(&fl)
		if err != nil {
			return 0, err
		}
		return fl, nil
	}
}

func (s DBStorage) GetCounterByID(ctx context.Context, id string) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		var val int64
		row := s.DB.QueryRowContext(ctx, `SELECT Delta FROM counter WHERE name = $1;`, id)
		err := row.Scan(&val)
		if err != nil {
			return 0, err
		}
		return val, nil
	}
}

func (s DBStorage) SetGuage(ctx context.Context, id string, v *float64) {
	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
	default:
		_, err := s.DB.ExecContext(ctx,
			`INSERT INTO gauge (name, Value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET Value = $3`,
			id, Gauge(*v), Gauge(*v))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s DBStorage) SetCounter(ctx context.Context, id string, v *int64) {
	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
	default:
		_, err := s.DB.ExecContext(ctx,
			`INSERT INTO counter (name, Delta) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET Delta = counter.Delta + $3`,
			id, Gauge(*v), Gauge(*v))

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s DBStorage) SetGaugeFromMemStats(ctx context.Context, id string, value float64) {
	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
	default:
		_, err := s.DB.ExecContext(ctx,
			`INSERT INTO gauge (name, Value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET Value = $3`,
			id, Gauge(value), Gauge(value))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s DBStorage) SetCounterFromMemStats(ctx context.Context, id string, value int64) {
	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
	default:
		_, err := s.DB.ExecContext(ctx,
			`INSERT INTO counter (name, Delta) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET Delta = counter.Delta + $3`,
			id, Counter(value), Counter(value))

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s DBStorage) GetGauges(ctx context.Context) (map[string]Gauge, error) {
	select {
	case <-ctx.Done():
		return map[string]Gauge{}, ctx.Err()
	default:
		rows, err := s.DB.QueryContext(ctx, `SELECT * FROM gauge;`)

		if err != nil {
			log.Fatal(err)
		}

		if rows.Err() != nil {
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
		return s.GaugeField, nil
	}
}

func (s DBStorage) GetCounters(ctx context.Context) (map[string]Counter, error) {
	select {
	case <-ctx.Done():
		return map[string]Counter{}, ctx.Err()
	default:
		rows, err := s.DB.QueryContext(ctx, `SELECT * FROM counter;`)
		if err != nil {
			log.Fatal(err)
		}

		if rows.Err() != nil {
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
		return s.CounterField, nil
	}
}

func (s DBStorage) CheckGaugeByName(ctx context.Context, name string) bool {
	select {
	case <-ctx.Done():
		return ctx.Err() == nil
	default:
		row, err := s.DB.QueryContext(ctx, `SELECT Value FROM gauge WHERE name = $1;`, name)
		if row.Err() != nil {
			log.Fatal(err)
		}

		return err == nil
	}
}

func (s DBStorage) CheckCounterByName(ctx context.Context, name string) bool {
	select {
	case <-ctx.Done():
		return ctx.Err() == nil
	default:
		row, err := s.DB.QueryContext(ctx, `SELECT Delta FROM counter WHERE name = $1;`, name)
		if row.Err() != nil {
			log.Fatal(err)
		}

		return err == nil
	}
}

func (s DBStorage) SaveMany(ctx context.Context, metricsList []metric.Metrics) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		tx, err := s.DB.Begin()
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
}
