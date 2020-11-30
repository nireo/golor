package api

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/nireo/golor/utils"
	"github.com/valyala/fasthttp"
)

// TODO: add different html template serve pages

// UploadImage is a fasthttp post request handler, which is the starting point for finding the
// most prominant colors in a image. This method only handles saving the file and redirecting another route.
func UploadImage(ctx *fasthttp.RequestCtx) {
	// parse the file header from the request body
	header, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
		return
	}

	// get the file extension, since this is needed for proper decoding of the image
	splitted := strings.Split(header.Filename, ".")
	extension := splitted[len(splitted)-1]

	// save the file in a temporary folder, in which the image is stored under a unique id, with it's
	// file extension
	uid := utils.GenUUID()
	if err := fasthttp.SaveMultipartFile(header, fmt.Sprintf("./tmp/%s.%s", uid, extension)); err != nil {
		ctx.Error(
			fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
			fasthttp.StatusInternalServerError)
		return
	}

	// redirect to the route, which handles finding the most prominant colors (GetImageColors)
	ctx.Redirect(fmt.Sprintf("/api/colors?file=%s.%s", uid, extension), fasthttp.StatusMovedPermanently)
}

// ColorResultPage holds all the values that are the most prominent colors on an image.
// We need to store it in a struct, since golang templates need structs as data.
type ColorResultPage struct {
	// Colors holds all the index values and hex values of the most prominent colors
	Colors   []ColorItem
	Filename string
}

// ColorItem stores the index and the color of a single color in a image.
type ColorItem struct {
	Index int
	Color string
}

// GetImageColors is a fasthttp get request handler, which gets redirected to by
// the UploadImage handler. GetImageColors handles the decoding an image and finding the most prominant
// colors in an image.
func GetImageColors(ctx *fasthttp.RequestCtx) {
	// Load the filename which appended as a query from the /api/img request
	filename := string(ctx.QueryArgs().Peek("file"))

	// load the image, by first checking if the image exists and then decoding it!
	img, err := utils.LoadImage(fmt.Sprintf("tmp/%s", filename))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	// use a library to easily find all the color values fast
	colors, err := prominentcolor.KmeansWithArgs(prominentcolor.ArgumentNoCropping, img)
	if err != nil {
		// do not display a fasthttp.StatusMessage, since the error message is quite important,
		// for example for debugging purposes.
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	// Construct the ColorResultPage struct to store all the needed data.
	var results ColorResultPage
	for index, color := range colors {
		results.Colors = append(results.Colors, ColorItem{
			Index: index + 1,
			Color: color.AsString(),
		})
	}

	results.Filename = filename

	// Set the correct Content-Type so that the html template renders corrently.
	ctx.Response.Header.Set("Content-Type", "text/html")

	// Parse the template file and add the ColorResultPage struct data to the template,
	// which the results.html file specifically renders.
	tmpl := template.Must(template.ParseFiles("./pages/results.html"))
	if err := tmpl.Execute(ctx, results); err != nil {
		ctx.Error(
			fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
			fasthttp.StatusInternalServerError,
		)
		return
	}

	// Remove the file since we don't want to store it.
	if err := os.Remove(fmt.Sprintf("./tmp/%s", filename)); err != nil {
		ctx.Error(
			fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
			fasthttp.StatusInternalServerError,
		)
		return
	}
}

// ServeImageUploadPage servers the html of the page, from which the user can upload images,
// and then get redirected to the color results page.
func ServeImageUploadPage(ctx *fasthttp.RequestCtx) {
	// Set the right content type, since without the right content type the html won't render
	// properly.
	ctx.Response.Header.Set("Content-Type", "text/html")

	// Parse the upload page into golang template, which we can add more data to
	tmpl := template.Must(template.ParseFiles("./pages/upload.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		// Since we cannot execute the template we just return a internal server error, since
		// something must be very wrong.
		ctx.Error(
			fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
			fasthttp.StatusInternalServerError,
		)
		return
	}
}
