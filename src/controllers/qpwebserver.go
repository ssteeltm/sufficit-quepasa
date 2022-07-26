package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sufficit/sufficit-quepasa-fork/models"
)

func QPWebServerStart() {
	r := newRouter()
	webAPIPort := os.Getenv("WEBAPIPORT")
	webAPIHost := os.Getenv("WEBAPIHOST")
	if len(webAPIPort) == 0 {
		webAPIPort = "31000"
	}

	log.Printf("Starting Web Server on Port: %s", webAPIPort)
	err := http.ListenAndServe(webAPIHost+":"+webAPIPort, r)
	if err != nil {
		log.Fatal(err)
	}
}

func NormalizePathsToLower(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "" {
			r.URL.Path = strings.ToLower(r.URL.Path)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func newRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(NormalizePathsToLower)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	shouldLog, _ := models.GetEnvBool("HTTPLOGS", false)
	if shouldLog {
		r.Use(middleware.Logger)
	}

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// web routes
	// authenticated web routes
	r.Group(RegisterFormAuthenticatedControllers)

	// unauthenticated web routes
	r.Group(RegisterFormControllers)

	// api routes
	addAPIRoutes(r)

	// static files
	workDir, _ := os.Getwd()
	assetsDir := filepath.Join(workDir, "assets")
	fileServer(r, "/assets", http.Dir(assetsDir))

	return r
}

func addAPIRoutes(r chi.Router) {
	r.Group(RegisterAPIControllers)
	r.Group(RegisterAPIV2Controllers)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
