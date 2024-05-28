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



#include <ATen/ops/_lu_with_info_ops.h>

namespace at {


// aten::_lu_with_info(Tensor self, bool pivot=True, bool check_errors=True) -> (Tensor LU, Tensor pivots, Tensor info)
inline ::std::tuple<at::Tensor,at::Tensor,at::Tensor> _lu_with_info(const at::Tensor & self, bool pivot=true, bool check_errors=true) {
    return at::_ops::_lu_with_info::call(self, pivot, check_errors);
}

}
