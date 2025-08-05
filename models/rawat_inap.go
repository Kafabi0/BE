package models

import "time"

type RawatInap struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	NamaPasien  string     `json:"nama_pasien"`
		Keluhan       string     `json:"keluhan"` // ✅ Tambahkan ini

	TanggalMasuk time.Time  `json:"tanggal_masuk"`
	TanggalKeluar *time.Time `json:"tanggal_keluar,omitempty"`
	Ruangan    string     `json:"ruangan"`
	Biaya      int        `json:"biaya"`
	Status     string     `json:"status"`       // contoh: "menunggu", "disetujui", "ditolak", "disetujui_user", "ditolak_user"
	Catatan    string     `json:"catatan"`
	IsUserSetuju *bool     `json:"is_user_setuju"`  // pointer karena nullable, nil = belum setuju/tolak
}

func (RawatInap) TableName() string {
	return "rawat_inap" // sesuai nama tabel di database kamu
}