package testfuck

import (
	"errors"
	"fmt"

	"../Utils"
)

func Div(num1, num2 int) (int, error) {
	if num2 == 0 {
		return 0, errors.New("qa")
	}
	return num1 / num2, nil
}
func main() {
	res := Utils.IsfileExists("../view/error.html")
	if res {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
}
