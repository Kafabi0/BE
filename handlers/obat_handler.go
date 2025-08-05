package handlers

import (
    "encoding/json" 
    "net/http"
    "loginapp/database"
    "loginapp/models"
)

// Handler untuk membuat data obat
func CreateObat(w http.ResponseWriter, r *http.Request) {
    var obat models.Obat
    err := json.NewDecoder(r.Body).Decode(&obat)
    if err != nil {
        http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
        return
    }

    result := database.DB.Create(&obat)
    if result.Error != nil {
        http.Error(w, "Database error: "+result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(obat)
}

// Handler untuk menampilkan semua data obat
func GetAllObat(w http.ResponseWriter, r *http.Request) {
    var obats []models.Obat
    result := database.DB.Find(&obats)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(obats)
}

// Handler untuk mendapatkan satu data obat berdasarkan ID
func GetObatByID(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    var obat models.Obat
    result := database.DB.First(&obat, id)
    if result.Error != nil {
        http.Error(w, "Obat tidak ditemukan", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(obat)
}

// Handler untuk mengupdate data obat
func UpdateObat(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    var obat models.Obat
    result := database.DB.First(&obat, id)
    if result.Error != nil {
        http.Error(w, "Obat tidak ditemukan", http.StatusNotFound)
        return
    }

    var updatedObat models.Obat
    err := json.NewDecoder(r.Body).Decode(&updatedObat)
    if err != nil {
        http.Error(w, "Input tidak valid: "+err.Error(), http.StatusBadRequest)
        return
    }

    database.DB.Model(&obat).Updates(updatedObat)
    json.NewEncoder(w).Encode(obat)
}

// Handler untuk menghapus data obat
func DeleteObat(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    result := database.DB.Delete(&models.Obat{}, id)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Data berhasil dihapus"))
}
