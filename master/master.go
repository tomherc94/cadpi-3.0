package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var nameDBlist []string

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
func UploadFile(file, filename string, nameDB string) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	conn := InitiateMongoClient()
	bucket, err := gridfs.NewBucket(
		conn.Database(nameDB),
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

func up(file string, wg *sync.WaitGroup, nameDB string) {

	defer wg.Done()

	/*files, err := ioutil.ReadDir("./masterInput")
	if err != nil {
		log.Fatal(err)
	}*/

	/*for _, f := range files {
		wg.Add(1)

		filename := "./masterInput/" + path.Base(f)

		//fmt.Println(filename + " -> MongoDB")
		go UploadFile(filename, f, &wg)

	}*/

	filename := "./masterInput/" + file

	UploadFile(filename, file, nameDB)

	//fmt.Printf("Numero de goroutines: %d\n", runtime.NumGoroutine())

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

	for _, nameDB := range nameDBlist {

		db := conn.Database(nameDB)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		db.Drop(ctx)
	}

	db := conn.Database("convertedImages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db.Drop(ctx)

}

func analyzeDB() int {

	var wg sync.WaitGroup

	//var listTotal []string
	listTotal := make(map[string]int64)

	files, err := ioutil.ReadDir("./masterInput")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		listTotal[f.Name()] = f.Size()
		//listTotal = append(listTotal, f.Name())
	}

	listTotalOrdenado := sortList(listTotal)

	fmt.Println("Quantidade total de imagens: " + strconv.Itoa(len(listTotal)))

	qtdWorkers := len(listTotal) / 10

	//qtdWorkers := 1

	if qtdWorkers == 0 {
		qtdWorkers = 1
	}

	fmt.Println("Quantidade de workers: " + strconv.Itoa(qtdWorkers))

	i := 1

	//Distribui a base de imagens entre os buckets
	for _, imageName := range listTotalOrdenado {

		wg.Add(1)

		nameDBcurrent := "worker" + strconv.Itoa(i)

		nameDBlist = append(nameDBlist, nameDBcurrent)

		go up(imageName, &wg, nameDBcurrent)

		i++

		if i > qtdWorkers {
			i = 1
		}

	}

	wg.Wait()

	return (qtdWorkers)
}

func createWorkers(qtdWorkers int) {

	var wg sync.WaitGroup

	for i := 1; i <= qtdWorkers; i++ {

		wg.Add(1)

		nameDBcurrent := "worker" + strconv.Itoa(i)

		arg0 := "./executerWorkerContainer.sh"

		arg1 := nameDBcurrent

		cmd := exec.Command(arg0, arg1)
		//fmt.Println(cmd.String())

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			errCmd := cmd.Run()

			if errCmd != nil {
				fmt.Println(errCmd)
			}
		}(&wg)

		fmt.Println("CRIADO WORKER" + strconv.Itoa(i))

	}
	wg.Wait()
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

		qtdWorkers := analyzeDB()

		createWorkers(qtdWorkers)

	case "down":
		down()
		deleteDatabases()
		if err := zipSource("masterOutput", "banco.zip"); err != nil {
			log.Fatal(err)
		}
		clearMaster()

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
