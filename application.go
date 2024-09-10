package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

type DbConfig struct {
	Name      string
	Host      string
	Port      int
	User      string
	Password  string
	Container bool
}

type Application struct {
	context  context.Context
	Postgres DbConfig
}

// DbConnURL returns string URL, which may be used for connect to postgres database.
func (c *DbConfig) ConnURL() string {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
	return url
}

func (a *Application) PostgresConnect(db *DbConfig) error {
	conn, err := pgx.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	var name string
	var weight int64
	err = conn.QueryRow("select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(name, weight)
	return nil
}

func (a *Application) MysqlConnect(db *DbConfig) error {
	return nil
}

func (a *Application) SQLiteConnect(db *DbConfig) error {
	return nil
}

func (a *Application) PostgresRestore(db *DbConfig) error {
	return nil
}

func (a *Application) MysqlRestore(db *DbConfig) error {
	return nil
}

func (a *Application) SQLiteRestore(db *DbConfig) error {
	return nil
}

func (a *Application) PostgresBackup(db *DbConfig) error {
	return nil
}
func (a *Application) MysqlBackup(db *DbConfig) error {
	return nil
}
func (a *Application) SQLiteBackup(db *DbConfig) error {
	return nil
}
