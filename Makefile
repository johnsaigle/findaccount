make:
	clean
	build

build:
	go build

build-debug:
	go build -gcflags="all=-N -l"

clean:
	rm findaccount
