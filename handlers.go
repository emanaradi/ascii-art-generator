package asciiartweb

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"image"
	"image/draw"
	"image/png"
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

func ExportPNGHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed, "405 Method Not Allowed")
		return
	}

	art := r.FormValue("ascii-art")

	lines := strings.Split(art, "\n")

	charWidth:=8
	charHeight:=16

	padding:=20

	maxLength :=0

	for _, line := range lines {
		if len(line) > maxLength {
			maxLength = len(line)
		}
	}

	width := maxLength* charWidth +padding*2
	height := len(lines)*charHeight + padding*2

	img := image.NewRGBA(
		image.Rect(0,0,width,height),
	)

	draw.Draw(
		img,
		img.Bounds(),
		image.White,
		image.Point{},
		draw.Src,
	)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(art)))
	w.Header().Set("Content-Disposition", `attachment; filename="ascii-art.png"`)

	err := png.Encode(w, img)
	if err != nil {
		return
	}
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