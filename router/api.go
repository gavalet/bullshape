package router

import (
	"bullshape/confs"
	"fmt"
	"net/http"
	"strconv"
)

func Run() {

	router := LoadRoutes()
	server := &http.Server{
		Addr:    ":" + strconv.FormatInt(confs.Conf.ServerPort, 10),
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Println("Server error:", err)
	}
}
