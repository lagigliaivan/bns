package main

type Stringifiable interface  {
	ToJsonString() string

}