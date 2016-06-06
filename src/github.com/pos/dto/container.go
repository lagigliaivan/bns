package dto

type Container interface {

	Add(value interface{})
	IsEmpty() bool
}