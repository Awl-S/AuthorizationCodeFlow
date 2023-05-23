package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OAuthServer struct {
	clients      map[string]string
	authCodes    map[string]string
	accessTokens map[string]string
}

func NewOAuthServer() *OAuthServer {
	return &OAuthServer{
		clients:      make(map[string]string),
		authCodes:    make(map[string]string),
		accessTokens: make(map[string]string),
	}
}

func (s *OAuthServer) AuthorizeHandler(c *gin.Context) {
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")

	if _, ok := s.clients[clientID]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_client"})
		return
	}

	// Пропускаем процесс проверки и авторизации пользователя для упрощения
	authCode := randSeq(10)
	s.authCodes[authCode] = clientID

	c.Redirect(http.StatusFound, redirectURI+"?code="+authCode)
}

func (s *OAuthServer) TokenHandler(c *gin.Context) {
	authCode := c.PostForm("code")
	clientID := c.PostForm("client_id")

	storedClientID, ok := s.authCodes[authCode]
	if !ok || storedClientID != clientID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant"})
		return
	}

	// Возвращаем access token (в реальности должен быть JSON Web Token или аналогичным типом токена)
	accessToken := randSeq(20)
	s.accessTokens[accessToken] = clientID

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "token_type": "Bearer"})
}

func (s *OAuthServer) ResourceHandler(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")

	if _, ok := s.accessTokens[accessToken]; !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "secure data"})
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *OAuthServer) CallbackHandler(c *gin.Context) {
	code := c.Query("code")
	c.JSON(http.StatusOK, gin.H{"code": code})
}

func main() {
	rand.Seed(time.Now().UnixNano())

	oauthServer := NewOAuthServer()
	oauthServer.clients["my_client_id"] = "my_client_secret"

	r := gin.Default()

	r.GET("/authorize", oauthServer.AuthorizeHandler)
	r.POST("/token", oauthServer.TokenHandler)
	r.GET("/resource", oauthServer.ResourceHandler)
	r.GET("/callback", oauthServer.CallbackHandler)

	r.Run(":8080")
}
