// Copyright ©2016 The ev3go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ev3dev

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

var execCommand = exec.Command

// Console provides an simple API for text printing on EV3 LCD console.
//  cons, _ := NewConsole("Lat15-TerminusBold16")
//  cons.PutCursorAt(4, 4)
//  fmt.Print("Hello World!")
//  time.Sleep(5*time.Second)
type Console struct {
	// Number of columns on the EV3 LCD console supported by the current font.
	columns int
	// Number of rows on the EV3 LCD console supported by the current font.
	rows int
}

// NewConsole returns a new Console and set specified EV3 LCD console font.
func NewConsole(font string) (Console, error) {
	c := Console{}
	err := c.SetFont(font)
	if err != nil {
		return Console{}, fmt.Errorf("ev3dev: could not set font %s: %s", font, err)
	}
	return c, nil
}

// SetFont sets the EV3 LCD console font to specified one.
// List of available fonts can be found in /usr/share/consolefonts/
func (c *Console) SetFont(font string) (err error) {
	if font != "" {
		err = execCommand("setfont", font).Run()
		if err != nil {
			return
		}
	}

	c.rows, c.columns, err = c.readConsoleSize()
	if err != nil {
		return
	}
	if c.rows == 0 || c.columns == 0 {
		return fmt.Errorf("ev3dev: invalid size of console - %v rows, %v columns", c.rows, c.columns)
	}

	return nil
}

// Rows returns number of rows of EV3 LCD console.
func (c *Console) Rows() int {
	return c.rows
}

// Columns returns number of columns of EV3 LCD console.
func (c *Console) Columns() int {
	return c.columns
}

// PutCursorAt puts cursor on EV3 LCD console at specified position.
// It does nothing if position is out of range of LCD console size.
func (c *Console) PutCursorAt(columns, rows int) {
	if columns > 0 && columns <= c.columns && rows > 0 && rows <= c.rows {
		fmt.Printf("\x1b[%d;%dH", rows, columns)
	}
}

// HideCursor hides cursor on EV3 LCD console
func (c *Console) HideCursor() {
	fmt.Print("\x1b[?25l")
}

// ShowCursor shows cursor on EV3 LCD console
func (c *Console) ShowCursor() {
	fmt.Print("\x1b[?25h")
}

// Clear clears the EV3 LCD console using ANSI codes, and move the cursor to 1,1
func (c *Console) Clear() {
	fmt.Print("\x1b[2J\x1b[H")
}

// ClearToEOL clears to the end of line from specified position
// on the EV3 LCD console. Default to current cursor position.
func (c *Console) ClearToEOL(columns, rows int) {
	c.PutCursorAt(columns, rows)
	fmt.Print("\x1b[K")
}

// readConsoleSize reads and stores EV3 LCD console size.
func (c *Console) readConsoleSize() (rows, cols int, err error) {
	cmd := execCommand("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	s := bytes.Split(chomp(out), []byte(" "))
	if len(s) < 2 {
		return 0, 0, fmt.Errorf("invalid output of 'stty size' command")
	}

	return parseBytesToInt(s[0]), parseBytesToInt(s[1]), nil
}

func parseBytesToInt(b []byte) int {
	i, err := strconv.ParseInt(string(b), 10, 32)
	if err != nil {
		return 0
	}
	return int(i)
}
