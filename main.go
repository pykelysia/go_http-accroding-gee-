package main

import (
	"gee"
	"net/http"
)

func main() {
	server := gee.New()
	server.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee<h1>")
	})
	server.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello\n")
	})
	server.GET("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	e := server.Run(":9999")
	if e != nil {
		panic(e)
	}
}
