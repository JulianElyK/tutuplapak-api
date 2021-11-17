package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaksi struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at,omitempty" form:"created_at" json:"created_at"`
	Email     string             `bson:"email,omitempty" form:"email" json:"email"`
	Barang    []BarangTransaksi  `bson:"barang,omitempty" form:"barang" json:"barang"`
}

type BarangTransaksi struct {
	Item   Barang `bson:"item,omitempty" json:"item"`
	Jumlah int    `bson:"jumlah,omitempty" form:"barang" json:"jumlah"`
}
