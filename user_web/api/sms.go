package api

import (
	"context"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"project/user_web/forms"
	"project/user_web/global"
	"strings"
	"time"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient (accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func GenerateSmsCode(witdh int) string {
	//生成width长度的短信验证码

	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendSms(c *gin.Context)  {
	SendSmsForm:=forms.SendSmsForm{}
	if err := c.ShouldBind(&SendSmsForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	//ServerConfig
	client, err := CreateClient(tea.String(global.ServerConfig.AliyunInfo.AccessKeyId), tea.String(global.ServerConfig.AliyunInfo.AccessKeySecret))
	if err != nil {
		panic(err)
	}
	smsCode := GenerateSmsCode(6)

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers: tea.String(SendSmsForm.Mobile),
		SignName: tea.String(global.ServerConfig.AliyunInfo.SignName),
		TemplateCode: tea.String(global.ServerConfig.AliyunInfo.TemplateCode),
		TemplateParam: tea.String("{\"code\":"+smsCode+"}"),
	}
	// 复制代码运行请自行打印 API 的返回值
	_, err = client.SendSms(sendSmsRequest)

	rerr := global.Rdb.Set(context.Background(), SendSmsForm.Mobile, smsCode, 300*time.Second).Err()
	if rerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":"服务器內部出错",
		})
		zap.S().Error("保存redis出错:", rerr.Error())
		return
	}else {
		c.JSON(http.StatusOK, gin.H{
			"msg":"发送成功",
		})
	}
}
