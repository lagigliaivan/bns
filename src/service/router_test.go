package main

import (
	"testing"
	"golang.org/x/oauth2"
	"log"
	"golang.org/x/oauth2/google"
	"net/http"
)



func Test_Gmail_Outh(t *testing.T) {

	code := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjZmYWE0ZTllYzMwMDMwNzg0Yjg5NDI2MDZmYjYxNzYyYWRhOTcyNTMifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhdWQiOiI3NzE4NzUzNzktcWJkdnJxcmpkaWkwZ2ltczl1cG51cWNxcmY2NzUzZWkuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMTY0MzA1NjY2ODEwMjQ3OTA3ODYiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXpwIjoiNzcxODc1Mzc5LTlqbGl2ZWRmODkybTVncmZxM3ZlYzk1cWRjMXBoZGFhLmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiZW1haWwiOiJsYWdpZ2xpYWl2YW5AZ21haWwuY29tIiwiaWF0IjoxNDcwNzk3ODE0LCJleHAiOjE0NzA4MDE0MTQsIm5hbWUiOiJJdmFuIExhZ2lnbGlhIiwiZ2l2ZW5fbmFtZSI6Ikl2YW4iLCJmYW1pbHlfbmFtZSI6IkxhZ2lnbGlhIiwibG9jYWxlIjoiZXMtNDE5In0.hlBo1uvRErEm85Ca1b7LywblFYZBrgSnY_RMrPnQwjyCzq-NWesNidxOKP-YiAYD1fwcXr5nCqpEVLriKuNfpW8N030JBanGS2JtEn43pdXMfEXlOedhhkA8Glvxpjacyvqty6vEvKvumi6ZTvtZs96lOWE4nwui8J7XEHEUG_R1cCahzaMwJeWxK7ygbEE2msZe0liXgLlqcqtcvO4jOzu1T--x3zY9L9hEDN5Mft4p7r3udu9WqE6TJSvaYbNhZV8q0TdQ32SSQ7rFLDXtf0VvxmXANmZed_Nj9lBSvRKewhytEHBxXHlNOl9jDy-aQkZi3luAvi8XJs-Vtww_og"

	res, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=eyJhbGciOiJSUzI1NiIsImtpZCI6IjZmYWE0ZTllYzMwMDMwNzg0Yjg5NDI2MDZmYjYxNzYyYWRhOTcyNTMifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhdWQiOiI3NzE4NzUzNzktcWJkdnJxcmpkaWkwZ2ltczl1cG51cWNxcmY2NzUzZWkuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMTY0MzA1NjY2ODEwMjQ3OTA3ODYiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXpwIjoiNzcxODc1Mzc5LTlqbGl2ZWRmODkybTVncmZxM3ZlYzk1cWRjMXBoZGFhLmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiZW1haWwiOiJsYWdpZ2xpYWl2YW5AZ21haWwuY29tIiwiaWF0IjoxNDcwNzk3ODE0LCJleHAiOjE0NzA4MDE0MTQsIm5hbWUiOiJJdmFuIExhZ2lnbGlhIiwiZ2l2ZW5fbmFtZSI6Ikl2YW4iLCJmYW1pbHlfbmFtZSI6IkxhZ2lnbGlhIiwibG9jYWxlIjoiZXMtNDE5In0.hlBo1uvRErEm85Ca1b7LywblFYZBrgSnY_RMrPnQwjyCzq-NWesNidxOKP-YiAYD1fwcXr5nCqpEVLriKuNfpW8N030JBanGS2JtEn43pdXMfEXlOedhhkA8Glvxpjacyvqty6vEvKvumi6ZTvtZs96lOWE4nwui8J7XEHEUG_R1cCahzaMwJeWxK7ygbEE2msZe0liXgLlqcqtcvO4jOzu1T--x3zY9L9hEDN5Mft4p7r3udu9WqE6TJSvaYbNhZV8q0TdQ32SSQ7rFLDXtf0VvxmXANmZed_Nj9lBSvRKewhytEHBxXHlNOl9jDy-aQkZi3luAvi8XJs-Vtww_og" + code)


	if err != nil {
		log.Printf("Err: %s", err)
		t.FailNow()
	}


	log.Printf("Resp: %s", res)



	return


	type Token struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		IdToken     string `json:"id_token"`
	}





	conf := &oauth2.Config{
		ClientID:     "312141132183-2f1n8ldrf76s197pp5uj73g9bd4b312e.apps.googleusercontent.com",
		ClientSecret: "Xa2uOEACex11bgOamkXkNAtm",
		RedirectURL:  "http://localhost",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},

		Endpoint: google.Endpoint,
	}

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	log.Printf("url %s", url)

	return

	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}

	log.Printf("token: %s", tok)



	id := tok.Extra("id_token").(string)


	userid, err:= decodeIdToken(id)

	/*client := conf.Client(oauth2.NoContext, tok)
	resp, err := client.Get("...")
*/
	/*conf := getJWTConfigFromJson()
	client := conf.Client(oauth2.NoContext)
	resp,err := client.Get("...")
*/
	/*if err != nil {
		log.Fatal("Error: %s" , err)
		t.FailNow()
	}
*/
	log.Printf("%s", userid)

}
