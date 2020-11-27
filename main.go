package main

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/gofrs/uuid"
	"github.com/valyala/fasthttp"
)

func GenUUID() string {
	v4, _ := uuid.NewV4()
	return v4.String()
}

func UploadImage(ctx *fasthttp.RequestCtx) {
	header, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
		return
	}

	if err := fasthttp.SaveMultipartFile(header, fmt.Sprintf("./temp/%s", GenUUID())); err != nil {
		ctx.Error(
			fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
			fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

func main() {
	router := fasthttprouter.New()
	router.POST("/api/img", UploadImage)

	if err := fasthttp.ListenAndServe("localhost:8080", router.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe")
	}
}
