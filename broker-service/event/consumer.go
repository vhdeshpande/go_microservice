package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

// type used for receiving events from the queue
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// create instance of the consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	// setup the consumer, open up the channel and declare an exchange
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil

}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	// return the result of declaring the exchange
	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Listens to the queue, listens to specific topics
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	// get a random queue
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		// bind our channel to each of this topics
		ch.QueueBind(
			q.Name,
			s, // topic
			"logs_topic",
			false, //no-wait
			nil,   //arguements?
		)

		if err != nil {
			return err
		}
	}

	// look for messages
	messages, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer
		true,   //auto acknowledge
		false,  // is it exclusive
		false,  // is it internal
		false,  //no-wait
		nil,    //arguements?
	)
	if err != nil {
		return err
	}

	// consumer forever until we exit
	// make a channel
	// keeep running in its go routine
	forever := make(chan bool)
	go func() {
		// current iteration
		for d := range messages {
			var payload Payload
			// read json into payload, current iteration of the messages
			// body is unmarshalled into payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	// keep the channel going forever
	<-forever

	return nil

}

// take an action based on the name of the event that we get pushed to us from the queue
func handlePayload(payload Payload) {
	switch payload.Name {

	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	case "auth":
		// authenticate

	// can have as many cases as you want, as long as you write the logic

	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}
	return nil
}
