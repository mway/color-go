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
	_fd       = os.Stderr.Fd()
	_stdout   = bufio.NewWriter(os.Stdout)
	_hasColor = !isset("NO_COLOR") && os.Getenv("TERM") != "dumb" &&
		(isatty.IsTerminal(_fd) || isatty.IsCygwinTerminal(_fd))
	_builders = pool.NewWithReleaser(
		func() *bytes.Buffer { return bytes.NewBuffer(make([]byte, 0, 256)) },
		func(x *bytes.Buffer) { x.Reset() },
	)
	_strings = [math.MaxUint8]string{
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
)

// A Color is a terminal color.
type Color uint8

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

// Wrap wraps str with c.
func (c Color) Wrap(str string) string {
	return c.Escape() + str + c.Reset()
}

// WrapN wraps strs, joined by spaces, with c.
func (c Color) WrapN(strs ...string) string {
	buf := _builders.Get()
	defer _builders.Put(buf)

	for i := range strs {
		if i > 0 {
			buf.WriteByte(' ') //nolint:errcheck
		}
		buf.WriteString(strs[i]) //nolint:errcheck
	}

	return c.Escape() + buf.String() + c.Reset()
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

// Colors are a group of [Color]s.
type Colors []Color

// Combine combines colors into [Colors].
func Combine(colors ...Color) Colors {
	return append([]Color(nil), colors...)
}

// Escape returns a concatenation of c's escape codes.
func (c Colors) Escape() string {
	buf := _builders.Get()
	defer _builders.Put(buf)

	for i := range c {
		buf.WriteString(c[i].Escape()) //nolint:errcheck
	}

	return buf.String()
}

// Reset returns the escape code the reset output after c.
func (c Colors) Reset() string {
	if len(c) == 0 {
		return ""
	}

	return c[len(c)-1].Reset()
}

// String returns c's escape code, regardless of whether color is enabled.
func (c Colors) String() string {
	buf := _builders.Get()
	defer _builders.Put(buf)

	for i := range c {
		buf.WriteString(c[i].String()) //nolint:errcheck
	}

	return buf.String()
}

// Wrap wraps str with c.
func (c Colors) Wrap(str string) string {
	return c.Escape() + str + c.Reset()
}

// WrapN wraps strs, joined by spaces, with c.
func (c Colors) WrapN(strs ...string) string {
	buf := _builders.Get()
	defer _builders.Put(buf)

	for i := range strs {
		if i > 0 {
			buf.WriteByte(' ') //nolint:errcheck
		}
		buf.WriteString(strs[i]) //nolint:errcheck
	}

	return c.Escape() + buf.String() + c.Reset()
}

// Print prints args as in fmt.Print, but wrapped in c.
func (c Colors) Print(args ...any) {
	defer _stdout.Flush() //nolint:errcheck

	_stdout.WriteString(c.Escape()) //nolint:errcheck
	fmt.Fprint(_stdout, args...)    //nolint:errcheck
	_stdout.WriteString(c.Reset())  //nolint:errcheck
}

// Printf prints msg and args as in fmt.Printf, but wrapped in c.
func (c Colors) Printf(msg string, args ...any) {
	defer _stdout.Flush() //nolint:errcheck

	_stdout.WriteString(c.Escape())    //nolint:errcheck
	fmt.Fprintf(_stdout, msg, args...) //nolint:errcheck
	_stdout.WriteString(c.Reset())     //nolint:errcheck
}

// Println prints args as in fmt.Println, but wrapped in c.
func (c Colors) Println(args ...any) {
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
func (c Colors) Sprint(args ...any) string {
	return c.Escape() + fmt.Sprint(args...) + c.Reset()
}

// Sprintf returns a string containing msg and args as in fmt.Sprintf, but
// wrapped in c.
func (c Colors) Sprintf(msg string, args ...any) string {
	return c.Escape() + fmt.Sprintf(msg, args...) + c.Reset()
}

// Sprintln returns a string containing args as in fmt.Sprintln, but wrapped
// in c.
func (c Colors) Sprintln(args ...any) string {
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
func (c Colors) Fprint(w io.Writer, args ...any) (int, error) {
	esc := c.Escape()
	fmt.Fprint(w, esc)                 //nolint:errcheck
	n, _ := fmt.Fprint(w, args...)     //nolint:errcheck
	m, err := fmt.Fprint(w, c.Reset()) //nolint:errcheck
	return n + m + len(esc), err
}

// Fprintf prints msg and args to w as in [fmt.Fprintf], but wrapped in c.
func (c Colors) Fprintf(w io.Writer, msg string, args ...any) (int, error) {
	esc := c.Escape()
	fmt.Fprint(w, esc)                   //nolint:errcheck
	n, _ := fmt.Fprintf(w, msg, args...) //nolint:errcheck
	m, err := fmt.Fprint(w, c.Reset())   //nolint:errcheck
	return n + m + len(esc), err
}

// Fprintln prints args to w as in [fmt.Fprintln], but wrapped in c.
func (c Colors) Fprintln(w io.Writer, args ...any) (int, error) {
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

// Enabled returns whether color is enabled based on terminal settings.
func Enabled() bool {
	return _hasColor
}

func isset(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}
