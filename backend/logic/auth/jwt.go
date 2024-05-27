package auth

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 用户token秘钥
const SignToken = "go_build_smart_darwin"
const ExpireTime = 86400 * 7

type TokenInfo struct {
	Uid      string `json:"uid"`
	UserName string `json:"user_name"` // 用户名
	Name     string `json:"name"`      // 真实名称
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Appid    string `json:"appid"`
	Exp      int64  `json:"exp"`
}

func CreateToken(key string, tokenInfo TokenInfo) string {

	claims := make(jwt.MapClaims)

	t := reflect.TypeOf(tokenInfo)
	v := reflect.ValueOf(tokenInfo)

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")

		switch v.Field(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			claims[tag] = v.Field(i).Int()
		case reflect.String:
			claims[tag] = v.Field(i).String()
		}
	}
	claims["iat"] = time.Now().Unix()

	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(key))

	return tokenString
}

func ParseToken(key string, tokenString string) (tokenInfo *TokenInfo, err error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})

	// 校验
	if token == nil || token.Claims == nil {
		return nil, err
	}

	// 校验
	if err := token.Claims.Valid(); err != nil {
		return nil, err
	}

	return GetTokenInfo(token.Claims), nil
}

func GetTokenInfo(claims interface{}) *TokenInfo {

	tokenInfo := &TokenInfo{}
	claimsValue := reflect.ValueOf(claims)
	if claimsValue.Kind() != reflect.Map {
		return tokenInfo
	}

	t := reflect.TypeOf(*tokenInfo)
	v := reflect.ValueOf(tokenInfo).Elem()

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")

		value, ok := claims.(jwt.MapClaims)[tag]
		if !ok {
			continue
		}

		switch v.Field(i).Kind() {
		case reflect.Int64, reflect.Float64:
			vInt := int64(value.(float64))
			v.Field(i).Set(reflect.ValueOf(vInt))
		case reflect.String:
			v.Field(i).Set(reflect.ValueOf(value))
		case reflect.Slice:
			v.Field(i).Set(reflect.ValueOf(value.([]string)))
		}
	}

	return tokenInfo
}
