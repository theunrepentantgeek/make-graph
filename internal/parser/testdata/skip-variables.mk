CC := gcc
CFLAGS = -Wall -O2
SRCS ?= main.c
OBJS += main.o

ifeq ($(DEBUG),1)
CFLAGS += -g
endif

build: main.o
	$(CC) $(CFLAGS) -o build main.o

clean:
	rm -f build
