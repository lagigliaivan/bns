package main

import (
	"testing"
	"golang.org/x/oauth2"
	"log"
	"golang.org/x/oauth2/google"
)



func Test_Gmail_Outh(t *testing.T) {



	type Token struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		IdToken     string `json:"id_token"`
	}


	code := ""

	conf := &oauth2.Config{
		ClientID:     "771875379-9jlivedf892m5grfq3vec95qdc1phdaa.apps.googleusercontent.com",
		ClientSecret: "86:37:98:89:A5:36:57:09:BB:9D:72:0A:78:2E:1D:F3:16:12:F3:8D",
		RedirectURL:  "http://localhost",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}


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
