package main

/*
#cgo pkg-config: python-3.8
#include <stdlib.h>
#include <Python.h>
*/
import "C"

// func main() {
// 	C.Py_Initialize()
// 	fmt.Println(C.GoString(C.Py_GetVersion()))

// 	fooModule := C.PyImport_ImportModule(C.CString("foo"))
// 	if fooModule == nil {
// 		panic("Error importing module")
// 	}

// 	helloFunc := fooModule.GetAttrString("hello")
// 	if helloFunc == nil {
// 		panic("Error importing function")
// 	}

// 	// to convert from string to c string :
// 	// cs := C.CString(s)
// 	// defer C.free(unsafe.Pointer(cs))

// 	helloFunc.Call(C.PyTuple_New(0), C.PyDict_New())
// 	C.Py_Finalize()
// }
