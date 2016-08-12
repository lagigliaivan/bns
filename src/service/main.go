/**
This Package starts up a server which has the following APIs:
GET /purchases
**/
package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"encoding/json"
	"io/ioutil"
	"crypto/sha1"
	"io"
	"fmt"
)

const (
	BNS_DB = "BNS_DB"
	LOCALDB = "LOCALDB"
	MEMDB = "MEMDB"
)

var androidAppClientID = os.Getenv("ANDROID_APP_ID")

func main() {

	dbType := os.Getenv("BNS_DB")


	var db DB

	if strings.Compare(dbType,LOCALDB) == 0 {

		log.Print("Using LOCALDB")
		db, _ = NewDynamoDB("http://localhost:8000", "us-west-2")

	} else if strings.Compare(dbType, MEMDB) == 0 {

		db = NewMemDb()
		log.Print("Using MEMDB")
	} else {

		db, _ = NewDynamoDB("", "us-west-2")
		log.Print("Using DYNAMODB")
	}

	purchasesService := NewPurchaseService(db)
	preRouter := NewPreRouter(isAValidUser).AddService(purchasesService)

	log.Fatal(http.ListenAndServe(":8080", preRouter))
}

const(
	HEADER = "Authorization"
	USER_ID = "User-ID"
	GOOGLE_TOKEN_INFO_URL = "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=";

)

func isAValidUser(request *http.Request) bool{

	userToken := request.Header.Get(HEADER)

	if len(userToken) == 0 {

		log.Printf("Security Header needs to be present")
		return false

	}else {
		log.Printf("Validating token agains " + GOOGLE_TOKEN_INFO_URL + userToken)
		res, err := http.Get(GOOGLE_TOKEN_INFO_URL + userToken)

		if isHTTPStatus(http.StatusBadRequest, res, err){
			log.Printf("Error while validating user token")
			return false;
		}

		googleDto := new(GoogleSignInDto)
		body, _ := ioutil.ReadAll(res.Body)
		if err := json.Unmarshal(body, googleDto); err != nil {
			log.Printf("Error while parsing google user token validation response: %s", err)
			return false
		}

		userEmail := googleDto.Email

		if googleDto.Email == "" || googleDto.Azp != androidAppClientID {
			return false
		}

		sha := sha1.New()
		io.WriteString(sha, userEmail)

		request.Header.Add(USER_ID,  fmt.Sprintf("%x", sha.Sum(nil)))

		log.Printf(fmt.Sprintf("%x", sha.Sum(nil)))

		return true;
	}

	return false
}