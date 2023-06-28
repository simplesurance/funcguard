package testdata

import (
	"io"
	"net/http"
)

var httpClt http.Client

func httpPost() {
	http.Post("", "", io.Reader(nil))    // want "use context-aware http.NewRequestWithContext method instead"
	httpClt.Post("", "", io.Reader(nil)) // want "use context-aware http.NewRequestWithContext method instead"
}

func httpPostForm() {
	http.PostForm("", nil)    // want "use context-aware http.NewRequestWithContext method instead"
	httpClt.PostForm("", nil) // want "use context-aware http.NewRequestWithContext method instead"
}

func httpGet() {
	http.Get("")    // want "use context-aware http.NewRequestWithContext method instead"
	httpClt.Get("") // want "use context-aware http.NewRequestWithContext method instead"
}

func httpHead() {
	http.Head("")    // want "use context-aware http.NewRequestWithContext method instead"
	httpClt.Head("") // want "use context-aware http.NewRequestWithContext method instead"
}
