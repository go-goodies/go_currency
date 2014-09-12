package go_currency

import (
	"testing"
	"log"
	"github.com/remogatto/prettytest"
	jc "github.com/go-goodies/go_jsoncfg"
	u "github.com/go-goodies/go_utils"
	"fmt"
)

type mySuite struct {
	prettytest.Suite
}

func TestRunner(t *testing.T) {
	prettytest.Run(t, new(mySuite))
}

func (s *mySuite) TestReadConfigFile() {

	//------------------------------------------------
	//               Read JSON file
	//------------------------------------------------
	/*
	{
	  "lineitem1": {
		  "description": "SSD Drive",
		  "count": 2,
		  "priceUSD": "250.00"
	  },
	  "lineitem2": {
		  "description": "RAM Chip",
		  "count": 4,
		  "priceUSD": "125.50"
	  }
	}
	 */

	obj, err := jc.ReadFile("testdata/invoice.json")
	if err != nil {
		log.Fatal(err)
	}
	lineitem1 := obj.RequiredObject("lineitem1")

	//------------------------------------------------
	//          description == "SSD Drive"
	//------------------------------------------------
	s.Equal(lineitem1["description"], "SSD Drive")

	//------------------------------------------------
	//                count == 2
	//------------------------------------------------
	lineitem1_count, err := u.ConvNumToInt(lineitem1["count"])
	if err != nil {
		log.Fatal(err)
	}
	s.Equal(lineitem1_count, int(2))

	//------------------------------------------------
	//             priceUSD == 250.00
	//------------------------------------------------
	item1PriceUSD, err := ParseUSD(lineitem1["priceUSD"].(string))
	if err != nil {
		log.Fatal(err)
	}
	usd250, err := ParseUSD("250.00")
	if err != nil {
		log.Fatal(err)
	}
	s.Equal(item1PriceUSD, usd250)

	s.Equal(u.TypeOf(usd250), "go_currency.USD")
}

func (s *mySuite) TestErrors() {
	emptyUSD, err := ParseUSD("")
	s.Equal(fmt.Sprintf("%s", err), "go_utils.ParseUSD: parsing \"\": value out of range")
	s.Equal(u.TypeOf(emptyUSD), "go_currency.USD")

	unnecessaryDeciaml, err := ParseUSD("2.")
	s.Equal(fmt.Sprintf("%s", err), "go_utils.ParseUSD: parsing \"2.\": value out of range")
	s.Equal(u.TypeOf(unnecessaryDeciaml), "go_currency.USD")

	usd250, err := ParseUSD(".")
	s.Equal(fmt.Sprintf("%s", err), "go_utils.ParseUSD: parsing \".\": value out of range")
	s.Equal(u.TypeOf(usd250), "go_currency.USD")
	s.Equal(fmt.Sprintf("%s", usd250), "0.00")
}

func initLineItemPrices(jsonFile string) (item1PriceUSD USD, item2PriceUSD USD) {
	obj, err := jc.ReadFile(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	lineitem1 := obj.RequiredObject("lineitem1")
	item1PriceUSD, _ = ParseUSD(lineitem1["priceUSD"].(string))
	lineitem2 := obj.RequiredObject("lineitem2")
	item2PriceUSD, _ = ParseUSD(lineitem2["priceUSD"].(string))
	return item1PriceUSD, item2PriceUSD
}

func (s *mySuite) TestAdd() {
	item1PriceUSD, item2PriceUSD := initLineItemPrices("testdata/invoice.json")
	s.Equal(item1PriceUSD.ToString(), "250.00")
	s.Equal(item2PriceUSD.ToString(), "125.50")
	sum, _ := item1PriceUSD.Add(item2PriceUSD)
	usd375_50, _ := ParseUSD("375.50")
	s.Equal(sum, usd375_50)
}

func (s *mySuite) TestSubtract() {
	item1PriceUSD, item2PriceUSD := initLineItemPrices("testdata/invoice.json")
	s.Equal(item1PriceUSD.ToString(), "250.00")
	s.Equal(item2PriceUSD.ToString(), "125.50")
	diff, _ := item1PriceUSD.Subtract(item2PriceUSD)
	usd124_50, _ := ParseUSD("124.50")
	s.Equal(diff, usd124_50)
}

func (s *mySuite) TestMultiply() {
	s.Pending()
}
