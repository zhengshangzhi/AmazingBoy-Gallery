package photo

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func ReadExif(fname string) {
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
	_, shutter, _ := ExposureTime.Rat2(0)
	Iso, _ := ISOSpeedRatings.Int(0)
	FocalLength, _, _ := Focal.Rat2(0)

	fmt.Println(Date)
	fmt.Println(Camera)
	fmt.Println(Aperture)
	fmt.Println(shutter)
	fmt.Println(Iso)
	fmt.Println(FocalLength)
	fmt.Println()
}
