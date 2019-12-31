package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

// 一些常量
var (
	// TokenExpired token过期
	TokenExpired error = errors.New("Token is expired")
	// TokenNotValidYet token验证错误
	TokenNotValidYet error = errors.New("Token not active yet")
	// TokenMalformed token未携带
	TokenMalformed error = errors.New("That's not even a token")
	// TokenInvalid token无效
	TokenInvalid error = errors.New("Couldn't handle this token:")
	// SignKey 签名
	SignKey string = "1233444"
)

// CustomClaims 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	FatherId string `json:"father_id"`
	jwt.StandardClaims
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// NewJWT 新建一个jwt实例
func NewJWT(sign string) *JWT {
	return &JWT{
		[]byte(sign),
	}
}

// ParseToken 解析Tokne
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}
