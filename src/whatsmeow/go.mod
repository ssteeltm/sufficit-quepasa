module github.com/sufficit/sufficit-quepasa/whatsmeow

require github.com/sufficit/sufficit-quepasa/whatsapp v0.0.0-00010101000000-000000000000

require (
	github.com/gosimple/slug v1.13.1 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
)

require (
	filippo.io/edwards25519 v1.0.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.11
	github.com/sirupsen/logrus v1.8.1
	go.mau.fi/libsignal v0.0.0-20220628090436-4d18b66b087e // indirect
	go.mau.fi/whatsmeow v0.0.0-20220811191500-f650c10b0068
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/sufficit/sufficit-quepasa/whatsmeow => ./

replace github.com/sufficit/sufficit-quepasa/whatsapp => ../whatsapp

go 1.17
