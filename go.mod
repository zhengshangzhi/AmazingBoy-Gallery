module Z-Gallery

go 1.18

require (
	gorm.io/driver/sqlite v1.3.6 // indirect
	gorm.io/gorm v1.23.8
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.14 // indirect
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd // indirect
)

require pkg/model v0.0.0

replace pkg/model => ./pkg/model
