package persistence

import (
	"encoding/json"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var (
	ErrOrderNotFound = errors.New("order_not_found")
	ErrInvalidOrder  = errors.New("invalid_order")
	ErrOrderExists   = errors.New("order_already_exists")
)

const ordersBucket = "orders"

type Store struct {
	db *bolt.DB
}

func Open(path string) (*Store, error) {
	db, err := bolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists([]byte(ordersBucket))
		return e
	})
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) CreateOrder(rec OrderRecord) error {
	if err := validateNewOrder(rec); err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ordersBucket))
		if b == nil {
			return fmt.Errorf("missing_bucket:%s", ordersBucket)
		}

		key := []byte(rec.EngineID)
		if existing := b.Get(key); existing != nil {
			return ErrOrderExists
		}

		payload, err := json.Marshal(rec)
		if err != nil {
			return err
		}

		return b.Put(key, payload)
	})
}

func (s *Store) GetOrder(engineID string) (OrderRecord, error) {
	var out OrderRecord

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ordersBucket))
		if b == nil {
			return fmt.Errorf("missing_bucket:%s", ordersBucket)
		}

		raw := b.Get([]byte(engineID))
		if raw == nil {
			return ErrOrderNotFound
		}

		return json.Unmarshal(raw, &out)
	})

	return out, err
}

func (s *Store) ListOrders() ([]OrderRecord, error) {
	orders := make([]OrderRecord, 0)

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ordersBucket))
		if b == nil {
			return fmt.Errorf("missing_bucket:%s", ordersBucket)
		}

		return b.ForEach(func(_, v []byte) error {
			var rec OrderRecord
			if err := json.Unmarshal(v, &rec); err != nil {
				return err
			}
			orders = append(orders, rec)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Store) UpdateOrderState(engineID string, nextState OrderState, updatedAt int64, exchangeID string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ordersBucket))
		if b == nil {
			return fmt.Errorf("missing_bucket:%s", ordersBucket)
		}

		key := []byte(engineID)
		raw := b.Get(key)
		if raw == nil {
			return ErrOrderNotFound
		}

		var rec OrderRecord
		if err := json.Unmarshal(raw, &rec); err != nil {
			return err
		}

		if err := ValidateTransition(rec.State, nextState); err != nil {
			return err
		}

		rec.State = nextState
		rec.UpdatedAt = updatedAt

		if exchangeID != "" {
			rec.ExchangeID = exchangeID
		}

		payload, err := json.Marshal(rec)
		if err != nil {
			return err
		}

		return b.Put(key, payload)
	})
}

func validateNewOrder(rec OrderRecord) error {
	if rec.EngineID == "" {
		return ErrInvalidOrder
	}
	if rec.Pair == "" {
		return ErrInvalidOrder
	}
	if rec.Amount <= 0 {
		return ErrInvalidOrder
	}
	if rec.State == "" {
		return ErrInvalidOrder
	}
	if rec.State != OrderStatePending {
		return fmt.Errorf("initial_state_must_be_pending")
	}
	return nil
}
