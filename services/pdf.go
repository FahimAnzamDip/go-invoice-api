package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (p *PDFService) GenerateInvoicePDF(data interface{}) ([]byte, error) {
	var templ *template.Template
	var err error

	// use Go's default HTML template generation tools to generate your HTML
	if templ, err = template.ParseFiles("./templates/pdf/invoice.html"); err != nil {
		return nil, err
	}

	// apply the parsed HTML template data and keep the result in a Buffer
	var body bytes.Buffer
	if err = templ.Execute(&body, data); err != nil {
		return nil, err
	}

	// set path to wkhtmltopdf
	wkhtmltopdf.SetPath("C:\\wkhtmltopdf\\bin\\wkhtmltopdf.exe")

	// initalize a wkhtmltopdf generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	// read the HTML page as a PDF page
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(body.Bytes()))

	// enable this if the HTML file contains local references such as images, CSS, etc.
	page.EnableLocalFileAccess.Set(true)

	// add the page to your generator
	pdfg.AddPage(page)

	// manipulate page attributes as needed
	pdfg.MarginLeft.Set(20)
	pdfg.MarginRight.Set(20)
	pdfg.MarginTop.Set(20)
	pdfg.MarginBottom.Set(20)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	// creates the pdf ad keeps it in the buffer
	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	// It will make the file name with the invoice reference
	err = pdfg.WriteFile(fmt.Sprintf("./public/pdfs/invoice-%s.pdf", strconv.Itoa(int(time.Now().Unix()))))
	if err != nil {
		log.Println(err.Error())
	}

	return pdfg.Bytes(), nil
}
