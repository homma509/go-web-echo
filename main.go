package main

import (
	"net/http"

	"github.com/labstack/echo"
)

// Comment is a struct to hold unit of request
type Comment struct {
	ID   int64  `json:"id"`
	Name string `json:"name" form:"name"`
	Text string `json:"text" form:"text"`
}

func main() {
	e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	name := c.QueryParam("name")
	// 	return c.String(http.StatusOK, "Hello " + name)
	// })
	// e.GET("/:name", func(c echo.Context) error {
	// 	name := c.Param("name")
	// 	return c.String(http.StatusOK, "Hello " + name)
	// })
	e.POST("api/comments", func(c echo.Context) error {
		var comment Comment
		if err := c.Bind(&comment); err != nil {
			return c.String(http.StatusBadRequest, "Bind: "+err.Error())
		}
		return c.String(http.StatusOK, "name " + comment.Name)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
