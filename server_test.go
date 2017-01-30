/*
 * MIT License
 *
 * Copyright (c) 2017 Hiroaki Mizuguchi
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
package main

import (
	"bytes"
	"github.com/jarcoal/httpmock"
	"github.com/pin/tftp"
	"testing"
)

func TestHTTPGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"http://localhost/hoge",
		httpmock.NewStringResponder(200, "hoge"))

	s := NewServer("http://localhost")
	buf, err := s.getContent("hoge", "127.0.0.1")
	if buf.Len() != 4 {
		t.Error("length mismatch")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestHTTPPut(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"PUT",
		"http://localhost/fuga",
		httpmock.NewStringResponder(201, "fuga"))

	s := NewServer("http://localhost")
	err := s.putContent("fuga", "127.0.0.1", bytes.NewBufferString("hogehoge"))
	if err != nil {
		t.Error(err)
	}
}

func TestE2E(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := NewServer("http://localhost")
	go func() {
		s.ListenAndServe("127.0.0.1:1069")
	}()
	defer s.Shutdown()

	c, err := tftp.NewClient("127.0.0.1:1069")
	if err != nil {
		t.Error(err)
	}

	{
		httpmock.RegisterResponder(
			"GET",
			"http://localhost/hoge",
			httpmock.NewStringResponder(200, "hoge"))
		wt, err := c.Receive("hoge", "octet")
		if err != nil {
			t.Error(err)
		}
		var buf bytes.Buffer
		n, err := wt.WriteTo(&buf)
		if bytes.Compare(buf.Bytes(), []byte("hoge")) != 0 {
			t.Error("not mismatch")
		}
		t.Logf("%d bytes received\n", n)
	}

	{
		httpmock.RegisterResponder(
			"PUT",
			"http://localhost/fuga",
		httpmock.NewStringResponder(201, ""))
		rf, err := c.Send("fuga", "octet")
		if err != nil {
			t.Error(err)
		}
		b := bytes.NewBufferString("hogehogehoge")
		n, err := rf.ReadFrom(b)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%d bytes sent\n", n)
	}
}
