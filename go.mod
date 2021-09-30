module project

go 1.16

require (
	github.com/alibabacloud-go/darabonba-openapi v0.1.7
	github.com/alibabacloud-go/dysmsapi-20170525/v2 v2.0.2
	github.com/alibabacloud-go/tea v1.1.17
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/v8 v8.11.3
	github.com/hashicorp/consul/api v1.3.0
	github.com/mbobakov/grpc-consul-resolver v1.4.4
	github.com/mojocn/base64Captcha v1.3.5
	github.com/nacos-group/nacos-sdk-go v1.0.9
	github.com/spf13/viper v1.8.1
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.14
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.40.0
