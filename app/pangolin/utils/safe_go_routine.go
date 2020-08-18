package utils

import (
	"context"
	"google.golang.org/appengine/log"
	"runtime"
)

// recoveryCallBack : 发生panic时回调
func SafeGoroutine(f func(), recoveryCallBack ...func(interface{})) {
	defer Recovery(recoveryCallBack...)
	f()
}

func Recovery(funcs ...func(interface{})) {
	if r := recover(); r != nil {
		recovered := false
		if len(funcs) > 0 {
			for _, fun := range funcs {
				if fun != nil {
					fun(r)
					recovered = true
				}
			}
		}
		if !recovered {
			buf := make([]byte, 1<<18)
			n := runtime.Stack(buf, false)
			log.Errorf(context.Background(), "%v, STACK: %s", r, buf[0:n])
		}
	}
}


