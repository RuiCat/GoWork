//Copyright 2013 Vastech SA (PTY) LTD
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package golang

import (
	"bytes"
	"fmt"
	"math"
	"path"
	"text/template"

	"ui/util/gocc/internal/ast"
	"ui/util/gocc/internal/io"
	"ui/util/gocc/internal/parser/lr1/action"
	"ui/util/gocc/internal/parser/lr1/items"
	"ui/util/gocc/internal/token"
)

func GenActionTable(outDir string, prods ast.SyntaxProdList, itemSets *items.ItemSets, tokMap *token.TokenMap, zip bool) map[int]items.RowConflicts {
	if zip {
		return GenCompActionTable(outDir, prods, itemSets, tokMap)
	}
	tmpl, err := template.New("parser action table").Parse(actionTableSrc[1:])
	if err != nil {
		panic(err)
	}
	wr := new(bytes.Buffer)
	data, conflicts := getActionTableData(prods, itemSets, tokMap)
	if err := tmpl.Execute(wr, data); err != nil {
		panic(err)
	}
	io.WriteFile(path.Join(outDir, "parser", "actiontable.go"), wr.Bytes())
	return conflicts
}

type actionTableData struct {
	Rows []*actRow
}

func getActionTableData(prods ast.SyntaxProdList, itemSets *items.ItemSets,
	tokMap *token.TokenMap) (actTab *actionTableData, conflicts map[int]items.RowConflicts) {
	actTab = &actionTableData{
		Rows: make([]*actRow, itemSets.Size()),
	}
	conflicts = make(map[int]items.RowConflicts)
	var cnflcts items.RowConflicts
	var row *actRow
	for i := range actTab.Rows {
		if row, cnflcts = getActionRowData(prods, itemSets.Set(i), tokMap); len(cnflcts) > 0 {
			conflicts[i] = cnflcts
		}
		actTab.Rows[i] = row
	}
	return
}

type actRow struct {
	CanRecover bool
	Actions    []string
}

func getActionRowData(prods ast.SyntaxProdList, set *items.ItemSet, tokMap *token.TokenMap) (data *actRow, conflicts items.RowConflicts) {
	data = &actRow{
		CanRecover: set.CanRecover(),
		Actions:    make([]string, len(tokMap.TypeMap)),
	}
	conflicts = make(items.RowConflicts)
	var max int
	// calculate padding.
	for _, sym := range tokMap.TypeMap {
		act, _ := set.Action(sym)
		switch act1 := act.(type) {
		case action.Accept:
			n := len("accept(true),")
			if n > max {
				max = n
			}
		case action.Error:
			n := len("nil,")
			if n > max {
				max = n
			}
		case action.Reduce:
			n := len("reduce(") + nbytes(int(act1)) + len("),")
			if n > max {
				max = n
			}
		case action.Shift:
			n := len("shift(") + nbytes(int(act1)) + len("),")
			if n > max {
				max = n
			}
		default:
			panic(fmt.Sprintf("Unknown action type: %T", act1))
		}
	}
	for i, sym := range tokMap.TypeMap {
		act, symConflicts := set.Action(sym)
		if len(symConflicts) > 0 {
			conflicts[sym] = symConflicts
		}
		switch act1 := act.(type) {
		case action.Accept:
			pad := max + 1 - len("accept(true),")
			data.Actions[i] = fmt.Sprintf("accept(true),%*c// %s", pad, ' ', sym)
		case action.Error:
			pad := max + 1 - len("nil,")
			data.Actions[i] = fmt.Sprintf("nil,%*c// %s", pad, ' ', sym)
		case action.Reduce:
			pad := max + 1 - (len("reduce(") + nbytes(int(act1)) + len("),"))
			data.Actions[i] = fmt.Sprintf("reduce(%d),%*c// %s, reduce: %s", int(act1), pad, ' ', sym, prods[int(act1)].Id)
		case action.Shift:
			pad := max + 1 - (len("shift(") + nbytes(int(act1)) + len("),"))
			data.Actions[i] = fmt.Sprintf("shift(%d),%*c// %s", int(act1), pad, ' ', sym)
		default:
			panic(fmt.Sprintf("Unknown action type: %T", act1))
		}
	}
	return
}

const actionTableSrc = `
// Code generated by gocc; DO NOT EDIT.

package parser

type (
	actionTable [numStates]actionRow
	actionRow   struct {
		canRecover bool
		actions    [numSymbols]action
	}
)

var actionTab = actionTable{
	{{- range $i, $r := .Rows }}
	actionRow{ // S{{$i}}
		canRecover: {{printf "%t" .CanRecover}},
		actions: [numSymbols]action{
			{{- range $a := .Actions }}
			{{$a}}
			{{- end }}
		},
	},
	{{- end }}
}
`

func GenCompActionTable(outDir string, prods ast.SyntaxProdList, itemSets *items.ItemSets, tokMap *token.TokenMap) map[int]items.RowConflicts {
	tab := make([]struct {
		CanRecover bool
		Actions    []struct {
			Index  int
			Action int
			Amount int
		}
	}, itemSets.Size())
	conflictss := make(map[int]items.RowConflicts)
	for i := range tab {
		set := itemSets.Set(i)
		tab[i].CanRecover = set.CanRecover()
		conflicts := make(items.RowConflicts)
		for j, sym := range tokMap.TypeMap {
			act, symConflicts := set.Action(sym)
			if len(symConflicts) > 0 {
				conflicts[sym] = symConflicts
			}
			switch act1 := act.(type) {
			case action.Accept:
				tab[i].Actions = append(tab[i].Actions, struct {
					Index  int
					Action int
					Amount int
				}{Index: j, Action: 0, Amount: 0})
			case action.Error:
				// skip
			case action.Reduce:
				tab[i].Actions = append(tab[i].Actions, struct {
					Index  int
					Action int
					Amount int
				}{Index: j, Action: 1, Amount: int(act1)})
			case action.Shift:
				tab[i].Actions = append(tab[i].Actions, struct {
					Index  int
					Action int
					Amount int
				}{Index: j, Action: 2, Amount: int(act1)})
			default:
				panic(fmt.Sprintf("Unknown action type: %T", act1))
			}
		}
		if len(conflicts) > 0 {
			conflictss[i] = conflicts
		}
	}
	bytesStr := genEnc(tab)
	tmpl, err := template.New("parser action table").Parse(actionCompTableSrc[1:])
	if err != nil {
		panic(err)
	}
	wr := new(bytes.Buffer)
	if err := tmpl.Execute(wr, bytesStr); err != nil {
		panic(err)
	}
	io.WriteFile(path.Join(outDir, "parser", "actiontable.go"), wr.Bytes())
	return conflictss
}

const actionCompTableSrc = `
// Code generated by gocc; DO NOT EDIT.

package parser

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
)

type (
	actionTable [numStates]actionRow
	actionRow   struct {
		canRecover bool
		actions    [numSymbols]action
	}
)

var actionTab = actionTable{}

func init() {
	tab := []struct {
		CanRecover bool
		Actions    []struct {
			Index  int
			Action int
			Amount int
		}
	}{}
	data := {{.}}
	buf, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&tab); err != nil {
		panic(err)
	}

	for i, row := range tab {
		actionTab[i].canRecover = row.CanRecover
		for _, a := range row.Actions {
			switch a.Action {
			case 0:
				actionTab[i].actions[a.Index] = accept(true)
			case 1:
				actionTab[i].actions[a.Index] = reduce(a.Amount)
			case 2:
				actionTab[i].actions[a.Index] = shift(a.Amount)
			}
		}
	}
}
`

// nbytes returns the number of bytes required to output the integer x.
func nbytes(x int) int {
	if x == 0 {
		return 1
	}
	n := 0
	if x < 0 {
		x = -x
		n++
	}
	n += int(math.Log10(float64(x))) + 1
	return n
}
