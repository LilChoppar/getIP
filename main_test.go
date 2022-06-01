package main

import (
	"fmt"
	"testing"

	"github.com/gorilla/mux"
)

func Test_getIP(t *testing.T) {
	router := mux.NewRouter()
	fmt.Println(router)
}
