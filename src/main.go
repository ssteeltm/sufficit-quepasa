package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	controllers "github.com/sufficit/sufficit-quepasa-fork/controllers"
	models "github.com/sufficit/sufficit-quepasa-fork/models"
	whatsmeow "github.com/sufficit/sufficit-quepasa-fork/whatsmeow"

	log "github.com/sirupsen/logrus"
)

func main() {

	// Carregando variaveis de ambiente apartir de arquivo .env
	godotenv.Load()

	if models.ENV.DEBUGJsonMessages() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Verifica se é necessario realizar alguma migração de base de dados
	err := models.MigrateToLatest()
	if err != nil {
		log.Fatalf("Database migration error: %s", err.Error())
	}

	whatsmeow.WhatsmeowService.Start()

	// Inicializando serviço de controle do whatsapp
	// De forma assíncrona
	err = models.QPWhatsappStart()
	if err != nil {
		panic(err.Error())
	}

	go func() {
		m := chi.NewRouter()
		m.Handle("/metrics", promhttp.Handler())
		host := fmt.Sprintf("%s:%s", os.Getenv("METRICS_HOST"), os.Getenv("METRICS_PORT"))

		log.Println("Starting Metrics Service")
		err := http.ListenAndServe(host, m)
		if err != nil {
			log.Fatal(err)
		}
	}()

	controllers.QPWebServerStart()
}
