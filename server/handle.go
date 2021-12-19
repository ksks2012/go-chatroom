package server

import (
	"net/http"
	"os"
	"path/filepath"
)

func RegisterHandle() {
	inferRootDir()

	// Processing broadcast message
	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}

// infers the project root directory
func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var infer func(d string) string
	infer = func(d string) string {
		// make sure that the template directory exists in the project root directory
		if exists(d + "/template") {
			return d
		}

		return infer(filepath.Dir(d))
	}

	rootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
