package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/infytvcode/pdf-generator/pdf"

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

/*
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
		fmt.Printf("err: %v\n", err)
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
*/

func create_pdf(w http.ResponseWriter, r *http.Request) {
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

	pdfc := pdf.NewPDFProvider()
	html_b := bytes.NewBuffer([]byte(generatePDFRequest.HTML))
	pdfg, err := pdfc.CreatePDF(*html_b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileName := "infytv.pdf"
	if generatePDFRequest.Filename != "" {
		fileName = generatePDFRequest.Filename
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Write(pdfg)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Infy.TV PDF Service"))
	})
	// r.Post("/generate-pdf", generatePDF)
	r.Post("/create-pdf", create_pdf)
	http.ListenAndServe(":8080", r)
}
