import ctypes as ct
from typing import List
import json

_lib = ct.cdll.LoadLibrary("./tokenizer.so")

_lib.TokenizeEng.argtypes = [ct.c_char_p, ct.c_int]
_lib.TokenizeEng.restype = ct.POINTER(ct.c_ubyte*8)
_lib.Free.argtypes = ct.c_void_p,
_lib.Free.restype = None

tokenize_eng = _lib.TokenizeEng
free = _lib.Free

def tokenize(urls: List[str]):
    try:
        data = json.dumps(urls).encode('utf-8')
        ptr = tokenize_eng(data, len(data))
        length = int.from_bytes(ptr.contents, byteorder='little')
        data = bytes(ct.cast(ptr,
                ct.POINTER(ct.c_ubyte*(8 + length))
                ).contents[8:])
        return json.loads(data.decode('utf-8'))
    finally:
        free(ptr)



print(tokenize(["https://www.google.com/hallo/essen", "https://www.facebook.com/autos/geld/news"]))