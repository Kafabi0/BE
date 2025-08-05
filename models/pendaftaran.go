package models

type Pendaftaran struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Nama     string `json:"nama"`
    Email    string `json:"email"`
    Gender   string `json:"gender"`   // Laki-laki / Perempuan
    Jenis    string `json:"jenis"`    // konsultasi / operasi
    Kategori string `json:"kategori"` // sakit gigi / jantung / dll
    Tanggal  string `json:"tanggal"`
    Catatan  string `json:"catatan"`
    
}
func (Pendaftaran) TableName() string {
	return "pendaftaran" // sesuai nama tabel di database kamu
}