package horcrux

import "fmt"

func logf(format string, v ...interface{}) {
	fmt.Printf(format, v...)

}
func logln(v ...interface{}) {
	fmt.Println(v...)
}
