package main

import (
	"bufio"
	"fmt"
	"go-demo/leetcode"
	"os"
	"strings"
)

func main() {
	// go_call_dll.Run()

	// call_interface.Run()

	// dispatcher.Run()

	// happends_before.Run()

	// var a int32
	// atomic.LoadInt32(&a)
	// println(a)

	// log.Println()

	// zero_copy.Run()

	// sort.Run()

	// iterator.Run()

	// unique_use.Run()

	// priority_select.Run()

	// go_eval.Run()

	// go_sqlite.Run2()

	// exercise.Run()

	// child_process_stdio.Run()

	leetcode.Run()
}

func data() {
	source, _ := os.Open("ttt.txt")
	defer source.Close()

	dst, _ := os.Create("ttt_result.txt")
	defer dst.Close()

	result := make(map[string][]string)
	result["数据处理"] = make([]string, 0)

	cate := []string{"读书积分", "微展评", "共读活动", "健康打卡", "健康之屋", "镜像"}
	reader := bufio.NewReader(source)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		var get bool
		line = line[:len(line)-1]
		for _, c := range cate {
			var contain bool
			if c == "镜像" && (strings.Contains(line, "resmake") || strings.Contains(line, "tarwork") || strings.Contains(line, "downwork") || strings.Contains(line, "resmirror") || strings.Contains(line, "镜像")) {
				contain = true
			} else if c == "镜像" && (strings.Contains(line, "健康之屋") || strings.Contains(line, "新活动平台")) {
				contain = true
			} else if strings.Contains(line, c) {
				contain = true
			}
			if contain {
				if sli, ok := result[c]; ok {
					result[c] = append(sli, line)
				} else {
					result[c] = []string{line}
				}
				get = true
				break
			}
		}
		if !get {
			result["数据处理"] = append(result["数据处理"], line)
		}
	}
	for cat, list := range result {
		dst.WriteString("----------------------" + cat + "----------------------\n")
		for _, vv := range list {
			dst.WriteString(vv + "\n")
		}
		dst.WriteString("\n")
	}
}

func cal() {
	sum := 60.0
	add := 10.0

	for range 5 {
		sum = sum + sum * 0.1 + add
	}

	fmt.Println(sum) // 157.6816
}