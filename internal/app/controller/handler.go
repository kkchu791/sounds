package controller

import (
	"github.com/kkchu791/sounds/internal/domain/model"
)

type Server struct {
	Controller *model.Controller
}

func NewServer() *Server {
	return &Server{
		Controller: model.NewController(),
	}
}
