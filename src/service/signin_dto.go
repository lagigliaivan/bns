package main



type GoogleSignInDto struct {
	Iss string
	Aud string
	Sub string
	Email_verified string
	Azp string
	Email string
	Iat string
	Exp string
	Name string
	Given_name string
	Family_name string
	Locale string
	Alg string
	Kid string
}



//Example
/*
{
"iss": "https://accounts.google.com",
"aud": "771875379-fadfadafaf.apps.googleusercontent.com",
"sub": "116430566681024790786",
"email_verified": "true",
"azp": "771875379-fadfadfa.apps.googleusercontent.com",
"email": "fsdfaf@gmail.com",
"iat": "1470880462",
"exp": "1470884062",
"name": "Ivan Lagiglia",
"given_name": "Ivan",
"family_name": "Lagiglia",
"locale": "es-419",
"alg": "RS256",
"kid": "fadfadfadfadfaf"
}
*/
type FacebookDataDto struct {
	App_id string 		`json:"app_id"`
	Application string	`json:"application"`
	Expires_at int64	`json:"expires_at"`
	Is_valid bool		`json:"is_valid"`
	Issued_at int64		`json:"issued_at"`

}

type FacebookSignInDto struct {
	Data FacebookDataDto  `json:"data"`
	User_id string	      `json:"user_id"`
	Scopes [] string      `json:"scopes"`

}

/*
{
"data": {
"app_id": "fadfafd",
"application": "fadfafd",
"expires_at": 1476374820,
"is_valid": true,
"issued_at": 1471190820,
"metadata": {
"auth_type": "rerequest"
},
"scopes": [
"email",
"public_profile"
],
"user_id": "23232323"
}
}*/
