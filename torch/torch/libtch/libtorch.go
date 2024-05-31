package libtch

// #cgo pkg-config: --define-variable=prefix=. ${SRCDIR}/libtorch.pc
// #cgo CFLAGS: -O3 -Wall -Wno-unused-variable -Wno-deprecated-declarations -Wno-c++11-narrowing -g -Wno-sign-compare -Wno-unused-function -D_GLIBCXX_USE_CXX11_ABI=1
// #cgo LDFLAGS: -Wl,-rpath=./lib -L${SRCDIR}/lib/ -lcuda -lstdc++ -ltorch -ltorch_cpu -ltorch_cuda -lc10
// #cgo CXXFLAGS: -std=c++17
import "C"
