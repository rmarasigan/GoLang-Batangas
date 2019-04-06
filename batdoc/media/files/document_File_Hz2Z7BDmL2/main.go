package main

import (
	"fmt"
	"net/http"
	"strings"
)

func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func PopHandler(w http.ResponseWriter, r *http.Request) {
	t := formatRequest(r)
	fmt.Println(t)
	fmt.Println("==============\n\n\n")
	guid := r.FormValue("ExtVar1")
	tickno := r.FormValue("ExtVar2")

	urlTicket1 := "http://localhost:4569/screenpop?t=report&id=115557&incidents.ref_no=%s-%s"
	urlTicket2 := "http://localhost:4569/screenpop?t=report&id=115439&contacts.c$p_gudid=%s"
	urlGuid := "http://localhost:4569/screenpop?t=report&id=115439&contacts.c$p_gudid=%s"

	if tickno != "" {
		if len(tickno) == 12 {
			http.Redirect(w, r, fmt.Sprintf(urlTicket1, tickno[0:6], tickno[6:12]), 303)
		} else {
			http.Redirect(w, r, fmt.Sprintf(urlTicket2, tickno), 303)
		}
		return
	}
	if guid != "" {
		http.Redirect(w, r, fmt.Sprintf(urlGuid, guid), 303)
		return
	}
	w.Write([]byte(t))
}

func main() {
	http.HandleFunc("/", PopHandler)
	http.ListenAndServeTLS(":8898", "/etc/letsencrypt/live/uadmin.io-0001/fullchain.pem", "/etc/letsencrypt/live/uadmin.io-0001/privkey.pem", nil)
}
