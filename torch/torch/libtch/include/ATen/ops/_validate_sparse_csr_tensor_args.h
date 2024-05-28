#pragma once

// @generated by torchgen/gen.py from Function.h

#include <ATen/Context.h>
#include <ATen/DeviceGuard.h>
#include <ATen/TensorUtils.h>
#include <ATen/TracerMode.h>
#include <ATen/core/Generator.h>
#include <ATen/core/Reduction.h>
#include <ATen/core/Tensor.h>
#include <c10/core/Scalar.h>
#include <c10/core/Storage.h>
#include <c10/core/TensorOptions.h>
#include <c10/util/Deprecated.h>
#include <c10/util/Optional.h>



#include <ATen/ops/_validate_sparse_csr_tensor_args_ops.h>

namespace at {


// aten::_validate_sparse_csr_tensor_args(Tensor crow_indices, Tensor col_indices, Tensor values, int[] size) -> ()
inline void _validate_sparse_csr_tensor_args(const at::Tensor & crow_indices, const at::Tensor & col_indices, const at::Tensor & values, at::IntArrayRef size) {
    return at::_ops::_validate_sparse_csr_tensor_args::call(crow_indices, col_indices, values, size);
}

}
