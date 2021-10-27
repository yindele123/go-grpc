package config

type MysqlConfig struct {
	Host        string `mapstructure:"host"  json:"host"`
	Port        int    `mapstructure:"port"   json:"port"`
	User        string `mapstructure:"user"  json:"user"`
	Password    string `mapstructure:"password"  json:"password"`
	DbName      string `mapstructure:"dbName"  json:"dbName"`
	TablePrefix string `mapstructure:"tablePrefix"  json:"tablePrefix"`
}

type consulConfig struct {
	Host string `mapstructure:"host"  json:"host"`
	Port int    `mapstructure:"port"  json:"port"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	ServiceName string       `mapstructure:"name"  json:"name"`
	Port        int          `mapstructure:"port"  json:"port"`
	Host        string       `mapstructure:"host"  json:"host"`
	MysqlInfo   MysqlConfig  `mapstructure:"mysql"  json:"mysql"`
	RedisInfo   RedisConfig  `mapstructure:"redisInfo" json:"redisInfo"`
	ConsulInfo  consulConfig `mapstructure:"consul"  json:"consul"`
}

type nacosInfo struct {
	Host        string `mapstructure:"host"`
	Port        uint64 `mapstructure:"port"`
	NamespaceId string `mapstructure:"namespaceId"`
	DataId      string `mapstructure:"dataId"`
	Group       string `mapstructure:"group"`
}

type NacosConfig struct {
	NacosInfo nacosInfo `mapstructure:"nacos"`
}
