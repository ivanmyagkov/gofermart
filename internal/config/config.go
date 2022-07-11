package config

var Flags struct {
	A string
	D string
	R string
}

var EnvVar struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:":8080"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI" envDefault:"postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable"`
}

type Config struct {
	RunAddress           string
	AccrualSystemAddress string
	DatabaseURI          string
}

func (c Config) GetRunAddress() string {
	return c.RunAddress
}

func (c Config) GetAccrualSystemAddress() string {
	return c.AccrualSystemAddress
}
func (c Config) GetDatabaseURI() string {
	return c.DatabaseURI
}

func NewConfig(runAddress, database, accrualAddress string) *Config {
	return &Config{
		RunAddress:           runAddress,
		DatabaseURI:          database,
		AccrualSystemAddress: accrualAddress,
	}
}

var TokenKey = []byte("secret")
