package net

func Htons16(n int) int {
	return (n & 0xFF) << 8 + (n >> 8) & 0xFF
}
