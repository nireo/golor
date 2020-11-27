package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/buaazp/fasthttprouter"
	uuid "github.com/satori/go.uuid"
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

	// take the file extension
	splitted := strings.Split(header.Filename, ".")
	extension := splitted[len(splitted)-1]

	uid := GenUUID()
	if err := fasthttp.SaveMultipartFile(header, fmt.Sprintf("./tmp/%s.%s", uid, extension)); err != nil {
		ctx.Error(
			fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
			fasthttp.StatusInternalServerError)
		return
	}


	ctx.Redirect(fmt.Sprintf("/api/colors?file=%s.%s", uid, extension), fasthttp.StatusMovedPermanently)
}

func LoadImage(filename string) (image.Image, error) {
	fmt.Println(filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("could not load file")
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func GetColorsInImage(ctx *fasthttp.RequestCtx) {
	// load image
	filename := string(ctx.QueryArgs().Peek("file"))
	img, err := LoadImage(fmt.Sprintf("tmp/%s", filename))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	// load all the prominent colors in an image
	colors, err := prominentcolor.KmeansWithArgs(prominentcolor.ArgumentNoCropping, img)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	var content string
	for index, color := range colors {
		content += fmt.Sprintf("%d. #%s\n", index, color.AsString())
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.AppendBody([]byte(content))
}

func main() {
	router := fasthttprouter.New()
	router.POST("/api/img", UploadImage)
	router.GET("/api/colors", GetColorsInImage)

	if err := fasthttp.ListenAndServe("localhost:8080", router.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe")
	}
}
