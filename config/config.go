package config

// Config contains server configs
type Config struct {
	DBConfig
	ServerPort int
}

// DBConfig contains DB configs
type DBConfig struct {
	Host     string
	Port     string
	Password string
}
