package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"loginapp/database"
	"loginapp/models"
	"github.com/gorilla/mux"
)

// CreateRawatInap user mengisi rawat inap
func CreateRawatInap(w http.ResponseWriter, r *http.Request) {
	var input models.RawatInap
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	input.Status = "menunggu"   // default status
	input.IsUserSetuju = nil    // belum setuju

	if err := database.DB.Create(&input).Error; err != nil {
		http.Error(w, "Gagal simpan data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(input)
}
func GetRawatInapByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["nama_pasien"]

	var rawats []models.RawatInap
	if err := database.DB.Where("nama_pasien = ?", username).Find(&rawats).Error; err != nil {
		http.Error(w, "Gagal mengambil data rawat inap", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rawats)
}

// GetRawatInap mengambil semua rawat inap (admin atau user)
func GetRawatInap(w http.ResponseWriter, r *http.Request) {
	var raws []models.RawatInap
	if err := database.DB.Find(&raws).Error; err != nil {
		http.Error(w, "Gagal ambil data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(raws)
}

// UpdateKonfirmasi admin mengonfirmasi rawat inap: menentukan ruangan & biaya
func UpdateKonfirmasi(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var input struct {
		Ruangan string `json:"ruangan"`
		Biaya   int    `json:"biaya"`
		Status  string `json:"status"` // biasanya "disetujui" atau "ditolak"
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var rawat models.RawatInap
	if err := database.DB.First(&rawat, id).Error; err != nil {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	rawat.Ruangan = input.Ruangan
	rawat.Biaya = input.Biaya
	rawat.Status = input.Status

	if err := database.DB.Save(&rawat).Error; err != nil {
		http.Error(w, "Gagal update data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rawat)
}

// SetujuiRawatInap user menyetujui atau menolak setelah konfirmasi admin
func SetujuiRawatInap(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var input struct {
		IsUserSetuju bool `json:"is_user_setuju"` // true = setuju, false = tolak
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var rawat models.RawatInap
	if err := database.DB.First(&rawat, id).Error; err != nil {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	rawat.IsUserSetuju = &input.IsUserSetuju
	if input.IsUserSetuju {
		rawat.Status = "disetujui_user"
	} else {
		rawat.Status = "ditolak_user"
	}

	if err := database.DB.Save(&rawat).Error; err != nil {
		http.Error(w, "Gagal update status", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rawat)
}
