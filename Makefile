export

# Includes


.PHONY: docker-build build shell_completion

build:
	@go build -o build/dbconntest

docker-build:
	@docker build -f build/Dockerfile -t titan_builder .
	docker run --rm \
	  -v ${HOME}/.ssh:/root/temp_ssh:ro \
	  -v ${PWD}:/usr/src/tempfolder \
	  -v ${PWD}/../godror:/usr/src/godror \
	  -v ${PWD}/../go_ibm_db:/usr/src/go_ibm_db \
	  -v ${GOPATH}/pkg/mod:/go/pkg/mod:delegated \
	  -w /usr/src/tempfolder \
	  titan_builder "build/build"

shell_completion:
	# this commands do not work from make (dyld[44503]: Library not loaded: libdb2.dylib)
	# but you can run them from the terminal
	rm -f shell_completion/_dbconntest_temp
	rm -f shell_completion/dbconntest_temp
	go run ./main.go completion bash > shell_completion/dbconntest_temp && mv -f shell_completion/dbconntest_temp shell_completion/dbconntest
	go run ./main.go completion zsh > shell_completion/_dbconntest_temp && mv -f shell_completion/_dbconntest_temp shell_completion/_dbconntest
