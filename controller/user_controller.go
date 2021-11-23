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

// Login user pertama kali.
// Data user disimpan di cookie.
// Pemanggilan data user di cookie pakai getEmail()
func LoginUser(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	email, _ := getEmailType(r)
	if email != "" {
		sendResponseData(w, 200, "An account already logged in! Log out first!", nil)
		return
	}

	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	checkErr(err)

	err = db.Database("tutuplapak").Collection("user").FindOne(ctx, bson.M{"email": user.Email, "password": user.Password}).Decode(&user)
	if checkErr(err) {
		sendResponseData(w, 400, "Login Failed: Email / Password is not correct!", nil)
	} else {
		if r.URL.Path == "/admin/login" {
			if user.Tipe != "A" {
				sendUnAuthorizedResponse(w)
				return
			}
		}
		generateToken(w, user.Email, user.Nama, user.Tipe)
		sendResponseData(w, 200, "Login Success!", nil)
		AddUserLog(email, "LI")
	}
}

// Logout user. Untuk logout
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	email, _ := getEmailType(r)
	if email == "" {
		sendUnAuthorizedResponse(w)
		return
	}
	AddUserLog(email, "LO")
	resetUsersToken(w)
	sendResponseData(w, 200, "Logout Success!", nil)
}

// Get user.
// Prioritas pakai param.
// Jika tidak ada pakai cookie.
// Default get all users
func GetUser(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	var filter interface{}
	// By Param
	id := r.URL.Query()["id"]
	email := r.URL.Query()["email"]
	if id != nil {
		id, err := primitive.ObjectIDFromHex(id[0])
		checkErr(err)

		filter = bson.M{"_id": id}
	} else if email != nil {
		filter = bson.M{"email": email[0]}
	} else {
		// By Cookie
		emailC, _ := getEmailType(r)
		filter = bson.M{"email": emailC}
	}
	var data interface{}
	var err error
	if filter != nil {
		var user model.User
		err = db.Database("tutuplapak").Collection("user").FindOne(ctx, filter).Decode(&user)
		data = user
	} else {
		// Get All
		cursor, err := db.Database("tutuplapak").Collection("user").Find(ctx, bson.D{})
		checkErr(err)
		var users []model.User
		err = cursor.All(ctx, &users)
		checkErr(err)
		data = users
	}

	if checkErr(err) {
		sendResponseData(w, 400, "Get Failed!", nil)
	} else {
		sendResponseData(w, 200, "Get Success!", data)
	}
}

// Insert user untuk registrasi user
// Alamat dalam bentuk objek
func InsertUser(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	var userData model.User
	err := json.NewDecoder(r.Body).Decode(&userData)
	checkErr(err)

	// Check email if exist
	var user model.User
	var opt options.FindOneOptions
	opt.SetProjection(bson.M{"email": 1})
	err = db.Database("tutuplapak").Collection("user").FindOne(ctx, bson.M{"email": userData.Email}, &opt).Decode(&user)
	if err == nil {
		if userData.Email == user.Email {
			sendResponseData(w, 200, "Email already exist!", nil)
			return
		}
	}

	_, err = db.Database("tutuplapak").Collection("user").InsertOne(ctx, user)
	if checkErr(err) {
		sendResponseData(w, 400, "Registration Failed!", nil)
	} else {
		sendResponseData(w, 200, "Registration Success!", nil)
		AddUserLog(userData.Email, "R")
	}

}

// Update data user. Harus login dulu.
// Alamat selalu terupdate, data sebelumnya harus dimasukkin ke form
// biar kalo alamatnya tidak diupdate sama user ga akan kosong
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	email, _ := getEmailType(r)
	if email == "" {
		sendUnAuthorizedResponse(w)
		return
	}

	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)

	_, err = db.Database("tutuplapak").Collection("user").UpdateOne(ctx, bson.M{"email": email}, bson.D{{"$set", user}})
	if checkErr(err) {
		sendResponseData(w, 400, "Update Failed!", nil)
	} else {
		sendResponseData(w, 200, "Update Success!", nil)
		AddUserLog(email, "U")
	}
}

// Post log_user tentang:
// LI - login
// LO - logout
// R - register
// U - update
func AddUserLog(email string, tipe string) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	log_user := model.UserLog{
		Email: email,
		Waktu: time.Now(),
		Tipe:  tipe,
	}
	_, err := db.Database("tutuplapak").Collection("log_user").InsertOne(ctx, log_user)
	checkErr(err)
}

// Read log_user hanya dari admin.
// Dapat mem-filter tanggal waktu <from>,<to> dan/atau tipe <tipe1>,<tipe2>,...
func ReadUserLog(w http.ResponseWriter, r *http.Request) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	filter := bson.D{}
	waktu := r.URL.Query().Get("waktu")
	tipe := r.URL.Query().Get("tipe")
	if waktu != "" {
		wArr := strings.Split(waktu, ",")
		w1, _ := time.Parse("2006-01-02 15:04:05", wArr[0])
		w2, _ := time.Parse("2006-01-02 15:04:05", wArr[1]+" 23:59:59")
		filter = append(filter, bson.E{"waktu", bson.M{"$gte": w1, "$lte": w2}})
	}
	if tipe != "" {
		tArr := strings.Split(tipe, ",")
		filter = append(filter, bson.E{"tipe", bson.M{"$in": tArr}})
	}

	var logs []model.UserLog
	cursor, err := db.Database("tutuplapak").Collection("log_user").Find(ctx, filter)
	if err == nil {
		err = cursor.All(ctx, &logs)
	}
	if checkErr(err) {
		sendResponseData(w, 400, "Get Failed!", nil)
	} else {
		sendResponseData(w, 200, "Get Success!", logs)
	}
}
