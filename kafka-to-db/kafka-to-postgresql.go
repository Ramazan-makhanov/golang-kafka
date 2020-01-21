package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "karajon4565"
	dbname   = "billings"
)

//MessageSMS struct
type MessageSMS struct {
	To   string `json:"to"`
	From string `json:"from"`
	text string `json:"text"`
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	receiveFromKafka()
}

func receiveFromKafka() {

	fmt.Println("Start receiving from Kafka")
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "group-id-1",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"sms-topic1"}, nil)

	for {
		msg, err := c.ReadMessage(-1)

		if err == nil {
			fmt.Printf("Received from Kafka %s: %s\n", msg.TopicPartition, string(msg.Value))
			sms := msg.Value
			saveSMSToDB(sms)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}

	c.Close()

}

func saveSMSToDB(sms MessageSMS) {
	//  inserting values
	_, err = db.Exec("INSERT INTO billings VALUES ($1, $2, $3)", sms.To, sms.From, sms.Text)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Message saved!")

}
