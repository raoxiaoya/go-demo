package happends_before

// var A int
// var B int

// func f5() {
// 	A = B + 1
// 	B = 1
// }

// var a, b int

// func f3() {
// 	a = 1
// 	b = 2
// }

// func g() {
// 	print(b)
// 	print(a)
// }

// func f0() {
// 	go f()
// 	g()
// }

// func f1() {
// 	var cc int32
// 	atomic.LoadInt32(&cc)

// 	var dd atomic.Int32
// 	dd.Add(1)
// }

// var s string

// func f() {
// 	print(s)
// }

// func Run() {
// 	s = "hello, world"
// 	go f()

// 	// time.Sleep(3 * time.Second)
// }

// var a, b int

// func f() {
// 	a = 1
// 	b = 2
// }

// func g() {
// 	println(b)
// 	println(a)
// }

// func Run() {
// 	go f()
// 	g()
// 	// time.Sleep(2 * time.Second)
// }

// ----------------------

// var mu sync.Mutex
// var num int
// func Run() {
// 	mu.Lock()
// 	mu.Unlock()
// 	num++

// 	ab := mu
// 	ab.Lock()
// 	num++
// 	ab.Unlock()
// }

