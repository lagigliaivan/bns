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
	"errors"
)

const (
	DB_TYPE = "DB_TYPE"
	DYNAMO_LOCAL_DB = "LOCALDB"
	MEM_DB = "MEMDB"

	DB_URL = "DB_URL"

	ANDROID_APP_ID = "ANDROID_APP_ID"
	FACEBOOK_ACCESS_TOKEN = "FACEBOOK_ACCESS_TOKEN"

	GOOGLE = "GOOGLE"
	FACEBOOK = "FACEBOOK"
	MAIL = "MAIL"
)

var androidAppClientID = os.Getenv(ANDROID_APP_ID)


var dbs map[string] func() *DB = make(map[string] func() *DB)
var getUser map[string] func(string) (*string, error) = make(map[string] func(string) (*string, error))

const(
	TOKEN_HEADER = "Authorization"
	TOKEN_TYPE_HEADER = "TokenType"

	USER_ID = "User-ID"
	GOOGLE_TOKEN_INFO_URL = "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token="
	FACEBOOK_TOKEN_INFO_URL = "https://graph.facebook.com/debug_token?input_token="
)

func init() {
	dbs[DYNAMO_LOCAL_DB] = getLocalDB
	dbs[MEM_DB] = getMemDB

	getUser[GOOGLE] = getGoogleUserId
	getUser[FACEBOOK] = getFacebookUserId
	getUser[MAIL] = getMailUserId
}

func main() {

	dbType := os.Getenv(DB_TYPE)

	createDb := dbs[dbType]

	if createDb == nil {
		createDb = getDynamoDB
	}

	purchasesService := NewPurchaseService(*createDb())

	preRouter := NewPreRouter(getValidUserId).AddService(purchasesService)

	log.Fatal(http.ListenAndServe(":8080", preRouter))
}

func getMemDB() *DB {
	log.Print("Using MEMDB")
	db := NewMemDb()
	cast_db := DB(*db)
	return &cast_db
}

func getDynamoDB() *DB {
	log.Print("Using DYNAMODB")
	db, _ := NewDynamoDB("", "us-west-2")
	cast_db := DB(*db)
	return &cast_db
}

func getLocalDB() *DB{

	url := os.Getenv(DB_URL)

	if strings.Compare(url, "") == 0 {
		log.Printf("URL_DB env variable is missing, so using localhost:8080")
		url = "http://localhost:8000"
	}

	log.Printf("Using LOCALDB url: %s", url)

	db, _ := NewDynamoDB(url, "us-west-2")
	cast_db := DB(*db)
	return &cast_db

}

func getValidUserId(request *http.Request) bool{

	userToken := request.Header.Get(TOKEN_HEADER)
	userTokenType := request.Header.Get(TOKEN_TYPE_HEADER)

	if len(userToken) == 0 {

		log.Printf("Security Header needs to be present")
		return false
	}

	getUserId := getUser[userTokenType]

	if getUserId == nil {
		log.Printf("Either %s http header is not present or invalid. Using %s by default", TOKEN_TYPE_HEADER, MAIL)
		getUserId = getUser[MAIL]
	}

	userId , err := getUserId(userToken)

	if err != nil {
		log.Printf("An error occurred while validating user [%s] %s ", request.RemoteAddr, userToken)
		return false
	}

	sha := sha1.New()
	io.WriteString(sha, *userId)

	//For the moment there is not a more practical way to use, later,
	//the user email as ID in DB. So, what I'm doing is to add it in a http header :(
	log.Printf("[%s]Token provided was successfully validated.", userTokenType)
	request.Header.Add(USER_ID,  fmt.Sprintf("%x", sha.Sum(nil)))
	return true
}

var users [2]string =  [...]string{"ce9e901ac338e02ce5ba652816718e6e54fed4cf", "a9aca56773bf9aa5af6955e49ea89a4eac762b8d"}

func getMailUserId(userToken string) (*string, error){

	if strings.Compare(userToken, users[0]) == 0 || strings.Compare(userToken, users[1]) == 0{
		return &userToken, nil
	}

	return nil, errors.New("It's not a valid user")
}

func getGoogleUserId(userToken string) (*string, error) {

	if strings.Compare(userToken, "") == 0{
		return nil, errors.New("user token cannot be emtpy, sorry.")
	}
	
	res, err := http.Get(GOOGLE_TOKEN_INFO_URL + userToken)

	if isHTTPStatus(http.StatusBadRequest, res, err){
		log.Printf("Error while validating user token")
		return nil, err;
	}

	googleDto := new(GoogleSignInDto)
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, googleDto); err != nil {
		log.Printf("Error while parsing google user token validation response: %s", err)
		return nil, err
	}

	userEmail := googleDto.Email

	if googleDto.Email == "" || googleDto.Aud != androidAppClientID {
		log.Printf("Either email is empty or aud does not match appclientid")
		return nil, errors.New("Either email is empty or aud does not match appclientid")
	}

	return &userEmail, nil

}


func getFacebookUserId(userToken string) (*string, error) {

	accessToken := os.Getenv(FACEBOOK_ACCESS_TOKEN)

	if strings.Compare(accessToken, "") == 0 {
		return nil, errors.New("FACEBOOK_ACCESS_TOKEN was not prived as env variable, so facebook validation will not be available")
	}

	res, err := http.Get(FACEBOOK_TOKEN_INFO_URL + userToken + "&access_token=" + accessToken)

	if isHTTPStatus(http.StatusBadRequest, res, err){
		log.Printf("Error while validating user token")
		return nil, err;
	}

	facebookDto := new(FacebookSignInDto)
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, facebookDto); err != nil {
		log.Printf("Error while parsing google user token validation response: %s", err)
		return nil, err
	}

	if facebookDto.Data.Is_valid != true {
		log.Printf("Facebook user token is not valid")
		return nil, errors.New("Facebook user token is not valid")
	}

	userId := facebookDto.User_id

	return &userId, nil
}
