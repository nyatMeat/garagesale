package database

import (
	"context"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
)

// COnfig for database connection
type Config struct {
	Host       string
	Name       string
	User       string
	Password   string
	DisableTLS bool
}

// Open knows how to open a database connection.
func Open(conf Config) (*sqlx.DB, error) {
	// Query parameters.
	q := make(url.Values)
	q.Set("sslmode", "require")
	if conf.DisableTLS {
		q.Set("sslmode", "disable")
	}
	q.Set("timezone", "utc")

	// Construct url.
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(conf.User, conf.Password),
		Host:     conf.Host,
		Path:     "postgres",
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}

func StatusCheck(ctxt context.Context, db *sqlx.DB) error {

	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctxt, q).Scan(&tmp)
}
