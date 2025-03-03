package shared

import (
	"encoding/json"
	"fmt"
	"os"
)

var Flags struct {
	Json bool
}

func PrintProgress(msg string) {
	if !Flags.Json {
		fmt.Printf("* %s\n", msg)
	}
}

func PrintItem(msg string) {
	if !Flags.Json {
		fmt.Printf("  - %s\n", msg)
	}
}

func PrintIndented(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func ExitWithResult(result any, key string, message string) {
	if Flags.Json {
		PrintIndented(map[string]any{"ok": true, key: result})
	} else {
		b, _ := json.Marshal(result)
		fmt.Printf("* %s: %s\n", message, string(b))
	}
	os.Exit(0)
}

func ExitWithMap(result map[string]any, message string) {
	if Flags.Json {
		result["ok"] = true
		PrintIndented(result)
	} else {
		b, _ := json.Marshal(result)
		fmt.Printf("* %s: %s\n", message, string(b))
	}
	os.Exit(0)
}

func ExitWithSlice[T any](results []T, key, noun, verb, singular, plural string) {
	if Flags.Json {
		PrintIndented(map[string]any{"ok": true, "count": len(results), key: results})
	} else {
		if len(results) == 1 {
			fmt.Printf("* %s 1 %s\n", verb, singular)
		} else {
			fmt.Printf("* %s %d %s\n", verb, len(results), plural)
		}
		for i, result := range results {
			b, _ := json.Marshal(result)
			fmt.Printf("  - %s %d: %s\n", noun, i, string(b))
		}
	}
	os.Exit(0)
}

func ExitWithStatus(err error) {
	if Flags.Json {
		if err != nil {
			PrintIndented(map[string]any{"ok": false, "error": err.Error()})
		} else {
			PrintIndented(map[string]any{"ok": true})
		}
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "* Error: %s\n", err.Error())
	}
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
