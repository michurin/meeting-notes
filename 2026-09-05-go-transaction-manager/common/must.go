package common

import "log"

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Must2[T any](value T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return value
}
