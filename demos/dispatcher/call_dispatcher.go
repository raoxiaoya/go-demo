package dispatcher

import (
	"log"
	"strconv"
	"time"
)

func Run() {
	dispatcher := New(10, 100)
	dispatcher.Run()

	for i := 1; i <= 100; i++ {
		dispatcher.Dispatch(Dog{
			Name: "dog" + strconv.Itoa(i),
		})
	}

	time.Sleep(time.Second * 5)
}

type Dog struct {
	Name string
}

func (d Dog) Do() error {
	log.Println("Dog Eat", d.Name)
	return nil
}
