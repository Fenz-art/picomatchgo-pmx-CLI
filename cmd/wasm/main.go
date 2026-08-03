//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	picomatch "github.com/debayansamal/port-mortem-picomatch-go"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("picomatchScan", js.FuncOf(picomatchScan))
	js.Global().Set("picomatchParse", js.FuncOf(picomatchParse))
	js.Global().Set("picomatchIsMatch", js.FuncOf(picomatchIsMatch))
	js.Global().Set("picomatchCompile", js.FuncOf(picomatchCompile))

	<-c
}

func picomatchScan(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "error: missing pattern"
	}
	pattern := args[0].String()

	var opts *picomatch.Options
	if len(args) > 1 && !args[1].IsNull() && !args[1].IsUndefined() {
		optsJSON := args[1].String()
		var o picomatch.Options
		if err := json.Unmarshal([]byte(optsJSON), &o); err == nil {
			opts = &o
		}
	}

	state := picomatch.Scan(pattern, opts)

	res := map[string]interface{}{
		"prefix":         state.Prefix,
		"input":          state.Input,
		"start":          state.Start,
		"base":           state.Base,
		"glob":           state.Glob,
		"isBrace":        state.IsBrace,
		"isBracket":      state.IsBracket,
		"isGlob":         state.IsGlob,
		"isExtglob":      state.IsExtglob,
		"isGlobstar":     state.IsGlobstar,
		"negated":        state.Negated,
		"negatedExtglob": state.NegatedExtglob,
		"maxDepth":       state.MaxDepth,
		"slashes":        state.Slashes,
		"parts":          state.Parts,
	}

	resJSON, err := json.Marshal(res)
	if err != nil {
		return "error: " + err.Error()
	}
	return string(resJSON)
}

func picomatchParse(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "error: missing pattern"
	}
	pattern := args[0].String()

	var opts *picomatch.Options
	if len(args) > 1 && !args[1].IsNull() && !args[1].IsUndefined() {
		optsJSON := args[1].String()
		var o picomatch.Options
		if err := json.Unmarshal([]byte(optsJSON), &o); err == nil {
			opts = &o
		}
	}

	state, err := picomatch.Parse(pattern, opts)
	if err != nil {
		return "error: " + err.Error()
	}

	res := map[string]interface{}{
		"input":   state.Input,
		"output":  state.Output,
		"negated": state.Negated,
	}

	resJSON, err := json.Marshal(res)
	if err != nil {
		return "error: " + err.Error()
	}
	return string(resJSON)
}

func picomatchIsMatch(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return false
	}
	input := args[0].String()
	pattern := args[1].String()

	var opts *picomatch.Options
	if len(args) > 2 && !args[2].IsNull() && !args[2].IsUndefined() {
		optsJSON := args[2].String()
		var o picomatch.Options
		if err := json.Unmarshal([]byte(optsJSON), &o); err == nil {
			opts = &o
		}
	}

	matched, err := picomatch.IsMatch(input, pattern, opts)
	if err != nil {
		return false
	}
	return matched
}

func picomatchCompile(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "error: missing pattern"
	}
	pattern := args[0].String()

	var opts *picomatch.Options
	if len(args) > 1 && !args[1].IsNull() && !args[1].IsUndefined() {
		optsJSON := args[1].String()
		var o picomatch.Options
		if err := json.Unmarshal([]byte(optsJSON), &o); err == nil {
			opts = &o
		}
	}

	re, err := picomatch.MakeRe(pattern, opts)
	if err != nil {
		return "error: " + err.Error()
	}
	return re.String()
}
