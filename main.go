package main

import (
	"fmt"
	"log"
	"net/http"
	"tutuplapak-api/controller"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()
	// Tipe User
	// A - Admin, B - Buyer/Normal User, S - Seller

	// Admin
	router.HandleFunc("/admin/login", controller.LoginUser).Methods("POST")
	router.HandleFunc("/admin/logout", controller.Authenticate(controller.LogoutUser, "A")).Methods("GET")
	router.HandleFunc("/admin/log", controller.Authenticate(controller.ReadUserLog, "A")).Methods("GET")
	router.HandleFunc("/admin/transaksi", controller.Authenticate(controller.ReadTransaksi, "A")).Methods("GET")

	// Non-admin
	router.HandleFunc("/user/register", controller.InsertUser).Methods("POST")
	router.HandleFunc("/login", controller.LoginUser).Methods("POST")
	router.HandleFunc("/user/update", controller.UpdateUser).Methods("PUT")
	router.HandleFunc("/users", controller.GetUser).Methods("GET")
	router.HandleFunc("/logout", controller.LogoutUser).Methods("GET")

	// Barang
	router.HandleFunc("/barang/jual", controller.Authenticate(controller.InsertBarang, "S")).Methods("POST")
	router.HandleFunc("/barang", controller.GetBarang).Methods("GET")
	router.HandleFunc("/barang/update", controller.Authenticate(controller.UpdateBarang, "S")).Methods("PUT")
	router.HandleFunc("/barang/hapus", controller.Authenticate(controller.DeleteBarang, "S")).Methods("PUT")

	// Transaksi
	router.HandleFunc("/transaksi/create", controller.CreateTransaksi).Methods("POST")
	router.HandleFunc("/transaksi", controller.ReadTransaksi).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})
	handler := corsHandler.Handler(router)

	http.Handle("/", handler)
	fmt.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
