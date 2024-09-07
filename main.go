/*
Copyright © 2024 Allan Muiruri <mwaranguallan345@gmail.com>
*/
package main

import (
	"fmt"
	"time"

	"github.com/AllanM007/snap/cmd"
)

func main() {
	loc, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		fmt.Println(err)
	}
	// handle err
	time.Local = loc // -> this is setting the global timezone

	cmd.Execute()
}
