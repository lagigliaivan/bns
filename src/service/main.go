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
	DB_TYPE = "DB_TYPE"
	DYNAMOLOCAL = "LOCALDB"
	MEMDB = "MEMDB"
	ANDROID_APP_ID = "ANDROID_APP_ID"
	USER_VALIDATION = "USER_VALIDATION"
	NON_G_USER = "NON_G_USER"
	DB_URL = "DB_URL"

)

var androidAppClientID = os.Getenv(ANDROID_APP_ID)


var dbs map[string] func() *DB = make(map[string] func() *DB)

const(
	HEADER = "Authorization"
	USER_ID = "User-ID"
	GOOGLE_TOKEN_INFO_URL = "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=";
)

func init() {
	dbs[DYNAMOLOCAL] = getLocalDB
	dbs[MEMDB] = getMemDB
}

func main() {

	dbType := os.Getenv(DB_TYPE)

	createDb := dbs[dbType]

	if createDb == nil {
		createDb = getDynamoDB
	}

	purchasesService := NewPurchaseService(*createDb())

	validateUser := isAValidGoogleUser

	if strings.Compare(os.Getenv(USER_VALIDATION), NON_G_USER) == 0 {
		log.Printf("Using simple user validation (NonGoogle sign in)")
		log.Printf("You can use the following users to send them out in the Authorization header %s %s", users[0], users[1] )
		validateUser = nonGoogleUserValidation
	}

	preRouter := NewPreRouter(validateUser).AddService(purchasesService)

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

func isAValidGoogleUser(request *http.Request) bool{

	userToken := request.Header.Get(HEADER)

	if len(userToken) == 0 {

		log.Printf("Security Header needs to be present")
		return false

	}else {
		log.Printf("Validating token")
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

		if googleDto.Email == "" || googleDto.Aud != androidAppClientID {
			log.Printf("Either email is empty or aud does not match appclientid")
			return false
		}

		sha := sha1.New()
		io.WriteString(sha, userEmail)

		//For the moment there is not a more practical way to use, later,
		//the user email as ID in DB. So, what I'm doing is to add it in a http header :(
		request.Header.Add(USER_ID,  fmt.Sprintf("%x", sha.Sum(nil)))

		return true;
	}

	return false
}

var users [2]string =  [...]string{"d563af2d08b4f672a11b3ed9065b7890a6412cab", "107cbb20a1d1e156beac1a9a7a331b36321300d4"}

func nonGoogleUserValidation(request *http.Request) bool{

	userToken := request.Header.Get(HEADER)

	if strings.Compare(userToken, users[0]) == 0 || strings.Compare(userToken, users[1]) == 0{
		request.Header.Add(USER_ID, userToken)
		return true;
	}else {
		return false;
	}
}
