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
	"errors"
	"github.com/pin/tftp"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func NewServer(baseURL string) *Server {
	s := &Server{
		BaseURL: baseURL,
	}
	s.server = tftp.NewServer(s.readHandler, s.writeHandler)
	return s
}

type Server struct {
	BaseURL string
	server  *tftp.Server
}

func (s *Server) readHandler(filename string, rf io.ReaderFrom) error {
	raddr := rf.(tftp.OutgoingTransfer).RemoteAddr()
	addr := raddr.String()
	log.Printf("RRQ from %s: %s\n", addr, filename)
	buf, err := s.getContent(filename, addr)
	if err != nil {
		return err
	}
	rf.(tftp.OutgoingTransfer).SetSize(int64(buf.Len()))
	if _, err = rf.ReadFrom(&buf); err != nil {
		return err
	}
	log.Printf("RRQ Complete from %s: %s\n", addr, filename)
	return nil
}

func (s *Server) writeHandler(filename string, wt io.WriterTo) error {
	raddr := wt.(tftp.IncomingTransfer).RemoteAddr()
	addr := raddr.String()
	log.Printf("WRQ from %s: %s\n", addr, filename)
	var buf bytes.Buffer
	if _, err := wt.WriteTo(&buf); err != nil {
		return err
	}
	if err := s.putContent(filename, addr, &buf); err != nil {
		return err
	}
	log.Printf("WRQ Complete from %s: %s\n", addr, filename)
	return nil
}

func (s *Server) getContent(filename string, raddr string) (buf bytes.Buffer, err error) {
	url := s.BaseURL + filepath.Clean("/"+filename)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Forwarded-For", raddr)
	res, err := client.Do(req)
	if err != nil {
		return
	}
	buf.ReadFrom(res.Body)
	return
}

func (s *Server) putContent(filename string, raddr string, buf *bytes.Buffer) error {
	url := s.BaseURL + filepath.Clean("/"+filename)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, buf)
	req.Header.Add("X-Forwarded-For", raddr)
	req.Header.Add("Content-Length", string(buf.Len()))
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == 200 || res.StatusCode == 201 {
		return nil
	} else {
		return errors.New("WRQ Fail from %s: %s: http put")
	}
}

func (s *Server) ListenAndServe(listen string) error {
	return s.server.ListenAndServe(listen)
}

func (s *Server) Shutdown() {
	s.server.Shutdown()
}
