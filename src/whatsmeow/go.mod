module github.com/sufficit/sufficit-quepasa-fork/whatsmeow

require (
	github.com/sufficit/sufficit-quepasa-fork/whatsapp v0.0.0-00010101000000-000000000000
	go.mau.fi/whatsmeow v0.0.0-20220215120744-a1550ccceb70
)

require (
	filippo.io/edwards25519 v1.0.0-rc.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/mattn/go-sqlite3 v1.14.11
	github.com/sirupsen/logrus v1.8.1
	go.mau.fi/libsignal v0.0.0-20211109153248-a67163214910 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/sufficit/sufficit-quepasa-fork/whatsmeow => ./

replace github.com/sufficit/sufficit-quepasa-fork/whatsapp => ../whatsapp

go 1.17
