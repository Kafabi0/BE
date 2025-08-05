// models/pasien.go
package models

type Pasien struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Nama        string `json:"nama"`
	Umur        int    `json:"umur"`
	JenisKelamin string `json:"jenis_kelamin"`
	Alamat      string `json:"alamat"`
	NoTelepon   string `json:"no_telepon"`
}
func (Pasien) TableName() string {
	return "pasien" // sesuai nama tabel di database kamu
}
