// Copyright (c) 2019 Andy Pan
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package gnet

import "github.com/Esbiya/gnet/internal/netpoll"

func (el *eventloop) handleEvent(fd int, ev uint32) error {
	if c, ok := el.connections[fd]; ok {
		// Don't change the ordering of processing EPOLLOUT | EPOLLRDHUP / EPOLLIN unless you're 100%
		// sure what you're doing!
		// Re-ordering can easily introduce bugs and bad side-effects, as I found out painfully in the past.

		// We should always check for the EPOLLOUT event first, as we must try to send the leftover data back to
		// client when any error occurs on a connection.
		//
		// Either an EPOLLOUT or EPOLLERR event may be fired when a connection is refused.
		// In either case loopWrite() should take care of it properly:
		// 1) writing data back,
		// 2) closing the connection.
		if ev&netpoll.OutEvents != 0 {
			if err := el.loopWrite(c); err != nil {
				return err
			}
		}
		if ev&netpoll.InEvents != 0 {
			return el.loopRead(c)
		}
		return nil
	}
	return el.loopAccept(fd)
}
