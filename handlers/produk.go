package handlers

import (
	"encoding/json"
	"io"
	"loginapp/database"
	"loginapp/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

// Get semua produk
func GetProduks(w http.ResponseWriter, r *http.Request) {
	var produks []models.Produk
	if err := database.DB.Find(&produks).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(produks)
}

// Tambah produk dengan upload gambar
func CreateProduk(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // max 10MB
	if err != nil {
		http.Error(w, "Gagal memproses form", http.StatusBadRequest)
		return
	}

	nama := r.FormValue("nama_produk")
	harga, _ := strconv.Atoi(r.FormValue("harga"))
	kategori := r.FormValue("kategori")
	stok, _ := strconv.Atoi(r.FormValue("stok"))

	file, handler, err := r.FormFile("gambar")
	if err != nil {
		http.Error(w, "Gagal membaca file gambar", http.StatusBadRequest)
		return
	}
	defer file.Close()

	os.MkdirAll("./uploads", os.ModePerm)
	filePath := filepath.Join("uploads", handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Gagal menyimpan file gambar", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Gagal menyimpan file gambar", http.StatusInternalServerError)
		return
	}

	produk := models.Produk{
		Nama_Produk: nama,
		Harga:       harga,
		Kategori:    kategori,
		Stok:        stok,
		Gambar:      handler.Filename,
	}

	if err := database.DB.Create(&produk).Error; err != nil {
		http.Error(w, "Gagal menyimpan produk", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Produk berhasil ditambahkan"})
}

// Hapus produk
func DeleteProduk(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var produk models.Produk
	if err := database.DB.First(&produk, id).Error; err != nil {
		http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
		return
	}

	// Hapus file gambar
	if produk.Gambar != "" {
		filePath := filepath.Join("uploads", produk.Gambar)
		if _, err := os.Stat(filePath); err == nil {
			os.Remove(filePath)
		}
	}

	if err := database.DB.Delete(&produk).Error; err != nil {
		http.Error(w, "Gagal menghapus produk", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Produk berhasil dihapus"})
}

// Get produk by ID
func GetProdukByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var produk models.Produk
	if err := database.DB.First(&produk, id).Error; err != nil {
		http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(produk)
}

// Edit produk dengan kondisi gambar tidak berubah jika tidak upload gambar baru
func EditProduk(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var produk models.Produk
	if err := database.DB.First(&produk, id).Error; err != nil {
		http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Gagal memproses form", http.StatusBadRequest)
		return
	}

	nama := r.FormValue("nama_produk")
	harga, _ := strconv.Atoi(r.FormValue("harga"))
	kategori := r.FormValue("kategori")
	stok, _ := strconv.Atoi(r.FormValue("stok"))

	produk.Nama_Produk = nama
	produk.Harga = harga
	produk.Kategori = kategori
	produk.Stok = stok

	file, handler, err := r.FormFile("gambar")
	if err == nil {
		defer file.Close()

		// Simpan file gambar baru
		path := filepath.Join("uploads", handler.Filename)
		dst, err := os.Create(path)
		if err != nil {
			http.Error(w, "Gagal menyimpan file gambar", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Gagal menyimpan file gambar", http.StatusInternalServerError)
			return
		}

		// Hapus gambar lama jika beda nama dan bukan kosong
		if produk.Gambar != "" && produk.Gambar != handler.Filename {
			oldPath := filepath.Join("uploads", produk.Gambar)
			if _, err := os.Stat(oldPath); err == nil {
				os.Remove(oldPath)
			}
		}

		produk.Gambar = handler.Filename
	}

	if err := database.DB.Save(&produk).Error; err != nil {
		http.Error(w, "Gagal memperbarui produk", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Produk diperbarui"})
}
