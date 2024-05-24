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

package color_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/color"
)

func TestNop(t *testing.T) {
	require.Zero(t, color.Nop.Code())
	require.Zero(t, color.Nop.Escape())
	require.Zero(t, color.Nop.Reset())
	require.Zero(t, color.Nop.String())
	require.Equal(t, t.Name(), color.Nop.Wrap(t.Name()))
	require.Zero(t, color.Nop, color.Nop.With(color.FgRed))
	require.Equal(
		t,
		fmt.Sprintf("%s %s", t.Name(), t.Name()),
		color.Nop.WrapN(t.Name(), t.Name()),
	)

	var buf bytes.Buffer
	copyWritten, err := color.Nop.Copy(&buf, bytes.NewBufferString(t.Name()))
	require.NoError(t, err)
	require.EqualValues(t, len(t.Name()), copyWritten)
	require.Equal(t, t.Name(), buf.String())

	require.Equal(t, t.Name(), color.Nop.Sprint(t.Name()))
	require.Equal(t, t.Name(), color.Nop.Sprintf("%s", t.Name()))
	require.Equal(t, t.Name()+"\n", color.Nop.Sprintln(t.Name()))

	buf.Reset()
	fprintWritten, err := color.Nop.Fprint(&buf, t.Name())
	require.NoError(t, err)
	require.Equal(t, len(t.Name()), fprintWritten)
	require.Equal(t, t.Name(), buf.String())

	buf.Reset()
	fprintfWritten, err := color.Nop.Fprintf(&buf, "%s", t.Name())
	require.NoError(t, err)
	require.Equal(t, len(t.Name()), fprintfWritten)
	require.Equal(t, t.Name(), buf.String())

	buf.Reset()
	fprintlnWritten, err := color.Nop.Fprintln(&buf, t.Name())
	require.NoError(t, err)
	require.Equal(t, len(t.Name())+1, fprintlnWritten)
	require.Equal(t, t.Name()+"\n", buf.String())
}
