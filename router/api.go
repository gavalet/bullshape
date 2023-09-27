package router

import (
	"bullshape/confs"
	"bullshape/utils/logger"
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func (s *Server) Run() {

	lg := logger.GetLogger()
	//lg.SetExtra("Req:", "Main things")
	lg.Slogger.Info("Eimai sto Run!!!")
	s.Logger.Info("lalalala")
	log := s.Logger
	router := s.loadRoutes()
	server := &http.Server{
		Addr:    ":" + strconv.FormatInt(confs.Conf.ServerPort, 10),
		Handler: router,
	}
	go s.handleSystemSignals(server)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Info("Server error:", err)
	}
}

func (s *Server) handleSystemSignals(server *http.Server) {
	log := s.Logger
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exit := make(chan int)
	go func() {

		s := <-signalChanel
		switch s {
		case syscall.SIGHUP:
			log.Info("SIGHUP received.")
			exit <- 0
		case syscall.SIGINT:
			log.Info("SIGINT received.")
			exit <- 0
		case syscall.SIGTERM:
			log.Info("SIGTERM received, put virtupian into maintenance mode before proceeding.")
			exit <- 0
		case syscall.SIGQUIT:
			log.Info("SIGQUIT received.")
			exit <- 0
		default:
			log.Info("Unknown signal received.")
			exit <- 1
		}
	}()
	exitCode := <-exit
	s.shutdownServer(server)
	os.Exit(exitCode)
}

func (s *Server) shutdownServer(server *http.Server) {
	log := s.Logger
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Stop accepting new connections.
	server.SetKeepAlivesEnabled(false)
	err := server.Shutdown(ctx)
	if err != nil {
		log.Error("Error during server shutdown: %v", err)
	}
}
