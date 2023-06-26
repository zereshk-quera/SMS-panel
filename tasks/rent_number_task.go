package tasks

import (
	"fmt"
)

func RentNumberTask() TaskFunc {
	return func() {
		fmt.Println("Yoooooo")
	}
}
