package analyzers

import (
	"fmt"
	"sort"

	"github.com/antonmamonov/IntrinsicValueBot/financehandlers"
)

// var CurrentStockPortfolio []string = []string{"ENB", "T", "VZ", "MFC", "TRP", "BRK-B", "BNS", "ADBE", "FFH", "BLK", "KO", "PG", "MET", "AP-UN", "CAR-UN", "K"}
// var CurrentStockPortfolio []string = []string{"T", "BNS", "BLK", "MET", "K"}

var CurrentStockPortfolio []string = []string{"BRK-B"}

func GetMostUndervaluedStock() (string, error) {

	intrinsicVals := []financehandlers.IntrinsicValue{}

	for _, stockSymbol := range CurrentStockPortfolio {
		finHandler, finHandlerErr := financehandlers.CreateNewFmpFinanceHandler(stockSymbol)

		if finHandlerErr != nil {
			panic(finHandlerErr)
		}

		intrinsicVal, getIntrinsicValueErr := finHandler.GetIntrinsicValue()

		if getIntrinsicValueErr == nil {
			intrinsicVals = append(intrinsicVals, intrinsicVal)
		} else {
			fmt.Println("ERROR", stockSymbol)
			fmt.Println(getIntrinsicValueErr)
		}

	}

	sort.Slice(intrinsicVals, func(i, j int) bool {
		return intrinsicVals[i].MarginOfSafety > intrinsicVals[j].MarginOfSafety
	})

	for _, intrinsicVal := range intrinsicVals {
		fmt.Println(intrinsicVal.StockSymbol)
		fmt.Println("Price:", intrinsicVal.Price)
		fmt.Println("Value:", intrinsicVal.Value)
		fmt.Println("MarginOfSafety:", intrinsicVal.MarginOfSafety)
		fmt.Println("assetToLiabilityRatio:", intrinsicVal.AssetToLiabilityRatio)
		fmt.Println("----------------------")
	}

	return "", nil
}
