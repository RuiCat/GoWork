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



#include <ATen/ops/unique_dim_consecutive_ops.h>

namespace at {


// aten::unique_dim_consecutive(Tensor self, int dim, bool return_inverse=False, bool return_counts=False) -> (Tensor, Tensor, Tensor)
inline ::std::tuple<at::Tensor,at::Tensor,at::Tensor> unique_dim_consecutive(const at::Tensor & self, int64_t dim, bool return_inverse=false, bool return_counts=false) {
    return at::_ops::unique_dim_consecutive::call(self, dim, return_inverse, return_counts);
}

// aten::unique_dim_consecutive.out(Tensor self, int dim, bool return_inverse=False, bool return_counts=False, *, Tensor(a!) out0, Tensor(b!) out1, Tensor(c!) out2) -> (Tensor(a!), Tensor(b!), Tensor(c!))
inline ::std::tuple<at::Tensor &,at::Tensor &,at::Tensor &> unique_dim_consecutive_out(at::Tensor & out0, at::Tensor & out1, at::Tensor & out2, const at::Tensor & self, int64_t dim, bool return_inverse=false, bool return_counts=false) {
    return at::_ops::unique_dim_consecutive_out::call(self, dim, return_inverse, return_counts, out0, out1, out2);
}
// aten::unique_dim_consecutive.out(Tensor self, int dim, bool return_inverse=False, bool return_counts=False, *, Tensor(a!) out0, Tensor(b!) out1, Tensor(c!) out2) -> (Tensor(a!), Tensor(b!), Tensor(c!))
inline ::std::tuple<at::Tensor &,at::Tensor &,at::Tensor &> unique_dim_consecutive_outf(const at::Tensor & self, int64_t dim, bool return_inverse, bool return_counts, at::Tensor & out0, at::Tensor & out1, at::Tensor & out2) {
    return at::_ops::unique_dim_consecutive_out::call(self, dim, return_inverse, return_counts, out0, out1, out2);
}

}
