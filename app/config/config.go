package config

type Configurations struct {
	Database DatabaseConfigurations
}

type DatabaseConfigurations struct {
	Hostname string
	Port     int
	DBName   string
	DBUser   string
	Password string
	Sslmode  string
}
