export

# Includes


.PHONY: docker-build build

build:
	@go build -o build/dbconntest

docker-build:
	@docker build -f build/Dockerfile -t titan_builder .
	docker run --rm \
	  -v ${HOME}/.ssh:/root/temp_ssh:ro \
	  -v ${PWD}:/usr/src/tempfolder \
	  -v ${GOPATH}/pkg/mod:/go/pkg/mod:delegated \
	  -w /usr/src/tempfolder \
	  titan_builder "build/build"

install_completion: build
	rm -f shell_completion/_dbconntest_temp
	rm -f shell_completion/dbconntest_temp
	go run ./main.go completion bash > shell_completion/dbconntest_temp & mv -f shell_completion/dbconntest_temp shell_completion/dbconntest
	go run ./main.go completion zsh > shell_completion/_dbconntest_temp & mv -f shell_completion/_dbconntest_temp shell_completion/_dbconntest
