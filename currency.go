package go_currency

import (
	"fmt"
	"errors"
	"strings"
	"strconv"
)

// A CurrencyError records a failed currency conversion.
type CurrencyError struct {
	Func 	string // the failing function (ParseUSD, parseEUR, ...)
	Amount  string // the input
	Err  	error  // the reason the conversion failed (ErrRange, ErrSyntax)
}

func (e *CurrencyError) Error() string {
	return "go_utils." + e.Func + ": " + "parsing " + strconv.Quote(e.Amount) + ": " + e.Err.Error()
}

type USD struct {
	Numerator 	int
	Denominator	int
}
func (usd USD) String() string {
	return fmt.Sprintf("%d.%02d", usd.Numerator, usd.Denominator)
}
// Helper to convert USD to string, ex:  s.Equal(item1PriceUSD.ToString(), "250.00")
func (usd USD) ToString() string {
	return fmt.Sprintf("%s", usd)
}
func (usd USD) ToCents() int {
	return usd.Numerator * 100 + usd.Denominator
}
func (usd1 USD) Add(usd2 USD) (amountUSD USD, err error) {
	usd1_cents := usd1.ToCents()
	usd2_cents := usd2.ToCents()
	sum := usd1_cents + usd2_cents
	sum_str := fmt.Sprintf("%v", sum)
	whole := sum_str[0:len(sum_str)-2]
	decimal := sum_str[len(sum_str)-2:]
	return ParseUSD(fmt.Sprintf("%s.%s", whole, decimal))
}
func (usd1 USD) Subtract(usd2 USD) (amountUSD USD, err error) {
	usd1_cents := usd1.ToCents()
	usd2_cents := usd2.ToCents()
	diff := usd1_cents - usd2_cents
	diff_str := fmt.Sprintf("%v", diff)
	whole := diff_str[0:len(diff_str)-2]
	decimal := diff_str[len(diff_str)-2:]
	return ParseUSD(fmt.Sprintf("%s.%s", whole, decimal))
}
func (usd1 USD) Multiply(multiplier int) (amountUSD USD, err error) {
	usd1_cents := usd1.ToCents()
	product := usd1_cents * multiplier
	product_str := fmt.Sprintf("%v", product)
	whole := product_str[0:len(product_str)-2]
	decimal := product_str[len(product_str)-2:]
	return ParseUSD(fmt.Sprintf("%s.%s", whole, decimal))
}

const FnParseUSD = "ParseUSD"

// ErrRange indicates that a value is out of range for the target type.
var ErrRange = errors.New("value out of range")

func CurrencyErrorFn(fn, str string) *CurrencyError {
	return &CurrencyError{fn, str, ErrRange}
}

func ParseUSD(priceStr string) (amountUSD USD, err error) {
	nilUSD := USD{0, 0}
	if priceStr == "0.00" || priceStr == ".00" || priceStr == "0" {
		return nilUSD, err
	}
	var wholeNumStr string
	decimalIdx := strings.Index(priceStr, ".")
	if decimalIdx < 0 && len(priceStr) == 0 || priceStr == "." || decimalIdx == len(priceStr) {
		err = CurrencyErrorFn(FnParseUSD, priceStr)
		// Following is true on error: fmt.Sprintf("%s", usd250) == "0.00"
		return amountUSD, err
	} else if decimalIdx > 0 {
		wholeNumStr = priceStr[0:decimalIdx]
		if err != nil{
			return nilUSD, err
		}
	}
	var decimalStr string
	if decimalIdx < 0 {
		wholeNumStr = priceStr
		decimalStr = "00"
	}  else {
		decimalStr = priceStr[decimalIdx + 1:len(priceStr)]
	}
	wholeNumInt, err := strconv.Atoi(wholeNumStr)
	decimalInt, err := strconv.Atoi(decimalStr)
	if err != nil {
		err = CurrencyErrorFn(FnParseUSD, priceStr)
		return nilUSD, err
	}
	if decimalInt > 100 {
		wholeNumInt += decimalInt / 100
		decimalInt = decimalInt % 100
	}
	amountUsd := new(USD)
	amountUsd.Numerator = wholeNumInt
	amountUsd.Denominator = decimalInt
	return *amountUsd, err
}


