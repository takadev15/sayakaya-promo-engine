package sql

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type ConnectionOption struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type psqlClient struct {
    db *sqlx.DB
}

func NewPostgresClient (opts ConnectionOption) (*psqlClient, error) {
    dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", opts.User, opts.Password, opts.Host, opts.Port, opts.Database)
    fmt.Println(opts.Database)
    fmt.Println(dbUrl)
    dbx, err := sqlx.Connect("pgx", dbUrl)
    if err != nil {
        return nil, err
    }
    return &psqlClient{
        db: dbx,
    }, nil
}

func (c *psqlClient) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.db.GetContext(ctx, dest, query, args...)
}

func (c *psqlClient) Exec(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

func (c *psqlClient) Query(ctx context.Context, query string, args ...interface{}) Row {
	return c.db.QueryRowContext(ctx, query, args...)
}
