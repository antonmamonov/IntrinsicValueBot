package financehandlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/antihax/optional"
)

type FinnHubFinanceHandler struct {
	finnhubClient      *finnhub.DefaultApiService
	finnhubAuthContext *context.Context
	StockSymbol        string
}

func (finnHubFinanceHandler FinnHubFinanceHandler) SetStockSymbol(stockSymbol string) error {
	finnHubFinanceHandler.StockSymbol = stockSymbol

	return nil
}

func (finnHubFinanceHandler FinnHubFinanceHandler) GetCurrentStockSymbol() (string, error) {
	stockSymbol := finnHubFinanceHandler.StockSymbol

	if stockSymbol == "" {
		return "", errors.New("stockSymbol is not set! Please set it first.")
	}

	return finnHubFinanceHandler.StockSymbol, nil
}

type FinnHubFinancialsReportedReportCF struct {
	Concept string  `json:"concept"`
	Label   string  `json:"label"`
	Unit    string  `json:"Unit"`
	Value   float64 `json:"value"`
}

type FinnHubFinancialsReportedReport struct {
	Cf []FinnHubFinancialsReportedReportCF `json:"cf"`
}

type FinnHubFinancialsReported struct {
	Form   string                          `json:"form"`
	Year   float64                         `json:"year"`
	Report FinnHubFinancialsReportedReport `json:"report"`
}

func (finnHubFinanceHandler FinnHubFinanceHandler) GetEarningsPerShareGrowthForPeriod(periods int, periodType string) (EarningsPerShareGrowth, error) {

	earningsPerShareGrowth := EarningsPerShareGrowth{}

	if periodType == "year" {
		financialsReported, _, err := finnHubFinanceHandler.finnhubClient.FinancialsReported(*finnHubFinanceHandler.finnhubAuthContext, &finnhub.FinancialsReportedOpts{Symbol: optional.NewString("ENB")})
		// fmt.Println(financialsReported.Data)

		// b, err := json.Marshal(financialsReported.Data)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(string(b))

		// fmt.Println(err)

		// var finnHubFinancialsReports []FinnHubFinancialsReported
		// err := json.Unmarshal(financialsReported.Data, &finnHubFinancialsReports)

		for _, financialsReportedDataElement := range financialsReported.Data {
			var finnHubFinancialsReported FinnHubFinancialsReported

			financialsReportedDataElementBytes, _ := json.Marshal(financialsReportedDataElement)
			json.Unmarshal(financialsReportedDataElementBytes, &finnHubFinancialsReported)
			fmt.Println(finnHubFinancialsReported.Report.Cf[0])

			for _, cf := range finnHubFinancialsReported.Report.Cf {
				fmt.Println(`'` + cf.Label + `'`)
			}

		}

		return earningsPerShareGrowth, nil
	}

	return earningsPerShareGrowth, errors.New(periodType + " periodType is not a valid value!")
}

func (finnHubFinanceHandler FinnHubFinanceHandler) GetIntrinsicValue() (IntrinsicValue, error) {
	return IntrinsicValue{}, nil
}

func CreateNewFinnHubFinanceHandler() (FinanceHandler, error) {
	finnhubClient := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: "bu5fl7748v6qku33onjg",
	})

	financeHandler := FinnHubFinanceHandler{finnhubClient, &auth, ""}

	return financeHandler, nil
}
