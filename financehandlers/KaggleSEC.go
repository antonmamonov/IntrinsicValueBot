package financehandlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/antonmamonov/IntrinsicValueBot/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// https://www.kaggle.com/finnhub/reported-financials
type KaggleSEC struct {
	dirPath                string
	StockSymbol            string
	MongoDbConnectionURI   string
	MongoSECCollectionName string
	db                     *mongo.Database
}

func (financeHandler KaggleSEC) SetStockSymbol(stockSymbol string) error {
	financeHandler.StockSymbol = stockSymbol

	return nil
}

func (financeHandler KaggleSEC) GetCurrentStockSymbol() (string, error) {
	stockSymbol := financeHandler.StockSymbol

	if stockSymbol == "" {
		return "", errors.New("stockSymbol is not set! Please set it first.")
	}

	return stockSymbol, nil
}

func (financeHandler KaggleSEC) GetEarningsPerShareGrowthForPeriod(periods int, periodType string) (EarningsPerShareGrowth, error) {
	earningsPerShareGrowth := EarningsPerShareGrowth{}

	return earningsPerShareGrowth, errors.New("not supported yet")
}

func (financeHandler KaggleSEC) GetIntrinsicValue() (IntrinsicValue, error) {

	intrinsicValue := IntrinsicValue{}

	return intrinsicValue, nil
}

type KaggleSecEntryDataElement struct {
	Concept string  `bson:"concept"`
	Label   string  `bson:"label"`
	Unit    string  `bson:"Unit"`
	Value   float64 `bson:"value"`
}

type KaggleSecEntryData struct {
	Bs []KaggleSecEntryDataElement `bson:"bs"`
	Cf []KaggleSecEntryDataElement `bson:"cf"`
	Ic []KaggleSecEntryDataElement `bson:"ic"`
}

type KaggleSecEntry struct {
	ID        string             `bson:"_id"`
	StartDate string             `bson:"startDate"`
	EndDate   string             `bson:"endDate"`
	Year      string             `bson:"year"`
	Quarter   string             `bson:"quarter"`
	Symbol    string             `bson:"symbol"`
	Data      KaggleSecEntryData `bson:"data"`
}

func (financeHandler KaggleSEC) LoadDataIntoMongo() error {

	quarters, quarterErr := ioutil.ReadDir(financeHandler.dirPath)
	if quarterErr != nil {
		log.Fatal(quarterErr)
	}

	collection := financeHandler.db.Collection(financeHandler.MongoSECCollectionName, &options.CollectionOptions{})

	for _, quarter := range quarters {
		// fmt.Println(quarter.Name())
		// fmt.Println(quarter.IsDir())
		// fmt.Println("")

		if quarter.IsDir() {
			files, fileErr := ioutil.ReadDir(financeHandler.dirPath + "/" + quarter.Name())
			if fileErr != nil {
				log.Fatal(fileErr)
			}

			for _, file := range files {
				fileName := file.Name()
				if fileName[len(fileName)-5:] == ".json" {
					jsonReadFileData, jsonReadFileErr := ioutil.ReadFile(financeHandler.dirPath + "/" + quarter.Name() + "/" + fileName)

					if jsonReadFileErr != nil {
						log.Fatal(jsonReadFileErr)
					}

					var kaggleSecEntry KaggleSecEntry

					json.Unmarshal(jsonReadFileData, &kaggleSecEntry)

					fmt.Println(kaggleSecEntry)

					kaggleSecEntry.ID = kaggleSecEntry.Symbol + "__" + kaggleSecEntry.Quarter + "__" + kaggleSecEntry.Year

					insertResult, insertOneErr := collection.InsertOne(context.TODO(), kaggleSecEntry)

					fmt.Println(insertResult)

					if insertOneErr != nil {
						log.Print(insertOneErr)
					}
				}
			}
		}
	}

	return nil
}

func CreateNewKaggleSEC(stockSymbol string) (FinanceHandler, error) {

	defaultMongoConnectionURI := "mongodb://127.0.0.1:27017"

	dbWrapper := db.NewMongoWrapper(defaultMongoConnectionURI)

	financeHandler := KaggleSEC{"/Users/antonyaar/Downloads/archive", stockSymbol, defaultMongoConnectionURI, "sec", dbWrapper}

	financeHandler.LoadDataIntoMongo()

	return financeHandler, nil
}
