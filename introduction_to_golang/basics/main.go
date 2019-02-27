// Every GO program is made up of packages. Programs start running in package main.
package main

// We can write imports as
// import "fmt"
// import "time"
// import "github.com/rmarasigan/introduction_to_golang/basics/vat"
// or
// import (
// 	"fmt"
// 	"time"

// 	"github.com/rmarasigan/introduction_to_golang/basics/vat"
// )
import (
	"fmt"
	"time"

	"github.com/rmarasigan/introduction_to_golang/basics/vat"
)

// A function can take zero or more arguments.
func main() {
	// for loop is a repetition control structure that will allow you
	// to efficiently write a loop that needs to execute a specific numer of times.
	// There are no parentheses surrounding the 3 components of the for
	// statement and the braces {} are always required
	for i := 0; i < 5; i++ {
		fmt.Printf("Line %d\n", i+1)
	}

	// A variable definition means to tell the compiler where and how much to
	// create the storage for the variable.
	// Example: var variable_name optional_data_type
	counter := 0
	fmt.Print("\n")
	for counter < 5 {
		counter++
		fmt.Println(counter)
	}

	// Counting from 0-4
	fmt.Println("\nCounter 0 - 4")
	counter = 0
	for {
		fmt.Printf("%d\n", counter)
		counter++

		if counter == 5 {
			break
		}

		time.Sleep(time.Millisecond * 2)
	}

	// Print Even numbers
	fmt.Println("\nEven Numbers")
	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			fmt.Printf("%d\n", i)
		}
	}

	// Printing up to 100
	fmt.Println("\nPrint up to 100")
	for i := 0; i < 100; i++ {
		go fmt.Print(i)
	}

	time.Sleep(time.Millisecond * 100)

	// Print the sum of Add function.
	fmt.Printf("\n")
	fmt.Print(Add(1, 100))

	// Print the Vat function.
	fmt.Println(vat.Vat(500))
}

// Add takes two parameters of type int.
// Type comes after the variable name.
func Add(x int, y int) int {
	fmt.Print("\n")
	return x + y
}
