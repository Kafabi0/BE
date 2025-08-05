package routes

import (
	"net/http"
	"rumahsakit/db"
	"rumahsakit/models"

	"github.com/gin-gonic/gin"
)

func PasienRoutes(r *gin.Engine) {
	r.GET("/pasien", func(c *gin.Context) {
		var data []models.Pasien
		db.DB.Find(&data)
		c.JSON(http.StatusOK, data)
	})

	r.POST("/pasien", func(c *gin.Context) {
		var pasien models.Pasien
		if err := c.ShouldBindJSON(&pasien); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validasi apakah nama sudah ada
		var existing models.Pasien
		if err := db.DB.Where("nama = ?", pasien.Nama).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Pasien dengan nama tersebut sudah ada"})
			return
		}

		if err := db.DB.Create(&pasien).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan pasien"})
			return
		}
		c.JSON(http.StatusCreated, pasien)
	})

	r.PUT("/pasien/:id", func(c *gin.Context) {
		id := c.Param("id")
		var pasien models.Pasien

		// Cek apakah data ada
		if err := db.DB.First(&pasien, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data pasien tidak ditemukan"})
			return
		}

		// Bind data update
		var updateData models.Pasien
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validasi apakah nama baru sudah digunakan pasien lain
		var existing models.Pasien
		if err := db.DB.Where("nama = ? AND id != ?", updateData.Nama, id).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Nama pasien sudah digunakan"})
			return
		}

		// Update field
		pasien.Nama = updateData.Nama
		pasien.Umur = updateData.Umur
		pasien.Alamat = updateData.Alamat
		// tambahkan field lain sesuai kebutuhan

		if err := db.DB.Save(&pasien).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui pasien"})
			return
		}

		c.JSON(http.StatusOK, pasien)
	})

	r.DELETE("/pasien/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Cek apakah data ada sebelum hapus
		var pasien models.Pasien
		if err := db.DB.First(&pasien, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data pasien tidak ditemukan"})
			return
		}

		db.DB.Delete(&pasien)
		c.JSON(http.StatusOK, gin.H{"message": "Pasien berhasil dihapus"})
	})
}
