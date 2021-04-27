package integration_tests

import (
	"testing"

	"github.com/antonmamonov/IntrinsicValueBot/analyzers"
	// "github.com/antonmamonov/IntrinsicValueBot/financehandlers"
)

func TestFMPEnd2End(t *testing.T) {

	// financeHandler, createNewFmpFinanceHandlerErr := financehandlers.CreateNewFmpFinanceHandler("ENB")

	// if createNewFmpFinanceHandlerErr != nil {
	// 	t.Error(createNewFmpFinanceHandlerErr)
	// }

	// intrinsicValue, getIntrinsicValueErr := financeHandler.GetIntrinsicValue()

	// if getIntrinsicValueErr != nil {
	// 	t.Error(getIntrinsicValueErr)
	// }

	// fmt.Println(intrinsicValue)

	analyzers.GetMostUndervaluedStock()

}
