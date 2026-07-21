package handler

import (
	"net/http"
	"path/filepath"
)

// it handles and gives the frontend files in static folder
func NewStaticHandler(staticDir string) http.Handler {
	absPath, err := filepath.Abs(staticDir)

	if err != nil {
		panic("Invalid static directory: " + err.Error())
	}
	fs := http.FileServer(http.Dir(absPath))
	return fs
}
