package main

func TypeCh(line string, pos int, ch rune) string {
	return line[:pos] + string(ch) + line[pos:]
}

func RemoveCh(line string, pos int) string {
	return line[:pos] + line[pos+1:]
}
