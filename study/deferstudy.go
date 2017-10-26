package study

import (
	"log"
	"time"
)

func funcA() int { //5
	x := 5
	defer func() {
		x += 1
	}()
	return x
}
func funcB() (x int) {
	defer func() {
		x += 1
	}()
	return 5
}
func funcC() (y int) {
	x := 5
	defer func() {
		x += 1
	}()
	return x
}
func funcD() (x int) {
	defer func(x int) {
		x += 1
	}(x)
	return 5
}
func bigSlowwOperation() {
	f := teace("bigjajhdfhjdk")
	defer f()
	time.Sleep(5 * time.Second)

}
func teace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", msg)
	return func() {
		log.Printf("exit %s(%s)", msg, time.Since(start))
	}
}
