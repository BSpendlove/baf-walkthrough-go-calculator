package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/BSpendlove/baf-walkthrough-go-calculator/internal/history"
	"github.com/BSpendlove/baf-walkthrough-go-calculator/internal/parser"
)

// Version can be overridden at build time with -ldflags "-X main.Version=..."
var Version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: calc <expression> | calc history | calc clear | calc --version | calc --whoami")
		os.Exit(1)
	}

	arg := os.Args[1]

	switch {
	case arg == "--version" || arg == "--whoami":
		fmt.Printf("calc v%s\n", Version)

	case arg == "history":
		entries, err := history.Load(history.DefaultPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if len(entries) == 0 {
			fmt.Println("No history.")
			return
		}
		for _, e := range entries {
			fmt.Printf("  %s = %s\n", e.Expr, formatResult(e.Result))
		}

	case arg == "clear":
		if err := history.Clear(history.DefaultPath); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("History cleared.")

	default:
		expr := strings.Join(os.Args[1:], " ")
		result, err := parser.Eval(expr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(formatResult(result))
		if err := history.Store(history.DefaultPath, expr, result); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func formatResult(f float64) string {
	if f == math.Trunc(f) && !math.IsInf(f, 0) {
		return fmt.Sprintf("%g", f)
	}
	return fmt.Sprintf("%g", f)
}
