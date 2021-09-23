package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"signingKey" json:"signingKey"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"name" json:"name"`
	Host        string        `mapstructure:"host" json:"host"`
	Port        int           `mapstructure:"port" json:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo     JwtConfig     `mapstructure:"jwtinfo" json:"jwtinfo"`
}
