package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"ebase64/pkg/app"
	"ebase64/pkg/exec"
)

const MaxBodySize = 2 << 20

func Index(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tpl, err := template.ParseFiles("templates/base.html",
				"templates/head.html",
				"templates/topmenu.html",
				"templates/footer.html",
				"templates/index.html")
			if err != nil {
				log.Printf("error parsing template files: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = tpl.Execute(w, "base.html")
			if err != nil {
				log.Printf("error executing template files: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

func Exec(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)

			err := r.ParseForm()
			if err != nil {
				log.Printf("error form parsing: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			radioButton := r.FormValue("base_radio")
			if (radioButton != "encode") && (radioButton != "decode") {
				log.Printf("unknown radioButton choice: %s", err)
				w.WriteHeader(http.StatusTeapot)
				return
			}

			decodeAction := false
			if radioButton == "decode" {
				decodeAction = true
			}

			var str string

			if decodeAction {
				str, err = exec.DecodeURL(r.FormValue("base64_body"))
				if err != nil {
					log.Printf("error decoding string: %s", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			} else {
				str = exec.EncodeURL(r.FormValue("base64_body"))
			}

			tpl, err := template.ParseFiles("templates/exec.html")
			if err != nil {
				log.Printf("error parsing template files: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = tpl.Execute(w, str)
			if err != nil {
				log.Printf("error executing template files: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

func Static(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir("static"))
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	}
}

func Healthz(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func Readyz(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

// LoggingMiddleware wraps every request with logging data.
func LoggingMiddleware(app *app.Application, next http.Handler) http.Handler {
	start := time.Now()
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			app.Sugar.Infow("Request processed",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
			)
		})
}
