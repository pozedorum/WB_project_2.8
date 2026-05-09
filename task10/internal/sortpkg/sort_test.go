package sortpkg

import (
	"container/heap"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pozedorum/WB_project_2/task10/pkg/options"
)

func testFlags(mods ...func(*options.FlagStruct)) options.FlagStruct {
	k := 1
	n, r, u, m, b, c, h := false, false, false, false, false, false, false
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
	for _, mod := range mods {
		mod(&fs)
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
