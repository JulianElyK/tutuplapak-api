package controller

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tutuplapak-api/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Read history transaksi.
// Selain admin, return history dari user yang login.
// Admin dapat read semua transaksi / user tertentu &/ range tanggal yang dipisahkan "," tanpa spasi
func ReadTransaksi(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	err := r.ParseForm()
	if err != nil {
		return
	}

	filter := bson.D{}

	email, uType := getEmailType(r)
	waktu := r.Form.Get("waktu")
	if uType == "A" {
		email = r.Form.Get("email")
		if email != "" {
			filter = bson.D{{"email", email}}
		}
	} else {
		if email != "" {
			filter = bson.D{{"email", email}}
		} else {
			sendUnAuthorizedResponse(w)
			return
		}
	}
	if waktu != "" {
		wArr := strings.Split(waktu, ",")
		w1, _ := time.Parse("2006-01-02 15:04:05", wArr[0])
		w2, _ := time.Parse("2006-01-02 15:04:05", wArr[1]+" 23:59:59")
		filter = append(filter, bson.E{"created_at", bson.M{"$gte": w1, "$lte": w2}})
	}

	cursor, err := db.Database("tutuplapak").Collection("transaksi").Find(ctx, filter)
	if err != nil {
		log.Println(err)
		sendResponseData(w, 400, "Get Failed!", nil)
	} else {
		var transaksi []model.Transaksi
		err = cursor.All(ctx, &transaksi)
		if err != nil {
			log.Println(err)
			sendResponseData(w, 400, "Get Failed!", nil)
		} else {
			sendResponseData(w, 200, "Get Success!", transaksi)
		}
	}
}

// Create transaksi ketika beli barang.
// Input string id barang & jumlah beli yang dipisah "," tanpa spasi.
// Jumlah beli tiap barang harus urut sesuai dengan urutan id barang.
func CreateTransaksi(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	err := r.ParseForm()
	if err != nil {
		return
	}

	email, _ := getEmailType(r)
	if email == "" {
		sendResponseData(w, 200, "Log In required!", nil)
		return
	}

	id := r.Form.Get("id")
	jumlah := r.Form.Get("jumlah")

	iArr := strings.Split(id, ",")
	var convIArr []primitive.ObjectID
	for _, v := range iArr {
		convI, _ := primitive.ObjectIDFromHex(v)
		convIArr = append(convIArr, convI)
	}
	jArr := strings.Split(jumlah, ",")
	var convJArr []int
	for _, v := range jArr {
		convJ, _ := strconv.Atoi(v)
		convJArr = append(convJArr, convJ)
	}
	barang, msg, valid := getValidBarangAndUpdate(db, ctx, convIArr, convJArr)
	if !valid {
		sendResponseData(w, 400, msg, nil)
		return
	}

	transaksi := model.Transaksi{
		CreatedAt: time.Now(),
		Email:     email,
	}
	for i := 0; i < len(barang); i++ {
		transaksi.Barang = append(transaksi.Barang, model.BarangTransaksi{
			Item:   barang[i],
			Jumlah: convJArr[i],
		})
	}

	_, err = db.Database("tutuplapak").Collection("transaksi").InsertOne(ctx, transaksi)
	if err != nil {
		sendResponseData(w, 400, "Insert Failed!", nil)
	} else {
		if msg == "" {
			msg = "Insert Success!"
		} else {
			msg = "Insert Success with Warning: " + msg + "!"
		}
		sendResponseData(w, 200, msg, nil)
	}
}

func getValidBarangAndUpdate(db *mongo.Client, ctx context.Context, ids []primitive.ObjectID, jml []int) ([]model.Barang, string, bool) {
	msg := "Unmatched Length!"
	if len(ids) == len(jml) {
		msg = "Invalid barang not updated:"
		var opt options.FindOneAndUpdateOptions
		opt.SetProjection(bson.M{"nama": 1, "harga": 1, "penjual": 1})
		var barang []model.Barang
		isError := false
		for i := 0; i < len(ids); {
			var item model.Barang
			err := db.Database("tutuplapak").Collection("barang").FindOneAndUpdate(ctx, bson.M{"_id": ids[i]}, bson.M{"$inc": bson.M{"stok": -jml[i]}}, &opt).Decode(&item)
			if err == nil {
				barang = append(barang, item)
				i++
			} else {
				log.Println(err)
				ids = removeSliceID(ids, i)
				jml = removeSliceInt(jml, i)
				msg += " (" + ids[i].String() + ")"
				isError = true
			}
		}
		if !isError {
			msg = ""
		}
		return barang, msg, true
	}
	return nil, msg, false
}
