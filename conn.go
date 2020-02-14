// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logs

import (
	"bytes"
	"io"
	"net"
	"sync"
)

// ConnWriter implements LoggerInterface.
// it writes messages in keep-live tcp connection.
type ConnWriter struct {
	sync.Mutex
	writer         io.WriteCloser
	rec            chan *LogRecord
	format         string
	ReconnectOnMsg bool   `json:"reconnectOnMsg"`
	Reconnect      bool   `json:"reconnect"`
	Net            string `json:"net"`
	Addr           string `json:"addr"`
	Level          Level    `json:"level"`
}

// NewConn create new ConnWrite returning as LoggerInterface.
func NewConn(Net, Addr, format string, level Level) *ConnWriter {
	w := &ConnWriter{
		rec:       make(chan *LogRecord, LogBufferLength),
		format:    "[%D %T] [%L] (%S) %M",
		Net:Net,
		Addr:Addr,
		Level:level,
	}

	go func() {
		defer recoverPanic()
		defer func() {
			w.connect()
		}()

		for {
			select {
			case rec, ok := <-w.rec:
				if !ok {
					return
				}
				w.write(rec)
			}
		}
	}()
	return w
}

func (c *ConnWriter) LogWrite(rec *LogRecord) {
	c.rec <- rec
}

func (c *ConnWriter) SetFormat(format string) {
	c.format = format
}

func (c *ConnWriter) connect() error {
	if c.writer != nil {
		c.writer.Close()
		c.writer = nil
	}

	conn, err := net.Dial(c.Net, c.Addr)
	if err != nil {
		return err
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
	}

	c.writer = conn
	return nil
}

func (c *ConnWriter) needToConnectOnMsg() bool {
	if c.Reconnect {
		c.Reconnect = false
		return true
	}

	if c.writer == nil {
		return true
	}

	return c.ReconnectOnMsg
}

// This is the SocketLogWriter's output method
func (c *ConnWriter) write(rec *LogRecord) {
	c.Lock()
	if c.needToConnectOnMsg() {
		c.connect()
	}

	bt := bytes.NewBufferString(FormatLogRecord(c.format, rec))
	_, e := c.writer.Write(append(bt.Bytes(), '\n'))
	if e != nil {
		c.connect()
	}
	c.Unlock()
}

func (c *ConnWriter) Close() {
}
