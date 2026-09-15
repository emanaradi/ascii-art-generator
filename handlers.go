package asciiartweb

import (
	"html"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

type PageData struct {
	Result string
}

/**
 * Handles incoming HTTP requests home
 */
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusNotFound, "404 Not Found")
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, "405 Method Not Allowed")
		return
	}

	templ, err := template.ParseFiles("templates/index.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	err = templ.Execute(w, nil)

	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
}

func ASCIIArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art" {
		ErrorHandler(w, http.StatusNotFound, "404 Not Found")
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed, "405 Method Not Allowed")
		return
	}

	text := r.FormValue("text")
	banner, err := GetBanner(r.FormValue("banner"))

	if len(banner) < 95 || err != nil {
		ErrorHandler(w, http.StatusNotFound, "404 Not Found")
		return
	}

	if text == "" || !isValidText(text) {
		ErrorHandler(w, http.StatusBadRequest, "400 Bad Request")
		return
	}

	result := GenerateASCIIArt(text, banner)

	data := struct {
		Result string
	}{
		Result: result,
	}

	template, err := template.ParseFiles("templates/index.html")
	if err != nil {
		ErrorHandler(w, http.StatusNotFound, "404 Not Found")
		return
	}

	err = template.Execute(w, data)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
}

func ExportTXTHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed, "405 Method Not Allowed")
		return
	}

	art := r.FormValue("ascii-art")

	if art == "" {
		ErrorHandler(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(art)))
	w.Header().Set("Content-Disposition", `attachment; filename="ascii-art.txt"`)

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(art))
	if err != nil {
		return
	}
}

func ExportHTMLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed, "405 Method Not Allowed")
		return
	}

	art := r.FormValue("ascii-art")

	if art == "" {
		ErrorHandler(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	htmlContent := `<!DOCTYPE html>
	<html lang="eng">
	<head>
		<meta charset="UTF-8">
		<title>ASCII Art</title>
		<style>
		body {
			background-color: white;
			margin: 20px;
			}

		pre {
			font-family: monospace;
			font-size: 16px;
			}
		</style>
	</head>
	<body>
		<pre> ` + html.EscapeString(art) + `</pre>
	</body>
	</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Disposition", `attachment; filename="ascii-art.html"`)

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(htmlContent))
	if err != nil {
		return
	}
}

func ExportPNGHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed, "405 Method Not Allowed")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	if len(data) == 0 {
		ErrorHandler(w, http.StatusBadRequest, "400 Bad Request")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Disposition", `attachment; filename="ascii-art.png"`)

	w.Write(data)
}

// =========================== HELPERS ===========================

func ErrorHandler(w http.ResponseWriter, statusCode int, message string) {
	template, _ := template.ParseFiles("templates/error.html")

	w.WriteHeader(statusCode)

	err := template.Execute(w, message)

	if err != nil {
		return
	}
}

func isValidText(text string) bool {
	for _, char := range text {
		if char > ' ' && char < '~' {
			return true
		}
	}
	return false
}