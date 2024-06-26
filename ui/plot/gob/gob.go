// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gob // import "ui/plot/gob"

import (
	"encoding/gob"
	"image/color"

	"ui/plot/plotter"
	"ui/plot/text"

	"ui/plot"
)

func init() {
	// register types for proper gob-encoding/decoding
	gob.Register(color.Gray16{})

	// plot.Ticker
	gob.Register(plot.ConstantTicks{})
	gob.Register(plot.DefaultTicks{})
	gob.Register(plot.LogTicks{})

	// plot.Normalizer
	gob.Register(plot.LinearScale{})
	gob.Register(plot.LogScale{})

	// plot.Plotter
	gob.Register(plotter.BarChart{})
	gob.Register(plotter.Histogram{})
	gob.Register(plotter.BoxPlot{})
	gob.Register(plotter.YErrorBars{})
	gob.Register(plotter.XErrorBars{})
	gob.Register(plotter.Function{})
	gob.Register(plotter.GlyphBoxes{})
	gob.Register(plotter.Grid{})
	gob.Register(plotter.Labels{})
	gob.Register(plotter.Line{})
	gob.Register(plotter.QuartPlot{})
	gob.Register(plotter.Scatter{})

	// plotter.XYZer
	gob.Register(plotter.XYZs{})
	gob.Register(plotter.XYValues{})

	// text.Style
	gob.Register(&text.Plain{})
	gob.Register(&text.Latex{})
}
