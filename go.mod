module git.01.kood.tech/roosarula/forum

// +heroku goVersion go1.17
go 1.18

require github.com/mattn/go-sqlite3 v1.14.10

require golang.org/x/crypto v0.0.0-20220126234351-aa10faf2a1f8

require github.com/gorilla/websocket v1.5.0

require (
	github.com/golang-migrate/migrate v3.5.4+incompatible
	github.com/satori/go.uuid v1.2.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
