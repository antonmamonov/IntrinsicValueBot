package integration_tests

import (
	"fmt"
	"testing"

	"github.com/antonmamonov/IntrinsicValueBot/financehandlers"
)

func TestFinnHubEnd2End(t *testing.T) {
	// factory comes before the product
	finnHubFinanceHandlerFactory, finnHubFinanceHandlerFactoryError := financehandlers.CreateNewFinnHubFinanceHandlerFactory()

	if finnHubFinanceHandlerFactoryError != nil {
		t.Error(finnHubFinanceHandlerFactoryError)
	}

	fmt.Println(finnHubFinanceHandlerFactory)

	financeHandler, createNewFinnHubFinanceHandlerErr := finnHubFinanceHandlerFactory.CreateFinanceHandler("ENB")

	if createNewFinnHubFinanceHandlerErr != nil {
		t.Error(createNewFinnHubFinanceHandlerErr)
	}

	fmt.Println(financeHandler)

	periods, getEarningsPerShareGrowthForPeriodErr := financeHandler.GetEarningsPerShareGrowthForPeriod(5, "year")

	if getEarningsPerShareGrowthForPeriodErr != nil {
		t.Error(getEarningsPerShareGrowthForPeriodErr)
	}

	fmt.Println("periods", periods)

}
