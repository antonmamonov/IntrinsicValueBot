package financehandlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

type FMPFinanceHandler struct {
	apiKey      string
	StockSymbol string
}

func (finHandler FMPFinanceHandler) SetStockSymbol(stockSymbol string) error {
	finHandler.StockSymbol = stockSymbol

	return nil
}

func (finHandler FMPFinanceHandler) GetCurrentStockSymbol() (string, error) {
	stockSymbol := finHandler.StockSymbol

	if stockSymbol == "" {
		return "", errors.New("stockSymbol is not set! Please set it first.")
	}

	return stockSymbol, nil
}

type GetIncomeStatementJSONResponse struct {
	Symbol           string  `json:"symbol"`
	ReportedCurrency string  `json:"reportedCurrency"`
	Date             string  `json:"date"`
	EPSdiluted       float64 `json:"epsdiluted"`
}

// exchange -> regular - "" , canada - ".TO"
//
func (finHandler FMPFinanceHandler) GetIncomeStatementForExchange(periods int, periodType string, exchange string) ([]GetIncomeStatementJSONResponse, error) {

	getIncomeStatementJSONResponse := []GetIncomeStatementJSONResponse{}

	if periodType != "year" {
		return getIncomeStatementJSONResponse, errors.New("Sorry, we don't support that period type yet.")
	}

	url := "https://financialmodelingprep.com/api/v3/income-statement/" + finHandler.StockSymbol + exchange + "?apikey=" + finHandler.apiKey + "&limit=" + strconv.Itoa(periods)
	response, err := http.Get(url)
	// needed for rate-limiting
	time.Sleep(50 * time.Millisecond)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(string(responseData))

	var getIncomeStatementJSONResponses []GetIncomeStatementJSONResponse
	json.Unmarshal(responseData, &getIncomeStatementJSONResponses)

	if len(getIncomeStatementJSONResponses) == 0 {

		//try .TO
		if exchange == "" {
			return finHandler.GetIncomeStatementForExchange(periods, periodType, ".TO")
		}

		return getIncomeStatementJSONResponses, errors.New("couldn't find income statement for stock")
	}

	return getIncomeStatementJSONResponses, nil
}

type GetBalanceStatementJSONResponse struct {
	Symbol           string `json:"symbol"`
	ReportedCurrency string `json:"reportedCurrency"`
	TotalAssets      int64  `json:"totalAssets"`
	TotalLiabilities int64  `json:"totalLiabilities"`
}

// exchange -> regular - "" , canada - ".TO"
//
func (finHandler FMPFinanceHandler) GetBalanceStatementForExchange(exchange string) (GetBalanceStatementJSONResponse, error) {

	url := "https://financialmodelingprep.com/api/v3/balance-sheet-statement/" + finHandler.StockSymbol + exchange + "?apikey=" + finHandler.apiKey + "&limit=1"
	response, err := http.Get(url)
	// needed for rate-limiting
	time.Sleep(50 * time.Millisecond)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var getBalanceStatementJSONResponses []GetBalanceStatementJSONResponse
	json.Unmarshal(responseData, &getBalanceStatementJSONResponses)

	if len(getBalanceStatementJSONResponses) == 0 {

		//try .TO
		if exchange == "" {
			return finHandler.GetBalanceStatementForExchange(".TO")
		}

		return GetBalanceStatementJSONResponse{}, errors.New("couldn't find income statement for stock")
	}

	return getBalanceStatementJSONResponses[0], nil
}

func (finHandler FMPFinanceHandler) GetEarningsPerShareGrowthForPeriod(periods int, periodType string) (EarningsPerShareGrowth, error) {

	earningsPerShareGrowth := EarningsPerShareGrowth{}

	if periodType != "year" {
		return earningsPerShareGrowth, errors.New("Sorry, we don't support that period type yet.")
	}

	getIncomeStatementJSONResponses, getIncomeStatementJSONErr := finHandler.GetIncomeStatementForExchange(periods, periodType, "")

	if getIncomeStatementJSONErr != nil {
		return earningsPerShareGrowth, getIncomeStatementJSONErr
	}

	// first get growth
	latestEps := getIncomeStatementJSONResponses[0].EPSdiluted
	farestEps := getIncomeStatementJSONResponses[len(getIncomeStatementJSONResponses)-1].EPSdiluted

	totalGrowthMultiple := math.Abs(latestEps / farestEps)
	compoundGrowthRate := math.Pow(totalGrowthMultiple, float64(1.0/float64(periods-1))) - float64(1.0)

	earningsPerShareGrowth.LatestEPS = latestEps
	earningsPerShareGrowth.EPSGrowth = compoundGrowthRate

	return earningsPerShareGrowth, nil
}

type StockQuote struct {
	Symbol     string  `json:"symbol"`
	QuotePrice float64 `json:"price"`
}

func (finHandler FMPFinanceHandler) GetCurrentStockQuote() (StockQuote, error) {
	return finHandler.GetCurrentStockQuoteForExchange("")
}

func (finHandler FMPFinanceHandler) GetCurrentStockQuoteForExchange(exchange string) (StockQuote, error) {

	url := "https://financialmodelingprep.com/api/v3/quote/" + finHandler.StockSymbol + exchange + "?apikey=" + finHandler.apiKey
	response, err := http.Get(url)
	// needed for rate-limiting
	time.Sleep(50 * time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(responseData)

	var getStockQuoteJSONResponses []StockQuote

	json.Unmarshal(responseData, &getStockQuoteJSONResponses)
	// fmt.Println(getStockQuoteJSONResponses)

	if len(getStockQuoteJSONResponses) == 0 {

		//try .TO
		if exchange == "" {
			return finHandler.GetCurrentStockQuoteForExchange(".TO")
		}

		return StockQuote{}, errors.New("couldn't find stock quote")
	}

	return getStockQuoteJSONResponses[0], nil
}

func (finHandler FMPFinanceHandler) GetIntrinsicValue() (IntrinsicValue, error) {

	intrinsicValue := IntrinsicValue{}
	earningGrowth, getEarningGrowthErr := finHandler.GetEarningsPerShareGrowthForPeriod(10, "year")

	if getEarningGrowthErr != nil {
		return intrinsicValue, getEarningGrowthErr
	}

	getBalanceStatementJSONResponse, getBalanceStatementJSONErr := finHandler.GetBalanceStatementForExchange("")
	if getBalanceStatementJSONErr != nil {
		return intrinsicValue, getBalanceStatementJSONErr
	}

	assetToLiabilityRatio := float64(getBalanceStatementJSONResponse.TotalAssets) / float64(getBalanceStatementJSONResponse.TotalLiabilities)

	discountRate := 1.2
	valueOfShare := 0.0
	previousYearEPS := earningGrowth.LatestEPS
	yearsAhead := 1
	earningGrowthRate := earningGrowth.EPSGrowth

	for previousYearEPS > 0 {
		epsForYear := (previousYearEPS * (1.0 + float64(earningGrowthRate)))
		discountRateForYear := math.Pow(discountRate, float64(yearsAhead))

		thisYearAheadEPS := epsForYear / discountRateForYear
		// fmt.Println("stock:", finHandler.StockSymbol, "year:", yearsAhead, "earningsGrowthRate:", earningGrowthRate, "epsForYear:", epsForYear, "discountRateForYear", discountRateForYear, "yearEPS:", thisYearAheadEPS)
		yearsAhead += 1
		valueOfShare += thisYearAheadEPS
		previousYearEPS = epsForYear
		if yearsAhead > 100 {
			previousYearEPS = 0
		}
		earningGrowthRate = earningGrowthRate * 0.995
	}

	stockQuote, getStockQuoteErr := finHandler.GetCurrentStockQuoteForExchange("")

	if getStockQuoteErr != nil {
		return intrinsicValue, getStockQuoteErr
	}

	intrinsicValue.StockSymbol = finHandler.StockSymbol
	intrinsicValue.Value = valueOfShare
	intrinsicValue.MarginOfSafety = 1.0 - (stockQuote.QuotePrice / valueOfShare)
	intrinsicValue.Price = stockQuote.QuotePrice
	intrinsicValue.AssetToLiabilityRatio = assetToLiabilityRatio

	return intrinsicValue, nil
}

type FMP_DCF_JSONResponse struct {
	Symbol     string  `json:"symbol"`
	DCF        float64 `json:"dcf"`
	StockPrice float64 `json:"Stock Price"`
}

func (finHandler FMPFinanceHandler) GetIntrinsicValue_DiscountedCashFlow() (IntrinsicValue, error) {

	url := "https://financialmodelingprep.com/api/v3/discounted-cash-flow/" + finHandler.StockSymbol + "?apikey=" + finHandler.apiKey
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var getIntrinsicValueJSONResponses []FMP_DCF_JSONResponse

	json.Unmarshal(responseData, &getIntrinsicValueJSONResponses)

	intrinsicValue := IntrinsicValue{}

	for _, getIntrinsicValueJSONResponse := range getIntrinsicValueJSONResponses {
		if getIntrinsicValueJSONResponse.Symbol == finHandler.StockSymbol {
			intrinsicValue.Value = getIntrinsicValueJSONResponse.DCF
			intrinsicValue.MarginOfSafety = 1 - (getIntrinsicValueJSONResponse.StockPrice / getIntrinsicValueJSONResponse.DCF)
			intrinsicValue.StockSymbol = getIntrinsicValueJSONResponse.Symbol
		}
	}

	return intrinsicValue, nil
}

func CreateNewFmpFinanceHandler(stockSymbol string) (FinanceHandler, error) {
	financeHandler := FMPFinanceHandler{"221dc7ab2611611859ee89b055ea7c88", stockSymbol}

	return financeHandler, nil
}
