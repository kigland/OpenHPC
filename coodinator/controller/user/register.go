package user

import "github.com/gin-gonic/gin"

func register(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "register",
	})
}
