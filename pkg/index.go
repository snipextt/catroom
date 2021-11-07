package pkg

import "net/http"

func HandleIndex(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "./index.html")
}
