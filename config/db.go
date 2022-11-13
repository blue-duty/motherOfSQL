package config

type DB struct {
	Type     string `yaml:"type"`
	Host     string
	Port     int
	User     string
	Password string
}
