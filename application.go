package main

import (
	"context"
)

type ServerConfig struct {
	// DbType    Database
	Name      string
	Host      string
	Port      int
	User      string
	Password  string
	Container bool
}

type Application struct {
	context context.Context
	server  ServerConfig
}

const (
	MysqlPrefix    string = "mysql"
	PostgresPrefix string = "psql"
)

func (a *Application) SSHBackup() error {
	return nil
}

func (a *Application) SSHRestore() error {
	return nil
}

func LocalBackup(a *Application) error {
	return nil
}

func (a *Application) LocalRestore() error {

	return nil
}
