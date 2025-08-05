package db

import (
	"log"
	"rumahsakit/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=zerxy dbname=rumahsakit port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Menampilkan log query SQL
	})
	if err != nil {
		log.Fatal("Gagal konek database:", err)
	}

	// Auto migrate model
	err = database.AutoMigrate(&models.Pasien{})
	err = database.AutoMigrate(&models.User{})

	if err != nil {
		log.Fatal("Gagal migrate:", err)
	}

	DB = database
	log.Println("Koneksi ke database berhasil! [GORM]")
}
