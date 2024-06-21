package controller

import (
    "database/sql"
    "html/template"
    "net/http"
    "path/filepath"
)

type StrukBelanja struct {
    Id          int
    Namamenu    string
    Harga       int
    Namauser    string
    Jumlah      int
    MPembayaran string
    TotalHarga  int
}

type Pemesan struct {
    Nama            string
    MetodePembayaran string
}

func StrukPDF(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        nameUser := r.FormValue("nama-user")
        MPembayaran := r.FormValue("metode-pembayaran")

        pemesanModel := Pemesan{
            Nama:            nameUser,
            MetodePembayaran: MPembayaran,
        }

        Struk, err := db.Query("SELECT id, nama, harga, jumlah FROM `pesanan` ORDER BY id DESC")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer Struk.Close()

        StrukModel := []StrukBelanja{}
        totalHargaKeseluruhan := 0

        for Struk.Next() {
            var id, harga, jumlah int
            var namabarang string

            err = Struk.Scan(&id, &namabarang, &harga, &jumlah)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            totalHarga := harga * jumlah
            totalHargaKeseluruhan += totalHarga

            StrukList := StrukBelanja{
                Id:          id,
                Namamenu:    namabarang,
                Namauser:    nameUser,
                MPembayaran: MPembayaran,
                Harga:       harga,
                Jumlah:      jumlah,
                TotalHarga:  totalHarga,
            }

            StrukModel = append(StrukModel, StrukList)
        }

        if err = Struk.Err(); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        fp := filepath.Join("view", "struk.html")
        tmpl, err := template.ParseFiles(fp)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        data := map[string]interface{}{
            "struk":         StrukModel,
            "user":          pemesanModel,
            "totalHargaKeseluruhan": totalHargaKeseluruhan,
        }

        err = tmpl.Execute(w, data)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
}
