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

// func AddCoupons(c *gin.Context) {
// 	type data struct {
// 		CouponCode    string
// 		Year          uint
// 		Month         uint
// 		Day           uint
// 		DiscountPrice float64
// 		Expired       time.Time
// 	}

// 	var userEnterData data
// 	var couponData []models.Coupon
// 	db := database.DB

// 	if c.Bind(&userEnterData) != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"Error": "Could not bind the JSON Data",
// 		})
// 		return
// 	}
// 	specificTime := time.Date(int(userEnterData.Year), time.Month(userEnterData.Month), int(userEnterData.Day), 0, 0, 0, 0, time.UTC)
// 	userEnterData.Expired = specificTime

//		r := db.First(&couponData, "coupon_code = ?", userEnterData.CouponCode)
//		if r.Error != nil {
//			Data := models.Coupon{
//				CouponCode:    userEnterData.CouponCode,
//				DiscountPrice: userEnterData.DiscountPrice,
//				Expired:       specificTime,
//			}
//			r := db.Create(&Data)
//			if r.Error != nil {
//				c.JSON(http.StatusBadRequest, gin.H{
//					"Error": r.Error.Error(),
//				})
//				return
//			}
//			c.JSON(http.StatusOK, gin.H{
//				"Message": userEnterData,
//				"Success": "Coupon added successfully",
//			})
//		} else {
//			c.JSON(http.StatusBadRequest, gin.H{
//				"Message": "Coupon already exist",
//			})
//		}
//	}
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

// Invoice

// Template for creating pdf
// const invoiceTemplate = `
// <style>
// 	body{
// 		background-color: white;
// 	}
// 	table{
// 		border: 1px solid black;
// 		border-collapse: collapse;
// 	}
// 	  th{
// 		border: 1px solid black;
// 		border-collapse: collapse;
// 		padding-right: 15px;
// 		padding-left: 15px;
// 	 }
// 	 td{
// 		border: 1px solid black;
// 		border-collapse: collapse;
// 		padding-right: 15px;
// 		padding-left: 15px;
// 	 }
// 	 hr{
// 		color:solid black;
// 	 }
// </style>

// <b> TAX INVOICE- </b>
// Order ID : {{.OrderId}}<br>
// Order Date: {{.Date}} <br><hr>
// Name : {{.Name}} <br>
// Email : {{.Email}} <br>
// Billing Address<br>
// {{range .Address}}

// Phone number : {{.Phoneno}} <br>
// House name : {{.Housename}} <br>
// Area :{{.Area}} <br>
// Landmark : {{.Landmark}} <br>
// City : {{.City}} <br>
// Pincode : {{.Pincode}} <br>
// District : {{.District}} <br>
// State : {{.State}} <br>
// {{end}}
// <hr>
// Payment method : {{.PaymentMethod}} <br>
// <hr>
//  <br>

// <table>
// 	<tr>
// 		<th>Product</th>
// 		<th>Description</th>
// 		<th>Qty</th>
// 		<th>Price</th>
// 		<th>Discount</th>
// 		<th>Total Price </th>
// 	</tr>
// 	{{range .Items}}
// 	<tr>
// 		<td>{{.Product}}</td>
// 		<td>{{.Description}}</td>
// 		<td>{{.Qty}}</td>
// 		<td>{{.Price}}</td>
// 		<td>{{.Discount}}</td>
// 		<td>{{.Total}}</td>
// 	</tr>
// 	{{end}}
// </table>

// <br><hr>
// Total Amount : {{.TotalAmount}} <br><hr>`

// type Invoice struct {
// 	Name          string
// 	Email         string
// 	PaymentMethod string
// 	TotalAmount   int64
// 	Date          string
// 	OrderId       uint
// 	Address       []Address
// 	Items         []Item
// }
// type Address struct {
// 	Phoneno   uint
// 	Housename string
// 	Area      string
// 	Landmark  string
// 	City      string
// 	Pincode   uint
// 	District  string
// 	State     string
// }
// type Item struct {
// 	Product     string
// 	Description string
// 	Qty         uint
// 	Price       uint
// 	Discount    uint
// 	Total       uint
// }

// @Summary Generate Purchase Invoice
// @Description Generate a purchase invoice for a user's order
// @Tags Invoice, Orders, Users
// @Accept json
// @Security BearerToken
// @Produce json
// @Success 200 {string} html "HTML response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /user/invoice [get]

// func PurchaseInvoice(c *gin.Context) {
// 	email := c.GetString("user")
// 	db := database.DB
// 	id := getId(email, db)

// 	var user models.User
// 	var Payment models.Payment
// 	var orderData models.OrderDetails
// 	var address models.Address
// 	var orderItem models.OrderItem

// 	// Fetching the data from table OrderItems useing the id
// 	r := db.Last(&orderItem).Where("user_id_no = ?", id)
// 	if r.Error != nil {
// 		err := map[string]interface{}{
// 			"Error": r.Error.Error(),
// 		}
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}

// 	// Fetching data from order_details using userid and order id for fetching the order_item id
// 	r = db.Last(&orderData).Where("userid = ? AND order_item_id = ?", id, orderItem.ID)
// 	if r.Error != nil {
// 		err := map[string]interface{}{
// 			"Error": r.Error.Error(),
// 		}
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}

// 	// Fetching user data
// 	r = db.First(&user, id)
// 	if r.Error != nil {
// 		err := map[string]interface{}{
// 			"Error": r.Error.Error(),
// 		}
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}

// 	// Fetching user address
// 	r = db.First(&address, orderData.Addressid)
// 	if r.Error != nil {
// 		err := map[string]interface{}{
// 			"Error": r.Error.Error(),
// 		}
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}

// 	// Fetching payment details
// 	r = db.Last(&Payment, "userid = ?", id)
// 	if r.Error != nil {
// 		err := map[string]interface{}{
// 			"Error": r.Error.Error(),
// 		}
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}

// 	// Fetching the product data from products
// 	var products []models.Product //models.product
// 	err := db.Joins("JOIN order_details ON products.product_id = order_details.productid").Where("order_details.order_item_id = ?", orderData.OrderItemId).
// 		Find(&products).Error
// 	if err != nil {
// 		err := map[string]interface{}{
// 			"Error": "Somthing went wrong",
// 		}
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}

// 	// To list product details from products
// 	items := make([]Item, len(products))
// 	var qty, coupon []uint

// 	r = db.Table("order_details").Select("quantity").Where("userid = ? AND order_item_id = ?", id, orderItem.ID).Find(&qty)
// 	if r.Error != nil {
// 		err := map[string]interface{}{
// 			"Error": r.Error.Error(),
// 		}
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}
// 	r = db.Table("order_details").Select("couponid").Where("userid = ? AND order_item_id = ?", id, orderItem.ID).Find(&coupon)
// 	if r.Error != nil {
// 		err := map[string]interface{}{
// 			"Error": r.Error.Error(),
// 		}
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}
// 	for idx, c := range coupon {
// 		var couponDis uint
// 		if c > 0 {
// 			_ = db.Table("coupons").Select("discount_price").Where("id = ?", c).Scan(&couponDis)
// 		}
// 		coupon[idx] = couponDis
// 	}
// 	fmt.Println("Coupon:", coupon)
// 	for i, data := range products {

// 		items[i] = Item{
// 			Product:     data.ProductName,
// 			Price:       data.SpecialPrice,
// 			Description: data.Description,
// 			Qty:         qty[i],
// 			Discount:    coupon[i],
// 			Total:       qty[i]*data.SpecialPrice - coupon[i],
// 		}
// 	}

// 	timeString := Payment.Date.Format("02-01-2006")

// 	// Excuting the template invoice
// 	invoice := Invoice{
// 		Name:          user.First_Name,
// 		Date:          timeString,
// 		Email:         user.Email,
// 		OrderId:       orderItem.ID,
// 		PaymentMethod: Payment.PaymentMethod,
// 		TotalAmount:   int64(Payment.TotalAmount),
// 		Address: []Address{
// 			{
// 				Phoneno:   address.Phone,
// 				Housename: address.HouseName,
// 				Area:      address.Area,
// 				Landmark:  address.Landmark,
// 				City:      address.City,
// 				Pincode:   address.Pincode,
// 				District:  address.District,
// 				State:     address.State,
// 			},
// 		},
// 		Items: items,
// 	}

// 	tmplet, err := template.New("invoice").Parse(invoiceTemplate)
// 	if err != nil {
// 		err := map[string]interface{}{
// 			"Error": err.Error(),
// 		}
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}

// 	var buf bytes.Buffer
// 	err = tmplet.Execute(&buf, invoice)
// 	if err != nil {
// 		err := map[string]interface{}{
// 			"Error": err.Error(),
// 		}
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}
// 	cmd := exec.Command("wkhtmltopdf", "-", "./public/invoice.pdf")
// 	cmd.Stdin = &buf
// 	err = cmd.Run()
// 	if err != nil {
// 		err := map[string]interface{}{
// 			"Error": err.Error(),
// 		}
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}
// 	c.HTML(http.StatusOK, "invoice.html", gin.H{})
// }

// // @Summary Download Purchase Invoice
// // @Description Download the generated purchase invoice as a PDF
// // @Tags Invoice, Users
// // @Produce application/pdf
// // @Success 200 "PDF file"
// // @Router /user/invoice/download [get]
// func DownloadInvoice(c *gin.Context) {
// 	c.Header("content-Disposition", "attachment; filename=invoice.pdf")
// 	c.Header("Content-Type", "application/pdf")
// 	c.File("./public/invoice.pdf")
// }
