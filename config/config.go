package config

type Configuration struct {
	Server struct {
		Environment string `yaml:"environment"`
		Port        string `yaml:"port"`
		Host        string `yaml:"host"`
	} `yaml:"server"`
}
