package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"project/user_web/forms"
	"project/user_web/global"
	"project/user_web/global/reponse"
	"project/user_web/middlewares"
	"project/user_web/models"
	"project/user_web/proto"
	"strconv"
	"strings"
	"time"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": st.Message(),
				})
			case codes.AlreadyExists:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": st.Message(),
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": st.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg:": "内部错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			}
			return
		}
	}
	return

}

/*func HandleGrpcErrorInfoToHttp(err error, c *gin.Context) {
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": st.Message(),
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": st.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg:": "内部错误",
				})
			case codes.AlreadyExists:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": st.Message(),
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			}
			return
		}
	}
	return

}*/

func SetGender(gender uint32) string {
	var genderName string
	switch gender {
	case 1:
		genderName = "女"
	case 2:
		genderName = "男"
	default:
		genderName = "保密"
	}
	return genderName
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return

}
func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func GetUserList(c *gin.Context) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 链接grpc失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	pn := c.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := c.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	userServer := proto.NewUserClient(conn)
	res, err := userServer.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range res.Data {
		user := reponse.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Mobile:   value.Mobile,
			Birthday: reponse.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   SetGender(value.Gender),
		}

		result = append(result, user)
	}
	c.JSON(http.StatusOK, result)

}

func PassWordLogin(c *gin.Context) {
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest,gin.H{
			"msg":"请输入正确的验证码",
		})
		return
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 链接grpc失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	userServer := proto.NewUserClient(conn)
	if userMobileRes, err := userServer.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		zap.S().Errorw("[PassWordLogin] 登录用户失败")
		if e, ok := status.FromError(err); ok {
			fmt.Print(e.Code())
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "用户不存在",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "登录失败",
				})
			}
			return
		}
	} else {
		if pasRes, pasErr := userServer.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: userMobileRes.PassWord,
		}); pasErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "密码错误",
			})
		} else {
			if pasRes.Success {
				j:=middlewares.NewJWT()
				laims := models.CustomClaims{
					ID:             uint(userMobileRes.Id),
					NickName:       userMobileRes.NickName,
					AuthorityId:    uint(userMobileRes.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(), //签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
						Issuer: "xindele",
					},
				}
				token,jwtErr:=j.CreateToken(laims)

				if jwtErr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg":"生成token失败",
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"id": userMobileRes.Id,
					"nick_name": userMobileRes.NickName,
					"token": token,
					"expired_at": (time.Now().Unix() + 60*60*24*30)*1000,
					"msg": "登录成功",
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "密码错误",
				})
			}
		}
	}
	/*res, err := userServer.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
		Password:          uint32(pnInt),
		EncryptedPassword: uint32(pSizeInt),
	})*/
}
func Register(c *gin.Context)  {
	registerForm:=forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	value,err:=global.Rdb.Get(context.Background(),registerForm.Mobile).Result()
	if err == redis.Nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":"验证码错误",
		})
		return
	}else {
		if value != registerForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":"验证码错误",
			})
			return
		}
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 链接grpc失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	userServer := proto.NewUserClient(conn)
	userRes, userErr := userServer.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})
	if userErr!=nil {
		zap.S().Errorf("[Register] 查询 【新建用户失败】失败: %s", userErr.Error())
		HandleGrpcErrorToHttp(userErr, c)
		return
	}

	j:=middlewares.NewJWT()

	laims := models.CustomClaims{
		ID:             uint(userRes.Id),
		NickName:       userRes.NickName,
		AuthorityId:    uint(userRes.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(), //签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
			Issuer: "xindele",
		},
	}
	token,jwtErr:=j.CreateToken(laims)

	if jwtErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":"生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"id": userRes.Id,
		"nick_name": userRes.NickName,
		"token": token,
		"msg":"注册成功",
		"expired_at": (time.Now().Unix() + 60*60*24*30)*1000,
	})

}
