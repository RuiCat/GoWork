// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !minimal
// +build !minimal

package plot // import "ui/plot"

import (
	_ "ui/plot/vg/vgeps"
	_ "ui/plot/vg/vgimg"
	_ "ui/plot/vg/vgpdf"
	_ "ui/plot/vg/vgsvg"
	_ "ui/plot/vg/vgtex"
)
