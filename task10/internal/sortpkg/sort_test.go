package sortpkg

import (
	"bytes"
	"container/heap"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pozedorum/WB_project_2/task10/pkg/options"
)

func testFlags(modify func(*options.FlagStruct)) options.FlagStruct {
	k := 0
	n := false
	r := false
	u := false
	m := false
	b := false
	c := false
	h := false

	fs := options.FlagStruct{
		KFlag: &k,
		NFlag: &n,
		RFlag: &r,
		UFlag: &u,
		MFlag: &m,
		BFlag: &b,
		CFlag: &c,
		HFlag: &h,
	}

	if modify != nil {
		modify(&fs)
	}

	return fs
}

func TestStringsSortNumericReverseIsStrictAndCorrect(t *testing.T) {
	fs := testFlags(func(fs *options.FlagStruct) {
		*fs.NFlag = true
		*fs.RFlag = true
	})
	lines := []string{"2", "10", "2", "1"}

	ss := MakeSortStruct(lines, fs)
	if err := ss.StringsSort(); err != nil {
		t.Fatalf("StringsSort returned error: %v", err)
	}

	want := []string{"10", "2", "2", "1"}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("got %v, want %v", lines, want)
	}
}

func TestStringsSortMonthAndHuman(t *testing.T) {
	monthFS := testFlags(func(fs *options.FlagStruct) { *fs.MFlag = true })
	months := []string{"Dec", "jan", "Sep"}
	if err := MakeSortStruct(months, monthFS).StringsSort(); err != nil {
		t.Fatalf("month sort returned error: %v", err)
	}
	if want := []string{"jan", "Sep", "Dec"}; !reflect.DeepEqual(months, want) {
		t.Fatalf("got %v, want %v", months, want)
	}

	humanFS := testFlags(func(fs *options.FlagStruct) { *fs.HFlag = true })
	human := []string{"2M", "500K", "1G"}
	if err := MakeSortStruct(human, humanFS).StringsSort(); err != nil {
		t.Fatalf("human sort returned error: %v", err)
	}
	if want := []string{"500K", "2M", "1G"}; !reflect.DeepEqual(human, want) {
		t.Fatalf("got %v, want %v", human, want)
	}
}

func TestConflictingFlagsReturnError(t *testing.T) {
	fs := testFlags(func(fs *options.FlagStruct) {
		*fs.NFlag = true
		*fs.MFlag = true
	})

	if err := MakeSortStruct([]string{"Jan", "2"}, fs).StringsSort(); err == nil {
		t.Fatal("expected conflicting flags error")
	}
}

func TestHeapUsesSortFlags(t *testing.T) {
	fs := testFlags(func(fs *options.FlagStruct) {
		*fs.NFlag = true
	})
	h := &MinHeap{fs: fs}
	heap.Init(h)

	key10, err := getKey(fs, "10")
	if err != nil {
		t.Fatal(err)
	}

	key2, err := getKey(fs, "2")
	if err != nil {
		t.Fatal(err)
	}

	heap.Push(h, HeapItem{line: "10", key: key10, index: 0})
	heap.Push(h, HeapItem{line: "2", key: key2, index: 1})

	got := heap.Pop(h).(HeapItem).line
	if got != "2" {
		t.Fatalf("heap ignored numeric flag: got %q", got)
	}
}

func TestExternalSortMergeUsesFlagsUniqueAndTempDir(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.txt")
	output := filepath.Join(dir, "output.txt")

	if err := os.WriteFile(input, []byte("10\n2\n1\n2\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	fs := testFlags(func(fs *options.FlagStruct) {
		*fs.NFlag = true
		*fs.UFlag = true
	})

	ess := MakeExternalSortStruct(fs, input, output, 2)

	if err := ess.splitAndSort(); err != nil {
		t.Fatalf("splitAndSort returned error: %v", err)
	}

	for _, chunk := range ess.tempFilesList {
		if filepath.Dir(chunk) == "." {
			t.Fatalf("chunk created in current working directory: %s", chunk)
		}

		if filepath.Dir(chunk) == dir {
			t.Fatalf("chunk created near input/output instead of system temp dir: %s", chunk)
		}
	}

	if err := ess.mergeChunks(); err != nil {
		t.Fatalf("mergeChunks returned error: %v", err)
	}

	gotBytes, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}

	got := strings.Split(strings.TrimSpace(string(gotBytes)), "\n")
	want := []string{"1", "2", "10"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestIsSortedBasic(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "input.txt")

	if err := os.WriteFile(file, []byte("apple\nbanana\ncherry\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	fs := testFlags(nil)

	got, err := isSorted(file, fs)
	if err != nil {
		t.Fatal(err)
	}

	if !got {
		t.Fatalf("expected file to be sorted")
	}
}

func TestIsSortedReturnsFalse(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "input.txt")

	if err := os.WriteFile(file, []byte("banana\napple\ncherry\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	fs := testFlags(nil)

	got, err := isSorted(file, fs)
	if err != nil {
		t.Fatal(err)
	}

	if got {
		t.Fatalf("expected file to be unsorted")
	}
}

func TestIsSortedNumericByColumn(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "input.txt")

	if err := os.WriteFile(file, []byte("a 1\nb 2\nc 10\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	fs := testFlags(func(fs *options.FlagStruct) {
		*fs.NFlag = true
		*fs.KFlag = 2
	})

	got, err := isSorted(file, fs)
	if err != nil {
		t.Fatal(err)
	}

	if !got {
		t.Fatalf("expected file to be numerically sorted by column 2")
	}
}

func TestIsSortedReverse(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "input.txt")

	if err := os.WriteFile(file, []byte("cherry\nbanana\napple\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	fs := testFlags(func(fs *options.FlagStruct) {
		*fs.RFlag = true
	})

	got, err := isSorted(file, fs)
	if err != nil {
		t.Fatal(err)
	}

	if !got {
		t.Fatalf("expected file to be reverse sorted")
	}
}

func TestIsSortedEmptyFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "empty.txt")

	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}

	fs := testFlags(nil)

	got, err := isSorted(file, fs)
	if err != nil {
		t.Fatal(err)
	}

	if !got {
		t.Fatalf("empty file should be considered sorted")
	}
}

func TestIsSortedMissingFileReturnsError(t *testing.T) {
	fs := testFlags(nil)

	_, err := isSorted("missing-file.txt", fs)
	if err == nil {
		t.Fatalf("expected error for missing file")
	}
}

func TestExternalSortToStdout(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.txt")

	if err := os.WriteFile(input, []byte("banana\napple\ncherry\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	fs := testFlags(nil)

	got := captureStdout(t, func() {
		if err := ExternalSortToStdout(input, fs); err != nil {
			t.Fatal(err)
		}
	})

	want := "apple\nbanana\ncherry\n"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestProcessStdio(t *testing.T) {
	fs := testFlags(func(fs *options.FlagStruct) {
		*fs.NFlag = true
	})

	got := captureStdio(t, "10\n2\n1\n", func() {
		if err := ProcessStdio(fs); err != nil {
			t.Fatal(err)
		}
	})

	want := "1\n2\n10\n"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestProcessInteractiveInput(t *testing.T) {
	fs := testFlags(nil)

	got := captureStdio(t, "banana\napple\ncherry\n", func() {
		if err := ProcessInteractiveInput(fs); err != nil {
			t.Fatal(err)
		}
	})

	want := "apple\nbanana\ncherry\n"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}

	return buf.String()
}

func captureStdio(t *testing.T, input string, fn func()) string {
	t.Helper()

	oldStdin := os.Stdin
	oldStdout := os.Stdout

	stdinR, stdinW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := stdinW.WriteString(input); err != nil {
		t.Fatal(err)
	}
	if err := stdinW.Close(); err != nil {
		t.Fatal(err)
	}

	os.Stdin = stdinR
	os.Stdout = stdoutW

	fn()

	if err := stdoutW.Close(); err != nil {
		t.Fatal(err)
	}

	os.Stdin = oldStdin
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, stdoutR); err != nil {
		t.Fatal(err)
	}

	return buf.String()
}
