package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type Tests struct {
	name          string
	server        *httptest.Server
	response      *IPInfo
	expectedError error
}

func Test_consumeAPI(t *testing.T) {
	tests := []Tests{
		{
			name: "basic-request",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{ip_addres: "2a03:2880:f10d:83:face:b00c:0:25de",region: "California",country: "United States"}`))
			})),
			response: &IPInfo{
				Ip:      "2a03:2880:f10d:83:face:b00c:0:25de",
				State:   "California",
				Country: "United States",
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.server.Close()

			ip, country, state, err := consumeAPI(test.server.URL)
			resp := ip + country + state

			if !reflect.DeepEqual(resp, test.response) {
				t.Errorf("FAILED: expected %v. got %v,%v,%v", test.response, ip, country, state)
			}
			if !errors.Is(err, test.expectedError) {
				t.Errorf("Expected error FAILED: expected %v. got %v", test.expectedError, err)
			}
		})
	}
}
