package main

import (
	"os"
	"time"

	"github.com/Ksenofontovas/spread/internal/repository"
	"github.com/Ksenofontovas/spread/internal/service"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	s := gocron.NewScheduler(time.UTC)

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	var exchanges []*repository.Exchange
	pairs := viper.GetStringSlice("pairs")

	garantex := repository.NewExchange(
		viper.GetString("garantex.name"),
		viper.GetString("garantex.url"),
		viper.GetDuration("garantex.timeout")*time.Second,
	)

	binance := repository.NewExchange(
		viper.GetString("binance.name"),
		viper.GetString("binance.url"),
		viper.GetDuration("binance.timeout")*time.Second,
	)

	exchanges = append(exchanges, garantex, binance)

	repos := repository.NewRepository(db, exchanges)
	service := service.NewService(repos)

	s.Every(10).Second().Do(func() {
		tickers, err := service.GetTickers(pairs)
		if err != nil {
			logrus.Error(err)
		}
		for _, ticker := range tickers {
			_, err := service.SaveTicker(ticker, time.Now())
			if err != nil {
				logrus.Error(err)
			}
		}
	})
	s.StartBlocking()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
