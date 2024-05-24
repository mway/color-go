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
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	_hasColor = true
}

func TestColor_Escape(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		require.Equal(t, esc, Color(i).Escape())
		require.Equal(t, esc, Combine(Color(i)).Escape())
	}
}

func TestColor_Escape_Reset_NoColor(t *testing.T) {
	_hasColor = false
	defer func() {
		_hasColor = true
	}()

	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		require.Equal(t, "", Color(i).Escape())
		require.Equal(t, "", Color(i).Reset())
		require.Equal(t, "", Combine(Color(i)).Escape())
		require.Equal(t, "", Combine(Color(i)).Reset())
	}
}

func TestColor_Reset(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		if c := Color(i); c == Reset {
			require.Equal(t, "", c.Reset())
			require.Equal(t, "", Combine(c).Reset())
		} else {
			require.Equal(t, Reset.Escape(), c.Reset())
			require.Equal(t, Reset.Escape(), Combine(c).Reset())
		}
	}

	require.Equal(t, "", Combine().Reset())
}

func Test_String(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		styles := []Style{
			Color(i),
			Combine(Color(i)),
		}

		_hasColor = false
		for _, style := range styles {
			require.Equal(t, esc, style.String())
		}
		_hasColor = true
		for _, style := range styles {
			require.Equal(t, esc, style.String())
		}
	}
}

func TestWrap(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		styles := []Style{
			Color(i),
			Combine(Color(i)),
		}

		for _, style := range styles {
			want := style.Escape() + t.Name() + style.Reset()
			require.Equal(t, want, style.Wrap(t.Name()))
			require.Equal(t, want, Wrap(style, t.Name()))
		}
	}
}

func TestWrapN(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		styles := []Style{
			Color(i),
			Combine(Color(i)),
		}

		for _, style := range styles {
			want := style.Escape() + t.Name() + " " + t.Name() + style.Reset()
			require.Equal(t, want, style.WrapN(t.Name(), t.Name()))
			require.Equal(t, want, WrapN(style, t.Name(), t.Name()))
		}
	}
}

func TestPrint(t *testing.T) {
	var (
		stdout = _stdout
		buf    bytes.Buffer
	)

	_stdout = bufio.NewWriter(&buf)
	t.Cleanup(func() {
		_stdout = stdout
	})

	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		styles := []Style{
			Color(i),
			Combine(Color(i)),
		}

		for _, style := range styles {
			want := style.Escape() + t.Name() + style.Reset()

			buf.Reset()
			style.Print(t.Name())
			require.Equal(t, want, buf.String())

			buf.Reset()
			Print(style, t.Name())
			require.Equal(t, want, buf.String())
		}
	}
}

func TestPrintf(t *testing.T) {
	var (
		stdout = _stdout
		buf    bytes.Buffer
	)

	_stdout = bufio.NewWriter(&buf)
	t.Cleanup(func() {
		_stdout = stdout
	})

	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		styles := []Style{
			Color(i),
			Combine(Color(i)),
		}

		for _, style := range styles {
			want := style.Escape() + t.Name() + "1" + style.Reset()

			buf.Reset()
			style.Printf("%s%d", t.Name(), 1)
			require.Equal(t, want, buf.String())

			buf.Reset()
			Printf(style, "%s%d", t.Name(), 1)
			require.Equal(t, want, buf.String())
		}
	}
}

func TestPrintln(t *testing.T) {
	var (
		stdout = _stdout
		buf    bytes.Buffer
	)

	_stdout = bufio.NewWriter(&buf)
	t.Cleanup(func() {
		_stdout = stdout
	})

	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		styles := []Style{
			Color(i),
			Combine(Color(i)),
		}

		for _, style := range styles {
			want := style.Escape() + t.Name() + " " + t.Name() + style.Reset() + "\n"

			buf.Reset()
			style.Println(t.Name(), t.Name())
			require.Equal(t, want, buf.String())

			buf.Reset()
			Println(style, t.Name(), t.Name())
			require.Equal(t, want, buf.String())
		}
	}
}

func TestSprint(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		styles := []Style{
			Color(i),
			Combine(Color(i)),
		}

		for _, style := range styles {
			want := style.Escape() + t.Name() + style.Reset()
			require.Equal(t, want, style.Sprint(t.Name()))
			require.Equal(t, want, Sprint(style, t.Name()))
		}
	}
}

func TestSprintf(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		styles := []Style{
			Color(i),
			Combine(Color(i)),
		}

		for _, style := range styles {
			want := style.Escape() + t.Name() + "1" + style.Reset()
			require.Equal(t, want, style.Sprintf("%s%d", t.Name(), 1))
			require.Equal(t, want, Sprintf(style, "%s%d", t.Name(), 1))
		}
	}
}

func TestSprintln(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		styles := []Style{
			Color(i),
			Combine(Color(i)),
		}

		for _, style := range styles {
			want := style.Escape() + t.Name() + " " + t.Name() + style.Reset() + "\n"
			require.Equal(t, want, style.Sprintln(t.Name(), t.Name()))
			require.Equal(t, want, Sprintln(style, t.Name(), t.Name()))
		}
	}
}

func TestFprint(t *testing.T) {
	var buf bytes.Buffer

	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		styles := []Style{
			Color(i),
			Combine(Color(i)),
		}

		for _, style := range styles {
			var (
				want = style.Escape() + t.Name() + style.Reset()
				err  error
			)

			buf.Reset()
			_, err = style.Fprint(&buf, t.Name())
			require.NoError(t, err)
			require.Equal(t, want, buf.String())

			buf.Reset()
			_, err = Fprint(&buf, style, t.Name())
			require.NoError(t, err)
			require.Equal(t, want, buf.String())
		}
	}
}

func TestFprintf(t *testing.T) {
	var buf bytes.Buffer

	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		styles := []Style{
			Color(i),
			Combine(Color(i)),
		}

		for _, style := range styles {
			var (
				want = style.Escape() + t.Name() + "1" + style.Reset()
				err  error
			)

			buf.Reset()
			_, err = style.Fprintf(&buf, "%s%d", t.Name(), 1)
			require.NoError(t, err)
			require.Equal(t, want, buf.String())

			buf.Reset()
			_, err = Fprintf(&buf, style, "%s%d", t.Name(), 1)
			require.NoError(t, err)
			require.Equal(t, want, buf.String())
		}
	}
}

func TestFprintln(t *testing.T) {
	var buf bytes.Buffer

	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		styles := []Style{
			Color(i),
			Combine(Color(i)),
		}

		for _, style := range styles {
			var (
				want = style.Escape() + t.Name() + " " + t.Name() + style.Reset() + "\n"
				err  error
			)

			buf.Reset()
			_, err = style.Fprintln(&buf, t.Name(), t.Name())
			require.NoError(t, err)
			require.Equal(t, want, buf.String())

			buf.Reset()
			_, err = Fprintln(&buf, style, t.Name(), t.Name())
			require.NoError(t, err)
			require.Equal(t, want, buf.String())
		}
	}
}

func TestCopy(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		styles := []Style{
			Color(i),
			Combine(Color(i)),
		}

		for _, style := range styles {
			var (
				want = style.Escape() + t.Name() + style.Reset()
				src  *bytes.Buffer
				dst  bytes.Buffer
				err  error
			)

			src = bytes.NewBufferString(t.Name())
			dst.Reset()
			_, err = style.Copy(&dst, src)
			require.NoError(t, err)
			require.Equal(t, want, dst.String())

			src = bytes.NewBufferString(t.Name())
			dst.Reset()
			_, err = Copy(style, &dst, src)
			require.NoError(t, err)
			require.Equal(t, want, dst.String())
		}
	}
}

func TestParseFgColor(t *testing.T) {
	for name, want := range _fgNames {
		have, err := ParseFgColor(name)
		require.NoError(t, err)
		require.Equal(t, want, have)
	}

	_, err := ParseFgColor("unknown-color-name")
	require.ErrorIs(t, err, ErrInvalidColorName)
	require.ErrorContains(t, err, "unknown-color-name")
}

func TestParseBgColor(t *testing.T) {
	for name, want := range _bgNames {
		have, err := ParseBgColor(name)
		require.NoError(t, err)
		require.Equal(t, want, have)
	}

	_, err := ParseBgColor("unknown-color-name")
	require.ErrorIs(t, err, ErrInvalidColorName)
	require.ErrorContains(t, err, "unknown-color-name")
}
