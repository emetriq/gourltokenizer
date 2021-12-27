package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"encoding/binary"
	"encoding/json"
	"unsafe"

	tok "github.com/emetriq/gourltokenizer/tokenizer"
)

//export TokenizeEng
func TokenizeEng(urlsByte *C.char, size C.int) unsafe.Pointer {
	d := C.GoBytes(unsafe.Pointer(urlsByte), size)
	urls := make([]string, 0, size)
	_ = json.Unmarshal([]byte(d), &urls)
	result := make([][]string, 0, len(urls))
	for _, url := range urls {
		result = append(result, tok.TokenizeV2(url, tok.IsEnglishStopWord))
	}
	resultByte, _ := json.Marshal(result)
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(resultByte)))
	return C.CBytes(append(length, resultByte...))
}

//export Free
func Free(addr *C.char) {
	C.free(unsafe.Pointer(addr))
}

func main() {}
