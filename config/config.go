package config

type Configuration struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}

type Database struct {
	Username string `yaml:"DB_USERNAME"`
	Password string `yaml:"DB_PASSWORD"`
	Name     string `yaml:"DB_NAME"`
	Host     string `yaml:"DB_HOST"`
	Port     string `yaml:"DB_PORT"`
}

type Server struct {
	Environment string `yaml:"environment"`
	Port        string `yaml:"port"`
	Host        string `yaml:"host"`
	LogPath     string `yaml:"logpath"`
	DBName      string `yaml:"dbname"`
}
