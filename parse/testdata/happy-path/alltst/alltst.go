package alltst

import "log"

// Alltst is logging its execution.
func Alltst() {
	log.Printf("INFO - alltst.Alltst")
	helper()
}

func helper() {
	log.Printf("INFO - alltst.helper")
}
