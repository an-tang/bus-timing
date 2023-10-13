package jwt

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JWTClaims struct {
	UserID int `json:"client"`
	jwt.StandardClaims
}

func Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerAuthorization := c.GetHeader("Authorization")

		if headerAuthorization != "" {
			token, err := decodeToken(headerAuthorization)
			if err != nil {
				log.Println("JWT Authorized", err.Error())
				c.Abort()

				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
				return
			}

			if token.Valid {
				c.Next()
			}
		} else {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Token"})
			return
		}
	}
}

func decodeToken(tokenHeader string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		tokenHeader,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("config.Config.SecretKeyJWT"), nil
		},
	)
}

func GetClaims(token string) JWTClaims {
	jwtToken, err := decodeToken(token)
	if err != nil {
		return JWTClaims{}
	}

	claims, ok := jwtToken.Claims.(*JWTClaims)
	if !ok {
		return JWTClaims{}
	}

	return *claims
}

// GenerateJWT the Claims
func GenerateJWT(claims JWTClaims) (map[string]string, error) {
	claims.ExpiresAt = time.Now().Add(15 * time.Minute).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte("config.Config.SecretKeyJWT"))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rfToken, err := refreshToken.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": rfToken,
	}, nil
}
