package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Login struct {
	User     string `json:"user" form:"user" binding:"required"`
	Password string `json:"pass" form:"pass" binding:"required"`
}

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/hello/:name", func(context *gin.Context) {
		name := context.Param("name")
		fmt.Println("name: ", name)
		context.JSON(http.StatusOK, gin.H{
			"message": "your name input is: " + name,
		})
	})

	router.GET("/welcome", func(ctx *gin.Context) {
		name := ctx.DefaultQuery("name", "Guest")
		lastname := ctx.Query("lastname")
		fmt.Println("name: ", name, ", lastname: ", lastname)
	})

	router.POST("/loginJSON", func(ctx *gin.Context) {
		var json Login

		if ctx.BindJSON(&json) == nil {
			if json.User == "manu" && json.Password == "123" {
				ctx.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		}
	})

	router.POST("/loginForm", func(ctx *gin.Context) {
		var form Login

		if ctx.Bind(&form) == nil {
			if form.User == "manu" && form.Password == "123" {
				ctx.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		}
	})

	router.POST("/upload", func(ctx *gin.Context) {
		file, header, _ := ctx.Request.FormFile("upload")

		filename := header.Filename
		fmt.Println("filename: ", filename)

		// the dir must exsit
		out, err := os.Create("./tmp/" + filename)
		if err != nil {
			log.Fatal(err)
		}

		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
	})
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	// also: http.ListenAndServe(":8080", router)
	router.Run()
}
