package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	mu     sync.Mutex
	strings []string
}

func main() {
	server := &Server{strings: []string{}}

	http.HandleFunc("/", server.renderHTML)
	http.HandleFunc("/add", server.addString)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// renderHTML rendert die HTML-Seite und zeigt die Liste der Strings an.
func (s *Server) renderHTML(w http.ResponseWriter, r *http.Request) {
	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Log List</title>
	</head>
	<body>
		<h1>Log List</h1>
		<ul>
			{{range .}}
				<li>{{.}}</li>
			{{end}}
		</ul>
		<form action="/add" method="POST">
			<input type="text" name="newString" required>
			<button type="submit">Add String</button>
		</form>
	</body>
	</html>
	`
	t, err := template.New("page").Parse(tmpl)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	t.Execute(w, s.strings)
}

func (s *Server) addString(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	newString := r.FormValue("newString")
	s.mu.Lock()
	s.strings = append(s.strings, newString)
	s.mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
