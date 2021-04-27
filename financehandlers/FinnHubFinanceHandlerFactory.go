package financehandlers

type FinnHubFinanceHandlerFactory struct {
}

func CreateNewFinnHubFinanceHandlerFactory() (FinnHubFinanceHandlerFactory, error) {
	return FinnHubFinanceHandlerFactory{}, nil
}

func (f FinnHubFinanceHandlerFactory) CreateFinanceHandler(stockSymbol string) (FinanceHandler, error) {
	financeHandler, createNewFinnHubFinanceHandlerErr := CreateNewFinnHubFinanceHandler()

	if createNewFinnHubFinanceHandlerErr != nil {
		return financeHandler, createNewFinnHubFinanceHandlerErr
	}

	setStockSymbolErr := financeHandler.SetStockSymbol(stockSymbol)

	if setStockSymbolErr != nil {
		return financeHandler, setStockSymbolErr
	}

	return financeHandler, nil
}
