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

func TestColor_String(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		_hasColor = false
		require.Equal(t, esc, Color(i).String())
		require.Equal(t, esc, Combine(Color(i)).String())
		_hasColor = true
		require.Equal(t, esc, Color(i).String())
		require.Equal(t, esc, Combine(Color(i)).String())
	}
}

func TestColor_Wrap(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		{
			c := Color(i)
			require.Equal(t, c.Escape()+"foo"+c.Reset(), c.Wrap("foo"))
			require.Equal(t, c.Escape()+"foo"+c.Reset(), c.Wrap("foo"))
		}

		{
			c := Combine(Color(i))
			require.Equal(t, c.Escape()+"foo"+c.Reset(), c.Wrap("foo"))
			require.Equal(t, c.Escape()+"foo"+c.Reset(), c.Wrap("foo"))
		}
	}
}

func TestColor_WrapN(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		{
			c := Color(i)
			require.Equal(t, c.Escape()+"foo bar"+c.Reset(), c.WrapN("foo", "bar"))
			require.Equal(t, c.Escape()+"foo bar"+c.Reset(), c.WrapN("foo", "bar"))
		}

		{
			c := Combine(Color(i))
			require.Equal(t, c.Escape()+"foo bar"+c.Reset(), c.WrapN("foo", "bar"))
			require.Equal(t, c.Escape()+"foo bar"+c.Reset(), c.WrapN("foo", "bar"))
		}
	}
}

func TestColor_Print(t *testing.T) {
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

		{
			c := Color(i)
			c.Print("foo")
			require.Equal(t, c.Escape()+"foo"+c.Reset(), buf.String())
			buf.Reset()
		}

		{
			c := Combine(Color(i))
			c.Print("foo")
			require.Equal(t, c.Escape()+"foo"+c.Reset(), buf.String())
			buf.Reset()
		}
	}
}

func TestColor_Printf(t *testing.T) {
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

		{
			c := Color(i)
			c.Printf("foo%d", 1)
			require.Equal(t, c.Escape()+"foo1"+c.Reset(), buf.String())
			buf.Reset()
		}

		{
			c := Combine(Color(i))
			c.Printf("foo%d", 1)
			require.Equal(t, c.Escape()+"foo1"+c.Reset(), buf.String())
			buf.Reset()
		}
	}
}

func TestColor_Println(t *testing.T) {
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

		{
			c := Color(i)
			c.Println("foo", "bar")
			require.Equal(t, c.Escape()+"foo bar"+c.Reset()+"\n", buf.String())
			buf.Reset()
		}

		{
			c := Combine(Color(i))
			c.Println("foo", "bar")
			require.Equal(t, c.Escape()+"foo bar"+c.Reset()+"\n", buf.String())
			buf.Reset()
		}
	}
}

func TestColor_Sprint(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		{
			c := Color(i)
			require.Equal(t, c.Escape()+"foo"+c.Reset(), c.Sprint("foo"))
		}

		{
			c := Combine(Color(i))
			require.Equal(t, c.Escape()+"foo"+c.Reset(), c.Sprint("foo"))
		}
	}
}

func TestColor_Sprintf(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		{
			c := Color(i)
			require.Equal(t, c.Escape()+"foo1"+c.Reset(), c.Sprintf("foo%d", 1))
		}

		{
			c := Combine(Color(i))
			require.Equal(t, c.Escape()+"foo1"+c.Reset(), c.Sprintf("foo%d", 1))
		}
	}
}

func TestColor_Sprintln(t *testing.T) {
	for i, esc := range _strings {
		if len(esc) == 0 {
			continue
		}

		{
			c := Color(i)
			require.Equal(t, c.Escape()+"foo bar"+c.Reset()+"\n", c.Sprintln("foo", "bar"))
		}

		{
			c := Combine(Color(i))
			require.Equal(t, c.Escape()+"foo bar"+c.Reset()+"\n", c.Sprintln("foo", "bar"))
		}
	}
}

func TestColor_Fprint(t *testing.T) {
	var buf bytes.Buffer

	for i, esc := range _strings {
		c := Color(i)
		if len(esc) == 0 {
			continue
		}

		{
			_, err := c.Fprint(&buf, "foo")
			require.NoError(t, err)
			require.Equal(t, c.Escape()+"foo"+c.Reset(), buf.String())
			buf.Reset()
		}

		{
			c := Combine(c)
			_, err := c.Fprint(&buf, "foo")
			require.NoError(t, err)
			require.Equal(t, c.Escape()+"foo"+c.Reset(), buf.String())
			buf.Reset()
		}
	}
}

func TestColor_Fprintf(t *testing.T) {
	var buf bytes.Buffer

	for i, esc := range _strings {
		c := Color(i)
		if len(esc) == 0 {
			continue
		}

		{
			_, err := c.Fprintf(&buf, "foo%d", 1)
			require.NoError(t, err)
			require.Equal(t, c.Escape()+"foo1"+c.Reset(), buf.String())
			buf.Reset()
		}

		{
			c := Combine(c)
			_, err := c.Fprintf(&buf, "foo%d", 1)
			require.NoError(t, err)
			require.Equal(t, c.Escape()+"foo1"+c.Reset(), buf.String())
			buf.Reset()
		}
	}
}

func TestColor_Fprintln(t *testing.T) {
	var buf bytes.Buffer

	for i, esc := range _strings {
		c := Color(i)
		if len(esc) == 0 {
			continue
		}

		{
			_, err := c.Fprintln(&buf, "foo", "bar")
			require.NoError(t, err)
			require.Equal(t, c.Escape()+"foo bar"+c.Reset()+"\n", buf.String())
			buf.Reset()
		}

		{
			c := Combine(c)
			_, err := c.Fprintln(&buf, "foo", "bar")
			require.NoError(t, err)
			require.Equal(t, c.Escape()+"foo bar"+c.Reset()+"\n", buf.String())
			buf.Reset()
		}
	}
}
