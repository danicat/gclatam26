package levels

import (
	"regexp"
	"strings"
)

// Level represents a single bug-fixing puzzle level.
type Level struct {
	ID              int
	Title           string
	BugCategory     string
	PanicMessage    string
	Hint            string
	Explanation     string
	CodeLines       []string
	TargetLineIndex int
	TimeLimit       float64 // In seconds
	FallSpeed       float64 // Vertical pixels per second
	Validate        func(lineIdx int, text string) bool
}

// CleanCode removes leading/trailing whitespace and normalizes spaces for comparison.
func CleanCode(s string) string {
	s = strings.TrimSpace(s)
	// Replace multiple spaces with a single space
	spaceRegex := regexp.MustCompile(`\s+`)
	return spaceRegex.ReplaceAllString(s, " ")
}

// AllLevels contains all 10 game jam levels with increasing difficulty.
var AllLevels = []Level{
	{
		ID:           1,
		Title:        "Nil Pointer Dereference",
		BugCategory:  "Pointer / Allocation",
		PanicMessage: "panic: runtime error: invalid memory address or nil pointer dereference",
		Hint:         "Hint: Initialize 'gopher' with &Gopher{} or allocate before accessing Score.",
		Explanation:  "Dereferencing a nil pointer causes an immediate runtime crash.",
		TimeLimit:    40.0,
		FallSpeed:    5.0,
		TargetLineIndex: 2,
		CodeLines: []string{
			"func setupPlayer() {",
			"    var gopher *Gopher",
			"    gopher.Score = 100",
			"    println(gopher.Score)",
			"}",
		},
		Validate: func(lineIdx int, text string) bool {
			c := CleanCode(text)
			if lineIdx == 1 {
				// Player fixed line 1: var gopher *Gopher
				return strings.Contains(c, "&Gopher") ||
					strings.Contains(c, "new(Gopher)") ||
					strings.Contains(c, "gopher := &Gopher") ||
					strings.Contains(c, "gopher = &Gopher")
			}
			if lineIdx == 2 {
				// Player fixed line 2: gopher.Score = 100
				return strings.Contains(c, "&Gopher") ||
					strings.Contains(c, "new(Gopher)") ||
					strings.Contains(c, "if gopher != nil") ||
					strings.Contains(c, "gopher = &Gopher{Score: 100}") ||
					strings.Contains(c, "gopher := &Gopher{Score: 100}")
			}
			return false
		},
	},
	{
		ID:           2,
		Title:        "Integer Division by Zero",
		BugCategory:  "Math / Runtime",
		PanicMessage: "panic: runtime error: integer divide by zero",
		Hint:         "Hint: Guard against count == 0 before division, or ensure count > 0.",
		Explanation:  "Integer division by zero is not NaN/Inf in Go; it is an unrecoverable runtime panic unless caught.",
		TimeLimit:    38.0,
		FallSpeed:    6.0,
		TargetLineIndex: 2,
		CodeLines: []string{
			"func calcRatio(total, count int) int {",
			"    // count is currently 0!",
			"    ratio := total / count",
			"    return ratio",
			"}",
		},
		Validate: func(lineIdx int, text string) bool {
			c := CleanCode(text)
			if lineIdx == 2 {
				return strings.Contains(c, "count + 1") ||
					strings.Contains(c, "count != 0") ||
					strings.Contains(c, "if count") ||
					strings.Contains(c, "count = 1") ||
					strings.Contains(c, "ratio := 0") ||
					strings.Contains(c, "return 0")
			}
			if lineIdx == 1 {
				return strings.Contains(c, "count = 1") || strings.Contains(c, "if count == 0")
			}
			return false
		},
	},
	{
		ID:           3,
		Title:        "Slice Index Out of Range",
		BugCategory:  "Slice / Bounds Check",
		PanicMessage: "panic: runtime error: index out of range [3] with length 3",
		Hint:         "Hint: Go slices are 0-indexed. The last element is at len(items)-1.",
		Explanation:  "Accessing slice[len(slice)] always triggers an index out of range panic.",
		TimeLimit:    35.0,
		FallSpeed:    7.0,
		TargetLineIndex: 2,
		CodeLines: []string{
			"func fetchLast(items []string) string {",
			"    n := len(items)",
			"    last := items[n]",
			"    return last",
			"}",
		},
		Validate: func(lineIdx int, text string) bool {
			c := CleanCode(text)
			if lineIdx == 2 {
				return strings.Contains(c, "items[n-1]") ||
					strings.Contains(c, "items[n - 1]") ||
					strings.Contains(c, "items[len(items)-1]") ||
					strings.Contains(c, "items[len(items) - 1]") ||
					strings.Contains(c, "if n > 0")
			}
			return false
		},
	},
	{
		ID:           4,
		Title:        "Assignment to Nil Map",
		BugCategory:  "Map / Allocation",
		PanicMessage: "panic: assignment to entry in nil map",
		Hint:         "Hint: A declared map is nil. Allocate it with make(map[string]int) first.",
		Explanation:  "Reading from a nil map returns the zero-value, but writing to a nil map always panics.",
		TimeLimit:    33.0,
		FallSpeed:    8.0,
		TargetLineIndex: 1,
		CodeLines: []string{
			"func initInventory() {",
			"    var inventory map[string]int",
			"    inventory[\"potions\"] = 5",
			"    println(\"Inventory ready\")",
			"}",
		},
		Validate: func(lineIdx int, text string) bool {
			c := CleanCode(text)
			if lineIdx == 1 {
				return strings.Contains(c, "make(map[string]int)") ||
					strings.Contains(c, "map[string]int{}")
			}
			if lineIdx == 2 {
				return strings.Contains(c, "make(map[string]int)") ||
					strings.Contains(c, "inventory = make(")
			}
			return false
		},
	},
	{
		ID:           5,
		Title:        "Send on Closed Channel",
		BugCategory:  "Concurrency / Channels",
		PanicMessage: "panic: send on closed channel",
		Hint:         "Hint: Send the value before closing the channel, or defer closing.",
		Explanation:  "Sending to a closed channel panics immediately. Only sender should close channels.",
		TimeLimit:    30.0,
		FallSpeed:    9.0,
		TargetLineIndex: 2,
		CodeLines: []string{
			"func notify(ch chan int, msg int) {",
			"    close(ch)",
			"    ch <- msg",
			"    println(\"Notification sent\")",
			"}",
		},
		Validate: func(lineIdx int, text string) bool {
			c := CleanCode(text)
			if lineIdx == 1 {
				// Swap order or remove close(ch)
				return strings.Contains(c, "ch <- msg") ||
					strings.Contains(c, "defer close(ch)") ||
					strings.HasPrefix(c, "//")
			}
			if lineIdx == 2 {
				return strings.Contains(c, "close(ch)") ||
					strings.HasPrefix(c, "//") ||
					strings.Contains(c, "select")
			}
			return false
		},
	},
	{
		ID:           6,
		Title:        "Double Channel Close",
		BugCategory:  "Concurrency / Channels",
		PanicMessage: "panic: close of closed channel",
		Hint:         "Hint: Do not call close() twice on the same channel; remove the second close.",
		Explanation:  "Closing an already closed channel causes an immediate runtime panic in Go.",
		TimeLimit:    28.0,
		FallSpeed:    10.0,
		TargetLineIndex: 3,
		CodeLines: []string{
			"func shutdown(ch chan bool) {",
			"    close(ch)",
			"    // Redundant cleanup:",
			"    close(ch)",
			"}",
		},
		Validate: func(lineIdx int, text string) bool {
			c := CleanCode(text)
			if lineIdx == 3 {
				return strings.HasPrefix(c, "//") ||
					c == "" ||
					strings.Contains(c, "return") ||
					!strings.Contains(c, "close(ch)")
			}
			return false
		},
	},
	{
		ID:           7,
		Title:        "Infinite Recursion (Stack Overflow)",
		BugCategory:  "Stack / Call Graph",
		PanicMessage: "runtime: goroutine stack exceeds 1000000000-byte limit (fatal error)",
		Hint:         "Hint: Add a base case (e.g. if n <= 0 { return }) before the recursive call.",
		Explanation:  "Recursive calls without an exit condition exhaust the goroutine stack memory.",
		TimeLimit:    26.0,
		FallSpeed:    11.0,
		TargetLineIndex: 1,
		CodeLines: []string{
			"func countdown(n int) {",
			"    // missing base case!",
			"    println(n)",
			"    countdown(n - 1)",
			"}",
		},
		Validate: func(lineIdx int, text string) bool {
			c := CleanCode(text)
			if lineIdx == 1 || lineIdx == 3 {
				return strings.Contains(c, "if n <= 0") ||
					strings.Contains(c, "if n <=") ||
					strings.Contains(c, "if n == 0") ||
					strings.Contains(c, "if n > 0") ||
					strings.Contains(c, "return")
			}
			return false
		},
	},
	{
		ID:           8,
		Title:        "Unchecked Type Assertion",
		BugCategory:  "Type System / Interface",
		PanicMessage: "panic: interface conversion: interface {} is int, not string",
		Hint:         "Hint: Use the two-value type assertion: val, ok := data.(string).",
		Explanation:  "Single-value type assertion panics if the underlying type does not match.",
		TimeLimit:    24.0,
		FallSpeed:    12.0,
		TargetLineIndex: 1,
		CodeLines: []string{
			"func parseName(data any) string {",
			"    name := data.(string)",
			"    return strings.ToUpper(name)",
			"}",
		},
		Validate: func(lineIdx int, text string) bool {
			c := CleanCode(text)
			if lineIdx == 1 {
				return strings.Contains(c, ", ok") ||
					strings.Contains(c, ",ok") ||
					strings.Contains(c, "fmt.Sprint") ||
					strings.Contains(c, "switch")
			}
			return false
		},
	},
	{
		ID:           9,
		Title:        "Concurrent Map Write (Race Condition)",
		BugCategory:  "Concurrency / Data Race",
		PanicMessage: "fatal error: concurrent map writes",
		Hint:         "Hint: Synchronize map access with a sync.Mutex or use sync.Map.",
		Explanation:  "Go runtime detects unsynchronized concurrent map writes and panics with a fatal crash.",
		TimeLimit:    22.0,
		FallSpeed:    13.0,
		TargetLineIndex: 2,
		CodeLines: []string{
			"func hitCounter(hits map[string]int, key string) {",
			"    // Concurrent calls without mutex lock!",
			"    hits[key]++",
			"}",
		},
		Validate: func(lineIdx int, text string) bool {
			c := CleanCode(text)
			if lineIdx == 1 || lineIdx == 2 {
				return strings.Contains(c, "mu.Lock") ||
					strings.Contains(c, "lock") ||
					strings.Contains(c, "sync") ||
					strings.Contains(c, "atomic")
			}
			return false
		},
	},
	{
		ID:           10,
		Title:        "Unhandled Panic in Goroutine",
		BugCategory:  "Runtime / Defer & Recover",
		PanicMessage: "panic: critical payment worker crash (unhandled in goroutine)",
		Hint:         "Hint: Use defer func() { recover() }() to catch panics before exit.",
		Explanation:  "A panic inside a goroutine crashes the entire Go process unless caught by defer recover().",
		TimeLimit:    20.0,
		FallSpeed:    14.0,
		TargetLineIndex: 1,
		CodeLines: []string{
			"func paymentWorker() {",
			"    // missing defer recover!",
			"    panic(\"payment worker failed\")",
			"}",
		},
		Validate: func(lineIdx int, text string) bool {
			c := CleanCode(text)
			if lineIdx == 1 || lineIdx == 2 {
				return strings.Contains(c, "recover()") ||
					strings.Contains(c, "defer") ||
					strings.HasPrefix(c, "//")
			}
			return false
		},
	},
}
