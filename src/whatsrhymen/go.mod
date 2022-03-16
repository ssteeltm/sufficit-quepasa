module github.com/sufficit/sufficit-quepasa-fork/whatsrhymen

require (
	github.com/sufficit/sufficit-quepasa-fork/whatsapp v0.0.0-00010101000000-000000000000 // indirect
)

require (
	github.com/Rhymen/go-whatsapp v0.1.1
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/sufficit/sufficit-quepasa-fork/whatsrhymen => ./

replace github.com/sufficit/sufficit-quepasa-fork/whatsapp => ../whatsapp

go 1.17
