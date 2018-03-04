package queue

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

var (
	connection   *amqp.Connection
	channel      *amqp.Channel
	closeChannel chan *amqp.Error
	namespace    = "STEAM_"
)

func init() {
	closeChannel = make(chan *amqp.Error)
}

func RunConsumers() {

	//go func() {
	//
	//	for _ := range closeChannel {
	//		fmt.Println("Reconnecting")
	//		err := connect()
	//		if err != nil {
	//			fmt.Println(err.Error())
	//		}
	//	}
	//
	//}()

	go appConsumer()
	go changeConsumer()
	go packageConsumer()
	go playerConsumer()
}

func connect() (err error) {

	if connection != nil && channel != nil {
		return nil
	}

	connection, err = amqp.Dial(os.Getenv("STEAM_RABBIT"))
	connection.NotifyClose(closeChannel)
	if err != nil {
		return err
	}
	//defer connection.Close()

	channel, err = connection.Channel()
	if err != nil {
		return err
	}
	//defer channel.Close()

	//fmt.Println("connected")

	return nil
}

//// todo, Have consume queue and produce queue!!
//
//const (
//	QApp     = "app"
//	QPackage = "package"
//)
//
//var (
//	queues map[string]*queue
//)
//
//func setup() {
//
//	queues = make(map[string]*queue, 0)
//
//	newQueue(QApp)
//	newQueue(QPackage)
//}
//
//type queue struct {
//	connection   *amqp.Connection
//	channel      *amqp.Channel
//	queue        amqp.Queue
//	name         string
//	closeChannel chan *amqp.Error
//}
//
//func newQueue(key string) (q *queue) {
//
//	q = new(queue)
//	q.name = key
//	q.init()
//
//	queues[key] = q
//
//	return q
//}
//
//func (q *queue) runConsumer() {
//
//	go func() {
//
//	}()
//}
//
//func (q *queue) init() (err error) {
//
//	if q.connection != nil {
//		return nil
//	}
//
//	q.closeChannel = make(chan *amqp.Error)
//
//	q.connection, err = amqp.Dial(os.Getenv("STEAM_RABBIT"))
//	q.connection.NotifyClose(q.closeChannel)
//	if err != nil {
//		return err
//	}
//
//	q.channel, err = q.connection.Channel()
//	if err != nil {
//		return err
//	}
//
//	q.queue, err = q.channel.QueueDeclare(
//		"Steam_"+strings.Title(q.name), // name
//		true,                           // durable
//		false,                          // delete when unused
//		false,                          // exclusive
//		false,                          // no-wait
//		nil,                            // arguments
//	)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (q *queue) getQueue() (r amqp.Queue, err error) {
//
//	return r, nil
//}
//
//func (q *queue) produce(bytes []byte) (err error) {
//
//	if q.queue.Name == "" {
//		q.init()
//	}
//
//	err = q.channel.Publish(
//		"",           // exchange
//		q.queue.Name, // routing key
//		false,        // mandatory
//		false,        // immediate
//		amqp.Publishing{
//			DeliveryMode: amqp.Persistent,
//			ContentType:  "text/plain",
//			Body:         bytes,
//		})
//	if err != nil {
//		return err
//	}
//
//	return nil
//
//}
//
//func (q *queue) consume(func(msg amqp.Delivery) (err error)) {
//
//}

func GetQeueus() (resp []Queue, err error) {

	req, err := http.NewRequest("GET", "http://localhost:15672/api/queues", nil)
	req.SetBasicAuth("guest", "guest")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	defer response.Body.Close()

	// Convert to bytes
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return resp, err
	}

	// Unmarshal JSON
	if err := json.Unmarshal(bytes, &resp); err != nil {
		return resp, err
	}

	for k, v := range resp {
		resp[k].Name = strings.Replace(v.Name, namespace, "", 1)
	}

	return resp, nil
}

type Queue struct {
	MessagesDetails struct {
		Rate float64 `json:"rate"`
	} `json:"messages_details"`
	Messages int `json:"messages"`
	MessagesUnacknowledgedDetails struct {
		Rate float64 `json:"rate"`
	} `json:"messages_unacknowledged_details"`
	MessagesUnacknowledged int `json:"messages_unacknowledged"`
	MessagesReadyDetails struct {
		Rate float64 `json:"rate"`
	} `json:"messages_ready_details"`
	MessagesReady int `json:"messages_ready"`
	ReductionsDetails struct {
		Rate float64 `json:"rate"`
	} `json:"reductions_details"`
	Reductions int `json:"reductions"`
	MessageStats struct {
		DeliverGetDetails struct {
			Rate float64 `json:"rate"`
		} `json:"deliver_get_details"`
		DeliverGet int `json:"deliver_get"`
		AckDetails struct {
			Rate float64 `json:"rate"`
		} `json:"ack_details"`
		Ack int `json:"ack"`
		RedeliverDetails struct {
			Rate float64 `json:"rate"`
		} `json:"redeliver_details"`
		Redeliver int `json:"redeliver"`
		DeliverNoAckDetails struct {
			Rate float64 `json:"rate"`
		} `json:"deliver_no_ack_details"`
		DeliverNoAck int `json:"deliver_no_ack"`
		DeliverDetails struct {
			Rate float64 `json:"rate"`
		} `json:"deliver_details"`
		Deliver int `json:"deliver"`
		GetNoAckDetails struct {
			Rate float64 `json:"rate"`
		} `json:"get_no_ack_details"`
		GetNoAck int `json:"get_no_ack"`
		GetDetails struct {
			Rate float64 `json:"rate"`
		} `json:"get_details"`
		Get int `json:"get"`
		PublishDetails struct {
			Rate float64 `json:"rate"`
		} `json:"publish_details"`
		Publish int `json:"publish"`
	} `json:"message_stats"`
	Node string `json:"node"`
	Arguments struct {
	} `json:"arguments"`
	Exclusive            bool   `json:"exclusive"`
	AutoDelete           bool   `json:"auto_delete"`
	Durable              bool   `json:"durable"`
	Vhost                string `json:"vhost"`
	Name                 string `json:"name"`
	MessageBytesPagedOut int    `json:"message_bytes_paged_out"`
	MessagesPagedOut     int    `json:"messages_paged_out"`
	BackingQueueStatus struct {
		AvgAckEgressRate  float64       `json:"avg_ack_egress_rate"`
		AvgAckIngressRate float64       `json:"avg_ack_ingress_rate"`
		AvgEgressRate     float64       `json:"avg_egress_rate"`
		AvgIngressRate    float64       `json:"avg_ingress_rate"`
		Delta             []interface{} `json:"delta"`
		Len               int           `json:"len"`
		Mode              string        `json:"mode"`
		NextSeqID         int           `json:"next_seq_id"`
		Q1                int           `json:"q1"`
		Q2                int           `json:"q2"`
		Q3                int           `json:"q3"`
		Q4                int           `json:"q4"`
		TargetRAMCount    string        `json:"target_ram_count"`
	} `json:"backing_queue_status"`
	HeadMessageTimestamp       interface{} `json:"head_message_timestamp"`
	MessageBytesPersistent     int         `json:"message_bytes_persistent"`
	MessageBytesRAM            int         `json:"message_bytes_ram"`
	MessageBytesUnacknowledged int         `json:"message_bytes_unacknowledged"`
	MessageBytesReady          int         `json:"message_bytes_ready"`
	MessageBytes               int         `json:"message_bytes"`
	MessagesPersistent         int         `json:"messages_persistent"`
	MessagesUnacknowledgedRAM  int         `json:"messages_unacknowledged_ram"`
	MessagesReadyRAM           int         `json:"messages_ready_ram"`
	MessagesRAM                int         `json:"messages_ram"`
	GarbageCollection struct {
		MinorGcs        int `json:"minor_gcs"`
		FullsweepAfter  int `json:"fullsweep_after"`
		MinHeapSize     int `json:"min_heap_size"`
		MinBinVheapSize int `json:"min_bin_vheap_size"`
		MaxHeapSize     int `json:"max_heap_size"`
	} `json:"garbage_collection"`
	State                     string        `json:"state"`
	RecoverableSlaves         interface{}   `json:"recoverable_slaves"`
	Consumers                 int           `json:"consumers"`
	ExclusiveConsumerTag      interface{}   `json:"exclusive_consumer_tag"`
	EffectivePolicyDefinition []interface{} `json:"effective_policy_definition"`
	OperatorPolicy            interface{}   `json:"operator_policy"`
	Policy                    interface{}   `json:"policy"`
	ConsumerUtilisation       interface{}   `json:"consumer_utilisation"`
	IdleSince                 string        `json:"idle_since"`
	Memory                    int           `json:"memory"`
}
