package handler

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Tushar456/go-carzone/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.opentelemetry.io/otel"
)

// LoginHandler godoc
// @Summary      Login
// @Description  Authenticates user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      models.Credentials  true  "User credentials"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /login [post]
func LoginHandler(c *gin.Context) {
	_, span := otel.Tracer("loginservice").Start(c.Request.Context(), "LoginHandler")
	defer span.End()

	var credentials models.Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if credentials.Username != "admin" || credentials.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := GenerateToken(credentials.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func GenerateToken(username string) (string, error) {

	jwtSecretKey := os.Getenv("JWT_SECRET")
	if jwtSecretKey == "" {
		return "", jwt.ErrInvalidKey
	}
	expiryTime, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_TIME"))
	if err != nil || expiryTime <= 0 {
		expiryTime = 24 // default to 24 hours if not set or invalid
	}
	token := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(expiryTime)).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   username,
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, token).SignedString([]byte(jwtSecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
