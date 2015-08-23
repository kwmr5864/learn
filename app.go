package main
import (
	"net/http"
	"html/template"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"encoding/json"
	"fmt"
	mgo "gopkg.in/mgo.v2"
//	"gopkg.in/mgo.v2/bson"
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

type Response struct {
	Result bool
	Message string
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

	t, _ := template.ParseFiles(
		"templates/index.html",
		"templates/parts/search_form.html",
	)
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

	t, _ := template.ParseFiles(
		"templates/detail.html",
		"templates/parts/search_form.html",
	)
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

	t, _ := template.ParseFiles(
		"templates/detail.html",
		"templates/parts/search_form.html",
	)
	t.Execute(w, page)
}

func addWordApiHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := mgo.Dial("mongodb://localhost/kawamura")
	defer session.Close()
	db := session.DB("kawamura")

	var response = Response{
		Result:true,
		Message:r.FormValue("itemId"),
	}
	col := db.C("items")
	col.Insert(response)

	var data, _ = json.Marshal(response)
	fmt.Fprintf(w, string(data))
}

func main() {
	http.HandleFunc("/", indexViewHandler)
	http.HandleFunc("/word/", searchWordViewHandler)
	http.HandleFunc("/mean/", searchMeanViewHandler)
	http.HandleFunc("/api/add", addWordApiHandler)
	http.ListenAndServe(":8080", nil)
}
