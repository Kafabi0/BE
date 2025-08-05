package database

import (
    "log"
    "os"
    "time"

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

    loc, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        log.Fatal("Failed to load timezone:", err)
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        NowFunc: func() time.Time {
            return time.Now().In(loc)
        },
    })
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

    log.Println("Koneksi GORM ke database berhasil dengan timezone Asia/Jakarta!")
}
