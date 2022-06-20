package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

var filenames = []string{"public/upload.html", "public/download.html"}

// Compile templates on start of the application
var templates = template.Must(template.ParseFiles(filenames...))

// Display the named template
func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	dst, err := os.Create(handler.Filename)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Aguardando processamento ...\n")

	//CHAMAR O MASTER.GO
	descompactarMasterInput()
	master("up")

	http.Redirect(w, r, "/download", http.StatusFound)

}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	Openfile, err := os.Open("./banco.zip") //Open the file to be downloaded later
	//Close after function return

	if err != nil {
		http.Error(w, "File not found.", 404) //return 404 if file is not found
		return
	}
	defer Openfile.Close()

	tempBuffer := make([]byte, 512)                       //Create a byte array to read the file later
	Openfile.Read(tempBuffer)                             //Read the file into  byte
	FileContentType := http.DetectContentType(tempBuffer) //Get file header

	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	Filename := "demo_download"

	//Set the headers
	w.Header().Set("Content-Type", FileContentType+";"+Filename)
	w.Header().Set("Content-Length", FileSize)

	Openfile.Seek(0, 0)  //We read 512 bytes from the file already so we reset the offset back to 0
	io.Copy(w, Openfile) //'Copy' the file to the client

	fmt.Fprintf(w, "Successfully Download File\n")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "upload", nil)
	case "POST":
		uploadFile(w, r)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "download", nil)
	case "POST":
		downloadFile(w, r)
	}
}

func main() {
	// Upload route
	http.HandleFunc("/upload", uploadHandler)

	http.HandleFunc("/download", downloadHandler)

	println(("Escutando na porta 8080"))

	//Listen on port 8080
	http.ListenAndServe(":8080", nil)

}
