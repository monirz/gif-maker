package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/nfnt/resize"
)

type data struct {
	Success Success `json:"success"`
	Status  int     `json:"status"`
}

type Success struct {
	Message string
}

func main() {
	http.HandleFunc("/create", fileUpload)
	http.HandleFunc("/file", getFile)
	fmt.Println("Listening on :8090...")
	http.ListenAndServe(":8090", nil)
}

func getFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "output.gif")
}

//fileUpload uploads nultiple files from formdata
func fileUpload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(20000)

	if err != nil {
		log.Fatal(err)
	}

	formdata := r.MultipartForm

	var files []*multipart.FileHeader
	for k, v := range formdata.File {
		fmt.Println(k, v)
		files = v
	}

	d := formdata.Value
	w.Header().Set("Content-Type", "application/json")

	if len(d["delay"]) < 1 || len(files) < 1 {
		w.WriteHeader(400)
		w.Write([]byte(`{"error": "Delay field or files are required in formdata"}`))
		return
	}

	di, err := strconv.Atoi(d["delay"][0])

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	//send file names and delay
	err = createGif(files, di)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	dt := data{Success{Message: "Gif file created successfully"}, 201}

	w.WriteHeader(201)

	err = json.NewEncoder(w).Encode(&dt)

	if err != nil {
		w.Write([]byte(err.Error()))
	}

}

//
func createGif(files []*multipart.FileHeader, delay int) error {

	var frames []*image.Paletted
	var dx = []int{}
	var dy = []int{}

	var newTempImg image.Image

	for i := range files {

		file, err := files[i].Open()

		if err != nil {
			return err
		}
		img, err := jpeg.Decode(file)

		if err != nil {
			return errors.New("Failed decoding jpeg: " + err.Error())
		}

		buf := bytes.Buffer{}

		err = gif.Encode(&buf, img, nil)

		if err != nil {
			return err
		}

		tmpimg, err := gif.Decode(&buf)

		if err != nil {
			err = errors.New("error decoding gif file: " + err.Error())
			return err
		}

		r := tmpimg.Bounds()

		var newX, newY int
		if len(dx) > 0 {
			if dx[i-1] != r.Dx() {
				newX = dx[i-1]
			}
		}

		if len(dy) > 0 {
			if dy[i-1] != r.Dy() {
				newY = dy[i-1]
				// return errors.New("All image must be same height")
			}
		}

		if newX > 0 || newY > 0 {
			newTempImg = resize.Resize(uint(newX), uint(newY), tmpimg, resize.Lanczos3)
		}

		dx = append(dx, r.Dx())
		dy = append(dy, r.Dy())

		if newTempImg != nil {

			err = gif.Encode(&buf, newTempImg, nil)

			if err != nil {
				return errors.New("Failed encoding resized image: " + err.Error())
			}

			tempImg, err := gif.Decode(&buf)
			if err != nil {
				return errors.New("Failed decoding resized image: " + err.Error())
			}

			frames = append(frames, tempImg.(*image.Paletted))

		} else {

			frames = append(frames, tmpimg.(*image.Paletted))
		}

	}

	delays := make([]int, len(frames))
	for j := range delays {
		delays[j] = delay
	}

	opfile, err := os.Create("output.gif")

	if err != nil {
		return errors.New("Failed ceating .gif file on disk: " + err.Error())
	}

	err = gif.EncodeAll(opfile, &gif.GIF{Image: frames, Delay: delays, LoopCount: 0})

	if err != nil {
		return errors.New("Failed gif encoding: " + err.Error())
	}

	return nil

}
