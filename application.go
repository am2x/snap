package main

import "context"

type DbConfig struct {
	Name      string
	Host      string
	Port      int
	User      string
	Password  string
	Container bool
}

type Application struct {
	context context.Context
}

func (a *Application) PostgresConnect(db *DbConfig) error {
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
