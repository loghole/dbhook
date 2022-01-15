package main

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/loghole/dbhook"
)

const (
	dsn        = "postgresql://root@localhost:29999/defaultdb?sslmode=disable"
	driverName = "postgres-test"
)

type Hook struct{}

func (h *Hook) Before(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	log.Println(input.Caller, input.Query)

	return ctx, nil
}

func (h *Hook) After(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	log.Println(input.Caller, input.Query)

	return ctx, nil
}

func (h *Hook) Error(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	log.Println(input.Caller, input.Query)

	return ctx, input.Error
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	hooks := dbhook.NewHooks(dbhook.WithHook(&Hook{}))

	sql.Register(driverName, dbhook.Wrap(&pq.Driver{}, hooks))

	sqlxSQL()

	log.Println("all ok")
}

type TestToken struct {
	Token string `db:"token"`
}

func sqlxSQL() {
	dbSTD, err := sql.Open(driverName, dsn)
	if err != nil {
		log.Fatal(err)
	}

	db := sqlx.NewDb(dbSTD, "postgres")

	rows, err := db.QueryxContext(context.Background(), `SELECT token FROM tokens limit 2`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}

		log.Println("id", id)
	}

	if _, err = db.ExecContext(
		context.Background(),
		`INSERT INTO tokens(token) VALUES($1)`,
		strconv.FormatInt(time.Now().Unix(), 32),
	); err != nil {
		log.Panicln(err)
	}

	tt := TestToken{Token: strconv.FormatInt(time.Now().Unix(), 30)}

	if _, err := db.NamedExecContext(context.Background(), `INSERT INTO tokens(token) VALUES(:token)`, tt); err != nil {
		log.Panicln(err)
	}

	stmt, err := db.PrepareContext(context.Background(), `INSERT INTO tokens(token) VALUES($1)`)
	if err != nil {
		log.Panicln(err)
	}

	if _, err := stmt.Exec(strconv.FormatInt(time.Now().Unix(), 29)); err != nil {
		log.Panicln(err)
	}

	if _, err := stmt.Exec(strconv.FormatInt(time.Now().Unix(), 28)); err != nil {
		log.Panicln(err)
	}

	if err := stmt.Close(); err != nil {
		log.Panicln(err)
	}

	tx, err := db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Panicln(err)
	}

	if _, err = tx.ExecContext(
		context.Background(),
		`INSERT INTO tokens(token) values($1)`,
		strconv.FormatInt(time.Now().UnixNano(), 32),
	); err != nil {
		log.Panicln(err)
	}

	if err := tx.Commit(); err != nil {
		log.Panicln(err)
	}
}
