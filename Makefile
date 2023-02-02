install-packages:
	go get github.com/gin-contrib/cors
	go get -u github.com/gin-gonic/gin
	go get github.com/golang-jwt/jwt
	go get github.com/satori/go.uuid
	go get github.com/spf13/viper
	go get gorm.io/driver/sqlite
	go get -u gorm.io/gorm