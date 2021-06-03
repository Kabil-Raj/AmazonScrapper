module sellerapp.assignment/restApiAssignment

go 1.16

require (
	github.com/PuerkitoBio/goquery v1.6.1
	github.com/gocolly/colly/v2 v2.1.0
	github.com/gorilla/mux v1.8.0
	sellerapp.assignment/datahandler v0.0.0-00010101000000-000000000000
)

replace sellerapp.assignment/datahandler => ../DataHandler
