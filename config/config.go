package config

type Configuration struct {
	Server struct {
		Environment string `yaml:"environment"`
		Port        string `yaml:"port"`
		Host        string `yaml:"host"`
		LogPath     string `yaml:"logpath"`
		DBName      string `yaml:"dbname"`
	} `yaml:"server"`
}
