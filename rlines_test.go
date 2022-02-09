package rlines

import (
	"bufio"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRlines(t *testing.T) {
	f, err := os.Open("testdata/commonsense.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	pr, pw := io.Pipe()
	tr := io.TeeReader(f, pw)

	var lines1, lines2 []string

	go func() {
		sc := bufio.NewScanner(tr)
		for sc.Scan() {
			lines1 = append(lines1, sc.Text())
		}
		// Skip checking sc.Err()
		pw.Close()
	}()

	rl := NewReader(pr)
	for {
		l := rl.Next()
		if l == nil {
			break
		}
		lbytes, err := io.ReadAll(l)
		if err != nil {
			t.Fatal(err)
		}
		lines2 = append(lines2, string(lbytes))
	}

	if diff := cmp.Diff(lines1, lines2); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
