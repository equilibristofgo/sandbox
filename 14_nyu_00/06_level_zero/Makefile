OCLOC = /opt/intel/oneapi/vtune/2025.10/bin64/gma/GTPin/Profilers/ocloc/Bin/intel64/ocloc
CC = g++
CFLAGS = -fPIC -I/usr/include/level_zero
LDFLAGS = -shared -lze_loader

all: kernel.spv_bmg.spv libmatrix_ze.so matrix_go

kernel.spv_bmg.spv: kernel.cl
	$(OCLOC) compile -file $< -device bmg -output kernel.spv

libmatrix_ze.so: matrix_ze.cpp matrix_ze.h kernel.spv_bmg.spv
	$(CC) $(CFLAGS) $(LDFLAGS) -o $@ matrix_ze.cpp

matrix_go: main.go libmatrix_ze.so
	go build -o $@ main.go

clean:
	rm -f kernel.spv_bmg.spv kernel.spv_bmg.bin libmatrix_ze.so matrix_go *.log *.bin

.PHONY: all clean
