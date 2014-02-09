package main

import "fmt"
import "github.com/justingrammens/gdev/go-examples/mydate"
import "time"

func main() {
	now := time.Now()
	birthday := mydate.Birthdate()
	diff := birthday.YearDay() - now.YearDay()
	fmt.Printf("There are only %d days to my birthday!\n", diff )
}