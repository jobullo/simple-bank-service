package routes

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthController struct{}

// Login generates a JWT token
// @Summary Login generates a JWT token
// @Description Generates a JWT Token for use with Authorized endpoints
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Success 200 {object} AuthToken
// @Failure 401
// @Router /auth/login [post]
func (h AuthController) Login(c *gin.Context) {

	// Establish a 30 minute expiration
	expiresAt := time.Now().Add(time.Minute * 30).Unix()

	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = jwt.MapClaims{
		"exp": expiresAt,
		"iat": time.Now().Unix(),
	}

	// TODO: Load from config
	tokenString, err := token.SignedString([]byte("A14E45A7-D02B-4ADA-94BC-66DCBFD3181E"))
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, AuthToken{
		Token:     tokenString,
		TokenType: "Bearer",
		ExpiresIn: expiresAt,
	})
}

/** Auth struct for token management */
type AuthToken struct {
	TokenType string `json:"token_type"`
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}
