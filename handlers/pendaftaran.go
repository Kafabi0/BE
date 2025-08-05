package handlers

import (
	"encoding/json"
	"loginapp/database"
	"loginapp/models"
	"net/http"

	"github.com/gorilla/mux"
)

func GetPendaftarans(w http.ResponseWriter, r *http.Request) {
	var daftar []models.Pendaftaran

	if err := database.DB.Find(&daftar).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(daftar)
}

 
func GetPendaftaranByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var p models.Pendaftaran
	if err := database.DB.First(&p, id).Error; err != nil {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

 
func CreatePendaftaran(w http.ResponseWriter, r *http.Request) {
	var p models.Pendaftaran
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Input tidak valid", http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&p).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Pendaftaran berhasil dibuat"})
}

 
func DeletePendaftaran(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := database.DB.Delete(&models.Pendaftaran{}, id).Error; err != nil {
		http.Error(w, "Gagal menghapus pendaftaran", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Pendaftaran berhasil dihapus"})
}
