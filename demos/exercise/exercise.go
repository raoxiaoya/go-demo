package exercise

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
)

func Run() {
	KaiHeTiao()
	// angleTest()
}

func angleTest() {
	// fmt.Println(calculateAngleY(Point{
	// 	Y: 124.92188430450126,
	// 	X: 140.62501882960643,
	// }, Point{
	// 	Y: 207.65624975830502,
	// 	X: 136.8749997565486,
	// }))

	// fmt.Println(calculateAngleY(Point{
	// 	Y: 207.65624975830502,
	// 	X: 136.8749997565486,
	// }, Point{
	// 	Y: 279.99999975003914,
	// 	X: 132.34374975706632,
	// }))

	fmt.Println(calculateAngle3(
		Point{
			X: 242.40551380705966,
			Y: 254.18392354045025,
		},
		Point{
			X: 188.029114415816,
			Y: 192.39023238969463,
		},
		Point{
			X: 246.48420651451895,
			Y: 113.71941940066803,
		}))
}

type Point struct {
	X     float64
	Y     float64
	Score float64
	Name  string
}

type Output struct {
	Score     float64
	Keypoints []Point
}

var minPoseConfidence = 0.6  // 总体置信度区间
var minPointConfidence = 0.6 // 单个点的置信度区间

func GetOutput(filename string) []map[string]Point {
	f, _ := os.Open(filename)
	b, _ := io.ReadAll(f)
	var outputs []Output
	json.Unmarshal(b, &outputs)

	list := make([]map[string]Point, 0)
	for _, v := range outputs {
		if v.Score < minPoseConfidence {
			continue
		}
		data := make(map[string]Point)
		for _, v := range v.Keypoints {
			data[v.Name] = v
		}
		list = append(list, data)
	}

	return list
}

type Result struct {
	LastStatus     uint8
	LastStatusTime int64
	CurrentStatus  uint8
	Count          int
}

var dirname = "demos/exercise/posedata"

// 开合跳
func KaiHeTiao() {
	list := GetOutput(dirname + "/kaihetiao/1739416846502-data.json")
	result := Result{
		CurrentStatus: 0,
	}
	// close, open, close, close, close, open
	for _, v := range list {
		if close(v) {
			fmt.Println("close")
			if result.CurrentStatus == 0 {
				result.CurrentStatus = 1
			} else if result.CurrentStatus == 2 {
				result.CurrentStatus = 1
				result.Count++
			}
		} else if open(v) {
			fmt.Println("open")
			if result.CurrentStatus == 1 {
				result.CurrentStatus = 2
			}
		} else {
			fmt.Println("none")
		}
	}
	fmt.Print(result.Count)
}

var threshold float64 = 10 // 偏差
func checkAngle2(point0, point1 Point, angle float64) bool {
	return math.Abs(calculateAngle2(point0, point1)-angle) < threshold
}

func checkAngle3(point0, point1, point2 Point, angle float64) bool {
	return math.Abs(calculateAngle3(point0, point1, point2)-angle) < threshold
}

var closeAngle float64 = 0

func close(keypoints map[string]Point) bool {
	// right
	if !checkAngle3(keypoints["right_shoulder"], keypoints["right_elbow"], keypoints["right_wrist"], closeAngle) ||
		!checkAngle3(keypoints["right_hip"], keypoints["right_knee"], keypoints["right_ankle"], closeAngle) {
		return false
	}

	// left
	if !checkAngle3(keypoints["left_shoulder"], keypoints["left_elbow"], keypoints["left_wrist"], closeAngle) ||
		!checkAngle3(keypoints["left_hip"], keypoints["left_knee"], keypoints["left_ankle"], closeAngle) {
		return false
	}

	if !(keypoints["right_wrist"].Y > keypoints["right_elbow"].Y && keypoints["right_elbow"].Y > keypoints["right_shoulder"].Y) {
		return false
	}
	if !(keypoints["left_wrist"].Y > keypoints["left_elbow"].Y && keypoints["left_elbow"].Y > keypoints["left_shoulder"].Y) {
		return false
	}

	return true
}

var openAngle float64 = 10

func open(keypoints map[string]Point) bool {
	// right
	if !checkAngle2(keypoints["right_shoulder"], keypoints["right_elbow"], openAngle) ||
		!checkAngle2(keypoints["right_elbow"], keypoints["right_wrist"], openAngle) ||
		!checkAngle2(keypoints["right_hip"], keypoints["right_knee"], openAngle) ||
		!checkAngle2(keypoints["right_knee"], keypoints["right_ankle"], openAngle) {
		return false
	}

	// left
	if !checkAngle2(keypoints["left_shoulder"], keypoints["left_elbow"], openAngle) ||
		!checkAngle2(keypoints["left_elbow"], keypoints["left_wrist"], openAngle) ||
		!checkAngle2(keypoints["left_hip"], keypoints["left_knee"], openAngle) ||
		!checkAngle2(keypoints["left_knee"], keypoints["left_ankle"], openAngle) {
		return false
	}

	if !(keypoints["right_wrist"].Y < keypoints["right_elbow"].Y && keypoints["right_elbow"].Y < keypoints["right_shoulder"].Y) {
		return false
	}
	if !(keypoints["left_wrist"].Y < keypoints["left_elbow"].Y && keypoints["left_elbow"].Y < keypoints["left_shoulder"].Y) {
		return false
	}

	return true
}

func calculateAngle2(p, q Point) float64 {
	return (180 * math.Acos(math.Abs(p.Y-q.Y)/math.Sqrt(math.Pow(p.X-q.X, 2)+math.Pow(p.Y-q.Y, 2)))) / math.Pi
}

func calculateAngleX(p, q Point) float64 {
	return (180 * math.Acos(math.Abs(p.X-q.X)/math.Sqrt(math.Pow(p.X-q.X, 2)+math.Pow(p.Y-q.Y, 2)))) / math.Pi
}

func calculateAngle3(point0, point1, point2 Point) float64 {
	// 计算向量 point1 -> point0 和 point2 -> point0
	dx1 := point1.X - point0.X
	dy1 := point1.Y - point0.Y
	dx2 := point1.X - point2.X
	dy2 := point1.Y - point2.Y

	// 计算两个向量的点积和模长
	dotProduct := dx1*dx2 + dy1*dy2
	magnitude1 := math.Sqrt(dx1*dx1 + dy1*dy1)
	magnitude2 := math.Sqrt(dx2*dx2 + dy2*dy2)

	// 计算余弦值
	cosine := dotProduct / (magnitude1 * magnitude2)

	// 将余弦值转换为角度（以度为单位）
	angleInDegrees := (180 * math.Acos(cosine)) / math.Pi

	return angleInDegrees
}
