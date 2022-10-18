package configs

import (
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

// Db database object
type Db struct {
	Host            string `toml:"HOST" env:"DB_HOST"`
	Port            int    `toml:"PORT" env:"DB_PORT"`
	DbName          string `toml:"NAME" env:"DB_NAME"`
	SSLMode         string `toml:"SSL_MODE" env:"DB_SSLMODE"`
	StaffUsername   string `toml:"USERNAME_STAFF" env:"DB_USERNAME_STAFF"`
	StaffPassword   string `toml:"PASSWORD_STAFF" env:"DB_PASSWORD_STAFF"`
	DsrUsername     string `toml:"USERNAME_DSR" env:"DB_USERNAME_DSR"`
	DsrPassword     string `toml:"PASSWORD_DSR" env:"DB_PASSWORD_DSR"`
	CommonsUsername string `toml:"USERNAME_COMMONS" env:"DB_USERNAME_COMMONS"`
	CommonsPassword string `toml:"PASSWORD_COMMONS" env:"DB_PASSWORD_COMMONS"`
	MaxIdleConns    int    `toml:"MAX_IDLE_CONNS" env:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns    int    `toml:"MAX_OPEN_CONNS" env:"DB_MAX_OPEN_CONNS"`
}

// ConfigureDB configures db
func (db *Db) ConfigureDB(userType string) (*sqlx.DB, error) {
	var databaseUrl string
	switch userType {
	case "staff":
		databaseUrl = db.GetStaffConnectionString()
	case "dsr":
		databaseUrl = db.GetDsrConnectionString()
	case "commons":
		databaseUrl = db.GetCommonsConnectionString()
	default:
		return nil, errors.New("db name not found")
	}
	database, err := sqlx.Open("pgx", databaseUrl)
	if err != nil {
		return nil, err
	}
	if err = database.Ping(); err != nil {
		return nil, err
	}
	database.SetMaxIdleConns(db.MaxIdleConns)
	database.SetMaxOpenConns(db.MaxOpenConns)
	return database, nil
}

func (db *Db) GetStaffConnectionString() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s user=%s password=%s",
		db.Host,
		db.Port,
		db.DbName,
		db.SSLMode,
		db.StaffUsername,
		db.StaffPassword)
}

func (db *Db) GetDsrConnectionString() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s user=%s password=%s",
		db.Host,
		db.Port,
		db.DbName,
		db.SSLMode,
		db.DsrUsername,
		db.DsrPassword)
}

func (db *Db) GetCommonsConnectionString() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s user=%s password=%s",
		db.Host,
		db.Port,
		db.DbName,
		db.SSLMode,
		db.DsrUsername,
		db.DsrPassword)
}
