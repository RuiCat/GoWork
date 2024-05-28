// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package liberation exports the Liberation fonts as a font.Collection.
package liberation // import "ui/plot/font/liberation"

import (
	"fmt"
	"sync"

	"ui/util/liberation/liberationmonobold"
	"ui/util/liberation/liberationmonobolditalic"
	"ui/util/liberation/liberationmonoitalic"
	"ui/util/liberation/liberationmonoregular"
	"ui/util/liberation/liberationsansbold"
	"ui/util/liberation/liberationsansbolditalic"
	"ui/util/liberation/liberationsansitalic"
	"ui/util/liberation/liberationsansregular"
	"ui/util/liberation/liberationserifbold"
	"ui/util/liberation/liberationserifbolditalic"
	"ui/util/liberation/liberationserifitalic"
	"ui/util/liberation/liberationserifregular"

	stdfnt "ui/util/x/image/font"
	"ui/util/x/image/font/opentype"

	"ui/plot/font"
)

var (
	once       sync.Once
	collection font.Collection
)

func Collection() font.Collection {
	once.Do(func() {
		addColl(font.Font{}, liberationserifregular.TTF)
		addColl(font.Font{Style: stdfnt.StyleItalic}, liberationserifitalic.TTF)
		addColl(font.Font{Weight: stdfnt.WeightBold}, liberationserifbold.TTF)
		addColl(font.Font{
			Style:  stdfnt.StyleItalic,
			Weight: stdfnt.WeightBold,
		}, liberationserifbolditalic.TTF)

		// mono variant
		addColl(font.Font{Variant: "Mono"}, liberationmonoregular.TTF)
		addColl(font.Font{
			Variant: "Mono",
			Style:   stdfnt.StyleItalic,
		}, liberationmonoitalic.TTF)
		addColl(font.Font{
			Variant: "Mono",
			Weight:  stdfnt.WeightBold,
		}, liberationmonobold.TTF)
		addColl(font.Font{
			Variant: "Mono",
			Style:   stdfnt.StyleItalic,
			Weight:  stdfnt.WeightBold,
		}, liberationmonobolditalic.TTF)

		// sans-serif variant
		addColl(font.Font{Variant: "Sans"}, liberationsansregular.TTF)
		addColl(font.Font{
			Variant: "Sans",
			Style:   stdfnt.StyleItalic,
		}, liberationsansitalic.TTF)
		addColl(font.Font{
			Variant: "Sans",
			Weight:  stdfnt.WeightBold,
		}, liberationsansbold.TTF)
		addColl(font.Font{
			Variant: "Sans",
			Style:   stdfnt.StyleItalic,
			Weight:  stdfnt.WeightBold,
		}, liberationsansbolditalic.TTF)

		n := len(collection)
		collection = collection[:n:n]
	})

	return collection
}

func addColl(fnt font.Font, ttf []byte) {
	face, err := opentype.Parse(ttf)
	if err != nil {
		panic(fmt.Errorf("vg: could not parse font: %+v", err))
	}
	fnt.Typeface = "Liberation"
	if fnt.Variant == "" {
		fnt.Variant = "Serif"
	}
	collection = append(collection, font.Face{
		Font: fnt,
		Face: face,
	})
}
