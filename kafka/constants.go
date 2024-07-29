package kafka

var Brokers = []string{
	"localhost:9093",
	"localhost:9094",
	"localhost:9095",
}

const Topic = "test"

const Group = "my-group"

var WriteToKafka bool
