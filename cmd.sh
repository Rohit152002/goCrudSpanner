go get -u github.com/gin-gonic/gin
go get -u go.uber.org/zap
go get github.com/joho/godotenv

go get -u gorm.io/gorm

go get -u github.com/googleapis/go-sql-spanner
go get -u github.com/googleapis/go-gorm-spanner
go get -u golang.org/x/crypto/bcrypt
go get github.com/DATA-DOG/go-sqlmock
go get gorm.io/driver/mysql
go get github.com/stretchr/testify/assert
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
