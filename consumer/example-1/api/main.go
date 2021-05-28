package main

//this is a simple API server so that you can use curl on the console to push events into NATS that the consumer will process
import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/nats-io/nats.go"
	"github.com/stewelarend/consumer/message"
	"github.com/stewelarend/logger"
)

var log = logger.New()

func main() {
	s := apiServer{
		topic: os.Getenv("NATS_TOPIC"),
	}
	if s.topic == "" {
		panic("NATS_TOPIC is not defined")
	}

	//connect to NATS before starting the api server
	uri := os.Getenv("NATS_URI")
	var err error
	for i := 0; i < 5; i++ {
		s.nc, err = nats.Connect(uri)
		if err == nil {
			break
		}
		fmt.Println("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		panic(fmt.Errorf("Error establishing connection to NATS:", err))
	}

	fmt.Println("Connected to NATS at:", s.nc.ConnectedUrl())

	//connect to KAFKA as well
	s.kp, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka-01.stagealot.com:6667,kafka-02.stagealot.com:6667",
		"client.id":         "test-producer", //socket.gethostname(),
		"acks":              "all"})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	//create chan for producer push results - just to log errors
	s.kafkaDeliveryChan = make(chan kafka.Event, 10000)
	go func() {
		for {
			e := <-s.kafkaDeliveryChan
			m := e.(*kafka.Message)
			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				fmt.Printf("Delivered message to topic(%s).partition(%d).offset(%v)\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		}
	}()
	fmt.Println("Connected to Kafka\n")

	//setup HTTP routes
	http.HandleFunc("/", s.baseRoot)
	http.HandleFunc("/healthz", s.healthz)
	http.HandleFunc("/nats/request/", s.natsRequest)   //send request, wait for response
	http.HandleFunc("/nats/publish/", s.natsPublish)   //publish event, expect no reply
	http.HandleFunc("/kafka/produce/", s.kafkaProduce) //produce event, expect no reply

	//start HTTP server
	fmt.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(fmt.Errorf("HTTP Server failed: %v", err))
	}
}

type apiServer struct {
	nc                *nats.Conn
	kp                *kafka.Producer
	kafkaDeliveryChan chan kafka.Event
	topic             string
}

func (s apiServer) baseRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Basic NATS based microservice example v0.0.1")
}

func (s apiServer) healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

//this handler sends a request to NATS, waiting for a response
//the consumer will see the Reply path is set, and push a response
func (s apiServer) natsRequest(httpRes http.ResponseWriter, httpReq *http.Request) {
	requestAt := time.Now()
	if httpReq.Method != http.MethodPost {
		http.Error(httpRes, "expecting method POST only", http.StatusMethodNotAllowed)
		return
	}
	prefix := "/nats/request/"
	if !strings.HasPrefix(httpReq.URL.Path, prefix) || len(httpReq.URL.Path) <= len(prefix) {
		http.Error(httpRes, fmt.Sprintf("expecting event type in URL after %s ...", prefix), http.StatusBadRequest)
		return
	}

	var reqData map[string]interface{}
	if err := json.NewDecoder(httpReq.Body).Decode(&reqData); err != nil {
		http.Error(httpRes, fmt.Sprintf("cannot parse JSON body: %v", err), http.StatusBadRequest)
		return
	}
	request := message.Request{
		Message: message.Message{
			Timestamp: requestAt,
		},
		TTL:       5 * time.Second,
		Operation: httpReq.URL.Path[len(prefix):], //skip prefix
		Data:      reqData,
	}
	jsonRequest, _ := json.Marshal(request)
	natsResponse, err := s.nc.Request(s.topic, jsonRequest, 5*time.Second)
	if err != nil {
		http.Error(httpRes, fmt.Sprintf("NATS request to topic(%s) failed: %v", s.topic, err), http.StatusServiceUnavailable)
		return
	}
	duration := time.Since(requestAt)
	var response message.Response
	if err := json.Unmarshal(natsResponse.Data, &response); err != nil {
		http.Error(httpRes, fmt.Sprintf("failed to decode JSON response: %v\n", err), http.StatusInternalServerError)
		return
	}

	if err := response.Validate(); err != nil {
		http.Error(httpRes, fmt.Sprintf("invalid response: %v\n", err), http.StatusInternalServerError)
	}

	jsonResponse, _ := json.Marshal(response)
	httpRes.Header().Set("Content-Type", "application/json")
	httpRes.Write(jsonResponse)
	log.Debugf("Request(%+v)->Response(%+v) Duration(%+v)", request, response, duration)
} //apiServer.natsRequest()

//this handler publishes an event not expecting a response
//the consumer will see the Reply path is unset and will not push a response at all
func (s apiServer) natsPublish(httpRes http.ResponseWriter, httpReq *http.Request) {
	requestAt := time.Now()
	if httpReq.Method != http.MethodPost {
		http.Error(httpRes, "expecting method POST only", http.StatusMethodNotAllowed)
		return
	}
	prefix := "/nats/publish/"
	if !strings.HasPrefix(httpReq.URL.Path, prefix) || len(httpReq.URL.Path) <= len(prefix) {
		http.Error(httpRes, fmt.Sprintf("expecting event type in URL after %s ...", prefix), http.StatusBadRequest)
		return
	}

	var reqData map[string]interface{}
	if err := json.NewDecoder(httpReq.Body).Decode(&reqData); err != nil {
		http.Error(httpRes, fmt.Sprintf("cannot parse JSON body: %v", err), http.StatusBadRequest)
		return
	}

	message := message.Event{
		Message: message.Message{
			Timestamp: requestAt,
		},
		Type: httpReq.URL.Path[len(prefix):], //skip prefix
		Data: reqData,
	}
	jsonMessage, _ := json.Marshal(message)
	err := s.nc.Publish(s.topic, jsonMessage)
	if err != nil {
		http.Error(httpRes, fmt.Sprintf("NATS publish failed:", err), http.StatusServiceUnavailable)
		return
	}
	//success
	duration := time.Since(requestAt)
	fmt.Fprintf(httpRes, "NATS publish to topic(%s) type(%s) success. Duration(%+v)\n", s.topic, message.Type, duration)
} //apiServer.natsPublish()

//this handler publishes an event not expecting a response
//the consumer will see the Reply path is unset and will not push a response at all
func (s apiServer) kafkaProduce(httpRes http.ResponseWriter, httpReq *http.Request) {
	requestAt := time.Now()
	if httpReq.Method != http.MethodPost {
		http.Error(httpRes, "expecting method POST only", http.StatusMethodNotAllowed)
		return
	}
	prefix := "/kafka/produce/"
	if !strings.HasPrefix(httpReq.URL.Path, prefix) || len(httpReq.URL.Path) <= len(prefix) {
		http.Error(httpRes, fmt.Sprintf("expecting event type in URL after %s ...", prefix), http.StatusBadRequest)
		return
	}

	var reqData map[string]interface{}
	if err := json.NewDecoder(httpReq.Body).Decode(&reqData); err != nil {
		http.Error(httpRes, fmt.Sprintf("cannot parse JSON body: %v", err), http.StatusBadRequest)
		return
	}

	message := message.Event{
		Message: message.Message{
			Timestamp: requestAt,
		},
		Type: httpReq.URL.Path[len(prefix):], //skip prefix
		Data: reqData,
	}
	jsonMessage, _ := json.Marshal(message)

	//async push to kafka
	//delivery_chan will get the push result
	if err := s.kp.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &s.topic, Partition: kafka.PartitionAny},
			Value:          jsonMessage,
		},
		s.kafkaDeliveryChan,
	); err != nil {
		http.Error(httpRes, fmt.Errorf("failed to produce event: %s", err).Error(), http.StatusServiceUnavailable)
	}
	//success
	duration := time.Since(requestAt)
	fmt.Fprintf(httpRes, "Kafka produce success. Duration(%+v)\n", duration)
} //apiServer.kafkaProduce()
