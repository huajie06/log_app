package recap

import "fmt"

type Database interface {
	Connect() error
	Query(query string) string
}

type MySql struct{}

func (m MySql) Connect() error {
	fmt.Println("Connected to MySQL")
	return nil
}

func (m MySql) Query(query string) string {
	return "MySQL result for: " + query
}

// another database
type SQL01 struct{}

func (m SQL01) Connect() error {
	fmt.Println("=========")
	for i := 1; i <= 4; i++ {
		fmt.Println("now i'm in a different db!")
	}
	fmt.Println("=========")
	return nil
}

func (m SQL01) Query(query string) string {
	return "you won't get any connection really for your query\n" + query
}

func Recap03() {
	var db Database

	// the benefit is that db can be multiple types
	db = MySql{}
	db.Connect()
	queryResult := db.Query("select * from xxx")
	fmt.Println(queryResult)

	db = SQL01{}
	db.Connect()
	queryResult = db.Query("what ???")
	fmt.Println(queryResult)

	// since i defined a method for the INTERFACE, it can work on the db
	fmt.Println("------")
	MakeAnimalSound(db)

}

// now define something for the interface
func MakeAnimalSound(db Database) {
	fmt.Println("do some dumb stuff for this database")
	fmt.Println("you know what? i don't even want to think about the method!!")

	fmt.Println("\n------------------ variadic parameters ------------------")
	variadicparms("haha", "wawa", "gaga")
}

func variadicparms(mystring ...any) {
	for _, param := range mystring {
		fmt.Println("param:", param)
	}

	fmt.Println("\n------------------ variadic parameters example 2 ------------------")
	fmt.Println(mysum(1, 2, 3, 4))
	fmt.Println(mysum(10, 20, 30))
	fmt.Println(mysum())
}

func mysum(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}
