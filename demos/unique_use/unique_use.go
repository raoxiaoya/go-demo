package unique_use

import (
	"fmt"
	"strings"
	"unique"
	"unsafe"
)

func Run() {
	t3()
}

func t1() {
	h1 := unique.Make("hello")
	h2 := unique.Make("hello")
	h3 := unique.Make("world")
	fmt.Println(h1 == h2) // true
	fmt.Println(h1 == h3) // false
}

func t2() {
	s0 := "raoxiaoya"
	fmt.Println(unsafe.StringData(s0)) // 0xc3bf0a

	s1 := "raoxiaoya"
	fmt.Println(unsafe.StringData(s1)) // 0xc3bf0a

	s2 := s1
	fmt.Println(unsafe.StringData(s2)) // 0xc3bf0a

	s0 = "xiao"
	fmt.Println(unsafe.StringData(s0)) // 0xc3b854

	ss := "rao"
	fmt.Println(unsafe.StringData(ss)) // 0xc3b7e5

	s2 = s1[:3]
	fmt.Println(unsafe.StringData(s1)) // 0xc3bf0a
	fmt.Println(unsafe.StringData(s2)) // 0xc3bf0a

	s3 := strings.Clone(s1)
	fmt.Println(unsafe.StringData(s1)) // 0xc3bf0a
	fmt.Println(unsafe.StringData(s3)) // 0xc00008c0d0
}

type Person struct {
	Name string
	Age int
}
func t3() {
	p1 := Person{"rao", 12}
	p2 := Person{"rao", 12}
	fmt.Println(p1 == p2) // true

	h1 := unique.Make(p1)
	h2 := unique.Make(p2)
	fmt.Println(h1 == h2) // true
}