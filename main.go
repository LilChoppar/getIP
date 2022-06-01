package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

type IPInfo struct {
	Ip      string `json:"ip_address"`
	State   string `json:"region"`
	Country string `json:"country"`
}

func main() {
	fmt.Println("Starting endpoint on port: 2525")
	startServer()
}

func startServer() {
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/v1/ip-info/{address}", getIP).Methods("GET")
	http.ListenAndServe(":2525", myRouter)
}

func getIP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "No-Store")

	vars := mux.Vars(r)
	userInput := vars["address"]

	hostName := validateIP(userInput)
	ip, country, state, err := consumeAPI(hostName)
	if err != nil {
		fmt.Fprint(w, "There was an issue consuming API: ", err)
	}
	fmt.Fprint(w, "Ip Address: ", ip, "\n", "Country: ", country, "\n", "State: ", state)

}

func consumeAPI(address string) (Ip, Country, State string, err error) {
	response, err := http.Get("https://ipgeolocation.abstractapi.com/v1/?api_key=0f73d32a308f43859d18747f922df76f&ip_address=" + address)
	if err != nil {
		return "", "", "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", "", "", err
	}

	var info IPInfo
	json.Unmarshal(body, &info)

	return info.Ip, info.Country, info.State, nil
}

//function that checks if user input is ip or hostname.
//if hostname, convert to ip using net.lookupip
func validateIP(ip string) string {
	var convertedIP string
	if net.ParseIP(ip) == nil {
		convert, _ := net.LookupIP(ip)
		//Will only loop once. only need one IP doesn't matter if it's IPv6/4.
		for i := 1; i < len(convert); i++ {
			convertedIP = convert[i].String()
		}

		return convertedIP
	}
	return ip
}
