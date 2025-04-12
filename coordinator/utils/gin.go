package utils

import "github.com/gin-gonic/gin"

func BodyAs[T any](c *gin.Context) (T, error) {
	var body T
	if err := c.ShouldBindJSON(&body); err != nil {
		return body, err
	}
	return body, nil
}

func BodyAsF[T any](c *gin.Context) T {
	var body T
	if err := c.ShouldBindJSON(&body); err != nil {
		panic(err)
	}
	return body
}

func ErrorMsg(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"error": msg,
	})
}
