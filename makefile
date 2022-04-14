
run:
	go run examples/personal/20220311_space_junk/app.go
build-html5:
	cp -r asset examples/personal/20220311_space_junk/output/
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" examples/personal/20220311_space_junk/output/
	GOOS=js GOARCH=wasm go build -o examples/personal/20220311_space_junk/output/app.wasm examples/personal/20220311_space_junk/app.go
	light-server -s examples/personal/20220311_space_junk/output -p 8001
