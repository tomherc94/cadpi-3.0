package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitiateMongoClient() *mongo.Client {
	var err error
	var client *mongo.Client
	//uri := "mongodb://localhost:27017"
	uri := "mongodb://root:example@localhost:27017/"
	opts := options.Client()
	opts.ApplyURI(uri)
	opts.SetMaxPoolSize(5)
	if client, err = mongo.Connect(context.Background(), opts); err != nil {
		fmt.Println(err.Error())
	}
	return client
}
func UploadFile(file, filename string, wg *sync.WaitGroup) {

	defer wg.Done()

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	conn := InitiateMongoClient()
	bucket, err := gridfs.NewBucket(
		conn.Database("myfiles"),
	)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	uploadStream, err := bucket.OpenUploadStream(
		filename,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(data)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Printf("Write file to DB was successful. File size: %d M\n", fileSize)
}
func DownloadFile(fileName string) {
	conn := InitiateMongoClient()

	// For CRUD operations, here is an example
	db := conn.Database("myfiles")
	fsFiles := db.Collection("fs.files")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var results bson.M
	err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results)
	if err != nil {
		log.Fatal(err)
	}
	// you can print out the results
	fmt.Println(results)

	bucket, _ := gridfs.NewBucket(
		db,
	)
	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStreamByName(fileName, &buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File size to download: %v\n", dStream)
	ioutil.WriteFile(fileName, buf.Bytes(), 0600)

}

func main() {
	now := time.Now()
	var wg sync.WaitGroup

	files, err := ioutil.ReadDir("./masterInput")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		wg.Add(1)

		filename := "./masterInput/" + path.Base(f.Name())

		//fmt.Println(filename + " -> MongoDB")
		go UploadFile(filename, filename, &wg)

	}
	fmt.Printf("Numero de goroutines: %d\n", runtime.NumGoroutine())
	wg.Wait()

	defer func() {
		fmt.Println()
		fmt.Println("RELATÓRIO")
		fmt.Println("Quantidade de imagens: ", len(files))
		fmt.Println("Tempo de execução: ", time.Since(now))
	}()

	/*// Get os.Args values
	file := os.Args[1] //os.Args[1] = testfile.zip
	filename := path.Base(file)
	UploadFile(file, filename)
	// Uncomment the below line and comment the UploadFile above this line to download the file
	//DownloadFile(filename)
	*/
}
