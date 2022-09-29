module github.com/sufficit/sufficit-quepasa/models

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/sufficit/sufficit-quepasa/whatsapp v0.0.0-00010101000000-000000000000
	github.com/sufficit/sufficit-quepasa/whatsmeow v0.0.0-00010101000000-000000000000
)

require (
	github.com/Rhymen/go-whatsapp v0.1.1 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
)

require (
	filippo.io/edwards25519 v1.0.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-chi/jwtauth v4.0.4+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/jinzhu/copier v0.3.5
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.5.2
	github.com/mattn/go-sqlite3 v1.14.15
	github.com/sirupsen/logrus v1.8.1
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086
	github.com/stretchr/testify v1.4.0 // indirect
	go.mau.fi/libsignal v0.0.0-20220628090436-4d18b66b087e // indirect
	go.mau.fi/whatsmeow v0.0.0-20220811191500-f650c10b0068 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace github.com/sufficit/sufficit-quepasa/library => ../library

replace github.com/sufficit/sufficit-quepasa/whatsmeow => ../whatsmeow

replace github.com/sufficit/sufficit-quepasa/whatsapp => ../whatsapp

replace github.com/sufficit/sufficit-quepasa/models => ./

go 1.17
