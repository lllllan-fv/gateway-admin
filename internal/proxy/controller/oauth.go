package controller

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/conf"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/dao"
	"github.com/lllllan-fv/gateway-admin/public/consts"
	jjwt "github.com/lllllan-fv/gateway-admin/public/jwt"
	"github.com/lllllan-fv/gateway-admin/public/resp"
)

type TokensInput struct {
	GrantType string `json:"grant_type" form:"grant_type" comment:"授权类型" example:"client_credentials" validate:"required"` //授权类型
	Scope     string `json:"scope" form:"scope" comment:"权限范围" example:"read_write" validate:"required"`                   //权限范围
}

type TokensOutput struct {
	AccessToken string `json:"access_token" form:"access_token"` //access_token
	ExpiresIn   int    `json:"expires_in" form:"expires_in"`     //expires_in
	TokenType   string `json:"token_type" form:"token_type"`     //token_type
	Scope       string `json:"scope" form:"scope"`               //scope
}

func Tokens(c *gin.Context) {
	params := &TokensInput{}
	if err := c.ShouldBind(&params); err != nil {
		resp.Error(c, 2000, err)
		return
	}

	splits := strings.Split(c.GetHeader("Authorization"), " ")
	if len(splits) != 2 {
		resp.Error(c, 2001, errors.New("用户名或密码格式错误"))
		return
	}

	appSecret, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		resp.Error(c, 2002, err)
		return
	}
	//fmt.Println("appSecret", string(appSecret))

	//  取出 app_id secret
	//  生成 app_list
	//  匹配 app_id
	//  基于 jwt生成token
	//  生成 output
	parts := strings.Split(string(appSecret), ":")
	if len(parts) != 2 {
		resp.Error(c, 2003, errors.New("用户名或密码格式错误"))
		return
	}

	appList := dao.ListApp()
	for _, appInfo := range appList {
		if appInfo.AppID == parts[0] && appInfo.Secret == parts[1] {
			claims := jwt.StandardClaims{
				Issuer:    appInfo.AppID,
				ExpiresAt: time.Now().Add(consts.JwtExpires * time.Second).In(conf.TimeLocation).Unix(),
			}
			token, err := jjwt.Encode(claims)
			if err != nil {
				resp.Error(c, 2004, err)
				return
			}
			output := &TokensOutput{
				ExpiresIn:   consts.JwtExpires,
				TokenType:   "Bearer",
				AccessToken: token,
				Scope:       "read_write",
			}
			resp.Success(c, output)
			return
		}
	}
	resp.Error(c, 2005, errors.New("未匹配正确APP信息"))
}
