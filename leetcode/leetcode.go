package leetcode

import (
	"fmt"
	"strconv"
)

// 两数之和
func twoSum(nums []int, target int) []int {
	total := make(map[int]int)
	for i, v := range nums {
		if idx, ok := total[target-v]; ok {
			return []int{idx, i}
		}
		total[v] = i
	}
	return nil
}

////////////////////////////////////////////////

// 两数相加
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	var l1s, l2s string
	for {
		if l1 != nil {
			l1s += strconv.Itoa(l1.Val)
			l1 = l1.Next
		} else {
			l1s += "0"
		}
		if l2 != nil {
			l2s += strconv.Itoa(l2.Val)
			l2 = l2.Next
		} else {
			l2s += "0"
		}

		if l1 == nil && l2 == nil {
			break
		}
	}

	// println("l1s", l1s)
	// println("l2s", l2s)

	// 无论是int还是int64都是有上限的，但是链表却是无限的，因此不能转换成int类型来存储
	n := 0
	c := len(l1s)
	sums := make([]int, c)
	for i, v := range l1s {
		v1, _ := strconv.Atoi(string(v))
		v2, _ := strconv.Atoi(string(l2s[i]))
		add := v1 + v2 + n
		if add >= 10 {
			n = 1
			add -= 10
		} else {
			n = 0
		}
		sums[c-i-1] = add
	}
	if n == 1 {
		sums = append([]int{1}, sums...)
	}

	// fmt.Println("sums", sums)

	var res *ListNode
	for _, v := range sums {
		if res == nil {
			res = &ListNode{Val: v}
		} else {
			res = &ListNode{Val: v, Next: res}
		}
	}

	return res
}

type ListNode struct {
	Val  int
	Next *ListNode
}

// 链表的构造与 vals 顺序相同
func buildNodes(vals []int) *ListNode {
	var res *ListNode
	for i := len(vals) - 1; i >= 0; i-- {
		if res == nil {
			res = &ListNode{Val: vals[i]}
		} else {
			res = &ListNode{Val: vals[i], Next: res}
		}
	}

	return res
}

func OutputNodes(in *ListNode) {
	for {
		if in == nil {
			break
		}
		println(in.Val)
		in = in.Next
	}
}

/**

// l1 := buildNodes([]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1})

// l1 := buildNodes([]int{9, 9, 9, 9, 9, 9, 9})
// l2 := buildNodes([]int{9, 9, 9, 9})

l1 := buildNodes([]int{2, 4, 3})
l2 := buildNodes([]int{5, 6, 4})

// l1 := buildNodes([]int{0})
// l2 := buildNodes([]int{0})

res := addTwoNumbers(l1, l2)

// OutputNodes(l1)
// OutputNodes(l2)
OutputNodes(res)


*/

////////////////////////////////////////////////

// 无重复字符的最长子串的长度
// abcabcbb --> abc
// bbbbb --> b
// pwwkew --> wke
// b --> b
// dvdf --> vdf
func lengthOfLongestSubstring(s string) int {
	bt := []byte(s)
	var max int
	arr := make(map[byte]int)
	for key, val := range bt {
		before, ok := arr[val]
		if ok {
			if max < len(arr) {
				max = len(arr)
			}
			for v, k := range arr {
				if k <= before {
					delete(arr, v)
				}
			}
		}
		arr[val] = key
	}
	if max < len(arr) {
		max = len(arr)
	}

	return max
}

func lengthOfLongestSubstring2(s string) int {
	var res int
	m := map[byte]bool{}
	for l, r := 0, 0; r < len(s); r++ {
		for l < r && m[s[r]] {
			m[s[l]] = false
			l++
		}
		m[s[r]] = true
		res = max(res, r-l+1)
	}
	return res
}

// //////////////////////////////////////////////
// 最长公共前缀
// ["flower","flow","flight"] --> "fl"
// ["dog","racecar","car"] --> ""
func longestCommonPrefix(strs []string) string {
	for i := 0; i < len(strs[0]); i++ {
		for j := 1; j < len(strs); j++ {
			if i == len(strs[j]) || strs[0][i] != strs[j][i] {
				return strs[0][:i]
			}
		}
	}
	return strs[0]
}

// //////////////////////////////////////////////
// 有效的括号
// () --> true
// ()[]{} --> true
// (] --> false
// ([]) --> true
// (([]){}) --> true
// (([)]{}) --> false
func isValid(s string) bool {
	stack := make([]byte, 0)
	for i := 0; i < len(s); i++ {
		if len(stack) == 0 && (s[i] == ')' || s[i] == ']' || s[i] == '}') {
			return false
		}
		if s[i] == '(' || s[i] == '[' || s[i] == '{' {
			stack = append(stack, s[i])
		} else if s[i] == ')' && stack[len(stack)-1] == '(' {
			stack = stack[:len(stack)-1]
		} else if s[i] == ']' && stack[len(stack)-1] == '[' {
			stack = stack[:len(stack)-1]
		} else if s[i] == '}' && stack[len(stack)-1] == '{' {
			stack = stack[:len(stack)-1]
		} else {
			return false
		}
	}
	if len(stack) != 0 {
		return false
	}
	return true
}

////////////////////////////////////////////////
// 合并两个有序链表
// [1,2,4] [1,3,4] --> [1,1,2,3,4,4]
// [] [] --> []
// [] [0] --> [0]
// [1,6,9] [2,3,4] --> [1,2,3,4,6,9]

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	var list3, node, tmp *ListNode
	for {
		// fmt.Println("list1.Val:", list1.Val, "list2.Val:", list2.Val)
		if list1 != nil && list2 != nil {
			if list1.Val <= list2.Val {
				tmp = &ListNode{Val: list1.Val}
				list1 = list1.Next
			} else {
				tmp = &ListNode{Val: list2.Val}
				list2 = list2.Next
			}
		} else if list2 != nil {
			tmp = &ListNode{Val: list2.Val}
			list2 = list2.Next
		} else if list1 != nil {
			tmp = &ListNode{Val: list1.Val}
			list1 = list1.Next
		}
		if node == nil {
			node = tmp
		} else {
			node.Next = tmp
			node = tmp
		}
		if list3 == nil {
			list3 = node
		}
		if list1 == nil && list2 == nil {
			return list3
		}
	}
}

/*

l1 := buildNodes([]int{1, 6, 9})
l2 := buildNodes([]int{2, 3, 4})
res := mergeTwoLists(l1, l2)
OutputNodes(res)

*/

////////////////////////////////////////////////

// 删除有序数组中的重复项
// nums 是一个升序的数组，过滤掉重复值，然后返回数组长度
func removeDuplicates(nums []int) int {
	newlist := nums
	k := 0
	for _, v := range nums {
		if k > 0 && newlist[k-1] == v {
			continue
		}
		newlist[k] = v
		k++
	}
	nums = newlist[:k]
	return k
}

// //////////////////////////////////////////////
// 删除元素
func removeElement(nums []int, val int) int {
	newlist := nums
	k := 0
	for _, v := range nums {
		if val != v {
			newlist[k] = v
			k++
		}
	}
	nums = newlist[:k]
	return k
}

////////////////////////////////////////////////

// 找出字符串中第一个匹配项的下标
// 输入：haystack = "sadbutsad", needle = "sad"
// 输出：0
// 解释："sad" 在下标 0 和 6 处匹配。
// 第一个匹配项的下标是 0 ，所以返回 0 。
func strStr(haystack string, needle string) int {
	position := -1
	searchLen := len(needle)
	if searchLen == 0 {
		return position
	}
	for i := range len(haystack) {
		sub := haystack[i:]
		if len(sub) >= searchLen && sub[:searchLen] == needle {
			return i
		}
	}

	return position
}

// //////////////////////////////////////////////
// 搜索插入位置
// 给定一个排序数组和一个目标值，在数组中找到目标值，并返回其索引。如果目标值不存在于数组中，返回它将会被按顺序插入的位置。
// nums 为 无重复元素 的 升序 排列数组
// 请必须使用时间复杂度为 O(log n) 的算法。
// 输入: nums = [1,3,5,6], target = 5
// 输出: 2
// 输入: nums = [1,3,5,6], target = 2
// 输出: 1
// 输入: nums = [1,3,5,6], target = 7
// 输出: 4
func searchInsert(nums []int, target int) int {
	l := len(nums)
	left := 0
	right := l - 1
	res := l
	for {
		pos := (right-left)/2 + left
		if target <= nums[pos] {
			res = pos
			right = pos - 1
		} else {
			left = pos + 1
		}
		if left > right {
			break
		}
	}
	return res
}

// //////////////////////////////////////////////
// 最后一个单词的长度
// 给你一个字符串 s，由若干单词组成，单词前后用一些空格字符隔开。返回字符串中 最后一个 单词的长度。
// 单词 是指仅由字母组成、不包含任何空格字符的最大子字符串。
// s 仅有英文字母和空格 ' ' 组成
// s 中至少存在一个单词
func lengthOfLastWord(s string) int {
	l := len(s)
	var left, right = -1, 0
	for i := l - 1; i >= 0; i-- {
		if s[i] != ' ' && right == 0 {
			right = i
		}
		if s[i] == ' ' && right != 0 {
			left = i
			break
		}
	}
	// println("left:", left, "right:", right)
	return right - left
}

// //////////////////////////////////////////////
// 加一
// 给定一个由 整数 组成的 非空 数组所表示的非负整数，在该数的基础上加一。
// 最高位数字存放在数组的首位， 数组中每个元素只存储单个数字。
// 你可以假设除了整数 0 之外，这个整数不会以零开头。
// 输入：digits = [1,2,3]
// 输出：[1,2,4]
// 输入：digits = [4,3,2,1]
// 输出：[4,3,2,2]
// 输入：digits = [9]
// 输出：[1,0]
func plusOne(digits []int) []int {
	l := len(digits)
	for i := l - 1; i >= 0; i-- {
		digits[i] += 1
		if digits[i] == 10 {
			digits[i] = 0
		} else {
			break
		}
	}
	if digits[0] == 0 {
		digits = append([]int{1}, digits...)
	}
	return digits
}

// //////////////////////////////////////////////
// 二进制求和
// 给你两个二进制字符串 a 和 b ，以二进制字符串的形式返回它们的和。
// 输入:a = "11", b = "1"
// 输出："100"
// 输入：a = "1010", b = "1011"
// 输出："10101"
// a 和 b 仅由字符 '0' 或 '1' 组成
// 字符串如果不是 "0" ，就不含前导零
func addBinary(a string, b string) string {
	posa := len(a) - 1
	posb := len(b) - 1
	var stra, strb, ext, res string
	for {
		if posa < 0 && posb < 0 {
			if ext == "1" {
				res = "1" + res
			}
			break
		}
		if posa < 0 {
			stra = ""
		} else {
			stra = string(a[posa])
		}
		if posb < 0 {
			strb = ""
		} else {
			strb = string(b[posb])
		}
		sum := ""
		if stra == "1" {
			sum = sum + stra
		}
		if strb == "1" {
			sum = sum + strb
		}
		if ext == "1" {
			sum = sum + ext
		}
		if len(sum) == 0 {
			res = "0" + res
		} else if len(sum) == 1 {
			res = "1" + res
		} else if len(sum) == 2 {
			res = "0" + res
		} else if len(sum) == 3 {
			res = "1" + res
		}
		if len(sum) >= 2 {
			ext = "1"
		} else {
			ext = ""
		}
		posa--
		posb--
	}
	return res
}

// //////////////////////////////////////////////
// 爬楼梯
// 假设你正在爬楼梯。需要 n 阶你才能到达楼顶。
// 每次你可以爬 1 或 2 个台阶。你有多少种不同的方法可以爬到楼顶呢？
// 输入：n = 2
// 输出：2
// 输入：n = 3
// 输出：3
// 动态规划问题：f(x)=f(x−1)+f(x−2)
func climbStairs(n int) int {
	p, q, r := 0, 0, 1
	for i := 1; i <= n; i++ {
		p = q
		q = r
		r = p + q
	}
	return r
}

// //////////////////////////////////////////////
// 删除排序链表中的重复元素
// 给定一个已排序的链表的头 head ， 删除所有重复的元素，使每个元素只出现一次 。返回 已排序的链表 。
func deleteDuplicates(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	cur := head
	for cur.Next != nil {
		if cur.Val == cur.Next.Val {
			cur.Next = cur.Next.Next
		} else {
			cur = cur.Next
		}
	}

	return head
}

// //////////////////////////////////////////////
// 合并两个有序数组
// 输入：nums1 = [1,2,3,0,0,0], m = 3, nums2 = [2,5,6], n = 3
// 输出：[1,2,2,3,5,6]
// 输入：nums1 = [1], m = 1, nums2 = [], n = 0
// 输出：[1]
// 输入：nums1 = [0], m = 0, nums2 = [1], n = 1
// 输出：[1]
// 输入：nums1 = [1,2,3,7,9,10,0,0,0], m = 9, nums2 = [2,5,6], n = 3
// 输出：[1,2,2,3,5,6,7,9,10]
func merge(nums1 []int, m int, nums2 []int, n int) {
	sorted := make([]int, 0, m+n)
    p1, p2 := 0, 0
    for {
        if p1 == m {
            sorted = append(sorted, nums2[p2:]...)
            break
        }
        if p2 == n {
            sorted = append(sorted, nums1[p1:]...)
            break
        }
        if nums1[p1] < nums2[p2] {
            sorted = append(sorted, nums1[p1])
            p1++
        } else {
            sorted = append(sorted, nums2[p2])
            p2++
        }
    }
    copy(nums1, sorted)
}
// nums1 := []int{1, 2, 3, 7, 8, 10, 0, 0, 0}
// nums2 := []int{2, 5, 6}
// merge(nums1, 6, nums2, 3)

// nums1 := []int{1}
// nums2 := []int{}
// merge(nums1, 1, nums2, 0)

// nums1 := []int{0}
// nums2 := []int{1}
// merge(nums1, 0, nums2, 1)

// nums1 := []int{1,0}
// nums2 := []int{2}
// merge(nums1, 1, nums2, 1)

// nums1 := []int{4,0,0,0,0,0}
// nums2 := []int{1,2,3,5,6}
// merge(nums1, 1, nums2, 5)

////////////////////////////////////////////////

////////////////////////////////////////////////

func Run() {
	// res := climbStairs(5)
	// fmt.Println(res)

	


	fmt.Println(nums1)
}
