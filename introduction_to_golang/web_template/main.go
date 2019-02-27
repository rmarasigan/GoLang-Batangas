package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// http.HandleFunc("/", HelloHandler) : We tell the server to use the handler
// HelloHandler to every request that hits the server.
func main() {
	fmt.Println("Starting Web Server")
	http.HandleFunc("/", HelloHandler)
	http.HandleFunc("/main", MainHandler)
	http.HandleFunc("/api/", APIHandler)
	http.HandleFunc("/ajax/", AJAXHandler)

	// Set listen Port
	http.ListenAndServe(":8080", nil)

}

// In this little code, we are creating a handler called HelloHandler
// which retrieves the path of the URL (first line), removes the first
// slash (second line) and appends the "Hello" to the beginning of the
// sentence.
// We write the final message to the ResponseWriter converted to bytes.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello" + message

	w.Write([]byte(message))
}

// MainHandler
func MainHandler(w http.ResponseWriter, r *http.Request) {
	c := map[string]interface{}{}
	c["Name"] = "Your name here"
	c["Food"] = []string{
		"pizza",
		"burgers",
	}

	c["Active"] = true
	c["Department"] = map[string]string{
		"Name":     "Development",
		"Location": "Integ2",
		"Manager":  "Maenard",
	}

	tmpl, _ := template.ParseFiles("template/index.html")
	tmpl.Execute(w, c)
}

func AJAXHandler(w http.ResponseWriter, r *http.Request) {
	c := map[string]interface{}{}
	c["Name"] = ""
	c["Food"] = []string{}

	tmpl, _ := template.ParseFiles("template/ajax.html")
	tmpl.Execute(w, c)
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	c := map[string]interface{}{}

	c["Name"] = "Please your name here"
	c["Food"] = []string{
		"Pizza",
		"Burgers",
	}

	c["Department"] = map[string]string{
		"Name":     "Development Department",
		"Location": "IntegrityNet II",
		"Manager":  "Maenard",
	}

	buf, _ := json.Marshal(c)
	w.Write(buf)
}
