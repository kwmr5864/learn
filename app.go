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

func viewHandler(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Path[1:]

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

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, page)
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.ListenAndServe(":8080", nil)
}
