package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/fabianMendez/jarlib/core"
)

// Generate returns an http handler to generate a self-contained jar of a dependency
func Generate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)

		dependency := path.Base(r.URL.Path)
		q := r.URL.Query()
		javaVersion := q.Get("javaVersion")
		projectName := q.Get("projectName")

		tmpfile, err := os.CreateTemp("", "*")
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error" "could not create temp file"}`))
			return
		}
		defer func() {
			tmpfile.Close()
			os.Remove(tmpfile.Name())
		}()

		err = core.Generate(dependency, javaVersion, tmpfile)
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error" "could not generate library"}`))
			return
		}

		w.Header().Set("Content-Type", "application/java-archive")
		if projectName != "" {
			w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename="%s.jar"`, strings.ReplaceAll(projectName, `"`, "")))
		}
		w.WriteHeader(http.StatusOK)

		tmpfileread, err := os.Open(tmpfile.Name())
		if err != nil {
			log.Println(err)
		}
		_, err = io.Copy(w, tmpfileread)
		if err != nil {
			log.Println(err)
		}
	}
}
