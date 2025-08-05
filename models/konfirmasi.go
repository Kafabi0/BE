package models

type Konfirmasi struct {
  ID                uint      `json:"id" gorm:"primaryKey"`
  PendaftaranID     uint      `json:"pendaftaran_id"`
  Status            string    `json:"status"`
  Keterangan        string    `json:"keterangan"`
	TanggalKonfirmasi string `json:"tanggal_konfirmasi"`
}
func (Konfirmasi) TableName() string {
	return "konfirmasi" // sesuai nama tabel di database kamu
}
