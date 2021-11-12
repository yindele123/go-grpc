package config

type SrvInfo struct {
	Name string `mapstructure:"name" json:"name"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"signingKey" json:"signingKey"`
}

type AliyunConfig struct {
	AccessKeyId     string `mapstructure:"accessKeyId" json:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret" json:"accessKeySecret"`
	SignName        string `mapstructure:"SignName" json:"SignName"`
	TemplateCode    string `mapstructure:"TemplateCode" json:"TemplateCode"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type consulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name          string       `mapstructure:"name" json:"name"`
	Host          string       `mapstructure:"host" json:"host"`
	Port          int          `mapstructure:"port" json:"port"`
	GoodsSrvInfo  SrvInfo      `mapstructure:"goodsSrv" json:"goodsSrv"`
	UserOpSrvInfo SrvInfo      `mapstructure:"userOpSrv" json:"userOpSrv"`
	JWTInfo       JwtConfig    `mapstructure:"jwtinfo" json:"jwtinfo"`
	AliyunInfo    AliyunConfig `mapstructure:"aliyunInfo" json:"aliyunInfo"`
	RedisInfo     RedisConfig  `mapstructure:"redisInfo" json:"redisInfo"`
	ConsulInfo    consulConfig `mapstructure:"consul" json:"consul"`
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
