package controller

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func InsertToDb(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {

			err := r.ParseMultipartForm(10 << 20)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			formdata := r.MultipartForm

			nama := r.Form["name-item"][0]
			harga := r.Form["harga"][0]
			foto := formdata.File["foto"]

			for i := range foto {
				file, err := foto[i].Open()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer file.Close()

				temp, err := ioutil.TempFile("file/", "upload-*.jpg")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer temp.Close()

				fotoname := temp.Name()

				filebytes, err := ioutil.ReadAll(file)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				_, err = temp.Write(filebytes)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				_, err = db.Exec("INSERT INTO `menu` (foto, nama, harga) VALUES (?,?,?)", fotoname, nama, harga)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				log.Printf("Upload File Sukses\n")
			}

			http.Redirect(w, r, "/", http.StatusMovedPermanently)

		} else if r.Method == "GET" {

			fp := filepath.Join("view", "input.html")
			tmpl, err := template.ParseFiles(fp)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}
