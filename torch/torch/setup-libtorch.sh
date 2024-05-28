#!/bin/bash
LIBTORCH_VERSION="${LIBTORCH_VER:-2.1.0}"
CUDA_VERSION="${CUDA_VER:-11.8}"
if [ "${CUDA_VERSION}" == "cpu" ]; then
  CU_VERSION="cpu"
else
  CU_VERSION="cu${CUDA_VERSION//./}"
fi

# Install Libtorch
#=================
LIBTORCH_ZIP="libtorch-cxx11-abi-shared-with-deps-${LIBTORCH_VERSION}%2B${CU_VERSION}.zip"
LIBTORCH_URL="https://download.pytorch.org/libtorch/${CU_VERSION}/${LIBTORCH_ZIP}"
echo $LIBTORCH_URL
wget  -q --show-progress --progress=bar:force:noscroll "$LIBTORCH_URL"