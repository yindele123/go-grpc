package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
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

type ServerConfig struct {
	Name        string        `mapstructure:"name" json:"name"`
	Host        string        `mapstructure:"host" json:"host"`
	Port        int           `mapstructure:"port" json:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo     JwtConfig     `mapstructure:"jwtinfo" json:"jwtinfo"`
	AliyunInfo  AliyunConfig  `mapstructure:"aliyunInfo" json:"aliyunInfo"`
	RedisInfo   RedisConfig   `mapstructure:"redisInfo" json:"redisInfo"`
}
