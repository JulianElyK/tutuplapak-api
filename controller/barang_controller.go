package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tutuplapak-api/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get barang sesuai id yang dimasukkan atau kategori yang dipisah "," tanpa spasi
func GetBarang(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	err := r.ParseForm()
	if err != nil {
		return
	}

	id := r.Form.Get("id")
	email := r.Form.Get("email")
	kategori := r.Form.Get("kategori")

	filter := bson.D{{"deleted_at", bson.M{"$exists": false}}}
	if email != "" {
		filter = append(filter, bson.E{"penjual", email})
	}
	if id != "" {
		oId, _ := primitive.ObjectIDFromHex(id)
		filter = append(filter, bson.E{"_id", oId})
	} else if kategori != "" {
		kArr := strings.Split(kategori, ",")
		filter = append(filter, bson.E{"kategori", bson.D{{"$in", kArr}}})
	}

	var data interface{}
	if id != "" {
		var barang model.Barang
		err = db.Database("tutuplapak").Collection("barang").FindOne(ctx, filter).Decode(&barang)
		data = barang
	} else {
		cursor, err := db.Database("tutuplapak").Collection("barang").Find(ctx, filter)
		if err != nil {
			log.Println(err)
			return
		}
		var barang []model.Barang
		err = cursor.All(ctx, &barang)
		data = barang
	}

	if err != nil {
		log.Println(err)
		sendResponseData(w, 400, "Get Failed!", nil)
	} else {
		sendResponseData(w, 200, "Get Success!", data)
	}
}

// Insert barang oleh Seller yg login.
// Memasukkan string nama, int harga, int stok, string kategori yang dipisah "," tanpa spasi jika > 1 kategori
// dan file foto
func InsertBarang(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	err := r.ParseForm()
	if err != nil {
		return
	}

	email, _ := getEmailType(r)

	nama := r.Form.Get("nama")
	harga, _ := strconv.Atoi(r.Form.Get("harga"))
	stok, _ := strconv.Atoi(r.Form.Get("stok"))
	kategori := r.Form.Get("kategori")
	kArr := strings.Split(kategori, ",")
	// foto := r.Form.Get("foto")

	barang := model.Barang{
		CreatedAt: time.Now(),
		Nama:      nama,
		Harga:     harga,
		Penjual:   email,
		Stok:      stok,
		Kategori:  kArr,
		// Foto: foto,
	}

	_, err = db.Database("tutuplapak").Collection("barang").InsertOne(ctx, barang)
	if err != nil {
		sendResponseData(w, 400, "Insert Failed!", nil)
		log.Println(err)
	} else {
		sendResponseData(w, 200, "Insert Success!", nil)
	}

}

// Update data barang:
// nama baru, harga baru, menambah stok, mengubah kategori
//, dan mengganti foto
func UpdateBarang(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	err := r.ParseForm()
	if err != nil {
		return
	}

	id, err := primitive.ObjectIDFromHex(r.Form.Get("id"))
	if err != nil {
		log.Println(err)
	}

	// Get penjual barang
	var barang model.Barang
	var opt options.FindOneOptions
	opt.SetProjection(bson.M{"deleted_at": 1, "penjual": 1, "stok": 1})
	err = db.Database("tutuplapak").Collection("barang").FindOne(ctx, bson.M{"_id": id}, &opt).Decode(&barang)
	if err != nil {
		log.Println(err)
	}

	// Cek penjual barang dengan user login
	if !checkEmail(r, barang.Penjual) {
		sendUnAuthorizedResponse(w)
		return
	}
	// Cek barangnya dihapus
	if !barang.DeletedAt.IsZero() {
		sendResponseData(w, 400, "Barang telah dihapus!", nil)
		return
	}

	nama := r.Form.Get("nama")
	harga, _ := strconv.Atoi(r.Form.Get("harga"))
	stok, _ := strconv.Atoi(r.Form.Get("stok"))
	var kArr []string
	if kategori := r.Form.Get("kategori"); kategori != "" {
		kArr = strings.Split(kategori, ",")
	}
	// foto := r.Form.Get("foto")

	barang = model.Barang{
		Nama:     nama,
		Harga:    harga,
		Stok:     stok + barang.Stok,
		Kategori: kArr,
		// Foto: foto,
	}

	_, err = db.Database("tutuplapak").Collection("barang").UpdateOne(ctx, bson.M{"_id": id}, bson.D{{"$set", barang}})
	if err != nil {
		sendResponseData(w, 400, "Update Failed!", nil)
		log.Println(err)
	} else {
		sendResponseData(w, 200, "Update Success!", nil)
	}
}

// Soft Delete barang dengan menambah field deleted_at
func DeleteBarang(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	err := r.ParseForm()
	if err != nil {
		return
	}

	id, err := primitive.ObjectIDFromHex(r.Form.Get("id"))
	if err != nil {
		log.Println(err)
	}

	// Get penjual barang
	var barang model.Barang
	var opt options.FindOneOptions
	opt.SetProjection(bson.M{"penjual": 1})
	err = db.Database("tutuplapak").Collection("barang").FindOne(ctx, bson.M{"_id": id}, &opt).Decode(&barang)
	if err != nil {
		log.Println(err)
	}

	// Cek penjual barang dengan user
	if !checkEmail(r, barang.Penjual) {
		sendUnAuthorizedResponse(w)
		return
	}

	_, err = db.Database("tutuplapak").Collection("barang").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	if err != nil {
		sendResponseData(w, 400, "Delete Failed!", nil)
		log.Println(err)
	} else {
		sendResponseData(w, 200, "Delete Success!", nil)
	}
}

func checkEmail(r *http.Request, email string) bool {
	emailC, _ := getEmailType(r)
	return emailC == email
}
