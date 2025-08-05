package main

import (
	"log"
	"net/http"
"os"
	"loginapp/database"
	"loginapp/handlers"
	"loginapp/middleware" // kita buat middleware role-based

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.ConnectDB()

	r := mux.NewRouter()

	// Public endpoints (register, login, dll)
	r.HandleFunc("/api/register", handlers.Register).Methods("POST")
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")

	// Protected endpoints tanpa role khusus (cukup login valid)
	r.Handle("/api/protected", middleware.JWTMiddleware(http.HandlerFunc(handlers.Protected))).Methods("GET")

	// Produk (misal hanya admin yang bisa tambah, edit, hapus)
	r.HandleFunc("/api/produk", handlers.GetProduks).Methods("GET")

	r.Handle("/api/produk/tambah",
		middleware.JWTMiddleware(middleware.RoleAuthorization("admin")(http.HandlerFunc(handlers.CreateProduk))),
	).Methods("POST")

	r.Handle("/api/produk/edit/{id}",
		middleware.JWTMiddleware(middleware.RoleAuthorization("admin")(http.HandlerFunc(handlers.EditProduk))),
	).Methods("PUT")

	r.Handle("/api/produk/hapus/{id}",
		middleware.JWTMiddleware(middleware.RoleAuthorization("admin")(http.HandlerFunc(handlers.DeleteProduk))),
	).Methods("DELETE")

	r.HandleFunc("/api/produk/{id}", handlers.GetProdukByID).Methods("GET")

	// Pasien (bisa dibatasi role sesuai kebutuhan, misal admin dan dokter bisa akses)
	r.Handle("/admin/pasien",
		middleware.JWTMiddleware(middleware.RoleAuthorization("admin")(http.HandlerFunc(handlers.GetPasiens))),
	).Methods("GET")

	r.Handle("/admin/pasien/{id}",
		middleware.JWTMiddleware(middleware.RoleAuthorization("admin")(http.HandlerFunc(handlers.GetPasienByID))),
	).Methods("GET")

	r.Handle("/api/admin/pasien/create",
		middleware.JWTMiddleware(middleware.RoleAuthorization("admin")(http.HandlerFunc(handlers.CreatePasien))),
	).Methods("POST")

	r.Handle("/admin/pasien/{id}",
		middleware.JWTMiddleware(middleware.RoleAuthorization("admin")(http.HandlerFunc(handlers.UpdatePasien))),
	).Methods("PUT")

	r.Handle("/admin/pasien/{id}",
		middleware.JWTMiddleware(middleware.RoleAuthorization("admin")(http.HandlerFunc(handlers.DeletePasien))),
	).Methods("DELETE")

	// Serve gambar statis dari folder uploads
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))


	r.HandleFunc("/api/pendaftaran", handlers.CreatePendaftaran).Methods("POST")
r.HandleFunc("/api/pendaftaran", handlers.GetPendaftarans).Methods("GET")
r.HandleFunc("/api/pendaftaran/{id}", handlers.GetPendaftaranByID).Methods("GET")
r.HandleFunc("/api/pendaftaran/{id}", handlers.DeletePendaftaran).Methods("DELETE")

// r.HandleFunc("/api/konfirmasi", handlers.GetKonfirmasiHandler).Methods("GET")
// r.HandleFunc("/api/konfirmasi/{id:[0-9]+}", handlers.GetKonfirmasiByID).Methods("GET")
// r.HandleFunc("/api/konfirmasi", handlers.CreateKonfirmasi).Methods("POST")
// r.HandleFunc("/api/konfirmasi/{id:[0-9]+}", handlers.UpdateKonfirmasi).Methods("PUT")
// r.HandleFunc("/api/konfirmasi/{id:[0-9]+}", handlers.DeleteKonfirmasi).Methods("DELETE")

r.HandleFunc("/api/rawat-inap", handlers.CreateRawatInap).Methods("POST")
r.HandleFunc("/api/rawat-inap", handlers.GetRawatInap).Methods("GET")
r.HandleFunc("/api/rawat-inap/{id}/konfirmasi", handlers.UpdateKonfirmasi).Methods("PUT")
r.HandleFunc("/api/rawat-inap/{id}/setujui", handlers.SetujuiRawatInap).Methods("PUT")

r.HandleFunc("/api/rawat-inap/nama_pasien/{nama_pasien}", handlers.GetRawatInapByUserID).Methods("GET")
r.HandleFunc("/api/obat", handlers.CreateObat).Methods("POST")
r.HandleFunc("/api/obat", handlers.GetAllObat).Methods("GET")
r.HandleFunc("/api/obat/detail", handlers.GetObatByID).Methods("GET")
r.HandleFunc("/api/obat/update", handlers.UpdateObat).Methods("PUT")
r.HandleFunc("/api/obat/delete", handlers.DeleteObat).Methods("DELETE")



	// Middleware CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Ganti dengan asal frontend kamu di production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)
port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Jalankan server

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
