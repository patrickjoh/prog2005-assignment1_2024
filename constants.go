package prog2005_assignment1_2024

const (
	// Ports
	DEFAULT_PORT = "8080"

	// Paths
	DEFAULT_PATH    = "/"
	COUNT_PATH      = "/librarystats/v1/bookcount/"
	READERSHIP_PATH = "/librarystats/v1/readership/"
	STATUS_PATH     = "/librarystats/v1/status/"

	// APIs
	GUTENDEX_API = ""
	L2C_API      = ""
	COUNTRY_API  = ""

	// Status
	GUTENDEX_STATUS = "http://129.241.150.113:8000/books/"
	L2C_STATUS      = "http://129.241.150.113:8080/v3.1/name/norway"
	COUNTRY_STATUS  = "http://129.241.150.113:8080/v3.1/alpha?codes=nor,swe,fin,rus"
	STATUS_VERSION  = "v1"
)
