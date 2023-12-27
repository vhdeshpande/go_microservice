package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //name of the exchange
		"topic",      //type - topic of the exchange
		true,         //is the exchange durable?
		false,        // is this auto deleted?
		false,        //is this exchange just used internally?
		false,        //no-wait?
		nil,          //arguements?
	)

}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    //Name
		false, //is the exchange durable?
		false, //deleted when unused?
		true,  // is this exclusive
		false, //no-wait
		nil,   //arguements?
	)
}
