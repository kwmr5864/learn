package main
import (
	"net/http"
	"html/template"
)

type Page struct {
	Title string
	Body string
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, Page{
		Title:r.URL.Path[1:],
		Body:"こんにちは, " + r.URL.Path[1:],
	})
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.ListenAndServe(":8080", nil)
}
