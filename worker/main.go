package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var nameDB, _ = os.Hostname()
var conn = InitiateMongoClient()

func InitiateMongoClient() *mongo.Client {
	var err error
	var client *mongo.Client
	//uri := "mongodb://root:example@localhost:27017/"
	uri := "mongodb://root:example@mongo_container:27017/"
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
	//conn := InitiateMongoClient()
	bucket, err := gridfs.NewBucket(
		conn.Database("convertedImages"),
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
func DownloadFile(fileName string, wg *sync.WaitGroup) {

	defer wg.Done()

	//conn := InitiateMongoClient()

	// For CRUD operations, here is an example
	db := conn.Database(nameDB)
	fsFiles := db.Collection("fs.files")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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
	ioutil.WriteFile("./workerInput/"+fileName, buf.Bytes(), 0600)

}

func up() {

	var wg sync.WaitGroup

	files, err := ioutil.ReadDir("./workerOutput")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		wg.Add(1)

		filename := "./workerOutput/" + path.Base(f.Name())

		//fmt.Println(filename + " -> MongoDB")
		go UploadFile(filename, f.Name(), &wg)

	}

	wg.Wait()

}

func down(files []string) {

	var wg sync.WaitGroup

	//files := []string{"image_1.jpg", "image_2.jpg", "image_3.jpg", "image_4.jpg", "image_5.jpg"}

	for _, f := range files {
		wg.Add(1)

		go DownloadFile(f, &wg)
	}

	wg.Wait()

}

func workerApp() {

	arg0 := "./executeWorkerApp.sh"

	cmd := exec.Command(arg0)
	//fmt.Println(cmd.String())

	errCmd := cmd.Run()

	if errCmd != nil {
		fmt.Println(errCmd)
	}

}

func main() {

	//arg := os.Args[1]
	arg := "down"

	var files []string

	//conn := InitiateMongoClient()
	db := conn.Database(nameDB)

	fsFiles := db.Collection("fs.files")

	cursor, _ := fsFiles.Find(
		context.TODO(),
		bson.D{},
	)

	for cursor.Next(context.TODO()) {
		var result bson.D

		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}

		files = append(files, fmt.Sprint(result.Map()["filename"]))
		//fmt.Println(result.Map()["filename"])
		fmt.Println("Tamanho do banco: ")
		fmt.Print(len(files))

	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())

	switch arg {
	case "up":
		up()

	case "down":
		now := time.Now()
		fmt.Println("Download de imagens do BD ...")
		down(files)
		time.Sleep(1 * time.Second)
		fmt.Println("Aplicativo JAVA ...")
		workerApp()
		time.Sleep(1 * time.Second)
		fmt.Println("Upload de imagens para o BD ...")
		up()
		time.Sleep(1 * time.Second)
		fmt.Println("Tempo de execução: ", time.Since(now))

	default:
		log.Fatal("Parametro incorreto")
	}

	/*// Get os.Args values
	file := os.Args[1] //os.Args[1] = testfile.zip
	filename := path.Base(file)
	UploadFile(file, filename)
	// Uncomment the below line and comment the UploadFile above this line to download the file
	//DownloadFile(filename)
	*/
}
