package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"tutuplapak-api/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get barang sesuai id yang dimasukkan atau kategori yang dipisah "," tanpa spasi
// Untuk foto dalam bentuk string
func GetBarang(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	id := r.URL.Query()["id"]
	email := r.URL.Query()["email"]
	kategori := r.URL.Query()["kategori"]

	filter := bson.D{{"deleted_at", bson.M{"$exists": false}}}
	if email != nil {
		filter = append(filter, bson.E{"penjual", email[0]})
	}
	if id != nil {
		oId, _ := primitive.ObjectIDFromHex(id[0])
		filter = append(filter, bson.E{"_id", oId})
	} else if kategori != nil {
		kArr := strings.Split(kategori[0], ",")
		filter = append(filter, bson.E{"kategori", bson.D{{"$in", kArr}}})
	}

	var err error
	var data interface{}
	if id != nil {
		var barang model.Barang
		err = db.Database("tutuplapak").Collection("barang").FindOne(ctx, filter).Decode(&barang)
		data = barang
	} else {
		cursor, err := db.Database("tutuplapak").Collection("barang").Find(ctx, filter)
		if checkErr(err) {
			return
		}
		var barang []model.Barang
		err = cursor.All(ctx, &barang)
		data = barang
	}

	if checkErr(err) {
		sendResponseData(w, 400, "Get Failed!", nil)
	} else {
		sendResponseData(w, 200, "Get Success!", data)
	}
}

// Insert barang oleh Seller yg login.
// Memasukkan string nama, int harga, int stok, array kategori, dan string foto
func InsertBarang(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	email, _ := getEmailType(r)

	var barang model.Barang
	err := json.NewDecoder(r.Body).Decode(&barang)
	checkErr(err)
	barang.Penjual = email

	_, err = db.Database("tutuplapak").Collection("barang").InsertOne(ctx, barang)
	if checkErr(err) {
		sendResponseData(w, 400, "Insert Failed!", nil)
	} else {
		sendResponseData(w, 200, "Insert Success!", nil)
	}

}

// Update data barang:
// nama baru, harga baru, menambah stok, mengubah kategori
func UpdateBarang(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	qId := r.URL.Query()["id"]

	var id primitive.ObjectID
	var err error
	if qId == nil {
		id, err = primitive.ObjectIDFromHex(qId[0])
		checkErr(err)
	}

	// Get penjual barang
	var barang model.Barang
	var opt options.FindOneOptions
	opt.SetProjection(bson.M{"deleted_at": 1, "penjual": 1, "stok": 1})
	err = db.Database("tutuplapak").Collection("barang").FindOne(ctx, bson.M{"_id": id}, &opt).Decode(&barang)
	checkErr(err)

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

	err = json.NewDecoder(r.Body).Decode(&barang)
	checkErr(err)

	barangFlat, err := Flatten(barang)
	checkErr(err)

	_, err = db.Database("tutuplapak").Collection("barang").UpdateOne(ctx, bson.M{"_id": id}, bson.D{{"$set", barangFlat}})
	if checkErr(err) {
		sendResponseData(w, 400, "Update Failed!", nil)
	} else {
		sendResponseData(w, 200, "Update Success!", nil)
	}
}

// Soft Delete barang dengan menambah field deleted_at
func DeleteBarang(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	qId := r.URL.Query()["id"]

	var id primitive.ObjectID
	var err error
	if qId == nil {
		id, err = primitive.ObjectIDFromHex(qId[0])
		checkErr(err)
	}

	// Get penjual barang
	var barang model.Barang
	var opt options.FindOneOptions
	opt.SetProjection(bson.M{"penjual": 1})
	err = db.Database("tutuplapak").Collection("barang").FindOne(ctx, bson.M{"_id": id}, &opt).Decode(&barang)
	checkErr(err)

	// Cek penjual barang dengan user
	if !checkEmail(r, barang.Penjual) {
		sendUnAuthorizedResponse(w)
		return
	}

	_, err = db.Database("tutuplapak").Collection("barang").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	if checkErr(err) {
		sendResponseData(w, 400, "Delete Failed!", nil)
	} else {
		sendResponseData(w, 200, "Delete Success!", nil)
	}
}

func checkEmail(r *http.Request, email string) bool {
	emailC, _ := getEmailType(r)
	return emailC == email
}
