module github.com/bp1222/photoSync

go 1.18

require (
	github.com/bp1222/tinybeans-api/go-client v0.1.1
	github.com/joho/godotenv v1.4.0
	github.com/sirupsen/logrus v1.9.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/sqlite v1.2.6
	gorm.io/gorm v1.22.5
)

replace github.com/bp1222/tinybeans-api/go-client => ../tinybeans-api/go-client

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/mattn/go-sqlite3 v1.14.10 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	golang.org/x/net v0.0.0-20220114011407-0dd24b26b47d // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
