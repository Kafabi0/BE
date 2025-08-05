package models

type Pasien struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nama   string `gorm:"uniqueIndex" json:"nama"`
	Umur   string `json:"umur"`
	Alamat string `json:"alamat"`
}
