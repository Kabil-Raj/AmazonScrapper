package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	datamodel "github.com/Kabil-Raj/datamodel"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/home", homePage)
	myRouter.HandleFunc("/scrapproduct", scrapAmazonProduct).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

type ProductDetail struct {
	ProductName        string    `json:"ProductName"`
	ProductImageUrl    string    `json:"ProductImageUrl"`
	ProductDescription string    `json:"ProductDescription"`
	ProductPrice       string    `json:"ProductPrice"`
	ProductReviews     string    `json:"ProductReviews"`
	CreatedAt          time.Time `json:"CreatedAt"`
}

var ProductDetails []ProductDetail

func scrapAmazonProduct(w http.ResponseWriter, req *http.Request) {

	productUrl := req.URL.Query().Get("url")
	getProductDetails(productUrl)
	fmt.Println("Endpoint Hit : return product details")
	json.NewEncoder(w).Encode(ProductDetails)
}

func main() {

	datamodel.ConnectMySql()
	handleRequests()
}

func getProductDetails(productUrl string) {

	var productName string

	var productImageUrl string

	var productPrice string

	var productReviews string

	var productDescription string

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("#productDescription", func(e *colly.HTMLElement) {
		productDescription = strings.TrimSpace(e.DOM.Children().Text())
		fmt.Println(len(e.DOM.Children().Text()))
	})

	c.OnHTML("#acrCustomerReviewText", func(e *colly.HTMLElement) {
		productReviews = e.Text
	})

	c.OnHTML("#desktop_unifiedPrice", func(e *colly.HTMLElement) {
		// normal amazon price
		e.DOM.Find("#priceblock_ourprice").Each(func(i int, s *goquery.Selection) {
			productPrice = s.Text()
		})
		// deal price
		e.DOM.Find("#priceblock_dealprice").Each(func(i int, s *goquery.Selection) {
			productPrice = s.Text()
		})

	})

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		var getProductName string
		getProductName = getProductName + e.DOM.Find("span").Text()
		for _, character := range getProductName {
			if character != 10 {
				productName = productName + string(character)
			}
		}

	})

	c.OnHTML("#imgTagWrapperId", func(e *colly.HTMLElement) {
		productImageUrl = getProductImage(e.DOM.Children().Attr("src"))
	})

	c.Visit(productUrl)

	fmt.Println(productName)
	ProductDetails = []ProductDetail{
		{ProductName: productName, ProductImageUrl: productImageUrl, ProductDescription: productDescription, ProductPrice: productPrice, ProductReviews: productReviews, CreatedAt: time.Now()},
	}

	datamodel.SaveData(productName, productImageUrl, productDescription, productPrice, productReviews, time.Now())
}

func getProductImage(imgSource string, isProductImage bool) (soruce string) {
	if isProductImage {
		return imgSource
	}
	return "nothing"
}
