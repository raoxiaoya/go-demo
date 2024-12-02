package go_call_dll

import (
	"bytes"
	"encoding/json"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func Run() {
	GetStringAndInt()
	// GetInfo()
	// f0()
	// f1()
	// f2()
	// f3()
	// f4()
}

func test() {
	name := "raoxiaoya"

	fmt.Println(unsafe.Pointer(&name)) // 0xc000052090
	fmt.Printf("%p\n", &name)          // 0xc000052090

	fmt.Println(uintptr(unsafe.Pointer(&name))) // 824634056848
	fmt.Printf("%d\n", &name)                   // 824634056848

	fmt.Println(uintptr(3)) // 3
}

func GetStringAndInt() {
	nativeModule := windows.NewLazyDLL(`D:\dev\php\magook\trunk\server\go-demo\demos\go_call_dll\WebView2Loader.dll`)
	
	// 返回结果为string
	// public STDAPI GetAvailableCoreWebView2BrowserVersionString(PCWSTR browserExecutableFolder, LPWSTR * versionInfo)
	nativeGetAvailableCoreWebView2BrowserVersionString := nativeModule.NewProc("GetAvailableCoreWebView2BrowserVersionString")
	var versionResult *uint16
	var version string
	hr, _, _ := nativeGetAvailableCoreWebView2BrowserVersionString.Call(uintptr(unsafe.Pointer(nil)), uintptr(unsafe.Pointer(&versionResult)))
	defer windows.CoTaskMemFree(unsafe.Pointer(versionResult))
	if hr != uintptr(windows.S_OK) {
		fmt.Println(hr)
		return
	}
	version = windows.UTF16PtrToString(versionResult)
	fmt.Println("versionResult:", version)

	// 返回结果为int
	// public STDAPI CompareBrowserVersions(PCWSTR version1, PCWSTR version2, int * result)
	//	-1 = v1 < v2
	//	 0 = v1 == v2
	//	 1 = v1 > v2
	nativeCompareBrowserVersions := nativeModule.NewProc("CompareBrowserVersions")
	var compareResult int
	_v1 := "131.0.2903.63"
	_v2 := "131.0.2903.63"
	_, _, err := nativeCompareBrowserVersions.Call(uintptr(unsafe.Pointer(&_v1)), uintptr(unsafe.Pointer(&_v2)), uintptr(unsafe.Pointer(&compareResult)))
	if err != windows.ERROR_SUCCESS {
		compareResult = -2
	}
	fmt.Println("compareResult:", compareResult)
}

var f = "SetConsoleTextAttribute"

type Result struct {
	DwSize              Coord
	DwCursorPosition    Coord
	WAttributes         uint16
	SrWindow            SmallRect
	DwMaximumWindowSize Coord
}
type Coord struct {
	X int16
	Y int16
}
type SmallRect struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

func PrettyPrint(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println(v)
		return
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		fmt.Println(v)
		return
	}

	fmt.Println(out.String())
}

func GetInfo() {
	nativeModule := windows.NewLazyDLL("kernel32.dll")
	p := nativeModule.NewProc("GetConsoleScreenBufferInfo")
	result := Result{}
	r1, r2, lastError := p.Call(uintptr(syscall.Stdout), uintptr(unsafe.Pointer(&result)))
	fmt.Printf("r1: %v, r2: %v, lastError: %v\n", r1, r2, lastError)
	PrettyPrint(result)
}

func f0() {
	nativeModule := windows.NewLazyDLL("kernel32.dll")
	p := nativeModule.NewProc(f)
	// 设置 红底白字
	r1, r2, lastError := p.Call(uintptr(syscall.Stdout), uintptr(0x0040|0x0007))
	fmt.Printf("r1: %v, r2: %v, lastError: %v", r1, r2, lastError) // 1 0 The operation completed successfully.
	// 恢复到黑底白字
	p.Call(uintptr(syscall.Stdout), uintptr(0x0000|0x0007))
}

func f1() {
	// A LazyDLL implements access to a single DLL.
	// It will delay the load of the DLL until the first
	// call to its Handle method or to one of its
	// LazyProc's Addr method.
	//
	// LazyDLL is subject to the same DLL preloading attacks as documented
	// on LoadDLL.
	//
	// Use LazyDLL in golang.org/x/sys/windows for a secure way to
	// load system DLLs.
	dll := syscall.NewLazyDLL("kernel32.dll")
	p := dll.NewProc(f)
	// r1, r2, lastError := p.Call(uintptr(syscall.Stdout), uintptr(0x0002|0x0040))
	r1, r2, lastError := p.Call(uintptr(syscall.Stdout), uintptr(0x0000|0x0007))
	fmt.Println(r1, r2, lastError)
}

func f2() {
	// LoadDLL loads the named DLL file into memory.
	//
	// If name is not an absolute path and is not a known system DLL used by
	// Go, Windows will search for the named DLL in many locations, causing
	// potential DLL preloading attacks.
	//
	// Use LazyDLL in golang.org/x/sys/windows for a secure way to
	// load system DLLs.
	dll, _ := syscall.LoadDLL("kernel32.dll")
	p, _ := dll.FindProc(f)
	r1, r2, lastError := p.Call(uintptr(syscall.Stdout), uintptr(4))
	fmt.Println(r1, r2, lastError)
}

func f3() {
	// MustLoadDLL is like LoadDLL but panics if load operation fails.
	dll := syscall.MustLoadDLL("kernel32.dll")
	p := dll.MustFindProc(f)
	r1, r2, lastError := p.Call(uintptr(syscall.Stdout), uintptr(5))
	fmt.Println(r1, r2, lastError)
}

func f4() {
	handle, _ := syscall.LoadLibrary("kernel32.dll")
	defer syscall.FreeLibrary(handle)

	p, _ := syscall.GetProcAddress(handle, f)
	// r1, r2, errorNo := syscall.Syscall(p, 2, uintptr(syscall.Stdout), uintptr(6), 0)
	r1, r2, errorNo := syscall.SyscallN(p, uintptr(syscall.Stdout), uintptr(6))
	fmt.Println(r1, r2, errorNo)
	// 恢复到白色
	syscall.SyscallN(p, uintptr(syscall.Stdout), uintptr(7))
}
