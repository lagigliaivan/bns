package main



type GoogleSignInDto struct {
	Iss string
	Aud string
	Sub string
	Email_verified bool
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