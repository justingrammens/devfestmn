package mydate

import "time"

func Birthdate() time.Time {
	d := time.Date(2014, time.May, 6, 0, 0, 0, 0, time.UTC)	// calling the time.Date function
	return d
}