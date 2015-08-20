package main
import (
	"net/http"
	"html/template"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	ItemId int
	Word string
	Mean string
	Level int
}

type Page struct {
	Title string
	Word string
	Count int
	Items []Item
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func indexViewHandler(w http.ResponseWriter, r *http.Request) {
	count := 10

	db, err := sql.Open("sqlite3", "./ejdict.sqlite3")
	checkErr(err)

	rows, err := db.Query("SELECT item_id, word, mean, level FROM items ORDER BY random() LIMIT ?", count)
	checkErr(err)

	page := Page{
		Items: make([]Item, count),
	}
	index := 0
	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.ItemId, &item.Word, &item.Mean, &item.Level)
		checkErr(err)
		page.Items[index] = item
		index++
	}
	db.Close()

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, page)
}

func searchViewHandler(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Path[8:]

	db, err := sql.Open("sqlite3", "./ejdict.sqlite3")
	checkErr(err)

	row := db.QueryRow("SELECT COUNT(*) FROM items WHERE word LIKE ?", "%" + word + "%")
	var count int
	row.Scan(&count)

	rows, err := db.Query("SELECT item_id, word, mean, level FROM items WHERE word LIKE ?", "%" + word + "%")
	checkErr(err)

	page := Page{
		Title:word,
		Word:word,
		Count:count,
		Items: make([]Item, count),
	}
	index := 0
	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.ItemId, &item.Word, &item.Mean, &item.Level)
		checkErr(err)
		page.Items[index] = item
		index++
	}
	db.Close()

	t, _ := template.ParseFiles("templates/detail.html")
	t.Execute(w, page)
}

func main() {
	http.HandleFunc("/", indexViewHandler)
	http.HandleFunc("/search/", searchViewHandler)
	http.ListenAndServe(":8080", nil)
}
