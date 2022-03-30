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

install_completion:
	rm -f shell_completion/_dbconntest
	rm -f shell_completion/dbconntest
	./dbconntest completion zsh --URL=dummy --driver=dummy > shell_completion/_dbconntest
	./dbconntest completion bash --URL=dummy --driver=dummy > shell_completion/dbconntest
