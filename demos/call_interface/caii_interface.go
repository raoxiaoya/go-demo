package call_interface

import "fmt"

func Run() {
	dog := Dog{Name: "aaa"}
	cat := Cat{Name: "bbb"}

	AnimalEat(dog)
	AnimalEat(cat)
}

type Dog struct {
	Name string
}

func (d Dog) Eat() {
	fmt.Println("Dog Eat", d.Name)
}

type Cat struct {
	Name string
}

func (c Cat) Eat() {
	fmt.Println("Cat Eat", c.Name)
}

type Animal interface {
	Eat()
}

func AnimalEat(a Animal) {
	a.Eat()
}