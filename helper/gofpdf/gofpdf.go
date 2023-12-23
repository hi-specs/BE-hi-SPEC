package gofpdf

import (
	"fmt"
	"log"

	"github.com/jung-kurt/gofpdf"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm"
)

type TM struct {
	gorm.Model
	Nota       string
	ProductID  uint
	UserID     uint
	TotalPrice uint
	Status     string
	Token      string
	Url        string
	Pdf        string
}
type PM struct {
	gorm.Model
	Category  string
	Name      string
	CPU       string
	RAM       string
	Display   string
	Storage   string
	Thickness string
	Weight    string
	Bluetooth string
	HDMI      string
	Price     int
	Picture   string
}

type UM struct {
	gorm.Model
	Email       string `gorm:"unique"`
	Name        string `json:"name" form:"name"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Password    string `json:"password" form:"password"`
	Avatar      string `json:"avatar" form:"avatar"`
	Role        string `json:"role" form:"role"`
}

func GeneratePDF(tm TM, pm PM, um UM) {
	p := message.NewPrinter(language.English)
	// Create PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// header
	pdf.SetFont("Arial", "", 30)
	pdf.Cell(180, 10, fmt.Sprintln("hi-SPEC"))
	pdf.SetFont("Arial", "", 11)
	pdf.Image("helper/gofpdf/hispec.png", 10, 24, 0, 0, false, "", 0, "")
	pdf.SetX(-20)
	pdf.Cell(0, 10, fmt.Sprint("INVOICE"))
	pdf.Ln(3)
	pdf.SetX(-20)
	pdf.SetFont("Arial", "", 8)
	nota := tm.Nota
	pdf.Cell(100, 12, fmt.Sprint(nota))
	pdf.SetFont("Arial", "", 11)
	pdf.Ln(5)
	pdf.Cell(20, 10, fmt.Sprint("-----------------------------------------------------------------------------------------------------------------------------------------------"))

	// user
	pdf.Ln(5)
	pdf.Cell(20, 10, fmt.Sprint("USER"))
	pdf.Ln(5)
	user := um.Name
	pdf.Cell(20, 10, fmt.Sprintf("Name                : %s", user))
	pdf.Ln(5)
	add := um.Address
	pdf.Cell(20, 10, fmt.Sprintf("Address            : %s", add))
	pdf.Ln(5)
	phone := um.PhoneNumber
	pdf.Cell(20, 10, fmt.Sprintf("Phone               : %s", phone))
	pdf.Ln(5)
	pdf.Cell(20, 10, fmt.Sprint("-----------------------------------------------------------------------------------------------------------------------------------------------"))

	// product
	pdf.Ln(5)
	pdf.Cell(20, 10, fmt.Sprint("INFO PRODUCT"))
	pdf.Ln(5)
	prod := pm.Name
	pdf.Cell(20, 10, fmt.Sprintf("Name                : %s", prod))
	pdf.Ln(5)
	price := pm.Price
	pdf.Cell(20, 10, p.Sprintf("Price                 : Rp%d", price))
	pdf.Ln(5)
	pdf.Cell(20, 10, fmt.Sprint("-----------------------------------------------------------------------------------------------------------------------------------------------"))

	// detail transaction
	pdf.Ln(5)
	pdf.Cell(20, 10, fmt.Sprint("DETAIL TRANSACTION"))
	pdf.Ln(5)
	trans := tm.ID
	pdf.Cell(20, 10, fmt.Sprintf("Transaction ID  : %d", trans))
	pdf.Ln(5)
	status := tm.Status
	pdf.Cell(20, 10, fmt.Sprintf("Status               : %s", status))
	pdf.Ln(5)
	total := tm.TotalPrice
	pdf.Cell(20, 10, p.Sprintf("Total Price        : Rp%d", total))
	pdf.SetX(-83)
	DDMMYYYYhhmmss := "2006-01-02 15:04:05"
	time := tm.UpdatedAt
	pdf.Cell(20, 10, fmt.Sprintf("Last Updated: %s (GMT+11)", time.Format(DDMMYYYYhhmmss)))

	//Save PDF to a temporary file
	pdfPath := "helper/gofpdf/invoice.pdf"
	err2 := pdf.OutputFileAndClose(pdfPath)
	if err2 != nil {
		log.Fatal(err2)
	}

	return
}
