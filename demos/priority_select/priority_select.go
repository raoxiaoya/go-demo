package priority_select

import (
	"fmt"
	"time"
)

func worker(ch1, ch2 <-chan int, stopCh chan struct{}) {
	for {
		select {
		case <-stopCh:
			return
		case job1 := <-ch1:
			fmt.Println(job1)
		default:
			select {
			case job2 := <-ch2:
				fmt.Println(job2)
			default:
			}
		}
	}
}

/*
select语句中实现优先级

此示例要求 ch1 要先于 ch2 被处理。
*/
func worker2(ch1, ch2 <-chan int) {
	for {
		select {
		case job1 := <-ch1:
			println(job1)
		case job2 := <-ch2:
		label:
			// for 循环保证 ch1 被消费完采取处理 ch2
			for {
				select {
				case job1 := <-ch1:
					println(job1) // 执行 job1
				default:
					break label
				}
			}
			println(job2) // 执行 job2
		}
	}
}

func worker4() {
	fmt.Println(1)
	goto label
	fmt.Println(2)
label:
	fmt.Println(3)
}

func worker5() {
label:
	fmt.Println(1)
	time.Sleep(time.Second)
	goto label
	fmt.Println(2)
	fmt.Println(3)

}

func worker6() {
	for x := 1; x < 10; x++ {
	label:
		for i := 1; i < 10; i++ {
			for j := 1; j < 10; j++ {
				sum := i + j
				fmt.Println(sum)
				if sum == 3 {
					break label
				}
			}
		}
		fmt.Println("out1")
	}
	fmt.Println("out2")
}

func worker7() {
label:
	for i := 1; i < 3; i++ {
		for j := 1; j < 3; j++ {
			sum := i + j
			fmt.Println(sum)
			if sum >= 3 {
				continue label
			}
			fmt.Println("after")
		}
	}
	fmt.Println("out1")
}

func worker8() {
	ch1 := make(chan int, 0)
label:
	for {
		select {
		case job1 := <-ch1:
			fmt.Println(job1)
		default:
			break label
		}
		fmt.Println(1)
		time.Sleep(time.Second)
	}
}

func Run() {
	worker8()
}
