package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/tealeg/xlsx"
)

func Getsalesreport(c *gin.Context) {
	c.HTML(http.StatusOK, "salesreports.html", nil)
}
func SalesReport(c *gin.Context) {
	// generate := c.Query("generate")

	generate := c.PostForm("salesreport")
	generate = strings.ToLower(generate)

	//for specific to - from date
	if generate == "specific" {
		startDate := c.Query("from")
		endDate := c.Query("to")
		from, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Invalid start date",
			})
			return
		}

		to, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Invalid start date",
			})
			return
		}
		now := time.Now()
		from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, now.Location())
		to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 999999999, now.Location())
		fmt.Println("to:", to)
		fmt.Println("from: ", from)
		if generateSalesReport(from, to) != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			return
		}
	} else if generate == "monthly" {
		year, err := strconv.Atoi(c.Query("year"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			return
		}
		month_int, errr := strconv.Atoi(c.Query("month"))
		if errr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			return
		}
		month := time.Month(month_int)
		from := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		to := from.AddDate(0, 1, 0).Add(-time.Nanosecond)
		err = generateSalesReport(from, to)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": err.Error(),
			})
			return
		}
	} else if generate == "yearly" {
		year, err := strconv.Atoi(c.Query("year"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			return
		}
		from := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		to := from.AddDate(1, 0, 0).Add(-time.Nanosecond)
		err = generateSalesReport(from, to)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": err.Error(),
			})
			return
		}
	} else if generate == "daily" {
		day, err := strconv.Atoi(c.Query("day"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			return
		}
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), day, 0, 0, 0, 0, now.Location())
		to := time.Date(today.Year(), today.Month(), day, 23, 59, 59, 999999999, today.Location())
		fmt.Println("today: ", today)
		fmt.Println("to: ", to)
		err = generateSalesReport(today, to)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Sales report created successfully",
	})
	c.HTML(http.StatusOK, "SalesReport.html", nil)
}

func generateSalesReport(from time.Time, to time.Time) error {
	// Fetching the data from the order details table as per the date
	var orderDetail []models.OrderItem
	db := database.DB

	r := db.Preload("Product").Preload("payment").Where("created_at BETWEEN ? AND ?", from, to).Find(&orderDetail)
	if r.Error != nil {
		return r.Error
	}

	file := excelize.NewFile()

	// Create a new sheet
	SheetName := "Sheet1"
	index := file.NewSheet(SheetName)

	// Set the value of headers
	file.SetCellValue(SheetName, "A1", "Order Date")
	file.SetCellValue(SheetName, "B1", "Order ID")
	file.SetCellValue(SheetName, "C1", "Product Name")
	file.SetCellValue(SheetName, "D1", "Qty")
	file.SetCellValue(SheetName, "E1", "Price")
	file.SetCellValue(SheetName, "F1", "Total")
	file.SetCellValue(SheetName, "G1", "Pay Method")
	file.SetCellValue(SheetName, "H1", "Status")
	// Set the value of the cell
	for i, report := range orderDetail {
		var product string
		db.Table("products").Where("id=?", report.Product_ID).Scan(&product)
		var paymentid uint
		db.Table("orders").Select("payment_id").Where("order_id=?", report.Order_ID).Scan(&paymentid)
		var paymentType string
		db.Table("payments").Select("payment_type").Where("payment_id=?", paymentid).Scan(&paymentType)
		row := i + 2
		file.SetCellValue(SheetName, fmt.Sprintf("A%d", row), report.Created_at.Format("02/01/2006"))
		file.SetCellValue(SheetName, fmt.Sprintf("B%d", row), report.Order_ID)
		file.SetCellValue(SheetName, fmt.Sprintf("C%d", row), product)
		file.SetCellValue(SheetName, fmt.Sprintf("D%d", row), report.Quantity)
		file.SetCellValue(SheetName, fmt.Sprintf("E%d", row), report.Price)
		file.SetCellValue(SheetName, fmt.Sprintf("F%d", row), report.Total_Price)
		file.SetCellValue(SheetName, fmt.Sprintf("G%d", row), paymentType)
		file.SetCellValue(SheetName, fmt.Sprintf("H%d", row), report.Status)
	}

	// Set active sheet of the workbook
	file.SetActiveSheet(index)

	// Save the Excel file with the name "test.xlsx"
	if err := file.SaveAs("./public/SalesReport.xlsx"); err != nil {
		fmt.Println("Excel file save ERROR:", err)
	}
	// Convert excel to pdf
	ConvertExcelToPdf()

	return nil
}

func ConvertExcelToPdf() {
	xlFile, err := xlsx.OpenFile("./public/SalesReport.xlsx")
	if err != nil {
		fmt.Println("xlsx file open error:", err)
		//return
	}

	// Create a new pdf document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 10)
	// err := pdf.OutputFileAndClose("hello.pdf")

	// Converting each cell in the excel file to a pdf cell
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				if cell.Value == "" { // If there is any empty cell values then skiping that
					continue
				}
				pdf.Cell(25, 10, cell.Value)
			}
			pdf.Ln(-1)
		}
	}

	// Save the PDF document
	err = pdf.OutputFileAndClose("./public/SalesReport.pdf")
	if err != nil {
		fmt.Println("PDF saving error:", err)
	}
	fmt.Println("PDF saved successfully")
}

func DownloadExcel(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=SalesReport.xlsx")
	c.Header("Content-Type", "application/xlsx")
	c.File("./public/SalesReport.xlsx")
}

func DownloadPdf(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=SalesReport.pdf")
	c.Header("Content-Type", "application/pdf")
	c.File("./public/SalesReport.pdf")
}
