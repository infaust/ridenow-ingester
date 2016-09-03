package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"log"
	"ridenow/ingester"
	"ridenow/ingester/models"
	"ridenow/ingester/queue"
	"sync"
)

type Env struct {
	db    models.Datastore
	queue *queue.QueueProducer
}

func main() {
	// get config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}
	dbUser := viper.GetString("DB_USER")
	dbPass := viper.GetString("DB_PASSWORD")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/forecasts?sslmode=disable", dbUser, dbPass, dbHost, dbPort)

	qUser := viper.GetString("AMQP_USER")
	qPass := viper.GetString("AMQP_PASSWORD")
	qHost := viper.GetString("AMQP_HOST")
	qPort := viper.GetString("AMQP_PORT")
	qUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", qUser, qPass, qHost, qPort)

	// set up postgresql connection
	db, err := models.NewDB(dbUrl)
	if err != nil {
		log.Panic(err)
	}
	// set up rabbitmq connection
	queue, err := queue.NewQueueProducer(qUrl)
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db, queue}
	locations, err := env.db.RelevantLocations()
	if err != nil {
		log.Panic(err)
	}
	for _, loc := range locations {
		fmt.Printf("Location %d: %s\n", loc.Id, loc.Name)
	}

	var tokens = make(chan struct{}, 2) // IO semaphore
	ch := make(chan []msw.WeatherEntry, 3)
	var wg sync.WaitGroup
	for _, l := range locations {
		wg.Add(1)
		go func(l *models.Location) {
			defer wg.Done()
			tokens <- struct{}{} // acquire a token
			scrapper := msw.NewScrapper(l.Url)
			entries, _ := scrapper.GetForecastEntries()
			<-tokens // release the token
			ch <- entries
		}(l)
	}

	// closer
	go func() {
		wg.Wait()
		close(ch)
	}()

	for entries := range ch {
		for _, we := range entries {
			forecast := models.NewForecast(we.WaveHeightM, we.SwellPeriodSecs, we.Id, we.Time)
			changed, err := env.db.StoreForecast(forecast)
			if err != nil {
				log.Panic(err)
			}
			if changed || true {
				time := int64(forecast.Time.UnixNano())
				msg := msw.Forecast{
					Id:              &forecast.Id,
					LocationId:      &forecast.LocationId,
					WaveHeightM:     &forecast.WaveHeightM,
					SwellPeriodSecs: &forecast.SwellPeriodSecs,
					Time:            &time,
				}
				fmt.Printf("Forecast %d changed: %t\n", msg.LocationId, changed)
				bytes, err := proto.Marshal(&msg)
				err = queue.Send("ridenow.forecasts.update", bytes)
				if err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
