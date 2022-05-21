package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
	"github.com/vontikov/kafka-test-tool/internal/logging"
	"github.com/vontikov/kafka-test-tool/internal/util"
)

const (
	KafkaBootstrapServersEnv = "KAFKA_BOOTSTRAP_SERVERS"
	ListenAddressEnv         = "LISTEN_ADDRESS"
)

var (
	// App is the app name.
	App string = "N/A"
	// Version is the app version.
	Version string = "N/A"
)

var (
	logLevelFlag       = flag.String("log-level", "info", "Log level: trace|debug|info|warn|error|none")
	listenAddrFlag     = flag.String("listen-address", ":8080", "Listen address.")
	kafkaBootstrapFlag = flag.String("kafka-bootstrap", "localhost:9092", "Kafka bootstrap servers")
)

func main() {
	flag.Parse()

	logging.SetLevel(*logLevelFlag)
	logger := logging.NewLogger(App)

	hostname, err := os.Hostname()
	util.PanicOnError(err)

	logger.Info("starting", "name", App, "version", Version, "hostname", hostname)

	if v, ok := os.LookupEnv(ListenAddressEnv); ok {
		listenAddrFlag = &v
	}
	if v, ok := os.LookupEnv(KafkaBootstrapServersEnv); ok {
		kafkaBootstrapFlag = &v
	}

	logger.Info("Listen address", "value", *listenAddrFlag)
	logger.Info("Kafka bootstrap", "value", *kafkaBootstrapFlag)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	router := mux.NewRouter()
	router.PathPrefix("/topics/{topic}").Handler(topics(ctx)).Methods("POST")

	srv := &http.Server{
		Addr:    *listenAddrFlag,
		Handler: router,
	}

	logger.Info("listening", "address", *listenAddrFlag)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen", "error", err)
		}
		logger.Info("shutdown")
	}()

	go func() {
		<-ctx.Done()
		_ = srv.Close()
	}()

	sig := <-signals
	logger.Info("received", "signal", sig)
	cancel()
}

func topics(ctx context.Context) http.Handler {
	logger := logging.NewLogger("topics")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		topic := mux.Vars(r)["topic"]
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error("Could not read request body", "reason", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(data) == 0 {
			logger.Error("Request error", "reason", "empty request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		logger.Info("Sending message", "topic", topic, "data", string(data))
		kw := kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{*kafkaBootstrapFlag},
			Topic:   topic,
			Logger:  logger,
		})
		defer kw.Close()

		if err := kw.WriteMessages(ctx, kafka.Message{
			Key:   []byte(uuid.New().String()),
			Value: data,
		}); err != nil {
			logger.Error("Kafka error", "reason", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
