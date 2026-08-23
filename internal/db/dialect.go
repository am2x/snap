package db

import "strings"

type Dialect interface {
	QuoteIdentifier(name string) string
}

type AnsiDialect struct{}

func (d AnsiDialect) QuoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

type MysqlDialect struct{}

func (d MysqlDialect) QuoteIdentifier(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

type SqlDialect struct{}

func (d SqlDialect) QuoteIdentifier(name string) string {
	return "[" + strings.ReplaceAll(name, "]", "]]") + "]"
}

func GetDialect(driverName string) Dialect {
	switch driverName {
	case "mysql", "mariadb":
		return MysqlDialect{}
	case "postgresql":
		return AnsiDialect{}
	case "sqlserver":
		return SqlDialect{}
	default:
		return AnsiDialect{}
	}
}
