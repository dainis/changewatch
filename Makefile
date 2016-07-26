updatedeps:
	go get -u github.com/kardianos/govendor
	govendor fetch +vendor

build:
	go build -o bin/changewatch

test_run: build
	mkdir -p test_run/something/something
	bin/changewatch test_run/ echo change&
	sleep 0.1
	echo yolo > test_run/somefile
	sleep 0.1
	mkdir -p test_run/some_dir
	sleep 0.1
	touch test_run/some_dir/somefile
	sleep 0.1
	rm -r test_run/some_dir
	sleep 0.1
	killall changewatch

.PHONY: updatedeps build
