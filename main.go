package main

import (
	"os"

	"github.com/gin-contrib/cors"
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
	noteHandler, _ := NewNoteHandler()
	authHandler, _ := NewAuthHandler(os.Getenv("CLIENT_ID"), os.Getenv("AUTH_DOMAIN"))
	authMiddleware := func(c *gin.Context) {
		GetAuthHeaderOrError(authHandler, c)
	}

	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5000"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "authorization")
	router.Use(cors.New(corsConfig))

	router.GET("/login", func(c *gin.Context) {
		c.Redirect(302, authHandler.GetRedirectURL(os.Getenv("CALLBACK_URL")))
	})

	router.GET("/notes/", authMiddleware, func(c *gin.Context) {
		userData, _ := c.Get("userInfo")
		userInfo := userData.(*UserInfo)

		notes, err := noteHandler.GetNotesByAuthor(userInfo.Sub)
		if err != nil {
			c.AbortWithStatus(404)
		}

		c.JSON(200, notes)
	})
	router.GET("/notes/:id", authMiddleware, func(c *gin.Context) {
		userData, _ := c.Get("userInfo")
		userInfo := userData.(*UserInfo)
		id := c.Param("id")

		result, err := noteHandler.GetNoteByID(id)
		if err != nil {
			c.AbortWithStatus(404)
			return
		}
		if result.Author != userInfo.Sub {
			c.AbortWithStatus(401)
			return
		}

		c.JSON(200, result)
	})

	router.POST("/notes/", authMiddleware, func(c *gin.Context) {
		userData, _ := c.Get("userInfo")
		userInfo := userData.(*UserInfo)

		var note Note
		if c.ShouldBind(&note) == nil {
			note.Author = userInfo.Sub
			result, err := noteHandler.AddNote(&note)
			if err != nil {
				c.AbortWithStatus(402)
				return
			}

			c.JSON(200, result)
			return
		}
		c.AbortWithStatus(402)
	})
	router.POST("/notes/:id", authMiddleware, func(c *gin.Context) {
		userData, _ := c.Get("userInfo")
		userInfo := userData.(*UserInfo)

		var note Note
		if c.ShouldBind(&note) == nil {
			note.Author = userInfo.Sub
			note.ID = c.Param("id")
			result, err := noteHandler.UpdateNote(&note)
			if err != nil {
				c.AbortWithStatus(402)
				return
			}

			c.JSON(200, result)
			return
		}
		c.AbortWithStatus(402)
	})

	router.DELETE("/notes/:id", authMiddleware, func(c *gin.Context) {
		userData, _ := c.Get("userInfo")
		userInfo := userData.(*UserInfo)
		id := c.Param("id")

		currentNote, err := noteHandler.GetNoteByID(id)
		if err != nil {
			c.AbortWithStatus(404)
			return
		}
		if currentNote.Author != userInfo.Sub {
			c.AbortWithStatus(402)
			return
		}

		noteHandler.DeleteNote(id)
		c.String(200, "Success")
	})

	router.Run()
}
