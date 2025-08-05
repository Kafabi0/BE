package models

type Produk struct {
	ID        int    `json:"id"`
	Nama_Produk      string `json:"nama_produk"`
	Harga     int    `json:"harga"`
	Kategori string `json:"kategori"`
	Stok      int    `json:"stok"`
	Gambar    string `json:"gambar"`
}

func (Produk) TableName() string {
	return "produk" // sesuai nama tabel di database kamu
}