package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// Serve assets
func Fs(root http.Dir) http.Handler {
	fs := http.FileServer(root)
	fn := func(w http.ResponseWriter, r *http.Request) {
		cleanPath := filepath.Clean(r.URL.Path)
		fullPath := filepath.Join(string(root), cleanPath)

		fmt.Println("Serving file: ", fullPath)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			fmt.Println("file not found: ", err)
			http.NotFound(w, r)
			return
		}
		fs.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
