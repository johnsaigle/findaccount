make:
	clean
	build

build:
	go build -o findaccounts cmd/findaccounts/main.go

build-debug:
	go build -gcflags="all=-N -l" -o findaccounts cmd/findaccounts/main.go

clean:
	rm findaccounts
