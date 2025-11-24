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
		clean_path := filepath.Clean(r.URL.Path)
		full_path := filepath.Join(string(root), clean_path)

		fmt.Println("Serving file: ", full_path)

		if _, err := os.Stat(full_path); os.IsNotExist(err) {
			fmt.Println("file not found: ", err)
			http.NotFound(w, r)
			return
		}
		fs.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
