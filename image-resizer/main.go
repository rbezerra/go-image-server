package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
)

type ResizerParams struct {
	url    string
	height int
	width  int
}

func ResizeImage(w http.ResponseWriter, r *http.Request) {
	p, err := ParseQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	img, err := FetchAndResizeImage(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoded, err := EncodeImageToJpg(img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(encoded.Len()))

	cop, err := io.Copy(w, encoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ParseQuery(r *http.Request) (*ResizerParams, error) {
	var p ResizerParams
	query := r.URL.Query()
	url := query.Get("url")
	if url == "" {
		return &p, errors.New("Url param 'url' is missing")
	}

	width, _ := strconv.Atoi(query.Get("width"))
	height, _ := strconv.Atoi(query.Get("height"))

	if width == 0 && height == 0 {
		return &p, errors.New("Url param 'height' or 'width' must be set")
	}

	p = NewResizerParams(url, height, width)

	return &p, nil
}

func NewResizerParams(url string, height int, width int) ResizerParams {
	return ResizerParams{url, height, width}
}

func FetchAndResizeImage(p *ResizerParams) (*image.Image, error) {
	var dest image.Image

	response, err := http.Get(p.url)
	if err != nil {
		return &dest, err
	}
	defer response.Body.Close()

	src, _, err := image.Decode(response.Body)
	if err != nil {
		return &dest, err
	}

	dest = imaging.Resize(src, p.width, p.height, imaging.Lanczos)

	return &dest, nil
}

func EncodeImageToJpg(img *image.Image) (*bytes.Buffer, error) {
	encoded := &bytes.Buffer{}
	err := jpeg.Encode(encoded, *img, nil)
	return encoded, err
}

func main() {
	port := flag.Int("p", 8080, "server port")
	mux := http.NewServeMux()
	mux.HandleFunc("resize-image", ResizeImage)
	fmt.Printf("Startign local server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}
