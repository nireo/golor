package main

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/nireo/golor/api"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fasthttprouter.New()
	router.POST("/api/img", api.UploadImage)
	router.GET("/api/colors", api.GetImageColors)
	router.GET("/", api.ServeImageUploadPage)

	if err := fasthttp.ListenAndServe("localhost:8080", router.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe")
	}
}
