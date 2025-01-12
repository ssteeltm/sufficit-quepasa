module github.com/sufficit/sufficit-quepasa/main

require (
	github.com/joho/godotenv v1.4.0
	github.com/sirupsen/logrus v1.8.1
	github.com/sufficit/sufficit-quepasa/controllers v0.0.0-00010101000000-00000000000
	github.com/sufficit/sufficit-quepasa/models v0.0.0-00010101000000-000000000000
	github.com/sufficit/sufficit-quepasa/whatsmeow v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.0.0 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-chi/chi v4.1.2+incompatible // indirect
	github.com/go-chi/chi/v5 v5.0.7 // indirect
	github.com/go-chi/jwtauth v4.0.4+incompatible // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.7 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gosimple/slug v1.13.1 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/jinzhu/copier v0.3.5 // indirect
	github.com/jmoiron/sqlx v1.3.5 // indirect
	github.com/joncalhoun/migrate v0.0.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/lib/pq v1.5.2 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/nbutton23/zxcvbn-go v0.0.0-20210217022336-fa2cb2858354 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086 // indirect
	github.com/sufficit/sufficit-quepasa/library v0.0.0-00010101000000-000000000000 // indirect
	github.com/sufficit/sufficit-quepasa/metrics v0.0.0-00010101000000-000000000000 // indirect
	github.com/sufficit/sufficit-quepasa/whatsapp v0.0.0-00010101000000-000000000000 // indirect
	github.com/swaggo/files v0.0.0-20220728132757-551d4a08d97a // indirect
	github.com/swaggo/http-swagger v1.3.3 // indirect
	github.com/swaggo/swag v1.8.5 // indirect
	go.mau.fi/libsignal v0.0.0-20220628090436-4d18b66b087e // indirect
	go.mau.fi/whatsmeow v0.0.0-20221008133908-7f01b3072802 // indirect
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90 // indirect
	golang.org/x/net v0.0.0-20220907135653-1e95f45603a7 // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
	golang.org/x/tools v0.1.12 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/sufficit/sufficit-quepasa/controllers => ./controllers

replace github.com/sufficit/sufficit-quepasa/library => ./library

replace github.com/sufficit/sufficit-quepasa/metrics => ./metrics

replace github.com/sufficit/sufficit-quepasa/models => ./models

replace github.com/sufficit/sufficit-quepasa/whatsapp => ./whatsapp

replace github.com/sufficit/sufficit-quepasa/whatsmeow => ./whatsmeow

go 1.17
