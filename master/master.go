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
	"strconv"
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
		conn.Database("originalImages"),
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

	conn := InitiateMongoClient()

	// For CRUD operations, here is an example
	db := conn.Database("convertedImages")
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
	ioutil.WriteFile("./masterOutput/"+fileName, buf.Bytes(), 0600)

}

func up() int {

	var wg sync.WaitGroup

	files, err := ioutil.ReadDir("./masterInput")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		wg.Add(1)

		filename := "./masterInput/" + path.Base(f.Name())

		//fmt.Println(filename + " -> MongoDB")
		go UploadFile(filename, f.Name(), &wg)

	}
	//fmt.Printf("Numero de goroutines: %d\n", runtime.NumGoroutine())
	wg.Wait()

	return len(files)

}

func down() {

	var wg sync.WaitGroup

	//files := []string{"image_1.jpg", "image_2.jpg", "image_3.jpg", "image_4.jpg", "image_5.jpg"}

	files, err := ioutil.ReadDir("./masterInput")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		wg.Add(1)

		go DownloadFile(f.Name(), &wg)
	}

	//fmt.Printf("Numero de goroutines: %d\n", runtime.NumGoroutine())
	wg.Wait()

}

func deleteDatabases() {
	conn := InitiateMongoClient()
	db := conn.Database("originalImages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db.Drop(ctx)

	db = conn.Database("convertedImages")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db.Drop(ctx)

}

func createWorkers(qtdWorkers int) {

	arg0 := "./executerWorkerContainer.sh"

	cmd := exec.Command(arg0)
	//fmt.Println(cmd.String())

	errCmd := cmd.Run()

	if errCmd != nil {
		fmt.Println(errCmd)
	}
}

func clearMaster() {

	arg0 := "./clearMaster.sh"

	cmd := exec.Command(arg0)
	//fmt.Println(cmd.String())

	errCmd := cmd.Run()

	if errCmd != nil {
		fmt.Println(errCmd)
	}
}

func master(arg string, wg *sync.WaitGroup) {

	//arg := os.Args[1]
	defer wg.Done()

	switch arg {
	case "up":

		unzip("banco.zip")

		e := os.Remove("banco.zip")
		if e != nil {
			log.Fatal(e)
		}

		tamanhoBanco := up()

		fmt.Println("Quantidade de imagens: " + strconv.Itoa(tamanhoBanco))
		//cria os workers
		qtdWokers := tamanhoBanco / 10

		if qtdWokers == 0 {
			qtdWokers = 1
		}

		fmt.Println("Quantidade de workers: " + strconv.Itoa(qtdWokers))
		createWorkers(qtdWokers)

	case "down":
		down()
		deleteDatabases()
		if err := zipSource("masterOutput", "banco.zip"); err != nil {
			log.Fatal(err)
		}

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
