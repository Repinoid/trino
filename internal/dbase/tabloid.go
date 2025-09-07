package dbase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"triner/internal/models"
)

// Структура для базы данных.
type DBstruct struct {
	DB *pgxpool.Pool
	//	DB *pgx.Conn
}

func NewPostgresPool(ctx context.Context, DSN string) (*DBstruct, error) {

	poolConfig, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	return &DBstruct{DB: pool}, nil
}

// (p *pgxpool.Pool) Close()
// Close closes all connections in the pool and rejects future Acquire calls.
// Blocks until all connections are returned to pool and closed.
func (dataBase *DBstruct) Close() {
	dataBase.DB.Close()
}

func Ping(ctx context.Context, DSN string) error {
	dataBase, err := NewPostgresPool(ctx, DSN)
	if err != nil {
		return err
	}
	defer dataBase.DB.Close()

	err = dataBase.DB.Ping(ctx) // база то открыта ...
	if err != nil {
		models.Logger.Error("No PING ", "error", err.Error())
		return fmt.Errorf("no ping %w", err)
	}
	return nil
}

func (dataBase *DBstruct) CreateTable(ctx context.Context) (err error) {

	order := `
		CREATE TABLE IF NOT EXISTS tabl (
    		id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    		name VARCHAR(64) NOT NULL
		);    
	`
	_, err = dataBase.DB.Exec(ctx, order)

	return
}

func AddNameToTable(ctx context.Context, db *sql.DB, name string) (err error) {
	order := "INSERT INTO postgres.public.tabl (name) VALUES (?)"
	_, err = db.Exec(order, name)

	return
}

type Names struct {
	Id int
	Name string
}

func GetNamesFromTable(ctx context.Context, db *sql.DB) (nms []Names, err error) {
	order := "SELECT id, name from postgres.public.tabl"

	rows, err := db.Query(order)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//nms := []Names{}
	for rows.Next() {
		nm := Names{}
		
		err = rows.Scan(&nm.Id, &nm.Name)
		if err != nil {
			return nil, err
		}
		nms = append(nms, nm)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return
}
