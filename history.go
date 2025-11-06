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
	index   int // 現在の参照位置
}

func (h *history) Add(e historyEntry) {
	h.entries = append(h.entries, e)
	h.index = len(h.entries) // 最新に移動
}

func (h *history) Prev() *historyEntry {
	if h.index > 0 {
		h.index--
	}
	return &h.entries[h.index]
}

func (h *history) Next() *historyEntry {
	if h.index < len(h.entries)-1 {
		h.index++
		return &h.entries[h.index]
	}
	h.index = len(h.entries) // 最後に戻る
	return nil
}
func (h *history) Clear() {
	h.entries = []historyEntry{}
	h.index = 0
}

func (h *history) Pop() historyEntry {
	last := h.entries[len(h.entries)-1]
	h.entries = h.entries[:len(h.entries)-1]
	return last
}

func (h *history) Save() (string, error) {
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

	h.Clear()

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

	switch now.Day() {
	case 30:
		return "log_test_251030.csv"
	case 31:
		return "log_251031.csv"
	case 1:
		return "log_251101.csv"
	case 2:
		return "log_251102.csv"
	default:
		fmt.Println("warning: 2025年度文化祭は終了しています")
		return fmt.Sprintf("log_%v%v%v.csv", now.Year(), now.Month(), now.Day())
	}
}
