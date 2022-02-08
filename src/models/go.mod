module github.com/sufficit/sufficit-quepasa-fork/models

require (
	github.com/go-chi/jwtauth v4.0.4+incompatible
	github.com/golang-migrate/migrate/v4 v4.11.0
	github.com/google/uuid v1.1.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.5.2
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa
	filippo.io/edwards25519 v1.0.0-rc.1 // indirect
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/go-chi/chi v4.0.2+incompatible // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/mattn/go-sqlite3 v1.14.11 // indirect
	github.com/mongodb/mongo-go-driver v0.3.0 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	go.mau.fi/libsignal v0.0.0-20211109153248-a67163214910 // indirect
	go.mau.fi/whatsmeow v0.0.0-20220128124639-e64fb976bf15 // indirect
	golang.org/x/time v0.0.0-20190921001708-c4c64cad1fd0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)

replace "github.com/sufficit/sufficit-quepasa-fork/models" => "./"

go 1.17
