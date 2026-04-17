CC = icpx
CFLAGS = -fsycl -fPIC
LDFLAGS = -fsycl -shared

all: libmatrix_sycl.so

libmatrix_sycl.so: matrix_sycl.cpp matrix_sycl.h
	$(CC) $(CFLAGS) $(LDFLAGS) -o $@ $<

clean:
	rm -f libmatrix_sycl.so

.PHONY: all clean
