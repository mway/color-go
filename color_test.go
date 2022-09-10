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
			want := style.Escape() + "foo" + style.Reset()
			require.Equal(t, want, style.Wrap("foo"))
			require.Equal(t, want, Wrap(style, "foo"))
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
			want := style.Escape() + "foo bar" + style.Reset()
			require.Equal(t, want, style.WrapN("foo", "bar"))
			require.Equal(t, want, WrapN(style, "foo", "bar"))
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
			want := style.Escape() + "foo" + style.Reset()

			buf.Reset()
			style.Print("foo")
			require.Equal(t, want, buf.String())

			buf.Reset()
			Print(style, "foo")
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
			want := style.Escape() + "foo1" + style.Reset()

			buf.Reset()
			style.Printf("foo%d", 1)
			require.Equal(t, want, buf.String())

			buf.Reset()
			Printf(style, "foo%d", 1)
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
			want := style.Escape() + "foo bar" + style.Reset() + "\n"

			buf.Reset()
			style.Println("foo", "bar")
			require.Equal(t, want, buf.String())

			buf.Reset()
			Println(style, "foo", "bar")
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
			want := style.Escape() + "foo" + style.Reset()
			require.Equal(t, want, style.Sprint("foo"))
			require.Equal(t, want, Sprint(style, "foo"))
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
			want := style.Escape() + "foo1" + style.Reset()
			require.Equal(t, want, style.Sprintf("foo%d", 1))
			require.Equal(t, want, Sprintf(style, "foo%d", 1))
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
			want := style.Escape() + "foo bar" + style.Reset() + "\n"
			require.Equal(t, want, style.Sprintln("foo", "bar"))
			require.Equal(t, want, Sprintln(style, "foo", "bar"))
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
				want = style.Escape() + "foo" + style.Reset()
				err  error
			)

			buf.Reset()
			_, err = style.Fprint(&buf, "foo")
			require.NoError(t, err)
			require.Equal(t, want, buf.String())

			buf.Reset()
			_, err = Fprint(&buf, style, "foo")
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
				want = style.Escape() + "foo1" + style.Reset()
				err  error
			)

			buf.Reset()
			_, err = style.Fprintf(&buf, "foo%d", 1)
			require.NoError(t, err)
			require.Equal(t, want, buf.String())

			buf.Reset()
			_, err = Fprintf(&buf, style, "foo%d", 1)
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
				want = style.Escape() + "foo bar" + style.Reset() + "\n"
				err  error
			)

			buf.Reset()
			_, err = style.Fprintln(&buf, "foo", "bar")
			require.NoError(t, err)
			require.Equal(t, want, buf.String())

			buf.Reset()
			_, err = Fprintln(&buf, style, "foo", "bar")
			require.NoError(t, err)
			require.Equal(t, want, buf.String())
		}
	}
}
