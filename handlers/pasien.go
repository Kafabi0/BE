package handlers

import (
	"encoding/json"
	"loginapp/database"
	"loginapp/models"
	"net/http"

	"github.com/gorilla/mux"
)

// GET semua pasien
func GetPasiens(w http.ResponseWriter, r *http.Request) {
	var pasiens []models.Pasien

	if err := database.DB.Find(&pasiens).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pasiens)
}

// GET pasien by ID
func GetPasienByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var p models.Pasien
	if err := database.DB.First(&p, id).Error; err != nil {
		http.Error(w, "Pasien tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// POST pasien baru
func CreatePasien(w http.ResponseWriter, r *http.Request) {
	var p models.Pasien
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Input tidak valid", http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&p).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Pasien berhasil ditambahkan"})
}

// PUT update pasien
func UpdatePasien(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var p models.Pasien
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Input tidak valid", http.StatusBadRequest)
		return
	}

	// Pastikan pasien dengan ID ini ada dulu
	var existing models.Pasien
	if err := database.DB.First(&existing, id).Error; err != nil {
		http.Error(w, "Pasien tidak ditemukan", http.StatusNotFound)
		return
	}

	// Update data
	p.ID = existing.ID
	if err := database.DB.Save(&p).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Pasien berhasil diperbarui"})
}

// DELETE pasien
func DeletePasien(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := database.DB.Delete(&models.Pasien{}, id).Error; err != nil {
		http.Error(w, "Gagal menghapus pasien", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Pasien berhasil dihapus"})
}
