package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var t = template.Must(template.ParseFiles("index.html"))

func Handler(w http.ResponseWriter, r *http.Request) {
	mw := io.MultiWriter(os.Stdout, w)

	t.Execute(w, nil)

	if r.Method == "POST" {
		// ParseMultipartForm parses a request body as multipart/form-data.
		// The whole request body is parsed and up to a total of maxMemory bytes of
		// its file parts are stored in memory, with the remainder stored on
		// disk in temporary files.
		// ParseMultipartForm calls ParseForm if necessary.
		// After one call to ParseMultipartForm, subsequent calls have no effect.
		if err := r.ParseMultipartForm(200000); err != nil {
			fmt.Fprintln(mw, err)
			return
		}

		// Read the form data
		//
		// MultipartForm is the parsed multipart form, including file uploads.
		// This field is only available after ParseMultipartForm is called.
		// The HTTP client ignores MultipartForm and uses Body instead.
		formData := r.MultipartForm

		// Get the *FileHeaders
		//
		// A FileHeader describes a file part of a multipart request.
		fileHeaders := formData.File["multiplefiles"]

		// Iterate over all of the FileHeaders
		for i, _ := range fileHeaders {
			// Open the file associated with the current FileHeader
			//
			// Open opens and returns the FileHeader's associated File.
			file, err := fileHeaders[i].Open()
			if err != nil {
				fmt.Fprintln(mw, err)
				return
			}
			defer file.Close()

			fmt.Fprintf(os.Stdout, "Opened file `%s` successfully!\n", fileHeaders[i].Filename)

			// Create a file to store the file
			outFile, err := os.Create("uploads/" + fileHeaders[i].Filename)
			if err != nil {
				fmt.Fprintln(mw, "Unable to create the file for writing. Check your write access privilege.\n")
				w.Write([]byte("<br />"))
				return
			}
			defer outFile.Close()

			// Store the file
			_, err = io.Copy(outFile, file)
			if err != nil {
				fmt.Fprintln(mw, err)
				return
			}

			fmt.Fprintf(mw, "%s uploaded successfully!\n\n", fileHeaders[i].Filename)
			w.Write([]byte("<br />"))
		}
	}
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Create an uploads folder to store all uploads
	if err := os.Mkdir("uploads", os.ModePerm); err == nil {
		fmt.Println("Created uploads directory.\n")
	}

	go func() {
		<-c
		GracefulExit()
		fmt.Println("Finished exiting cleanly.")
		os.Exit(1)
	}()

	http.HandleFunc("/", Handler)
	log.Fatalln(http.ListenAndServe(":9000", nil))
}

func GracefulExit() {
	err := os.RemoveAll("uploads")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Removed uploads directory and all of it's contents.")
}
