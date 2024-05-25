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

// Package color provides color-based types and helpers.
package color

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/mattn/go-isatty"
	"go.mway.dev/errors"
	"go.mway.dev/pool"
)

// Style colors.
const (
	Reset Color = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground colors.
const (
	FgBlack Color = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground high-intensity colors.
const (
	FgHiBlack Color = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background colors.
const (
	BgBlack Color = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background high-intensity colors.
const (
	BgHiBlack Color = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

var (
	// ErrInvalidColorName is returned when attempting to parse a color using
	// an unknown name.
	ErrInvalidColorName = errors.New("invalid color name")

	_fd       = os.Stderr.Fd()
	_stdout   = bufio.NewWriter(os.Stdout)
	_hasColor = !isset("NO_COLOR") && os.Getenv("TERM") != "dumb" &&
		(isatty.IsTerminal(_fd) || isatty.IsCygwinTerminal(_fd))
	_builders = pool.NewWithReleaser(
		func() *bytes.Buffer { return bytes.NewBuffer(make([]byte, 0, 256)) },
		func(x *bytes.Buffer) { x.Reset() },
	)
	_strings = [math.MaxUint8]string{
		// Reset:        "\x1b[0m",
		Reset:        "\x1b[0m",
		Bold:         "\x1b[1m",
		Faint:        "\x1b[2m",
		Italic:       "\x1b[3m",
		Underline:    "\x1b[4m",
		BlinkSlow:    "\x1b[5m",
		BlinkRapid:   "\x1b[6m",
		ReverseVideo: "\x1b[7m",
		Concealed:    "\x1b[8m",
		CrossedOut:   "\x1b[9m",
		FgBlack:      "\x1b[30m",
		FgRed:        "\x1b[31m",
		FgGreen:      "\x1b[32m",
		FgYellow:     "\x1b[33m",
		FgBlue:       "\x1b[34m",
		FgMagenta:    "\x1b[35m",
		FgCyan:       "\x1b[36m",
		FgWhite:      "\x1b[37m",
		FgHiBlack:    "\x1b[90m",
		FgHiRed:      "\x1b[91m",
		FgHiGreen:    "\x1b[92m",
		FgHiYellow:   "\x1b[93m",
		FgHiBlue:     "\x1b[94m",
		FgHiMagenta:  "\x1b[95m",
		FgHiCyan:     "\x1b[96m",
		FgHiWhite:    "\x1b[97m",
		BgBlack:      "\x1b[40m",
		BgRed:        "\x1b[41m",
		BgGreen:      "\x1b[42m",
		BgYellow:     "\x1b[43m",
		BgBlue:       "\x1b[44m",
		BgMagenta:    "\x1b[45m",
		BgCyan:       "\x1b[46m",
		BgWhite:      "\x1b[47m",
		BgHiBlack:    "\x1b[100m",
		BgHiRed:      "\x1b[101m",
		BgHiGreen:    "\x1b[102m",
		BgHiYellow:   "\x1b[103m",
		BgHiBlue:     "\x1b[104m",
		BgHiMagenta:  "\x1b[105m",
		BgHiCyan:     "\x1b[106m",
		BgHiWhite:    "\x1b[107m",
	}
	_fgNames = map[string]Color{
		"black":      FgBlack,
		"red":        FgRed,
		"green":      FgGreen,
		"yellow":     FgYellow,
		"blue":       FgBlue,
		"magenta":    FgMagenta,
		"cyan":       FgCyan,
		"white":      FgWhite,
		"hi-black":   FgHiBlack,
		"hi-red":     FgHiRed,
		"hi-green":   FgHiGreen,
		"hi-yellow":  FgHiYellow,
		"hi-blue":    FgHiBlue,
		"hi-magenta": FgHiMagenta,
		"hi-cyan":    FgHiCyan,
		"hi-white":   FgHiWhite,
	}
	_bgNames = map[string]Color{
		"black":      BgBlack,
		"red":        BgRed,
		"green":      BgGreen,
		"yellow":     BgYellow,
		"blue":       BgBlue,
		"magenta":    BgMagenta,
		"cyan":       BgCyan,
		"white":      BgWhite,
		"hi-black":   BgHiBlack,
		"hi-red":     BgHiRed,
		"hi-green":   BgHiGreen,
		"hi-yellow":  BgHiYellow,
		"hi-blue":    BgHiBlue,
		"hi-magenta": BgHiMagenta,
		"hi-cyan":    BgHiCyan,
		"hi-white":   BgHiWhite,
	}

	_ Style = Color(0)
	_ Style = multiStyle{}
)

// A Style styles text.
type Style interface {
	Code() string
	Copy(io.Writer, io.Reader) (int64, error)
	Escape() string
	Fprint(io.Writer, ...any) (int, error)
	Fprintf(io.Writer, string, ...any) (int, error)
	Fprintln(io.Writer, ...any) (int, error)
	Join([]string, string) string
	Print(...any)
	Printf(string, ...any)
	Println(...any)
	Reset() string
	Sprint(...any) string
	Sprintf(string, ...any) string
	Sprintln(...any) string
	String() string
	With(...Style) Style
	Wrap(string) string
}

// A Color is a terminal color.
type Color uint8

// ParseFgColor parses the given name into a foreground [Color].
func ParseFgColor(name string) (Color, error) {
	return parseColor(_fgNames, name)
}

// ParseBgColor parses the given name into a background [Color].
func ParseBgColor(name string) (Color, error) {
	return parseColor(_bgNames, name)
}

// Escape returns c's escape code.
func (c Color) Escape() string {
	if !Enabled() {
		return ""
	}

	return _strings[c]
}

// Reset returns the escape code the reset output after c.
func (c Color) Reset() string {
	if c == Reset || !Enabled() {
		return ""
	}
	return Reset.Escape()
}

// String returns c's escape code, regardless of whether color is enabled.
func (c Color) String() string {
	return _strings[c]
}

// With returns a [Style] with the given styles amended to the current color.
func (c Color) With(styles ...Style) Style {
	switch len(styles) {
	case 0:
		return c
	case 1:
		if x, ok := styles[0].(Color); ok && x == c {
			return c
		}
	default:
	}

	return newMultiStyle(append([]Style{c}, styles...)...)
}

// Join joins each string (wrapped in this color) with the given delimeter.
func (c Color) Join(elems []string, sep string) string {
	buf := _builders.Get()
	defer _builders.Put(buf)

	for i, str := range elems {
		if i > 0 {
			buf.WriteString(sep) //nolint:errcheck
		}
		buf.WriteString(c.Wrap(str)) //nolint:errcheck
	}

	return buf.String()
}

// Wrap wraps str with c.
func (c Color) Wrap(str string) string {
	return c.Escape() + str + c.Reset()
}

// Code returns the ANSI code related to this color.
func (c Color) Code() string {
	x := _strings[c]
	if len(x) == 0 {
		return ""
	}
	return x[2 : len(x)-1]
}

// Copy copies src to dest as in io.Copy, but wrapped in c.
func (c Color) Copy(dst io.Writer, src io.Reader) (int64, error) {
	buf := _builders.Get()
	defer _builders.Put(buf)

	buf.WriteString(c.Escape()) //nolint:errcheck

	n, err := dst.Write(buf.Bytes())
	if err != nil {
		return int64(n), err
	}

	n64, err := io.Copy(dst, src)
	if err != nil {
		return int64(n) + n64, err
	}
	n += int(n64)

	buf.Reset()
	buf.WriteString(c.Reset()) //nolint:errcheck

	m, err := dst.Write(buf.Bytes())
	return int64(n + m), err
}

// Print prints args as in fmt.Print, but wrapped in c.
func (c Color) Print(args ...any) {
	defer _stdout.Flush() //nolint:errcheck

	_stdout.WriteString(c.Escape()) //nolint:errcheck
	fmt.Fprint(_stdout, args...)    //nolint:errcheck
	_stdout.WriteString(c.Reset())  //nolint:errcheck
}

// Printf prints msg and args as in fmt.Printf, but wrapped in c.
func (c Color) Printf(msg string, args ...any) {
	defer _stdout.Flush() //nolint:errcheck

	_stdout.WriteString(c.Escape())    //nolint:errcheck
	fmt.Fprintf(_stdout, msg, args...) //nolint:errcheck
	_stdout.WriteString(c.Reset())     //nolint:errcheck
}

// Println prints args as in fmt.Println, but wrapped in c.
func (c Color) Println(args ...any) {
	defer _stdout.Flush() //nolint:errcheck

	buf := _builders.Get()
	defer _builders.Put(buf)

	_stdout.WriteString(c.Escape()) //nolint:errcheck
	fmt.Fprintln(buf, args...)      //nolint:errcheck

	var (
		tmp = buf.Bytes()
		esc = c.Reset()
	)

	if len(esc) > 0 {
		tmp[len(tmp)-1] = esc[0]
		buf.WriteString(esc[1:]) //nolint:errcheck
		buf.WriteByte('\n')      //nolint:errcheck
	}

	_stdout.WriteString(buf.String()) //nolint:errcheck
}

// Sprint returns a string containing args as in fmt.Sprint, but wrapped in c.
func (c Color) Sprint(args ...any) string {
	return c.Escape() + fmt.Sprint(args...) + c.Reset()
}

// Sprintf returns a string containing msg and args as in fmt.Sprintf, but
// wrapped in c.
func (c Color) Sprintf(msg string, args ...any) string {
	return c.Escape() + fmt.Sprintf(msg, args...) + c.Reset()
}

// Sprintln returns a string containing args as in fmt.Sprintln, but wrapped
// in c.
func (c Color) Sprintln(args ...any) string {
	buf := _builders.Get()
	defer _builders.Put(buf)

	buf.WriteString(c.Escape()) //nolint:errcheck
	fmt.Fprintln(buf, args...)  //nolint:errcheck

	var (
		tmp = buf.Bytes()
		esc = c.Reset()
	)

	if len(esc) > 0 {
		tmp[len(tmp)-1] = esc[0]
		buf.WriteString(esc[1:]) //nolint:errcheck
		buf.WriteByte('\n')      //nolint:errcheck
	}

	return buf.String()
}

// Fprint prints args to w as in fmt.Fprint, but wrapped in c.
func (c Color) Fprint(w io.Writer, args ...any) (int, error) {
	esc := c.Escape()
	fmt.Fprint(w, esc)             //nolint:errcheck
	n, _ := fmt.Fprint(w, args...) //nolint:errcheck
	m, err := fmt.Fprint(w, c.Reset())
	return n + m + len(esc), err
}

// Fprintf prints msg and args to w as in [fmt.Fprintf], but wrapped in c.
func (c Color) Fprintf(w io.Writer, msg string, args ...any) (int, error) {
	esc := c.Escape()
	fmt.Fprint(w, esc)                   //nolint:errcheck
	n, _ := fmt.Fprintf(w, msg, args...) //nolint:errcheck
	m, err := fmt.Fprint(w, c.Reset())
	return n + m + len(esc), err
}

// Fprintln prints args to w as in [fmt.Fprintln], but wrapped in c.
func (c Color) Fprintln(w io.Writer, args ...any) (int, error) {
	buf := _builders.Get()
	defer _builders.Put(buf)

	buf.WriteString(c.Escape()) //nolint:errcheck
	fmt.Fprintln(buf, args...)  //nolint:errcheck

	var (
		tmp = buf.Bytes()
		esc = c.Reset()
	)

	if len(esc) > 0 {
		tmp[len(tmp)-1] = esc[0]
		buf.WriteString(esc[1:]) //nolint:errcheck
		buf.WriteByte('\n')      //nolint:errcheck
	}

	n, err := io.Copy(w, buf)
	return int(n), err
}

type multiStyle struct {
	escape string
	styles []Style
}

func newMultiStyle(s ...Style) Style {
	switch len(s) {
	case 0:
		return multiStyle{}
	case 1:
		// TODO(mway): Use s[0] directly
	default:
	}

	buf := _builders.Get()
	defer _builders.Put(buf)

	// TODO(mway): Potentially unpack multiStyles
	buf.WriteString("\x1b[") //nolint:errcheck
	var written int
	for _, style := range s {
		code := style.Code()
		if len(code) == 0 {
			continue
		}

		if written > 0 {
			buf.WriteByte(';') //nolint:errcheck
		}
		buf.WriteString(code) //nolint:errcheck
		written++
	}

	if written == 0 {
		return multiStyle{}
	}

	buf.WriteByte('m') //nolint:errcheck

	return multiStyle{
		styles: append([]Style(nil), s...),
		escape: buf.String(),
	}
}

// Combine combines colors into [Colors].
func Combine(s ...Style) Style {
	return newMultiStyle(s...)
}

func (s multiStyle) Escape() string {
	if !Enabled() {
		return ""
	}
	return s.escape
}

func (s multiStyle) Reset() string {
	if len(s.styles) == 0 {
		return ""
	}

	return s.styles[len(s.styles)-1].Reset()
}

func (s multiStyle) String() string {
	return s.escape
}

func (s multiStyle) With(styles ...Style) Style {
	if len(styles) == 0 {
		return s
	}

	return newMultiStyle(append([]Style{s}, styles...)...)
}

func (s multiStyle) Join(elems []string, sep string) string {
	buf := _builders.Get()
	defer _builders.Put(buf)

	for i, str := range elems {
		if i > 0 {
			buf.WriteString(sep) //nolint:errcheck
		}
		buf.WriteString(s.Wrap(str)) //nolint:errcheck
	}

	return buf.String()
}

func (s multiStyle) Wrap(str string) string {
	return s.Escape() + str + s.Reset()
}

func (s multiStyle) Code() string {
	if len(s.escape) == 0 {
		return ""
	}
	return s.escape[2 : len(s.escape)-1]
}

// Copy copies src to dest as in io.Copy, but wrapped in c.
func (s multiStyle) Copy(dst io.Writer, src io.Reader) (int64, error) {
	buf := _builders.Get()
	defer _builders.Put(buf)

	buf.WriteString(s.Escape()) //nolint:errcheck

	n, err := dst.Write(buf.Bytes())
	if err != nil {
		return int64(n), err
	}

	n64, err := io.Copy(dst, src)
	if err != nil {
		return int64(n) + n64, err
	}
	n += int(n64)

	buf.Reset()
	buf.WriteString(s.Reset()) //nolint:errcheck

	m, err := dst.Write(buf.Bytes())
	return int64(n + m), err
}

// Print prints args as in fmt.Print, but wrapped in c.
func (s multiStyle) Print(args ...any) {
	defer _stdout.Flush() //nolint:errcheck

	_stdout.WriteString(s.Escape()) //nolint:errcheck
	fmt.Fprint(_stdout, args...)    //nolint:errcheck
	_stdout.WriteString(s.Reset())  //nolint:errcheck
}

// Printf prints msg and args as in fmt.Printf, but wrapped in c.
func (s multiStyle) Printf(msg string, args ...any) {
	defer _stdout.Flush() //nolint:errcheck

	_stdout.WriteString(s.Escape())    //nolint:errcheck
	fmt.Fprintf(_stdout, msg, args...) //nolint:errcheck
	_stdout.WriteString(s.Reset())     //nolint:errcheck
}

// Println prints args as in fmt.Println, but wrapped in c.
func (s multiStyle) Println(args ...any) {
	defer _stdout.Flush() //nolint:errcheck

	buf := _builders.Get()
	defer _builders.Put(buf)

	_stdout.WriteString(s.Escape()) //nolint:errcheck
	fmt.Fprintln(buf, args...)      //nolint:errcheck

	var (
		tmp = buf.Bytes()
		esc = s.Reset()
	)

	if len(esc) > 0 {
		tmp[len(tmp)-1] = esc[0]
		buf.WriteString(esc[1:]) //nolint:errcheck
		buf.WriteByte('\n')      //nolint:errcheck
	}

	_stdout.WriteString(buf.String()) //nolint:errcheck
}

// Sprint returns a string containing args as in fmt.Sprint, but wrapped in c.
func (s multiStyle) Sprint(args ...any) string {
	return s.Escape() + fmt.Sprint(args...) + s.Reset()
}

// Sprintf returns a string containing msg and args as in fmt.Sprintf, but
// wrapped in c.
func (s multiStyle) Sprintf(msg string, args ...any) string {
	return s.Escape() + fmt.Sprintf(msg, args...) + s.Reset()
}

// Sprintln returns a string containing args as in fmt.Sprintln, but wrapped
// in c.
func (s multiStyle) Sprintln(args ...any) string {
	buf := _builders.Get()
	defer _builders.Put(buf)

	buf.WriteString(s.Escape()) //nolint:errcheck
	fmt.Fprintln(buf, args...)  //nolint:errcheck

	var (
		tmp = buf.Bytes()
		esc = s.Reset()
	)

	if len(esc) > 0 {
		tmp[len(tmp)-1] = esc[0]
		buf.WriteString(esc[1:]) //nolint:errcheck
		buf.WriteByte('\n')      //nolint:errcheck
	}

	return buf.String()
}

// Fprint prints args to w as in fmt.Fprint, but wrapped in c.
func (s multiStyle) Fprint(w io.Writer, args ...any) (int, error) {
	esc := s.Escape()
	fmt.Fprint(w, esc)                 //nolint:errcheck
	n, _ := fmt.Fprint(w, args...)     //nolint:errcheck
	m, err := fmt.Fprint(w, s.Reset()) //nolint:errcheck
	return n + m + len(esc), err
}

// Fprintf prints msg and args to w as in [fmt.Fprintf], but wrapped in c.
func (s multiStyle) Fprintf(w io.Writer, msg string, args ...any) (int, error) {
	esc := s.Escape()
	fmt.Fprint(w, esc)                   //nolint:errcheck
	n, _ := fmt.Fprintf(w, msg, args...) //nolint:errcheck
	m, err := fmt.Fprint(w, s.Reset())   //nolint:errcheck
	return n + m + len(esc), err
}

// Fprintln prints args to w as in [fmt.Fprintln], but wrapped in c.
func (s multiStyle) Fprintln(w io.Writer, args ...any) (int, error) {
	buf := _builders.Get()
	defer _builders.Put(buf)

	buf.WriteString(s.Escape()) //nolint:errcheck
	fmt.Fprintln(buf, args...)  //nolint:errcheck

	var (
		tmp = buf.Bytes()
		esc = s.Reset()
	)

	if len(esc) > 0 {
		tmp[len(tmp)-1] = esc[0]
		buf.WriteString(esc[1:]) //nolint:errcheck
		buf.WriteByte('\n')      //nolint:errcheck
	}

	n, err := io.Copy(w, buf)
	return int(n), err
}

// Enabled returns whether color is enabled based on terminal settings.
func Enabled() bool {
	return _hasColor
}

// Copy is a convenience function that calls s.Copy(dst, src).
func Copy(s Style, dst io.Writer, src io.Reader) (int64, error) {
	return s.Copy(dst, src)
}

// Join is a convenience function that calls s.Join(elems, sep).
func Join(s Style, elems []string, sep string) string {
	return s.Join(elems, sep)
}

// Wrap is a convenience function that calls s.Wrap(msg).
func Wrap(s Style, msg string) string {
	return s.Wrap(msg)
}

// Print is a convenience function that calls s.Print(args...).
func Print(s Style, args ...any) {
	s.Print(args...)
}

// Printf is a convenience function that calls s.Printf(msg, args...).
func Printf(s Style, msg string, args ...any) {
	s.Printf(msg, args...)
}

// Println is a convenience function that calls s.Println(args...).
func Println(s Style, args ...any) {
	s.Println(args...)
}

// Sprint is a convenience function that calls s.Sprint(args...).
func Sprint(s Style, args ...any) string {
	return s.Sprint(args...)
}

// Sprintf is a convenience function that calls s.Sprintf(msg, args...).
func Sprintf(s Style, msg string, args ...any) string {
	return s.Sprintf(msg, args...)
}

// Sprintln is a convenience function that calls s.Sprintln(args...).
func Sprintln(s Style, args ...any) string {
	return s.Sprintln(args...)
}

// Fprint is a convenience function that calls s.Fprint(dst, args...).
func Fprint(dst io.Writer, s Style, args ...any) (int, error) {
	return s.Fprint(dst, args...)
}

// Fprintf is a convenience function that calls s.Fprintf(dst, msg, args...).
func Fprintf(dst io.Writer, s Style, msg string, args ...any) (int, error) {
	return s.Fprintf(dst, msg, args...)
}

// Fprintln is a convenience function that calls s.Fprintln(dst, args...).
func Fprintln(dst io.Writer, s Style, args ...any) (int, error) {
	return s.Fprintln(dst, args...)
}

func isset(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

func parseColor(src map[string]Color, name string) (Color, error) {
	x, ok := src[name]
	if !ok {
		return Reset, errors.Wrap(ErrInvalidColorName, name)
	}
	return x, nil
}
