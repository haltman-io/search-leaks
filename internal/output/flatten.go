package output

import (
	"fmt"
	"sort"
	"strings"
)

type FlatLine struct {
	Brackets []string // e.g. ["stealer(1)"] or ["stealer(2)", "credentials(1)"]
	Key      string
	Value    string
}

func FlattenJSON(v any) []FlatLine {
	var out []FlatLine
	walk(&out, nil, "", v)
	return out
}

func walk(out *[]FlatLine, brackets []string, keyPath string, v any) {
	switch vv := v.(type) {
	case map[string]any:
		// stable order
		keys := make([]string, 0, len(vv))
		for k := range vv {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			child := vv[k]
			nextPath := k
			if keyPath != "" {
				nextPath = keyPath + "." + k
			}
			walk(out, brackets, nextPath, child)
		}

	case []any:
		// If array of primitives, print as one line
		if isPrimitiveArray(vv) {
			vals := make([]string, 0, len(vv))
			for _, it := range vv {
				vals = append(vals, fmtValue(it))
			}
			*out = append(*out, FlatLine{
				Brackets: brackets,
				Key:      keyPath,
				Value:    strings.Join(vals, ", "),
			})
			return
		}

		// Array of objects (or mixed): emit bracket segments with index
		base := arrayLabelFromKey(keyPath)
		for i, it := range vv {
			lbl := fmt.Sprintf("%s(%d)", base, i+1)
			nextBrackets := append(brackets, lbl)
			// For object items, we keep keyPath empty so keys don't repeat "stealers.x"
			switch it.(type) {
			case map[string]any:
				walk(out, nextBrackets, "", it)
			default:
				*out = append(*out, FlatLine{
					Brackets: nextBrackets,
					Key:      keyPath,
					Value:    fmtValue(it),
				})
			}
		}

	default:
		// Leaf
		if keyPath == "" {
			// Happens if an array object had primitive leaf with empty path; keep generic key
			keyPath = "value"
		}
		*out = append(*out, FlatLine{
			Brackets: brackets,
			Key:      keyPath,
			Value:    fmtValue(v),
		})
	}
}

func fmtValue(v any) string {
	if v == nil {
		return "null"
	}
	switch t := v.(type) {
	case string:
		return t
	case float64:
		// JSON numbers decode as float64
		if t == float64(int64(t)) {
			return fmt.Sprintf("%d", int64(t))
		}
		return fmt.Sprintf("%g", t)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", t)
	}
}

func isPrimitiveArray(a []any) bool {
	for _, it := range a {
		switch it.(type) {
		case map[string]any, []any:
			return false
		}
	}
	return true
}

func arrayLabelFromKey(keyPath string) string {
	// If keyPath is "stealers" => "stealer" (best-effort for nice output)
	// If nested path => take last segment
	last := keyPath
	if strings.Contains(keyPath, ".") {
		parts := strings.Split(keyPath, ".")
		last = parts[len(parts)-1]
	}
	last = strings.TrimSpace(last)
	if strings.HasSuffix(last, "s") && len(last) > 1 {
		return strings.TrimSuffix(last, "s")
	}
	return last
}
