package router

import (
	"time"

	"bullshape/utils/kafkala"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// handlerDefaultTimeout is the timeout the handlers pass to the inner layers.
const handlerDefaultTimeout = 5 * time.Second

// The Server is used as a container for the most important dependencies.
type Server struct {
	DB     *gorm.DB
	Logger *zap.SugaredLogger
	Kafka  *kafkala.CompaniesProducer
}

// NewServer returns a pointer to a new Server.
func NewServer(logger *zap.SugaredLogger, db *gorm.DB, kafka *kafkala.CompaniesProducer) *Server {
	server := &Server{
		DB:     db,
		Logger: logger,
		Kafka:  kafka,
	}
	return server
}
