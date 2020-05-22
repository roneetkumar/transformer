package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/roneetkumar/transformers/primitive"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<html>
			<body>
				<form action="/upload" method="post" enctype="multipart/form-data">
					<input type="file" name="image"/>
					<button type="submit">Upload Image</button>
				</form
			</body>
		</html>`

		fmt.Fprint(w, html)
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		ext := filepath.Ext(header.Filename)[1:]

		out, err := primitive.Transform(file, ext, 100, primitive.WithMode(primitive.ModeBeziers))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		outFile, err := tempFile("out_", ext)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer outFile.Close()

		io.Copy(outFile, out)

		rurl := fmt.Sprintf("/%s", outFile.Name())
		http.Redirect(w, r, rurl, http.StatusFound)
	})

	fs := http.FileServer(http.Dir("./img/"))

	mux.Handle("/img/", http.StripPrefix("/img/", fs))

	log.Fatal(http.ListenAndServe(":3000", mux))
}

func tempFile(prefix string, ext string) (*os.File, error) {
	in, err := ioutil.TempFile("./img/", prefix)
	if err != nil {
		return nil, errors.New("main : failed to create temp file")
	}
	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))
}
