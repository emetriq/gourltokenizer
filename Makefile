build: tokenizer.so
tokenizer.so:
	go build -o tokenizer.so -buildmode=c-shared tokenizer/cexport/tokenizer.go