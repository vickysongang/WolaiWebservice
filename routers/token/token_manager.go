package token

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	rsaService "WolaiWebservice/service/rsa"
)

type JWTManager struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

var tokenManager *JWTManager = nil

func GetTokenManager() *JWTManager {
	if tokenManager == nil {
		tokenManager = &JWTManager{
			privateKey: rsaService.GetPrivateKey(),
			PublicKey:  rsaService.GetPublicKey(),
		}
	}

	return tokenManager
}

func (m *JWTManager) GenerateToken(userId int64) (string, error) {
	var err error

	user, err := models.ReadUser(userId)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims["exp"] = time.Now().Add(time.Second * time.Duration(settings.TokenDuration())).Unix()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["userId"] = strconv.FormatInt(user.Id, 10)
	token.Claims["accessRight"] = strconv.FormatInt(user.AccessRight, 10)

	tokenString, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *JWTManager) TokenAuthenticate(userId int64, tokenString string) error {
	var err error

	_, err = models.ReadUser(userId)
	if err != nil {
		return errors.New("非法用户ID")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("非法的签名方法：%v", token.Header["alg"])
		} else {
			return m.PublicKey, nil
		}
	})

	if err != nil {
		return errors.New("非法安全令牌")
	}

	if token.Valid {
		if tokenId, ok := token.Claims["userId"]; ok {
			tokenIdStr, _ := tokenId.(string)
			userIdStr := strconv.FormatInt(userId, 10)

			if userIdStr == tokenIdStr {
				return nil
			} else {
				return errors.New("安全令牌信息不匹配")
			}

		} else {
			return errors.New("安全令牌信息格式错误")
		}

	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return errors.New("非法安全令牌")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return errors.New("安全令牌已失效")
		} else {
			return errors.New("非法安全令牌")
		}
	} else {
		return errors.New("非法安全令牌")
	}
}
