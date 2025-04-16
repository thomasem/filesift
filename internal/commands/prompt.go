package commands

import (
	"fmt"
)

func confirm() bool {
	var response string
	fmt.Print("Enter [y/N]: ")
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}
