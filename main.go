package main

import (
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func GetAuthHeaderOrError(authHandler AuthHandler, c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" {
		userInfo, err := authHandler.GetUserInfo(token)
		if err != nil {
			c.AbortWithStatus(401)
		} else {
			c.Set("userInfo", userInfo)
			c.Next()
		}
	} else {
		c.AbortWithStatus(401)
	}
}
func main() {
	router := gin.Default()
	authHandler, _ := NewAuthHandler(os.Getenv("CLIENT_ID"), os.Getenv("AUTH_DOMAIN"))
	authMiddleware := func(c *gin.Context) {
		GetAuthHeaderOrError(authHandler, c)
	}

	router.GET("/me", authMiddleware, func(c *gin.Context) {
		userData, _ := c.Get("userInfo")
		userInfo := userData.(*UserInfo)

		c.JSON(200, userInfo.Sub)
	})
	router.GET("/login", func(c *gin.Context) {
		c.Redirect(302, authHandler.GetRedirectURL(os.Getenv("CALLBACK_URL")))
	})
	router.Run()
}
