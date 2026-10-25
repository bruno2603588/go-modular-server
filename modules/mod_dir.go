package modules

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type DirConfig struct {
	StaticDir      string
	DirectoryIndex string
}

type DirHandler struct {
	Config DirConfig
}

func NewDirHandler(config DirConfig) *DirHandler {
	if config.DirectoryIndex == "" {
		config.DirectoryIndex = "index.html"
	}
	return &DirHandler{Config: config}
}

func (d DirHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cleanPath := filepath.Clean(r.URL.Path)

	fullPath := filepath.Join(d.Config.StaticDir, cleanPath)

	info, err := os.Stat(fullPath)

	if os.IsNotExist(err) {
		indexPath := filepath.Join(d.Config.StaticDir, d.Config.DirectoryIndex)
		http.ServeFile(w, r, indexPath)
		return
	} else if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
			return
		}
		indexPath := filepath.Join(fullPath, d.Config.DirectoryIndex)
		http.ServeFile(w, r, indexPath)

		if _, err := os.Stat(indexPath); err == nil {
			http.ServeFile(w, r, indexPath)
			return
		}
	}
	http.ServeFile(w, r, fullPath)
}
