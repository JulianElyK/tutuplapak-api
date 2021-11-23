package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Barang struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at,omitempty" form:"created_at" json:"created_at,omitempty"`
	DeletedAt time.Time          `bson:"deleted_at,omitempty" form:"deleted_at" json:"deleted_at,omitempty"`
	Nama      string             `bson:"nama,omitempty" form:"nama" json:"nama"`
	Harga     int                `bson:"harga,omitempty" form:"harga" json:"harga"`
	Penjual   string             `bson:"penjual,omitempty" form:"penjual" json:"penjual"`
	Stok      int                `bson:"stok,omitempty" form:"stok" json:"stok,omitempty"`
	Kategori  []string           `bson:"kategori,omitempty" form:"kategori" json:"kategori,omitempty"`
	Foto      string            `bson:"foto,omitempty" form:"foto" json:"foto,omitempty"`
}