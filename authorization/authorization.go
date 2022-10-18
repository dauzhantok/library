package authorization

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

// Claim for jwt token
type Claim struct {
	User *UserClaim `json:"user"`
	*jwt.StandardClaims
}

// UserClaim user object for claim
type UserClaim struct {
	User
	Roles   []*Role   `json:"roles"`
	Modules []*Module `json:"modules"`
}

// User user object
type User struct {
	Login      string  `json:"login"`
	Email      *string `json:"email"`
	FullName   *string `json:"fullName"`
	Title      *string `json:"title"`
	DepCode    *string `json:"depCode"`
	DepId      *string `json:"depId"`
	DepName    *string `json:"depName"`
	AuthType   *string `json:"authType"`
	Photo      *string `json:"photo"`
	FilialCode *string `json:"filialCode"`
	FilialName *string `json:"filialName"`
}

// Role role object for user
type Role struct {
	Authority   *string `json:"authority"`
	Code        *string `json:"code"`
	Comments    *string `json:"comments"`
	Description *string `json:"description"`
}

// Module module object for user profile role
type Module struct {
	Id          int     `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// SetCORSMiddleware configure CORS
func SetCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.GetHeader("origin"))
		c.Writer.Header().Set("Access-Control-Max-Age", "3600")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "x-requested-with, authorization, accesstoken, content-type, deviceid, ip, fingerprint")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusOK)
		}
	}
}

// AuthMiddleware verify token middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentToken := c.Request.Header.Get("Authorization")
		if currentToken == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "unauthorized"})
			return
		}

		jwtKey := c.MustGet("jwtKey").([]byte)
		jwtTokenMinute := c.MustGet("jwtTokenMinute").(int)

		claim := &Claim{}
		token, err := jwt.ParseWithClaims(currentToken, claim, func(token *jwt.Token) (i interface{}, e error) {
			return jwtKey, nil
		})
		if (token != nil && !token.Valid) || (err != nil && err == jwt.ErrSignatureInvalid) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "cant parse token with claim"})
			return
		}

		expTime := time.Now().Add(time.Duration(jwtTokenMinute) * time.Minute)
		claim.ExpiresAt = expTime.Unix()
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		refreshToken, err := token.SignedString(jwtKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "cant get jwtString"})
			return
		}

		c.Header("Authorization", refreshToken)
		c.Set("UserClaim", claim.User)
		c.Next()
		return
	}
}
