package main

import "fmt"
import "github.com/justingrammens/gdev/go-examples/mydate"
//import "time" // note: uncomment this line

func main() {
	//now := time.Now() // note: try uncommenting this line
	sum, difference := mydate.AddandSub(3,4);
	fmt.Printf("Here's the sum: %d and the difference: %d\n", sum, difference )
}