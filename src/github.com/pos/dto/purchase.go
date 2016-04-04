package dto

import "time"

type Purchase struct {
	Time time.Time //DateTime when this purchase was acquired.
	Items ItemsContainer
}

