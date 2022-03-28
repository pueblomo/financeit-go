package conf

import "log"

func CheckErrorFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}