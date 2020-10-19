package config

type Config struct {

}

type Database struct {
	Username string `env:"DATABASE_USERNAME"`
	Password string `env:"DATABASE_PASSWORD"`
	Hostname string `env:"DATABASE_HOSTNAME"`
	DBname string `env:"DATABASE_DBNAME"`
}
