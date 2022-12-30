package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/bars-squad/ais-user-query-service/config"
	es "github.com/bars-squad/ais-user-query-service/databases/elasticsearch"
	"github.com/bars-squad/ais-user-query-service/jwt"
	"github.com/bars-squad/ais-user-query-service/middleware"
	"github.com/bars-squad/ais-user-query-service/modules/admin"
	"github.com/bars-squad/ais-user-query-service/pubsub"
	"github.com/bars-squad/ais-user-query-service/responses"
	"github.com/bars-squad/ais-user-query-service/server"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload" //for development
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

var (
	cfg                 *config.Config
	httpResponse        = responses.HttpResponseStatusCodesImpl{}
	healthCheckMessage  = "Application running properly"
	pageNotFoundMessage = "You're lost, double check the endpoint"
)

const (
	createdAdministratorTopic = "ais-user.administrator"
)

func init() {
	cfg = config.Load()
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(cfg.Logger.Formatter)
	logger.SetReportCaller(true)

	validate := validator.New()

	// set elasctisearch
	esClient, err := es.NewClientAdapter().Connect(context.Background(), logger)
	if err != nil {
		logger.Fatal(err)
	}

	// set publisher object
	saramaAsyncProducer, err := sarama.NewAsyncProducer(
		cfg.SaramaKafka.Addresses,
		cfg.SaramaKafka.Config,
	)
	if err != nil {
		logger.Fatal(err)
	}
	publisher := pubsub.NewSaramaKafkaProducerAdapter(logger, &pubsub.SaramaKafkaProducerAdapterConfig{
		AsyncProducer: saramaAsyncProducer,
	})

	// set jwt object
	privateKey := jwt.GetRSAPrivateKey("./secret/private.pem")
	publicKey := jwt.GetRSAPublicKey("./secret/public.pem")
	jsonWebToken := jwt.NewJSONWebToken(privateKey, publicKey)

	// set basic auth
	basicAuth := middleware.NewBasicAuth(cfg.BasicAuth.Username, cfg.BasicAuth.Password)

	router := mux.NewRouter()
	router.HandleFunc("/", index)
	// http.Handle("/", router)
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	adminRepository := admin.NewRepository(logger, esClient)

	adminUsecase := admin.NewUsecase(&admin.Property{
		ServiceName:  cfg.Application.Name,
		Logger:       logger,
		Repository:   adminRepository,
		JSONWebToken: jsonWebToken,
		// Session:      sess,
		// Publisher:                publisher,
	})

	admin.NewHTTPHandler(logger, validate, router, basicAuth, adminUsecase)

	createdAccountEventHandler := admin.NewCreatedAdministratorEventHandler(logger, validate, adminUsecase)
	createdAccountConsumerHandler := pubsub.NewDefaultSaramaConsumerGroupHandler(cfg.Application.Name, createdAccountEventHandler, nil)

	updatedTPUserSubcriber, err := pubsub.NewSaramaKafkaConsumerGroupFullConfigAdapter(
		logger, cfg.SaramaKafka.Addresses, cfg.Application.Name, []string{createdAdministratorTopic},
		createdAccountConsumerHandler, cfg.SaramaKafka.Config,
	)

	if err != nil {
		logger.Fatal(err)
	}

	updatedTPUserSubcriber.Subscribe()

	handler := cors.New(cors.Options{
		AllowedOrigins:   cfg.Application.AllowedOrigins,
		AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization", "X-RECAPTCHA-TOKEN"},
		AllowCredentials: true,
	}).Handler(router)

	server := server.NewServer(logger, handler, cfg.Application.Port)
	server.Start()

	// When we run this program it will block waiting for a signal. By typing ctrl-C, we can send a SIGINT signal, causing the program to print interrupt and then exit.
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	// closing service for a gracefull shutdown.
	server.Close()
	publisher.Close()
}

func index(w http.ResponseWriter, r *http.Request) {
	responses.REST(w, httpResponse.Ok("").NewResponses(nil, healthCheckMessage))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	responses.REST(w, httpResponse.NotFound("PAGES_NOT_FOUND").NewResponses(nil, pageNotFoundMessage))
}
