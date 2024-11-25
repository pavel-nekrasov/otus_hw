package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

type Response struct {
	Data  interface{} `json:"data"`
	Error *Error      `json:"error"`
}

type HelloService struct {
	logger Logger
}

func NewHelloService(logger Logger) *HelloService {
	return &HelloService{logger: logger}
}

func (s *HelloService) GetHello(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	if r.Method != http.MethodGet {
		resp.Error = &Error{Message: fmt.Sprintf("method %s not supported on uri %s", r.Method, r.URL.Path)}
		w.WriteHeader(http.StatusMethodNotAllowed)
		s.WriteResponse(w, resp)
		return
	}

	resp.Data = "Hello world"

	w.WriteHeader(http.StatusOK)
	s.WriteResponse(w, resp)
}

func (s *HelloService) WriteResponse(w http.ResponseWriter, resp *Response) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
	_, err = w.Write(resBuf)
	if err != nil {
		s.logger.Error(fmt.Sprintf("response marshal error: %s", err))
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
