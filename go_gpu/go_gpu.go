package go_gpu

import (
	"math/rand"
	"syscall"
	"time"
	"unsafe"
)

const N = 1 << 20

type Lib struct {
    dll     *syscall.DLL
    dotProc *syscall.Proc
}

func LoadLib() (*Lib, error) {
    l := &Lib{}
    var err error
    defer func() {
        if nil != err {
            l.Release()
        }
    }()

    if l.dll, err = syscall.LoadDLL("cuda.dll"); nil != err {
        return nil, err
    }
    if l.dotProc, err = l.dll.FindProc("dot"); nil != err {
        return nil, err
    }
    return l, nil
}

func (l *Lib) Release() {
    if nil != l.dll {
        l.dll.Release()
    }
}

func (l *Lib) Dot(x, y []float32) float32 {
    var r float32
    l.dotProc.Call(
        uintptr(unsafe.Pointer(&x[0])),
        uintptr(unsafe.Pointer(&y[0])),
        uintptr(len(x)),
        uintptr(unsafe.Pointer(&r)),
    )
    return r
}

func main() {
    lib, err := LoadLib()
    if nil != err {
        println(err.Error())
        return
    }
    defer lib.Release()

    rand.Seed(time.Now().Unix())
    x, y := make([]float32, N), make([]float32, N)
    for i := 0; i < N; i++ {
        x[i], y[i] = rand.Float32(), rand.Float32()
    }

    t := time.Now()
    var r float32
    for i := 0; i < 100; i++ {
        r = 0
        for i := 0; i < N; i++ {
            r += x[i] * y[i]
        }
    }
    println(time.Now().Sub(t).Microseconds())
    println(r)

    t = time.Now()
    for i := 0; i < 100; i++ {
        r = lib.Dot(x, y)
    }
    println(time.Now().Sub(t).Microseconds())
    println(r)
}