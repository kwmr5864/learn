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
	Keyword string
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

func searchWordViewHandler(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Path[6:]

	db, err := sql.Open("sqlite3", "./ejdict.sqlite3")
	checkErr(err)

	row := db.QueryRow("SELECT COUNT(*) FROM items WHERE word LIKE ?", "%" + keyword + "%")
	var count int
	row.Scan(&count)

	rows, err := db.Query("SELECT item_id, word, mean, level FROM items WHERE word LIKE ?", "%" + keyword + "%")
	checkErr(err)

	page := Page{
		Title:keyword,
		Keyword:keyword,
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

func searchMeanViewHandler(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Path[6:]

	db, err := sql.Open("sqlite3", "./ejdict.sqlite3")
	checkErr(err)

	row := db.QueryRow("SELECT COUNT(*) FROM items WHERE mean LIKE ?", "%" + keyword + "%")
	var count int
	row.Scan(&count)

	rows, err := db.Query("SELECT item_id, word, mean, level FROM items WHERE mean LIKE ?", "%" + keyword + "%")
	checkErr(err)

	page := Page{
		Title:keyword,
		Keyword:keyword,
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
	http.HandleFunc("/word/", searchWordViewHandler)
	http.HandleFunc("/mean/", searchMeanViewHandler)
	http.ListenAndServe(":8080", nil)
}
