package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
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

func (h *history) save() (string, error) {
	fileName := getFileName()
	oldText, err := readOldText(fileName)
	if err != nil {
		return "", fmt.Errorf("[save] failed to read existing text: %w ", err)
	}

	f, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("[save] failed to create file: %w", err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	_, err = w.Write(oldText)
	if err != nil {
		return "", fmt.Errorf("[save] failed to write to buffer: %w", err)
	}

	for _, v := range h.entries {
		line := fmt.Sprintf("%v,%v\n", v.time.String(), v.kind.String())
		_, err := w.WriteString(line)
		if err != nil {
			return "", fmt.Errorf("[save] failed to write to buffer: %w", err)
		}
	}

	h.clear()

	err = w.Flush()
	if err != nil {
		return "", fmt.Errorf("[save] failed to flush buffer: %w", err)
	}

	return "Successfully saved.", nil
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
		return nil, fmt.Errorf("[readOldText] failed to read all text from %v: %w", path, err)
	}

	return oldText, nil
}

func getFileName() string {
	now := time.Now().Local()
	return fmt.Sprintf("logs/log_%v%v%v.csv", now.Year(), now.Month(), now.Day())
}
