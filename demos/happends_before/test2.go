package happends_before

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"
)

func Run() {
	// test1()
	// test2()
	// test3()
	test5()
}

type Counter struct {
	count int
}

func (c *Counter) Inc() {
	fmt.Printf("%p\n", c)
	c.count++
}

func test1() {
	c := Counter{}

	do := func() {
		for i := 0; i < 10; i++ {
			c.count++
		}
		fmt.Println("done")
	}

	go do()
	go do()
	time.Sleep(3 * time.Second)
	fmt.Println(c.count)
}
func test2() {
	c := &Counter{}
	do := func() {
		for i := 0; i < 10; i++ {
			c.Inc()
		}
		fmt.Println("done")
	}

	go do()
	go do()
	time.Sleep(3 * time.Second)
	fmt.Println(c.count)
}
func test3() {
	c := &Counter{}
	reflectValue := reflect.TypeOf(c)
	fmt.Println(reflectValue.Method(0).Type)
}

type MyError struct {

}

func (e *MyError) Error() string {
	return "bad"
}

func returnsError() error {
    var p *MyError = nil
    return p
}

func test4() {
	err := returnsError()
	fmt.Println(err)
	fmt.Println(err == nil)
}

func test5() {
	data := make(map[string]int)
	data["a"] = 1
	data["b"] = 2
	test6(data)
	fmt.Println(data)
}

func test6(d map[string]int) {
	d["c"] = 3

	var buf *bytes.Buffer
	io.Copy(buf, os.Stdin)
}