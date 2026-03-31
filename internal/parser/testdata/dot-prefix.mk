.PHONY: build test clean

build: main.o
	gcc -o build main.o

test: build
	./run-tests.sh

clean:
	rm -rf build/

.SUFFIXES: .c .o

.DEFAULT:
	echo "No rule for $@"
