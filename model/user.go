package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email,omitempty" form:"email" json:"email"`
	Password     string             `bson:"password,omitempty" form:"password"`
	Nama         string             `bson:"nama,omitempty" form:"nama" json:"nama"`
	Alamat       Alamat             `bson:"alamat,omitempty" json:"alamat"`
	JenisKelamin string             `bson:"jenis_kelamin,omitempty" form:"jenis_kelamin" json:"jenis_kelamin"`
	Telepon      string             `bson:"telepon,omitempty" form:"telepon" json:"telepon"`
	Tipe         string             `bson:"tipe,omitempty" form:"tipe" json:"tipe"`
}

type Alamat struct {
	Jalan     string `bson:"jalan,omitempty" form:"jalan" json:"jalan"`
	Rt        string `bson:"rt,omitempty" form:"rt" json:"rt"`
	Rw        string `bson:"rw,omitempty" form:"rw" json:"rw"`
	Desa      string `bson:"desa,omitempty" form:"desa" json:"desa"`
	Kelurahan string `bson:"kelurahan,omitempty" form:"kelurahan" json:"kelurahan"`
	Kecamatan string `bson:"kecamatan,omitempty" form:"kecamatan" json:"kecamatan"`
	Kota      string `bson:"kota,omitempty" form:"kota" json:"kota"`
	Provinsi  string `bson:"provinsi,omitempty" form:"provinsi" json:"provinsi"`
}

type UserLog struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email string             `bson:"email,omitempty" json:"email"`
	Waktu time.Time          `bson:"waktu,omitempty" json:"waktu"`
	Tipe  string             `bson:"tipe,omitempty" json:"tipe"`
}
