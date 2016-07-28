package main

import (
	"strings"
	"net/http"
	"log"
	"net/http/httptest"
)

func httpGet(user, url string) (*http.Response, error){

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error when creating PUT request %d.", err)
		return nil, err
	}
	req.Header.Add(HEADER, user)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error when creating PUT request %d.", err)
		return nil, err
	}
	return resp, err
}

func httpPost(user, url string, values Stringifiable) (*http.Response, error){

	body := strings.NewReader(values.ToJsonString())
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Printf("Error when creating POST request %d.", err)
		return nil, err
	}
	req.Header.Add(HEADER, user)
	resp, err := http.DefaultClient.Do(req)

	return resp, err
}


func httpDelete(user, url string) (*http.Response, error){

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Printf("Error when creating DELETE request %d.", err)
		return nil, err
	}
	req.Header.Add(HEADER, user)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error when creating DELETE request %d.", err)
		return nil, err
	}
	return resp, err
}

func getServer(service Service) *httptest.Server {

	router := NewRouter()
	service.ConfigureRouter(router)
	server := httptest.NewServer(router)

	return server
}

func isHTTPStatus(httpStatus int, res *http.Response, err error ) bool {
	return !( (err != nil) || (res.StatusCode != httpStatus) )
}
