# photoFS Makefile

info:
	@echo ""
	@echo  :: Avaiable make targets : $(shell grep -P "^[a-zA-Z0-9_]+:" Makefile | grep -Po "^[a-zA-Z0-9_]+")
	@echo ""

init: init_lib init_client init_server
fmt: fmt_lib fmt_client fmt_server
build: build_client build_server
all: init fmt build
run: run_server run_client


init_lib:
	@echo "[INFO] initializing photoFS lib Golang project"
	cd lib && if [ ! -f go.mod ];then go mod init photofs_lib;fi
	cd lib && go mod tidy

init_client:
	@echo "[INFO] initializing photoFS client Golang project"
	cd client && if [ ! -f go.mod ];then go mod init photofs_client;fi
	cd client && go mod tidy

init_server:
	@echo "[INFO] initializing photoFS server Golang project"
	cd server && if [ ! -f go.mod ];then go mod init photofs_server;fi
	cd server && go mod tidy


fmt_lib:
	@echo "[INFO] formating photoFS lib Golang code"
	cd lib && go fmt

fmt_client:
	@echo "[INFO] formating photoFS client Golang code"
	cd client && go fmt

fmt_server:
	@echo "[INFO] formating photoFS server Golang code"
	cd server && go fmt


build_client: fmt_lib fmt_client
	@echo "[INFO] building photoFS client Golang binary"
	cd client && go build -o ../bin/

build_server: fmt_lib fmt_server
	@echo "[INFO] building photoFS server Golang binary"
	cd server && go build -o ../bin/


test_data:
	@echo "[INFO] create test data files from server/tag_map.json"
	mkdir -p /tmp/test/src/ \
		&& cat server/tag_map.json  | jq "keys" | grep '^\s*"' | cut -d '"' -f2 | xargs touch


run_client:
	@echo "[INFO] running photoFS client binary"
	bin/photofs_client

run_server: test_data
	@echo "[INFO] running photoFS server binary"
	bin/photofs_server
