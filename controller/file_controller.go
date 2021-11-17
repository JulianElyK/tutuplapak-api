package controller

import (
	"io/ioutil"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func UploadFile(file, filename string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	db, ctx := connect()
	defer db.Disconnect(ctx)

	bucket, err := gridfs.NewBucket(db.Database("tutuplapak"))
	if err != nil {
		return err
	}

	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		return err
	}
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(data)
	if err != nil {
		return err
	}
	log.Printf("Write file was successful. File size: %d\n", fileSize)
	return nil
}

func DownloadFile(filename string) {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	var file bson.M
	if err := db.Database("tutuplapak").Collection("fs.files").FindOne(ctx, bson.M{}).Decode(&file); err != nil {
		log.Println(err)
	}
}
