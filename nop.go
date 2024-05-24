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
	"fmt"
	"io"
	"strings"
)

// Nop is a passthrough [Style] that does not add color.
var Nop Style = nop{}

type nop struct{}

func (nop) Code() string   { return "" }
func (nop) Escape() string { return "" }
func (nop) Reset() string  { return "" }
func (nop) String() string { return "" }

func (nop) With(...Style) Style         { return nop{} }
func (nop) Wrap(str string) string      { return str }
func (nop) WrapN(strs ...string) string { return strings.Join(strs, " ") }

func (nop) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

func (nop) Print(args ...any)              { fmt.Print(args...) }
func (nop) Printf(msg string, args ...any) { fmt.Printf(msg, args...) }
func (nop) Println(args ...any)            { fmt.Println(args...) }

func (nop) Sprint(args ...any) string              { return fmt.Sprint(args...) }
func (nop) Sprintf(msg string, args ...any) string { return fmt.Sprintf(msg, args...) }
func (nop) Sprintln(args ...any) string            { return fmt.Sprintln(args...) }

func (nop) Fprint(dst io.Writer, args ...any) (int, error) { return fmt.Fprint(dst, args...) }
func (nop) Fprintf(dst io.Writer, msg string, args ...any) (int, error) {
	return fmt.Fprintf(dst, msg, args...)
}
func (nop) Fprintln(dst io.Writer, args ...any) (int, error) { return fmt.Fprintln(dst, args...) }
