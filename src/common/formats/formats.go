package formats

import "fmt"

const formatPostgres = "host=%s port=%d dbname=%s user=%s password='%s' sslmode='%s' connect_timeout='%d'"

func Postgres(host string, port uint32, dbName, dbUser, dbPassword, dbSsl string, timeout uint32) string {
	return fmt.Sprintf(formatPostgres,
		host,
		port,
		dbName,
		dbUser,
		dbPassword,
		dbSsl,
		timeout,
	)
}

const formatAMQP = "amqp://%s:%s@%s:%d?connection_timeout=%d/"

func AMQP(host string, port uint32, user, password string, timeout uint32) string {
	return fmt.Sprintf(formatAMQP,
		user,
		password,
		host,
		port,
		timeout,
	)
}