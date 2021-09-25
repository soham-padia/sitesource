package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/solow-crypt/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//do noth

	default:
		t.Error(fmt.Sprintf("type is not chi.mux , typr is %T", v))
	}
}
