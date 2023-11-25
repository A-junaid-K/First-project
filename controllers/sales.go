package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/tealeg/xlsx"
)

func Sales(c *gin.Context) {
	db := database.DB

	var order []models.Order
	db.Find(&order)

	c.HTML(200, "salesreports.html", order)
}
func Salesreport(c *gin.Context) {

	layout := "2006-01-02"

	raw_starting_date := c.PostForm("startingdate")
	log.Println("raw : ", raw_starting_date)
	start, err := time.Parse(layout, raw_starting_date)
	if err != nil {
		log.Println("faild to get starting date : ", err)
		return
	}
	raw_expiry_date := c.PostForm("endingdate")
	log.Println("rew exp : ", raw_expiry_date)
	end, errr := time.Parse(layout, raw_expiry_date)
	end = end.AddDate(0, 0, 1)
	if errr != nil {
		log.Println("faild to get ending date : ", err)
		return
	}
	log.Println("start  : ", start)
	log.Println("end : ", end)

	//Fetching data from database and inner joins product table knowing product details
	var orders []models.Order
	err = database.DB.Table("orders").Where("date BETWEEN ? AND ? AND status =?", start, end, "delivered").Scan(&orders).Error
	if err != nil {
		log.Println("scanning error : ", err)
		c.HTML(400, "salesreports.html", gin.H{"error": "Scanning error"})
		return
	}
	log.Println("orders  : ", orders)

	f := excelize.NewFile()
	sheet := "sheet1"
	intex := f.NewSheet(sheet)

	f.SetCellValue(sheet, "A1", "User ID")
	f.SetCellValue(sheet, "B1", "Address ID")
	f.SetCellValue(sheet, "C1", "Order ID")
	f.SetCellValue(sheet, "D1", "Date")
	f.SetCellValue(sheet, "E1", "Payment Method")
	f.SetCellValue(sheet, "F1", "Payment ID")
	f.SetCellValue(sheet, "G1", "Grand Total")
	f.SetCellValue(sheet, "H1", "Status")

	for i, v := range orders {
		k := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", k), v.User_ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", k), v.Address_ID)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", k), v.Order_ID)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", k), v.Date.Format(layout))
		f.SetCellValue(sheet, fmt.Sprintf("E%d", k), v.Payment_Type)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", k), v.Payment_ID)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", k), v.Total_Price)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", k), v.Status)
	}
	f.SetActiveSheet(intex)

	if err := f.SaveAs("./public/salesreport.xlsx"); err != nil {
		fmt.Println(err)
		return
	}
	convertintoPdf(c)
	c.HTML(200, "salesreports.html", orders)

}

//---------------------------------------------

func convertintoPdf(c *gin.Context) {
	file, err := xlsx.OpenFile("./public/salesreport.xlsx")

	if err != nil {
		fmt.Println("File open error")
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(10, 10, 20)
	pdf.Ln(-1)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(5, 10, "Sales report")
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 8)
	// Convertig each cell in the Excel file to a PDF cell
	for _, sheet := range file.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				//if there is any empty cell values skiping that
				if cell.Value == "" {
					continue
				}

				pdf.CellFormat(24, 8, cell.Value, "1", 0, "C", false, 0, "")
			}
			pdf.Ln(-1)
		}
	}
	err = pdf.OutputFileAndClose("./public/salesreport.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("PDF saved successfully.")
}
func DownloadExel(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=salesreport.xlsx")
	c.Header("Content-Type", "application/xlsx")
	c.File("./public/salesreport.xlsx")
}

func Downloadpdf(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=salesreport.pdf")
	c.Header("Content-Type", "application/pdf")
	c.File("./public/salesreport.pdf")
}
