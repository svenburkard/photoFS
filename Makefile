# photoFS Makefile

info:
	@echo ""
	@echo  :: Avaiable make targets : $(shell grep -P "^[a-zA-Z0-9_]+:" Makefile | grep -Po "^[a-zA-Z0-9_]+")
	@echo ""

init: init_client init_server
fmt: fmt_client fmt_server
build: build_client build_server
all: init fmt build


init_client:
	@echo "[INFO] initializing photoFS client Golang project"
	cd client && if [ ! -f go.mod ];then go mod init photofs_client;fi
	cd client && go mod tidy

init_server:
	@echo "[INFO] initializing photoFS server Golang project"
	cd server && if [ ! -f go.mod ];then go mod init photofs_server;fi
	cd server && go mod tidy


fmt_client:
	@echo "[INFO] formating photoFS client Golang code"
	cd client && go fmt

fmt_server:
	@echo "[INFO] formating photoFS server Golang code"
	cd server && go fmt


build_client:
	@echo "[INFO] building photoFS client Golang binary"
	cd client && go build -o ../bin/

build_server:
	@echo "[INFO] building photoFS server Golang binary"
	cd server && go build -o ../bin/


run_client:
	@echo "[INFO] running photoFS client on-demand"
	cd client && go run .

run_server:
	@echo "[INFO] running photoFS server on-demand"
	cd server && go run .
