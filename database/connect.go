package database

import (
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"loginapp/models"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=zerxy dbname=coba port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal konek database:", err)
	}

	// Migrasi otomatis tabel RawatInap (dan model lain kalau ada)
	err = db.AutoMigrate(&models.RawatInap{},
		&models.Pasien{},
		&models.Produk{},
		&models.User{},
		&models.Pendaftaran{},
		&models.Obat{})
	
	if err != nil {
		log.Fatal("Gagal migrasi:", err)
	}

	DB = db

	log.Println("Koneksi GORM ke database berhasil!")
}
