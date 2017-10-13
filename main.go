package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type data struct {
	Success Success `json:"success"`
	Status  int     `json:"status"`
}

type Success struct {
	Message string
}

func main() {
	http.HandleFunc("/file", fileUpload)
	fmt.Println("Listening on :8090...")
	http.ListenAndServe(":8090", nil)
}

//fileUpload uploads nultiple files from formdata
func fileUpload(w http.ResponseWriter, r *http.Request) {
	fn := []string{}

	err := r.ParseMultipartForm(20000)

	if err != nil {
		log.Fatal(err)
	}

	formdata := r.MultipartForm

	files := formdata.File["file"]
	d := formdata.Value

	di, err := strconv.Atoi(d["delay"][0])

	if err != nil {
		log.Fatal(err)
	}

	for i := range files {
		file, err := files[i].Open()

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()
		num := strconv.Itoa(i)
		out, err := os.Create(num + ".jpg")

		fn = append(fn, num+".jpg")

		if err != nil {
			log.Fatal(err)
		}

		defer out.Close()

		_, err = io.Copy(out, file)

		if err != nil {
			log.Fatal("error copying bytes from file ", err)
		}

	}

	//send file names and delay
	createGif(fn, di)

	w.Header().Set("Content-Type", "application/json")

	dt := data{Success{Message: "Gif file created successfully"}, 2001}

	b, err := json.Marshal(dt)

	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(b))

}

//
func createGif(files []string, delay int) {

	var frames []*image.Paletted

	for _, name := range files {

		file, err := os.Open(name)

		if err != nil {
			log.Fatal(err)
		}
		img, err := jpeg.Decode(file)

		if err != nil {
			log.Fatal(err)
		}

		// out, err := os.Create("out.gif")
		//
		// if err != nil {
		// 	log.Fatal(err)
		// }

		buf := bytes.Buffer{}

		err = gif.Encode(&buf, img, nil)

		if err != nil {
			log.Fatal(err)
		}

		tmpimg, err := gif.Decode(&buf)

		if err != nil {
			log.Fatal("error decoding gif file ", err)
		}

		frames = append(frames, tmpimg.(*image.Paletted))

	}

	delays := make([]int, len(frames))
	for j := range delays {
		delays[j] = delay
	}

	opfile, err := os.Create("output.gif")

	if err != nil {
		log.Fatal(err)
	}

	err = gif.EncodeAll(opfile, &gif.GIF{Image: frames, Delay: delays, LoopCount: 0})

	if err != nil {
		log.Fatal(err)
	}

}
