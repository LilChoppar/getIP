package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type IPInfo struct {
	Ip      string `json:"ip_address"`
	State   string `json:"region"`
	Country string `json:"country"`
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	userInput := req.QueryStringParameters["address"]
	hostName, err := validateIP(userInput)
	if err != nil {
		return nil, err
	}

	ip, country, state, err := consumeAPI(hostName)
	if err != nil {
		return nil, err
	}

	info := IPInfo{Ip: ip, Country: country, State: state}
	ipBytes, _ := json.MarshalIndent(info, "", " ")

	res := &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(ipBytes),
	}

	return res, nil
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
func validateIP(ip string) (string, error) {
	var convertedIP string
	if net.ParseIP(ip) == nil {
		convert, err := net.LookupIP(ip)
		if err != nil {
			return "", err
		}
		//Will only loop once. only need one IP doesn't matter if it's IPv6/4.
		for i := 1; i < len(convert); i++ {
			convertedIP = convert[i].String()
		}

		return convertedIP, nil
	}
	return ip, nil
}
