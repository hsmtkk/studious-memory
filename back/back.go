package back

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HelloWorld", helloWorld)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}