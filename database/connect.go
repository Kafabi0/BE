package database

import (
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "loginapp/models"
)

var DB *gorm.DB

func ConnectDB() {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL environment variable not set")
    }
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Gagal konek database:", err)
    }

    err = db.AutoMigrate(&models.RawatInap{},
        &models.Pasien{},
        &models.Produk{},
        &models.User{},
        &models.Pendaftaran{},
        &models.Obat{},
    )

    if err != nil {
        log.Fatal("Gagal migrasi:", err)
    }

    DB = db

    log.Println("Koneksi GORM ke database berhasil!")
}
