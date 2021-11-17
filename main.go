package main

import (
	"fmt"
	"log"
	"net/http"
	"tutuplapak-api/controller"

	"github.com/gorilla/mux"
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

	http.Handle("/", router)
	fmt.Println("Connected to port 9090")
	log.Fatal(http.ListenAndServe(":9090", router))
}
