package main

//this is a simple API server so that you can use curl on the console to push events into NATS that the consumer will process
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stewelarend/consumer/message"
)

func main() {
	s := apiServer{
		topic: os.Getenv("NATS_TOPIC"),
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
		log.Fatal("Error establishing connection to NATS:", err)
	}

	fmt.Println("Connected to NATS at:", s.nc.ConnectedUrl())
	http.HandleFunc("/", s.baseRoot)
	http.HandleFunc("/healthz", s.healthz)
	http.HandleFunc("/request", s.request)
	http.HandleFunc("/publish/", s.publish)

	fmt.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type apiServer struct {
	nc    *nats.Conn
	topic string
}

func (s apiServer) baseRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Basic NATS based microservice example v0.0.1")
}

func (s apiServer) healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

//this handler sends a request to NATS, waiting for a response
//the consumer will see the Reply path is set, and push a response
func (s apiServer) request(w http.ResponseWriter, r *http.Request) {
	requestAt := time.Now()
	request := message.Request{
		Message: message.Message{
			Timestamp: requestAt,
		},
		TTL:  5 * time.Second,
		Data: r.URL.Query(),
	}
	jsonRequest, _ := json.Marshal(request)
	natsResponse, err := s.nc.Request(s.topic, jsonRequest, 5*time.Second)
	if err != nil {
		log.Println("NATS request failed:", err)
	} else {
		duration := time.Since(requestAt)
		fmt.Fprintf(w, "NATS request success. Duration(%+v) Response: %v\n", duration, string(natsResponse.Data))

		var response message.Response
		if err := json.Unmarshal(natsResponse.Data, &response); err != nil {
			fmt.Printf("failed to decode JSON response: %v\n", err)
		} else {
			if err := response.Validate(); err != nil {
				fmt.Printf("invalid response: %v\n", err)
			} else {
				fmt.Printf("Valid Response: %+v\n", response)
			}
		}
	}
}

//this handler publishes an event not expecting a response
//the consumer will see the Reply path is unset and will not push a response at all
func (s apiServer) publish(httpRes http.ResponseWriter, httpReq *http.Request) {
	requestAt := time.Now()
	if httpReq.Method != http.MethodPost {
		http.Error(httpRes, "expecting method POST only", http.StatusMethodNotAllowed)
		return
	}
	if !strings.HasPrefix(httpReq.URL.Path, "/publish/") || len(httpReq.URL.Path) <= len("/publish/") {
		http.Error(httpRes, "expecting event type in URL after /publish/...", http.StatusBadRequest)
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
		Type: httpReq.URL.Path[9:], //skip '/publish/'
		Data: reqData,
	}
	jsonMessage, _ := json.Marshal(message)
	err := s.nc.Publish(s.topic, jsonMessage)
	if err != nil {
		log.Println("NATS publish failed:", err)
	}
	duration := time.Since(requestAt)
	fmt.Fprintf(httpRes, "NATS publish success. Duration(%+v)\n", duration)
}
