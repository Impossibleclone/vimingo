package main

import (
	"os"
	"bufio"
)

type Buffer struct {
	Filename string
	Lines 	[]string
	Cursor  *Cursor
	Mode 	Mode
	Command []rune
	ScrollX int
	ScrollY int
	Register string
}

func LoadFile(filename string) (*Buffer, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644) //r/w,if not exist create and 0644 so owner can r/w and other can r 
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner :=bufio.NewScanner(file)	
	for scanner.Scan(){
		lines = append(lines, scanner.Text())
	}
	return &Buffer{Lines: lines}, scanner.Err()
}

func SaveFile(filename string, buffer *Buffer) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644) //r/w,if not exist create and 0644 so owner can r/w and other can r 
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _,line := range buffer.Lines {
		_,err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

