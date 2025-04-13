package events

/*
type SignalHandlerOrderBook struct{}


func (s *SignalHandlerOrderBook) HandleMessage(msg *sarama.ConsumerMessage) {
	log.Printf("[Signal] Processing: %s", msg.Topic)
	var payload dtos.OrderBook
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error unmarshalling message for order book signal",
			Message: "Error unmarshalling message into BSON: " + err.Error(),
			Type:    "error",
			Entity:  "orderbook",
			Data:    string(msg.Value),
		})
		return
	}

	var pgDb = database.PgClient()
	var intervals []entities.SignalInterval

	pgDb.Where("symbol = ?", payload.Symbol).Find(&intervals)

	for _, interval := range intervals {
		key50 := fmt.Sprintf("%s:%s:ma50", payload.Symbol, interval.Interval)
		key200 := fmt.Sprintf("%s:%s:ma200", payload.Symbol, interval.Interval)

		rdb := cache.RedisClient()

		pipe := rdb.Pipeline()

		pipe.RPush(context.Background(), key50)
		pipe.LTrim(context.Background(), key50, -50, -1)

		pipe.RPush(context.Background(), key200)
		pipe.LTrim(context.Background(), key200, -200, -1)

	}

	_, err = pipe.Exec(context.Background())
	if err != nil {
		log.Println("Redis pipeline error:", err)
	}

	signal, err := checkForSignal(rdb, payload.Symbol, timeframe)
	if err != nil {
		log.Println("Signal check error:", err)
		return
	}
	log.Printf("[%s][%s] Signal: %s", payload.Symbol, timeframe, signal)

	var indicators = map[string]interface{}{
		"MA50":  signal.MA50,
		"MA200": signal.MA200,
	}
	indData, _ := json.Marshal(indicators)
	lastTrade, _ := json.Marshal(payload)

	var signalData = &dtos.Signal{
		Symbol:        payload.Symbol,
		Timeframe:     timeframe,
		Signal:        signal.Signal,
		IndicatorData: string(indData),
		LastTrade:     string(lastTrade),
	}

	signalJson, _ := json.Marshal(signalData)
	client := KafkaClientNew()
	_, _, err = client.Produce("market-signal", signalData.Symbol, signalJson)
	if err != nil {
		log.Println("Kafka signal produce error:", err)
	}

}*/

type SignalHandlerAggTrade struct{}

/*
func (s *SignalHandlerAggTrade) HandleMessage(msg *sarama.ConsumerMessage) {
	log.Printf("[Signal] Processing: %s", msg.Topic)

	var payload dtos.AggTrade
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error unmarshalling message for order book signal",
			Message: "Error unmarshalling message into BSON: " + err.Error(),
			Type:    "error",
			Entity:  "aggTrade",
			Data:    string(msg.Value),
		})
		return
	}

	var pgDb = database.PgClient()
	var intervals []entities.SignalInterval

	err := pgDb.Where("symbol = ?", payload.Symbol).Find(&intervals).Error
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error fetching intervals from Postgres",
			Message: "Error fetching intervals from Postgres: " + err.Error(),
			Type:    "error",
			Entity:  "aggTrade",
			Data:    string(msg.Value),
		})
		return
	}

	for _, interval := range intervals {
		key50 := fmt.Sprintf("%s:%s:ma50", payload.Symbol, interval.Interval)
		key200 := fmt.Sprintf("%s:%s:ma200", payload.Symbol, interval.Interval)

		rdb := cache.RedisClient()
		pipe := rdb.Pipeline()
		price, _ := decimal.NewFromString(payload.Price)

		if exists, _ := rdb.Exists(context.Background(), key50).Result(); exists == 0 {
			if exists, _ := rdb.Exists(context.Background(), key200).Result(); exists > 0 {
				pipe.Del(context.Background(), key200)
			}

			candleSticks, err := candlestick.GetCandleSticksAndUpdate(interval.ExchangeID, payload.Symbol, interval.Interval, 200)
			if err != nil {
				log.Println("Error getting candlesticks:", err)
				return
			}
			var cnt = 0
			for _, candle := range candleSticks {
				if cnt < 50 {
					pipe.RPush(context.Background(), key50, candle.Close)
					cnt++
				}
				pipe.RPush(context.Background(), key200, candle.Close)
			}

		} else {
			pipe.RPush(context.Background(), key50, price.String())
			pipe.LTrim(context.Background(), key50, -50, -1)

			pipe.RPush(context.Background(), key200, price.String())
			pipe.LTrim(context.Background(), key200, -200, -1)
		}

		_, err := pipe.Exec(context.Background())
		if err != nil {
			log.Println("Redis pipeline error:", err)
		}

		signal, err := checkForSignal(rdb, payload.Symbol, interval.Interval)
		if err != nil {
			log.Println("Signal check error:", err)
			return
		}
		log.Printf("[%s][%s] Signal: %s", payload.Symbol, interval.Interval, signal)

	}
}*/

/*
func checkForSignal(rdb *redis.Client, symbol, timeframe string) (dtos.MaData, error) {
	var res dtos.MaData
	key50 := fmt.Sprintf("%s:%s:ma50", symbol, timeframe)
	key200 := fmt.Sprintf("%s:%s:ma200", symbol, timeframe)
	lastSignalKey := fmt.Sprintf("%s:%s:lastSignal", symbol, timeframe)

	pipe := rdb.Pipeline()
	res50 := pipe.LRange(context.Background(), key50, 0, -1)
	res200 := pipe.LRange(context.Background(), key200, 0, -1)

	_, err := pipe.Exec(context.Background())
	if err != nil {
		return res, fmt.Errorf("Redis read error: %v", err)
	}

	var lastSignal string
	lastSignalRes := rdb.Get(context.Background(), lastSignalKey)
	if lastSignalRes.Err() != nil {
		lastSignal = consts.HoldSignal
	} else {
		lastSignal = lastSignalRes.Val()
	}

	vals50, _ := res50.Result()
	vals200, _ := res200.Result()

	if len(vals50) < 50 || len(vals200) < 200 {
		return res, nil
	}

	ma50 := average(vals50)
	ma200 := average(vals200)

	res.MA50 = ma50.String()
	res.MA200 = ma200.String()

	if lastSignal == consts.BuySignal && ma50.GreaterThan(ma200) {
		return res, errors.New(consts.AlreadyInBuy)
	} else if lastSignal == consts.SellSignal && ma50.LessThan(ma200) {
		return res, errors.New(consts.AlreadyInSell)
	}

	jsonData, _ := json.Marshal(
		map[string]interface{}{
			"MA50":  res.MA50,
			"MA200": res.MA200,
		},
	)

	signal := entities.Signal{
		Symbol:    symbol,
		Timeframe: timeframe,
		Indicator: string(jsonData),
		LastTrade: "{}",
	}

	if ma50.GreaterThan(ma200) {
		rdb.Set(context.Background(), lastSignalKey, consts.BuySignal, 0)
		res.Signal = consts.BuySignal
		signal.Signal = consts.BuySignal
	}

	if ma50.LessThan(ma200) {
		rdb.Set(context.Background(), lastSignalKey, consts.SellSignal, 0)
		res.Signal = consts.SellSignal
		signal.Signal = consts.SellSignal
	}

	db := database.PgClient()
	if err := db.Create(&signal).Error; err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error inserting signal into Postgres",
			Message: "Error inserting signal into Postgres: " + err.Error(),
			Type:    "error",
			Entity:  "signal",
			Data:    fmt.Sprintf("Symbol: %s, Signal: %s", signal.Symbol, signal.Signal),
		})

		log.Printf("Error inserting signal into Postgres: %v", err)

	}
	log.Printf("Postgres updated for symbol %s with signal %s", signal.Symbol, signal.Signal)

	return res, errors.New(consts.AlreadyInHold)
}

func average(values []string) decimal.Decimal {
	var sum decimal.Decimal
	for _, val := range values {
		f, err := decimal.NewFromString(val)
		if err != nil {
			continue
		}
		sum = sum.Add(f)
	}
	return sum.Div(decimal.NewFromInt(int64(len(values))))
}
*/
