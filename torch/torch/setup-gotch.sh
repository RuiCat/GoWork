#!/bin/bash

GOTCH_VERSION="${GOTCH_VER:-v0.9.0}"
CUDA_VERSION="${CUDA_VER:-11.8}"

# Install gotch
#==============
echo "GOPATH:'$GOPATH'"
echo "GOTCH_VERSION: '$GOTCH_VERSION'"
echo "CUDA_VERSION: '$CUDA_VERSION'"

# Setup gotch for CUDA or non-CUDA device:
#=========================================
GOTCH_LIB_FILE="./libtch/lib.go"
if [ -f $GOTCH_LIB_FILE ]
then
  echo "$GOTCH_LIB_FILE existing. Deleting..."
  sudo rm $GOTCH_LIB_FILE
fi

# Create files for CUDA or non-CUDA device
if [ $CUDA_VERSION == "cpu" ]; then
  echo "creating $GOTCH_LIB_FILE for CPU"
  sudo tee -a $GOTCH_LIB_FILE > /dev/null <<EOT
package libtch

// #cgo CFLAGS: -I${SRCDIR} -O3 -Wall -Wno-unused-variable -Wno-deprecated-declarations -Wno-c++11-narrowing -g -Wno-sign-compare -Wno-unused-function
// #cgo CFLAGS: -I/usr/local/include
// #cgo CFLAGS: -D_GLIBCXX_USE_CXX11_ABI=1
// #cgo LDFLAGS: -lstdc++ -ltorch -lc10 -ltorch_cpu -L/lib64
// #cgo CXXFLAGS: -std=c++17
import "C"
EOT
else
  echo "creating $GOTCH_LIB_FILE for GPU"
  sudo tee -a $GOTCH_LIB_FILE > /dev/null <<EOT
package libtch

package libtch

// #cgo CFLAGS: -I${SRCDIR} -O3 -Wall -Wno-unused-variable -Wno-deprecated-declarations -Wno-c++11-narrowing -g -Wno-sign-compare -Wno-unused-function
// #cgo CFLAGS: -I/usr/local/cuda/include
// #cgo LDFLAGS: -lstdc++ -ltorch -ltorch_cpu -ltorch_cuda -lc10
// #cgo LDFLAGS: -L/usr/local/cuda/lib64 -lcuda
// #cgo CFLAGS: -D_GLIBCXX_USE_CXX11_ABI=1
// #cgo CXXFLAGS: -std=c++17
import "C"
EOT
fi

sudo ldconfig
