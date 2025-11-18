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
		fmt.Println("root", root)
		cleanPath := filepath.Clean(r.URL.Path)
		fmt.Println("trying:", cleanPath)
		fullPath := filepath.Join(string(root), cleanPath)
		fmt.Println("string(root) ", string(root))
		fmt.Println("cleanPath ", cleanPath)

		fmt.Println("fullPath ", fullPath)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			fmt.Println("file nto sound", err)
			http.NotFound(w, r)
			return
		}
		fs.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
