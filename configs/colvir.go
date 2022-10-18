package configs

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

// Colvir colvir object
type Colvir struct {
	Host         string `toml:"HOST" env:"DB_HOST_COLVIR,required"`
	Port         int    `toml:"PORT" env:"DB_PORT_COLVIR,required"`
	DbName       string `toml:"NAME" env:"DB_NAME_COLVIR,required"`
	Username     string `toml:"USERNAME" env:"DB_USERNAME_COLVIR,required"`
	Password     string `toml:"PASSWORD" env:"DB_PASSWORD_COLVIR,required"`
	TNSString    string `toml:"TNSSTRING" env:"DB_TNSSTRING_COLVIR,required"`
	MaxIdleConns int    `toml:"MAX_IDLE_CONNS" env:"DB_MAX_IDLE_CONNS,required"`
	MaxOpenConns int    `toml:"MAX_OPEN_CONNS" env:"DB_MAX_OPEN_CONNS,required"`
}

// ConfigureColvir configures colvir
func (colvir *Colvir) ConfigureColvir() (*sqlx.DB, error) {
	connectionString := colvir.GetConnectionString()
	db, err := sqlx.Open("godror", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(colvir.MaxIdleConns)
	db.SetMaxOpenConns(colvir.MaxOpenConns)
	db.SetConnMaxLifetime(1 * time.Hour)
	return db, nil
}

// GetConnectionString returns colvir connection string
func (colvir *Colvir) GetConnectionString() string {
	return fmt.Sprintf(`user="%s" password="%s" connectString="%s"`,
		colvir.Username,
		colvir.Password,
		colvir.TNSString)
}
