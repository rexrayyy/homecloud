package client

import (
	"testing"
	"unicode/utf16"
	"unsafe"
)

func encodeZ(parts ...string) []uint16 {
	var buf []uint16
	for _, part := range parts {
		buf = append(buf, utf16.Encode([]rune(part))...)
		buf = append(buf, 0)
	}
	return append(buf, 0)
}

func TestOFNSelections(t *testing.T) {
	one := ofnSelections(encodeZ(`C:\docs\a.txt`))
	if len(one) != 1 || one[0] != `C:\docs\a.txt` {
		t.Fatalf("single selection: %#v", one)
	}
	many := ofnSelections(encodeZ(`C:\docs`, "a.txt", "b.txt"))
	if len(many) != 2 || many[0] != `C:\docs\a.txt` || many[1] != `C:\docs\b.txt` {
		t.Fatalf("multi selection: %#v", many)
	}
	if got := ofnSelections(nil); len(got) != 0 {
		t.Fatalf("empty: %#v", got)
	}
}

func TestDialogStructLayout(t *testing.T) {
	if unsafe.Sizeof(uintptr(0)) != 8 {
		t.Skip("layout check is for 64-bit Windows")
	}
	if n := unsafe.Sizeof(openFileName{}); n != 152 {
		t.Fatalf("OPENFILENAMEW size %d, want 152", n)
	}
	if n := unsafe.Sizeof(browseInfo{}); n != 64 {
		t.Fatalf("BROWSEINFOW size %d, want 64", n)
	}
}
