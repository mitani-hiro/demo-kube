package app

import "fmt"

func Hoge(v []byte) {
	fmt.Printf("received message: %s\n", string(v))
}
