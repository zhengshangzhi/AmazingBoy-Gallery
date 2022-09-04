package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"pkg/model"
)

var db *gorm.DB

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/index.html", "template/footer.html")
	log.Println(t.Execute(w, nil))
}

func updateDB(w http.ResponseWriter, r *http.Request) {
	err := model.UpdateAlbums(db)
	if err != nil {
		fmt.Fprint(w, "Fail to update albums!!")
	}

	err = model.ScanFolder(db)
	if err != nil {
		fmt.Fprintf(w, "Fail to update photos!!")
	}

	fmt.Fprint(w, "Update Success!!")
}

func gallery(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		FirstAlbum  model.Album
		SecondAlbum model.Album
		ThirdAlbum  model.Album
		RestAlbums  []model.Album
	}
	var albums []model.Album
	_ = db.Find(&albums)
	data := Data{FirstAlbum: albums[0], SecondAlbum: albums[1], ThirdAlbum: albums[2], RestAlbums: albums[3:]}
	t, _ := template.ParseFiles("template/gallery.html")
	log.Println(t.Execute(w, data))
}

func album(w http.ResponseWriter, r *http.Request) {
	var photos []model.Photo
	albumID, _ := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	db.Where("album_id = ?", albumID).Find(&photos)
	for _, photo := range photos {
		photo.Path = photo.Path[1:]
	}
	t, _ := template.ParseFiles("template/album.html", "template/footer.html", "template/header.html")
	log.Println(t.Execute(w, photos))
}

func equipment(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/equipment.html", "template/footer.html", "template/header.html")
	log.Println(t.Execute(w, nil))
}

func research(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/research.html", "template/footer.html", "template/header.html")
	log.Println(t.Execute(w, nil))
}

func main() {
	db = model.ConnectDB()

	http.Handle("/static/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", index)
	http.HandleFunc("/update", updateDB)
	http.HandleFunc("/gallery", gallery)
	http.HandleFunc("/album", album)
	http.HandleFunc("/equipment", equipment)
	http.HandleFunc("/research", research)
	err := http.ListenAndServe(":3003", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
