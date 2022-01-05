# gourltokenizer

A powerful URL tokenizer

# install

`go get github.com/emetriq/gourltokenizer`

# usage

```golang
import (
   tok  "github.com/emetriq/gourltokenizer/tokenizer"
)
// set min token size
tok.MinWordSize = 3
// set default stop words
tok.DefaultStopWordFunc = IsGermanStopWord

reuslt := tok.TokenizeV2("mailto://www.Subdomain.example.com/HSV-fussbal%3asome/a")
// custom stop words
reuslt2 := tok.TokenizeV2("mailto://www.Subdomain.example.com/HSV-fussball%3asome/a", func(val string) bool {
	if val == "fussball" {
		return true
	}
	if val == "Subdomain" {
		return true
	}
	return false
})
```
# Benchmark Results

goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
|Benchmark|runs|time/op|B/op|allocs/op|
|---|---|---|---|---|
BenchmarkURLTokenizerV2-12|2026138|605.3 ns/op|256 B/op|1 allocs/op
BenchmarkURLTokenizerV2Fast-12|3609961|330.6 ns/op|256 B/op|1 allocs/op
BenchmarkURLTokenizerV1-12|1766235|676.3 ns/op|272 B/op|2 allocs/op