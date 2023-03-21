package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"

	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type GeneratePDFRequest struct {
	HTML     string `json:"html"`
	Filename string `json:"filename"`
}

type GeneratePDFResponse struct {
	Message string `json:"message"`
}

func generatePDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var generatePDFRequest GeneratePDFRequest
	err = json.Unmarshal(reqBody, &generatePDFRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate PDF
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error generating PDF", http.StatusInternalServerError)
		return
	}

	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/javascript", js.Minify)

	minifiedHTML, err := m.String("text/html", generatePDFRequest.HTML)
	if err != nil {
		http.Error(w, "Error minifying HTML", http.StatusInternalServerError)
		return
	}

	minifiedCSS, err := m.String("text/css", "")
	if err != nil {
		http.Error(w, "Error minifying CSS", http.StatusInternalServerError)
		return
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(minifiedHTML)))
	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(minifiedCSS)))

	err = pdfg.Create()
	if err != nil {
		http.Error(w, "Error generating PDF", http.StatusInternalServerError)
		return
	}

	err = pdfg.WriteFile(generatePDFRequest.Filename)
	if err != nil {
		http.Error(w, "Error saving PDF", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := GeneratePDFResponse{
		Message: fmt.Sprintf("PDF generated and saved as %s", generatePDFRequest.Filename),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	os.Setenv("WKHTMLTOPDF_PATH", os.Getenv("LAMBDA_TASK_ROOT"))
	// http.HandleFunc("/generate-pdf", generatePDF)
	// http.ListenAndServe(":8080", nil)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Post("/generate-pdf", generatePDF)
	http.ListenAndServe(":8080", r)
}
