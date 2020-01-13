package formats

import "fmt"

const (
	formatPostgres = "host=%s port=%d dbname=%s user=%s password='%s' sslmode='%s' connect_timeout='%d'"
	formatAMQP     = "amqp://%s:%s@%s:%d?connection_timeout=%d/"
)

type (
	ParamsPostgres struct {
		Host       string
		Port       uint32
		DbName     string
		DbUser     string
		DbPassword string
		DbSsl      string
		Timeout    uint32
	}

	ParamsAMQP struct {
		Host     string
		Port     uint32
		User     string
		Password string
		Timeout  uint32
	}
)

func Postgres(p ParamsPostgres) string {
	return fmt.Sprintf(formatPostgres,
		p.Host,
		p.Port,
		p.DbName,
		p.DbUser,
		p.DbPassword,
		p.DbSsl,
		p.Timeout,
	)
}

func AMQP(p ParamsAMQP) string {
	return fmt.Sprintf(formatAMQP,
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Timeout,
	)
}
