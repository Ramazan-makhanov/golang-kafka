package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

//MessageSMS struct
type MessageSMS struct {
	To   string `json:"to"`
	From string `json:"from"`
	text string `json:"text"`
}

func smsPostHandler(w http.ResponseWriter, r *http.Request) {

	//Save data into Sms struct
	var sms MessageSMS
	sms.To = r.FormValue("receiver")
	sms.From = r.FormValue("sender")
	sms.Text = r.FormValue("text")

	saveSmsToKafka(sms)

	//Convert sms struct into json
	jsonString, err := json.Marshal(sms)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//Set content-type http header
	w.Header().Set("content-type", "application/json")

	//Send back data as response
	w.Write(jsonString)

}

func saveJobToKafka(sms MessageSMS) {

	fmt.Println("save to kafka")

	jsonString, err := json.Marshal(sms)

	smsString := string(jsonString)
	fmt.Print(smsString)

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}

	// Produce messages to topic (asynchronously)
	topic := "sms-topic1"
	for _, word := range []string{string(smsString)} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}
}

func main() {

	http.HandleFunc("/", smsPostHandler)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
