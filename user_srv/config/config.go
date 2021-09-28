package config

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbName"`
}

type consulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type ServerConfig struct {
	ServiceName string       `mapstructure:"name"`
	Port        int          `mapstructure:"port"`
	Host        string       `mapstructure:"host"`
	MysqlInfo   MysqlConfig  `mapstructure:"mysql"`
	ConsulInfo  consulConfig `mapstructure:"consul"`
}
