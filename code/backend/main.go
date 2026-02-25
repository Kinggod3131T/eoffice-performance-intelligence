package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()

	r := gin.Default()

	r.POST("/goals", CreateGoal)
	r.GET("/goals", GetGoals)

	r.Run(":8080")
}
