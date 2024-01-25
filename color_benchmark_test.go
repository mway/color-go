// Copyright (c) 2024 Matt Way
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE THE SOFTWARE.

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
