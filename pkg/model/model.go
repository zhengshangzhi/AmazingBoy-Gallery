package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const PHOTOROOT = "./static/photo"

type Album struct {
	gorm.Model
	Name        string
	Type        string
	Description string
	Cover       string
}

type Albums struct {
	Albums []Album `gorm:"-"`
}

type Photo struct {
	gorm.Model
	Path        string
	Name        string
	Date        time.Time
	Camera      string
	Aperture    float32
	Shutter     int
	Iso         int
	FocalLength int
	Location    string `gorm:"default:"`
	AlbumID     uint
}

func (photo *Photo) ReadExif(fname string) {
	f, err := os.Open(fname)
	if err != nil {
		panic("can not open file")
	}

	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	DateTimeOriginal, _ := x.Get(exif.DateTimeOriginal)
	camModel, _ := x.Get(exif.Model)
	FNumber, _ := x.Get(exif.FNumber)
	ExposureTime, _ := x.Get(exif.ExposureTime)
	ISOSpeedRatings, _ := x.Get(exif.ISOSpeedRatings)
	Focal, _ := x.Get(exif.FocalLength)

	DateStr, _ := DateTimeOriginal.StringVal()
	Date, _ := time.Parse("2006:01:02 15:04:05", DateStr)
	Camera, _ := camModel.StringVal()
	ApertureNum, ApertureDen, _ := FNumber.Rat2(0)
	Aperture := float32(ApertureNum) / float32(ApertureDen)
	_, Shutter, _ := ExposureTime.Rat2(0)
	Iso, _ := ISOSpeedRatings.Int(0)
	FocalLength, _, _ := Focal.Rat2(0)

	photo.Date = Date
	photo.Camera = Camera
	photo.Aperture = Aperture
	photo.Shutter = int(Shutter)
	photo.Iso = Iso
	photo.FocalLength = int(FocalLength)
}

/*/func (photo *Photo) ReadBasicInfo(fname string) {
	s := strings.Split(fname, "_")
	photo.Path = fname
	photo.Name = strings.Split(s[1], ".")[0]
}/*/

func (photo *Photo) ReadInfo(fname string, albumName string) {
	//photo.ReadBasicInfo(fname)
	photo.Name = (strings.Split(fname, ".")[0])[3:]
	osPath := PHOTOROOT + "/" + albumName + "/" + fname
	photo.Path = osPath[1:]
	photo.ReadExif(osPath)
}

func ScanFolder(db *gorm.DB) error {
	folders, err := ioutil.ReadDir(PHOTOROOT)
	if err != nil {
		return err
	}

	var photoID uint = 1

	for _, folder := range folders {
		albumID, _ := strconv.ParseUint(folder.Name(), 10, 64)
		files, _ := ioutil.ReadDir(PHOTOROOT + "/" + folder.Name())
		for _, file := range files {
			var photo Photo
			photo.AlbumID = uint(albumID)
			photo.ID = photoID
			photoID++
			photo.ReadInfo(file.Name(), folder.Name())
			db.Save(&photo)
		}
	}
	return nil
}

func UpdateAlbums(db *gorm.DB) error {
	var albums Albums
	f, _ := os.ReadFile("album.json")
	err := json.Unmarshal(f, &albums)
	if err != nil {
		return err
	}

	for _, album := range albums.Albums {
		db.Save(&album)
	}
	return nil
}

func ConnectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("Gallery.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Album{}, &Photo{})

	return db
}
