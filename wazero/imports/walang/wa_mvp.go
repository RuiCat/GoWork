// 版权 @2023 凹语言 作者。保留所有权利。

// github.com/tetratelabs/wazero/imports/walang

package walang

import (
	"wazero/api"
	"wazero/internal/sys"
	"wazero/internal/wasm"
)

func ModCallContextSys(m api.Module) *sys.Context {
	return m.(*wasm.CallContext).Sys
}
