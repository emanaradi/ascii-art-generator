package asciiartweb

import (
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

	template, err := template.ParseFiles("templates/index.html")
	if err != nil {
		ErrorHandler(w, http.StatusNotFound, "404 Not Found")
		return
	}

	err = template.Execute(w, nil)

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

	if err != nil {
		banner, _ = GetBanner("standard")
	}

	if len(banner) < 95 {
		ErrorHandler(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	if text == "" || !isValidText(text) || len(text) > 100 {
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

func ErrorHandler(w http.ResponseWriter, statusCode int, message string) {
	template, _ := template.ParseFiles("templates/error.html")

	w.WriteHeader(statusCode)

	err := template.Execute(w, message)

	if err != nil {
		return
	}
}

func ExportTXTHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Disposition", `attachment; filename="ascii-art.txt"`)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
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

func isValidText(text string) bool {
    for _, char := range text {
        if char == '\n' || char == '\r' {
            continue
        }

        if char < ' ' || char > '~' {
            return false
        }
    }

    return true
}