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

reuslt := tok.TokenizeV3("mailto://www.Subdomain.example.com/HSV-fussbal%3asome/a")
// custom stop words
reuslt2 := tok.TokenizeV3("mailto://www.Subdomain.example.com/HSV-fussball%3asome/a", func(val string) bool {
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

goos: linux
goarch: amd64
pkg: github.com/emetriq/gourltokenizer/tokenizer
cpu: 11th Gen Intel(R) Core(TM) i5-11500H @ 2.90GHz
| Benchmark                         | runs    | time/op     | B/op     | allocs/op   |
|-----------------------------------|---------|-------------|----------|-------------|
| BenchmarkEscapedURLTokenizerV3-12 | 1000000 | 1080 ns/op  | 496 B/op | 3 allocs/op |
| BenchmarkURLTokenizerV3-12        | 4751826 | 255.5 ns/op | 256 B/op | 1 allocs/op |
| BenchmarkURLTokenizerV3Fast-12    | 6231590 | 191.6 ns/op | 256 B/op | 1 allocs/op |
| BenchmarkEscapedURLTokenizerV2-12 | 1000000 | 1042 ns/op  | 496 B/op | 3 allocs/op |
| BenchmarkURLTokenizerV2-12        | 3813273 | 484.2 ns/op | 256 B/op | 1 allocs/op |
| BenchmarkURLTokenizerV2Fast-12    | 5835351 | 199.6 ns/op | 256 B/op | 1 allocs/op |
| BenchmarkEscapedURLTokenizerV1-12 | 1942860 | 1084 ns/op  | 496 B/op | 3 allocs/op |
| BenchmarkURLTokenizerV1-12        | 2495599 | 510.7 ns/op | 272 B/op | 2 allocs/op |
| BenchmarkTokenizerV1-12           | 9431893 | 122.9 ns/op | 256 B/op | 1 allocs/op |
| BenchmarkTokenizerV2-12           | 7669710 | 157.0 ns/op | 256 B/op | 1 allocs/op |
| BenchmarkTokenizerV3-12           | 8120326 | 158.3 ns/op | 256 B/op | 1 allocs/op |