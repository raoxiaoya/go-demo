package sort

import (
	"container/heap"
	"fmt"
	"math/rand/v2"
	"slices"
	"time"
)

// 排序算法集合：https://github.com/MisterBooo/Article

func Run() {
	Sort()
	Sort0()
	Sort1()
	Sort2()
	Sort3()
}

var num = 10000

func RandSeries(n int) (series []int) {
	series = rand.Perm(n)
	return
}

// slices.Sort
func Sort() {
	series := RandSeries(50)
	// series := []int{6, 0, 1, 7, 9, 4, 3, 8, 2, 5}
	// fmt.Println(series)

	start := time.Now().UnixMicro()
	series = BucketSort(series)
	end := time.Now().UnixMicro()
	fmt.Println(series)
	fmt.Println(end - start)
}

func Sort0() {
	series := RandSeries(num)
	// series := []int{6, 0, 1, 7, 9, 4, 3, 8, 2, 5}
	// fmt.Println(series)

	start := time.Now().UnixMicro()
	slices.Sort(series)
	end := time.Now().UnixMicro()
	// fmt.Println(series)
	fmt.Println(end - start)
}

// 快速排序-基础版本
func Sort1() {
	// series := RandSeries(num)
	series := []int{6, 0, 1, 7, 9, 4, 3, 8, 2, 5}
	// fmt.Println(series)

	// start := time.Now().UnixMicro()
	series = QuickSort1(series)
	// end := time.Now().UnixMicro()
	fmt.Println(series)
	// fmt.Println(end - start)
}

// 快速排序-优化版本
func Sort2() {
	// series := RandSeries(num)
	series := []int{6, 0, 1, 7, 9, 4, 3, 8, 2, 5}
	// fmt.Println(series)

	// start := time.Now().UnixMicro()
	QuickSort2(series, 0, len(series)-1)
	// end := time.Now().UnixMicro()
	fmt.Println(series)
	// fmt.Println(end - start)
}

// 冒泡排序
func Sort3() {
	series := RandSeries(num)
	// series := []int{6, 0, 1, 7, 9, 4, 3, 8, 2, 5}
	// fmt.Println(series)

	start := time.Now().UnixMicro()
	BubbleSort(series)
	end := time.Now().UnixMicro()
	// fmt.Println(series)
	fmt.Println(end - start)
}

func QuickSort1(data []int) []int {
	if len(data) == 1 {
		return data
	}
	if len(data) == 2 {
		if data[0] > data[1] {
			data[0], data[1] = data[1], data[0]
		}
		return data
	}

	pivot := data[0]
	left := make([]int, 0)
	middle := make([]int, 0)
	right := make([]int, 0)
	for _, x := range data {
		if x > pivot {
			right = append(right, x)
		} else if x == pivot {
			middle = append(middle, x)
		} else {
			left = append(left, x)
		}
	}

	var rst []int
	if len(left) > 0 && len(right) > 0 {
		rst = slices.Concat(QuickSort1(left), middle, QuickSort1(right))
	} else if len(left) > 0 {
		rst = slices.Concat(QuickSort1(left), middle)
	} else if len(right) > 0 {
		rst = slices.Concat(middle, QuickSort1(right))
	}

	return rst
}

// https://www.geeksforgeeks.org/quick-sort-algorithm/
func QuickSort2(data []int, low, high int) {
	if low < high {
		pivot := data[high]
		i := low - 1
		for j := low; j <= high; j++ {
			if data[j] < pivot {
				i++
				data[i], data[j] = data[j], data[i]
			}
		}
		data[i+1], data[high] = data[high], data[i+1]

		pi := i + 1

		QuickSort2(data, low, pi-1)
		QuickSort2(data, pi+1, high)
	}
}

// 冒泡排序
// 疯狂的swap，效率低下
func BubbleSort(data []int) {
	length := len(data)
	for i := 0; i < length-1; i++ {
		for j := 0; j < length-i-1; j++ {
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
}

// 选择排序
//
// 首先在未排序序列中找到最小元素，存放到排序序列的起始位置。
// 再从剩余未排序元素中继续寻找最小（大）元素，然后放到已排序序列的末尾。
// 重复第二步，直到所有元素均排序完毕。
//
// 比冒泡排序减少了不必要的swap操作
func PickSort(data []int) {
	length := len(data)
	for i := 0; i < length-1; i++ {
		index := i
		for j := i + 1; j < length; j++ {
			if data[j] < data[index] {
				index = j
			}
		}
		if index != i {
			data[index], data[i] = data[i], data[index]
		}
	}
}

// 输入排序
//
// 将第一待排序序列第一个元素看做一个有序序列，把第二个元素到最后一个元素当成是未排序序列。
// 从头到尾依次扫描未排序序列，将扫描到的每个元素插入有序序列的适当位置。
//
// 从原理来看需要不停的对数组进行分割与合并，性能极差
func InsertSort(arr []int) []int {
	res := make([]int, 0)
	for _, v := range arr {
		if len(res) == 0 {
			res = append(res, v)
		} else {
			// 最左
			if v <= res[0] {
				res = append([]int{v}, res...)
				continue
			}
			// 最右
			if v > res[len(res)-1] {
				res = append(res, v)
				continue
			}
			// 中间
			for x, y := range res {
				if v < y {
					res = append(res[:x], append([]int{v}, res[x:]...)...)
					break
				}
			}
		}
	}
	return res
}

// 归并排序
// 将序列从中间分为两部分，递归拆分，直到序列长度为1。
// 将两个序列的元素逐个比较，将较小的放到新序列中。
func CombineSort(data []int) []int {
	length := len(data)
	if length <= 1 {
		return data
	}

	middle := length / 2
	// 递归拆分
	arr1 := CombineSort(data[:middle])
	arr2 := CombineSort(data[middle:])

	// 合并并排序，从小到大
	arr3 := make([]int, len(arr1)+len(arr2))
	var indexarr1, indexarr2, indexarr3, val int
	for indexarr1 < len(arr1) && indexarr2 < len(arr2) {
		if arr1[indexarr1] <= arr2[indexarr2] {
			val = arr1[indexarr1]
			indexarr1++
		} else {
			val = arr2[indexarr2]
			indexarr2++
		}
		arr3[indexarr3] = val
		indexarr3++
	}

	// 合并剩下的
	// arr1 和 arr2 至少有一个为空
	for _, v := range arr1[indexarr1:] {
		arr3[indexarr3] = v
		indexarr3++
	}
	for _, v := range arr2[indexarr2:] {
		arr3[indexarr3] = v
		indexarr3++
	}

	return arr3
}

// 堆排序
// golang中的container/heap包是最小堆（根节点最小），且是一颗完全二叉树。
func HeapSort(data []int) []int {
	hp := &myHeap{}
	for _, v := range data {
		heap.Push(hp, v)
	}
	heap.Init(hp)
	res := make([]int, hp.Len())
	for i := 0; hp.Len() > 0; i++ {
		res[i] = heap.Pop(hp).(int)
	}

	return res
}

// 计数排序
// 找到最大值和最小值，然后创建一个新的切片，其长度为 max-min+1，记录每个值出现的次数。
func CountSort(data []int) []int {
	maxValue := slices.Max(data)
	minValue := slices.Min(data)
	newSlice := make([]int, maxValue-minValue+1)
	for _, v := range data {
		newSlice[v-minValue]++
	}

	res := make([]int, len(data))
	var index int
	for k, v := range newSlice {
		if v > 0 {
			for i := 0; i < v; i++ {
				res[index] = k + minValue
				index++
			}
		}
	}

	return res
}

// 桶排序
// 设置固定数量的空桶。
// 把数据放到对应的桶中。
// 对每个不为空的桶中数据进行排序。
// 拼接不为空的桶中数据，得到结果
func BucketSort(data []int) []int {
	maxValue := slices.Max(data)
	minValue := slices.Min(data)
	bucketSize := 10                                  // 每个桶中有几个数
	bucketCount := (maxValue-minValue)/bucketSize + 1 // 桶的个数
	buckets := make([][]int, bucketCount)
	for _, v := range data {
		index := (v - minValue) / bucketSize
		if len(buckets[index]) == 0 {
			buckets[index] = make([]int, 0)
		}
		buckets[index] = append(buckets[index], v)
	}
	res := make([]int, len(data))
	var index int
	for _, bucket := range buckets {
		if len(bucket) <= 0 {
			continue
		}
		slices.Sort(bucket)
		for _, v := range bucket {
			res[index] = v
			index++
		}
	}

	return res
}
