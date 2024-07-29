package storage

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"homework-1/internal/cache"
	"homework-1/internal/constant"
	"homework-1/internal/model"
	"log"
)

var ErrFailedToBeginTransaction = errors.New("failed to begin transaction")

type Storage struct {
	Db    *sql.DB
	Cache *cache.Cache[uint64, *model.Order]
}

func New() (*Storage, error) {
	db, err := sql.Open("pgx", constant.DataBaseConnection)
	if err != nil {
		return nil, err
	}

	st := &Storage{
		Db:    db,
		Cache: cache.New[uint64, *model.Order](10),
	}
	return st, nil
}

func (s *Storage) CloseStorage() {
	err := s.Db.Close()
	if err != nil {
		log.Printf(err.Error())
	}
}

func (s *Storage) BeginTransaction() (*sql.Tx, error) {
	tx, err := s.Db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *Storage) MakeTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := s.Db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return ErrFailedToBeginTransaction
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
