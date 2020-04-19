// Presents a page to upload several files, and stores them in the downloads directory
package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

// compiling/caching the template
var templates = template.Must(template.New("tmpl").Parse(`
<html>
  <head>
    <title>File Upload Demo</title>
   <style>
   body {
        font-family: Sans-serif;
        padding-top: 40px;
        padding-bottom: 40px;
        background-color: #ffffff;
   }
   h1 {text-align: center; margin-bottom: 30px;}
   .message {font-weight:bold}
   fieldset {width:50%}
   </style>
  </head>
  <body>
    <div class="container">
      <h1>File Upload Demo</h1>
      {{range .}}<div class="message">{{.}} successful!</div>{{end}}
      <form class="form-signin" method="post" action="/upload" enctype="multipart/form-data">
          <fieldset>
            <input type="file" name="myfiles" id="myfiles" multiple="multiple">
            <input type="submit" name="submit" value="Submit">
        </fieldset>
      </form>
    </div>
  </body>
</html>
`))

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// GET to display the upload form.
	case "GET":
		err := templates.Execute(w, nil)
		if err != nil {
			log.Print(err)
		}
		// POST analyzes each part of the MultiPartReader (ie the uploaded file(s))
		// and saves them to disk.
	case "POST":
		// grab the request.MultipartReader
		reader, err := r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var successfulFiles []string
		// copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			// if part.FileName() is empty, skip this iteration.
			if part.FileName() == "" {
				continue
			}

			// prepare the dst
			dst, err := os.Create("./downloads/" + part.FileName())
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// copy the part to dst
			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			successfulFiles = append(successfulFiles, part.FileName())
			log.Print("Received " + part.FileName())
		}
		// displaying a success message.
		err = templates.Execute(w, successfulFiles)
		if err != nil {
			log.Print(err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func prepareDir() {
	if _, err := os.Stat("./downloads"); os.IsNotExist(err) {
		if err := os.MkdirAll("./downloads", 0755); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	prepareDir()
	http.HandleFunc("/", uploadHandler)
	log.Print("Listening on port:8000...")
	// Listen on port 8080
	http.ListenAndServe(":8000", nil)
}
