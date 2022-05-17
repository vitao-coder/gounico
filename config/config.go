package config

type Configuration struct {
	Server    Server    `yaml:"server"`
	Worker    Listener  `yaml:"listener"`
	Database  Database  `yaml:"database"`
	Messaging Messaging `yaml:"messaging"`
	Telemetry Telemetry `yaml:"telemetry"`
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
}

type Listener struct {
	Environment string `yaml:"environment"`
	Port        string `yaml:"port"`
	Host        string `yaml:"host"`
}

type Messaging struct {
	BrokerURL      string                    `yaml:"brokerURL"`
	ConsumerLimit  int                       `yaml:"channelConsumerLimit"`
	Configurations []MessagingConfigurations `yaml:"configurations"`
}

type MessagingConfigurations struct {
	Topic      string `yaml:"topic"`
	Subscriber string `yaml:"subscriber"`
	URL        string `yaml:"url"`
}

type Telemetry struct {
	JaegerEndpoint string `yaml:"jaegerEndpoint"`
	AppName        string `yaml:"appName"`
}
