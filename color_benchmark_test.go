package color

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

var (
	_tmpbuf bytes.Buffer
	_tmpstr string
)

func init() {
	_hasColor = true
}

func BenchmarkColor_Escape(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FgCyan.Escape()
	}
}

func BenchmarkColor_Reset(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FgCyan.Reset()
	}
}

func BenchmarkColor_Wrap(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_tmpstr = FgCyan.Wrap("x")
	}
}

func BenchmarkColor_Print(b *testing.B) {
	stdout := _stdout
	_stdout = bufio.NewWriter(&_tmpbuf)
	defer func() {
		_stdout = stdout
	}()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FgCyan.Print("x")
	}
}

func BenchmarkColor_Fprint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FgCyan.Fprint(io.Discard, "x") //nolint:errcheck
	}
}

func BenchmarkFprint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Fprint(io.Discard, FgCyan, "x") //nolint:errcheck
	}
}
