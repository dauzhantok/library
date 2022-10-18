package configs

import (
	"fmt"
)

// Mq RabbitMQ object
type Mq struct {
	Host     string `toml:"HOST" env:"RABBIT_CRM_HOST,required"`
	Port     int    `toml:"PORT" env:"RABBIT_CRM_PORT,required"`
	Username string `toml:"USERNAME" env:"RABBIT_CRM_USR,required"`
	Password string `toml:"PASSWORD" env:"RABBIT_CRM_PASSW,required"`
	Vhost    string `toml:"VHOST" env:"RABBIT_CRM_USR,required"`
}

// GetConnectionUrl returns RabbitMQ connection url
func (mq *Mq) GetConnectionUrl() string {
	if mq.Vhost != "" {
		return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", mq.Username, mq.Password, mq.Host, mq.Port, mq.Vhost)
	} else {
		return fmt.Sprintf("amqp://%s:%s@%s:%d", mq.Username, mq.Password, mq.Host, mq.Port)
	}
}
