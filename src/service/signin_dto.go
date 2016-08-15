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
"aud": "771875379-qbdvrqrjdii0gims9upnuqcqrf6753ei.apps.googleusercontent.com",
"sub": "116430566681024790786",
"email_verified": "true",
"azp": "771875379-9jlivedf892m5grfq3vec95qdc1phdaa.apps.googleusercontent.com",
"email": "lagigliaivan@gmail.com",
"iat": "1470880462",
"exp": "1470884062",
"name": "Ivan Lagiglia",
"given_name": "Ivan",
"family_name": "Lagiglia",
"locale": "es-419",
"alg": "RS256",
"kid": "104625465f6d4c7d214e3326913c5a5e4505699c"
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
"app_id": "1644514075861990",
"application": "AhorraYa",
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
"user_id": "115690035545503"
}
}*/
