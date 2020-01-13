package fatality

import "log"

// Panic panics if err != nil
func Panic(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Log prints message if err != nil
func Log(err error) {
	if err != nil {
		log.Println(err)
	}
}

// PanicMsg prints message and panics if err != nil
func PanicMsg(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %s\n", message, err)
	}
}

// PanicMsg prints message if err != nil
func LogMsg(message string, err error) {
	if err != nil {
		log.Printf("%s: %s\n", message, err)
	}
}
