package server

import (
	"database/sql"
	"fmt"
	"github.com/AgileBits/go-redis-queue/redisqueue"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/iseroukhov/brave-new-billing/pkg/entities/payment"
	"github.com/iseroukhov/brave-new-billing/pkg/http/handlers"
	"github.com/iseroukhov/brave-new-billing/pkg/http/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	config *Config
	router *mux.Router
	logger *logrus.Logger
}

func New(config *Config) *Server {
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./templates/static"))))
	logger := logrus.New()

	return &Server{
		config: config,
		router: router,
		logger: logger,
	}
}

func (s *Server) Run() error {
	if level, err := logrus.ParseLevel(s.config.Log.Level); err != nil {
		return err
	} else {
		s.logger.SetLevel(level)
	}

	mysql, err := s.MysqlDB()
	if err != nil {
		return err
	}

	queue, err := s.RedisQueue()
	if err != nil {
		return err
	}
	// repositories
	paymentRepo := payment.NewRepository(mysql, queue)

	// routing
	formHandler := handlers.NewFormHandler(s.logger, paymentRepo)
	s.router.HandleFunc("/payments/card/form", formHandler.Index()).Methods(http.MethodGet)
	s.router.HandleFunc("/payments/card/form", formHandler.Store()).Methods(http.MethodPost)

	registerHandler := handlers.NewRegisterHandler(s.logger, paymentRepo)
	s.router.HandleFunc("/register", registerHandler.Index()).Methods(http.MethodPost)

	paymentHandler := handlers.NewPaymentHandler(s.logger, paymentRepo)
	s.router.HandleFunc("/payments", paymentHandler.Index()).Methods(http.MethodGet)

	// middleware
	handler := middleware.API(s.router)
	handler = middleware.AccessLog(s.logger, handler)
	handler = middleware.Panic(s.logger, handler)

	s.logger.Infof("server started on %s port", s.config.Server.Addr)
	return http.ListenAndServe(s.config.Server.Addr, handler)
}

func (s *Server) MysqlDB() (*sql.DB, error) {
	c := s.config.DB
	dsn := ""
	switch c.Driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&interpolateParams=true&parseTime=true", c.User, c.Password, c.Host, c.Port, c.Database)
	default:
		return nil, fmt.Errorf("driver \"%s\" is not supported", c.Driver)
	}
	db, err := sql.Open(c.Driver, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (s *Server) RedisQueue() (*redisqueue.Queue, error) {
	conn, err := redis.Dial(s.config.Queue.Network, s.config.Queue.Host+":"+s.config.Queue.Port)
	if err != nil {
		return nil, err
	}
	q := redisqueue.New("queue", conn)
	return q, nil
}
