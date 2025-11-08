package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type historyEntry struct {
	time time.Time
	kind personKind
}

type history struct {
	entries []historyEntry
}

func (h *history) add(e historyEntry) {
	h.entries = append(h.entries, e)
}

func (h *history) clear() {
	h.entries = []historyEntry{}
}

// 一番後ろのhistoryEntryを取り出して返す
func (h *history) pop() historyEntry {
	last := h.entries[len(h.entries)-1]
	h.entries = h.entries[:len(h.entries)-1]
	return last
}

type writerWithFile struct {
	w *bufio.Writer
	f *os.File
}

func (h *history) save() error {
	writers := make(map[string]*writerWithFile)

	for _, v := range h.entries {
		path := getFilePath(v.time)
		var wtr *writerWithFile
		if wf, ok := writers[path]; !ok {
			newWtr, err := initFile(path)
			if err != nil {
				return err
			}

			writers[path] = newWtr
			wtr = newWtr
		} else {
			wtr = wf
		}
		line := fmt.Sprintf("%v,%v\n", v.time.String(), v.kind.string())
		_, err := wtr.w.WriteString(line)
		if err != nil {
			return fmt.Errorf("[save] failed to write to buffer: %w", err)
		}
	}

	// 履歴を空にする
	h.clear()

	for _, w := range writers {
		defer w.f.Close()
		err := w.w.Flush()
		if err != nil {
			return fmt.Errorf("[save] failed to flush buffer: %w", err)
		}
	}

	return nil
}

func initFile(path string) (*writerWithFile, error) {
	oldText, err := readOldText(path)
	if err != nil {
		return nil, fmt.Errorf("[save] failed to read existing text: %w ", err)
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("[save] failed to create file: %w", err)
	}
	w := bufio.NewWriter(f)

	_, err = w.Write(oldText)
	if err != nil {
		return nil, fmt.Errorf("[save] failed to write to buffer: %w", err)
	}
	return &writerWithFile{w, f}, nil
}

func readOldText(path string) ([]byte, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return []byte{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("[readOldText] failed to open file: %w", err)
	}
	defer f.Close()

	oldText, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("[readOldText] %w", err)
	}

	return oldText, nil
}

func getFilePath(time time.Time) string {
	time = time.Local()
	return fmt.Sprintf("logs/%s/log_%s.csv", time.Format("2006"), time.Format("060102"))
}
