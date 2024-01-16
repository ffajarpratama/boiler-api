package constant

type JwtKey int

const (
	UserIDKey JwtKey = iota
	RoleKey
)

const (
	FormatDate1        = "2006-01-02 15:04:05"
	FormatDate2        = "Monday, 02 January 2006 15:04"
	FormatDate3        = "02 January 2006"
	FormatDateExcelMUF = "02-Jan-2006"
	FormatDate4        = "020106" // DDMMYY
	FormatDate5        = "15:04:05"
	FormatDate6        = "2006-01-02 15:04"
	FormatYYYYMMDD     = "2006-01-02"
	FormatDDMMYYYY     = "02-01-2006"
	FormatOrderDateMuf = "2006/02/01 15:04:05"
	FormatYYYYMM       = "2006-01"
	FormatMMYYYY       = "01-2006"
	FormatYYYY         = "2006"
	FormatDatePDF      = "02-Jan-2006 15:04"
	FormatDate7        = "02/01/2006"
	FormatDateExcel    = "1/02/2006"
	FormatDateExcelHMM = "2/01/2006 15:04"
	FormatNameOfDay    = "Monday"
	FormatNameOfMonth  = "January"
	FormatDDMMYYYYHHMM = "2/1/2006 15:04"
	FormatYYYYMMDDHHMM = "2006-01-02 15:04"
)
