package sort

import (
	"slices"
	"testing"
)

// go test -v -run TestSort
// go test -v -bench=. -benchmem

func BenchmarkBubbleSort(b *testing.B) {
	BubbleSort(RandSeries(num))
}
func BenchmarkPickSort(b *testing.B) {
	PickSort(RandSeries(num))
}
func BenchmarkInsertSort(b *testing.B) {
	InsertSort(RandSeries(num))
}
func BenchmarkCombineSort(b *testing.B) {
	CombineSort(RandSeries(num))
}
func BenchmarkQuickSort1(b *testing.B) {
	QuickSort1(RandSeries(num))
}
func BenchmarkQuickSort2(b *testing.B) {
	series := RandSeries(num)
	QuickSort2(series, 0, len(series)-1)
}
func BenchmarkSlicesSort(b *testing.B) {
	slices.Sort(RandSeries(num))
}
func BenchmarkHeapSort(b *testing.B) {
	HeapSort(RandSeries(num))
}
func BenchmarkCountSort(b *testing.B) {
	CountSort(RandSeries(num))
}
func BenchmarkBucketSort(b *testing.B) {
	BucketSort(RandSeries(num))
}
