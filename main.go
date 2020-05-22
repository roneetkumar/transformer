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
	"text/template"

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

		a, err := generateImage(file, ext, 33, primitive.ModeBeziers)
		if err != nil {
			panic(err)
		}
		file.Seek(0, 0)
		b, err := generateImage(file, ext, 33, primitive.ModeTriangle)
		if err != nil {
			panic(err)
		}
		file.Seek(0, 0)
		c, err := generateImage(file, ext, 33, primitive.ModeCircle)
		if err != nil {
			panic(err)
		}
		file.Seek(0, 0)
		d, err := generateImage(file, ext, 33, primitive.ModeRect)
		if err != nil {
			panic(err)
		}

		html := `
			<html>
				<body>
				{{range .}}
					<img src="/{{.}}"/><br><br>
				{{end}}
				</body>
			</html>
		`

		tpl := template.Must(template.New("").Parse(html))

		images := []string{a, b, c, d}

		// for i, img := range images {
		// 	images[i] = "/" + img
		// }

		tpl.Execute(w, images)

		// rurl := fmt.Sprintf("/%s", b)
		// http.Redirect(w, r, rurl, http.StatusFound)
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

func generateImage(r io.Reader, ext string, shapes int, mode primitive.Mode) (string, error) {
	out, err := primitive.Transform(r, ext, shapes, primitive.WithMode(mode))
	if err != nil {
		return "", err
	}

	outFile, err := tempFile("out_", ext)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	io.Copy(outFile, out)

	return outFile.Name(), nil

}
