package config

type Configuration struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}

type Database struct {
	Maintable       string `yaml:"table"`
	Region          string `yaml:"region"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	SessionToken    string `yaml:"sessionToken"`
	EndpointURL     string `yaml:"endpointURL"`
}

type Server struct {
	Environment string `yaml:"environment"`
	Port        string `yaml:"port"`
	Host        string `yaml:"host"`
	LogPath     string `yaml:"logpath"`
	DBName      string `yaml:"dbname"`
}
