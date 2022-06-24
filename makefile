# game_project_dir = 20220311_space_junk
game_project_dir = 20220516_forest_zombies

run:
	go run examples/personal/${game_project_dir}/app.go
build-html5:
	cp -r asset examples/personal/${game_project_dir}/output/
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" examples/personal/${game_project_dir}/output/
	GOOS=js GOARCH=wasm go build -o examples/personal/${game_project_dir}/output/app.wasm examples/personal/${game_project_dir}/app.go
	light-server -s examples/personal/${game_project_dir}/output -p 8002
