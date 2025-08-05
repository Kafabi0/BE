package models

type Obat struct {
    ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Nama     string `json:"nama"`
    Jenis    string `json:"jenis"`
    Stok     int    `json:"stok"`
    Harga    int    `json:"harga"`
    Kegunaan string `json:"kegunaan"`
    Kategori string `json:"kategori"`
}

func (Obat) TableName() string {
	return "obats"
}
