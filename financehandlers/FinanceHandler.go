package financehandlers

type IntrinsicValue struct {
	StockSymbol           string
	Value                 float64
	Price                 float64
	MarginOfSafety        float64 // the closer to 1, the greater the safety (ie, 0.8 - 20% margin of safety)
	AssetToLiabilityRatio float64
}

type EarningsPerShareGrowth struct {
	EPSGrowth float64 // value between 0 & 1
	LatestEPS float64
}

type FinanceHandler interface {
	GetEarningsPerShareGrowthForPeriod(periods int, periodType string) (EarningsPerShareGrowth, error)
	SetStockSymbol(stockSymbol string) error
	GetCurrentStockSymbol() (string, error)
	GetIntrinsicValue() (IntrinsicValue, error)
}
