package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/fahimanzamdip/go-invoice-api/utils"
)

type PDFService struct{}

func init() {
	// set path to wkhtmltopdf
	if runtime.GOOS == "windows" {
		wkhtmltopdf.SetPath("./bin/wkhtmltopdf.exe")
	} else {
		wkhtmltopdf.SetPath("./bin/wkhtmltopdf-amd64")
	}
}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (p *PDFService) GenerateInvoicePDF(data interface{}) error {
	var err error

	templ := template.New("invoice.html")
	templ.Funcs(template.FuncMap{
		"Price": utils.FormatPrice,
		"Date":  utils.FormatToDate,
	})
	// use Go's default HTML template generation tools to generate your HTML
	templ, err = templ.ParseFiles("./templates/pdf/invoice.html")
	if err != nil {
		return err
	}

	// apply the parsed HTML template data and keep the result in a Buffer
	var body bytes.Buffer
	if err = templ.Execute(&body, data); err != nil {
		return err
	}

	// initalize a wkhtmltopdf generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return err
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
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeLetter)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	// creates the pdf ad keeps it in the buffer
	err = pdfg.Create()
	if err != nil {
		return err
	}

	// It will make the file name with the invoice reference
	err = pdfg.WriteFile(fmt.Sprintf("./public/pdfs/invoice-%s.pdf", strconv.Itoa(int(time.Now().Unix()))))
	if err != nil {
		log.Println(err.Error())
	}

	return nil
}

func (p *PDFService) GeneratePayReceiptPDF(data interface{}) error {
	var err error

	templ := template.New("payment.html")
	templ.Funcs(template.FuncMap{
		"Price": utils.FormatPrice,
		"Date":  utils.FormatToDate,
	})
	// use Go's default HTML template generation tools to generate your HTML
	templ, err = templ.ParseFiles("./templates/pdf/payment.html")
	if err != nil {
		return err
	}

	// apply the parsed HTML template data and keep the result in a Buffer
	var body bytes.Buffer
	if err = templ.Execute(&body, data); err != nil {
		return err
	}

	// initalize a wkhtmltopdf generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return err
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
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeLetter)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)

	// creates the pdf ad keeps it in the buffer
	err = pdfg.Create()
	if err != nil {
		return err
	}

	// It will make the file name with the invoice reference
	err = pdfg.WriteFile(fmt.Sprintf("./public/pdfs/pay-receipt-%s.pdf", strconv.Itoa(int(time.Now().Unix()))))
	if err != nil {
		log.Println(err.Error())
	}

	return nil
}
