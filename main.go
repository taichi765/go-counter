package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO: ドキュメントコメントを残す

// 種別を表すenum(的なもの)
type personKind int

// KIND
const (
	//小学生
	elemStudent personKind = iota
	//高校生男子
	hsBoy
	//高校生女子
	hsGirl
	//親
	parent
	//その他
	other
)

func (k personKind) string() string {
	// KIND
	switch k {
	case elemStudent:
		return "小学生"
	case hsBoy:
		return "中高生男子"
	case hsGirl:
		return "中高生女子"
	case parent:
		return "親"
	case other:
		return "その他"
	default:
		return "Unknown Kind"
	}
}

type saveResultMsg struct {
	success bool
	err     error
}

func save(h history) tea.Cmd {
	return func() tea.Msg {
		err := h.save()
		if err != nil {
			return saveResultMsg{false, err}
		}
		return saveResultMsg{true, nil}
	}
}

// TODO: mapにした方が変更可能性が高い
type counterModel struct {
	// KIND

	elemStudentCount int
	hsBoyCount       int
	hsGirlCount      int
	parentCount      int
	otherCount       int

	showAskSaveDialog bool

	history history
	message string
}

func (m counterModel) AddCount(kind personKind) counterModel {
	// KIND
	switch kind {
	case elemStudent:
		m.elemStudentCount++
	case hsBoy:
		m.hsBoyCount++
	case hsGirl:
		m.hsGirlCount++
	case parent:
		m.parentCount++
	case other:
		m.otherCount++
	}
	m.history.add(historyEntry{time: time.Now(), kind: kind})
	return m
}

func (m counterModel) Init() tea.Cmd {
	return nil
}

// TODO: 長いのでわける
func (m counterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.message = "" //最初にリセットしておく
	if m.showAskSaveDialog {
		switch msg := msg.(type) {
		case saveResultMsg:
			if msg.success {
				return m, tea.Quit
			} else {
				m.message = "Failed to save:" + msg.err.Error()
				return m, nil
			}
		case tea.KeyMsg:
			switch msg.String() {
			case "y":
				return m, save(m.history)
			case "n":
				return m, tea.Quit
			case "esc":
				m.showAskSaveDialog = false
				return m, nil
			}
		default:
			m.message = fmt.Sprintf("warning: unknown msg type when showing AskSaveDialog: %T", msg)
		}
	}
	switch msg := msg.(type) {
	case saveResultMsg:
		if msg.success {
			m.message = "successfully saved"
			return m, nil
		} else {
			m.message = "failed to save: " + msg.err.Error()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.showAskSaveDialog = true
		case "backspace":
			last := m.history.pop()
			switch last.kind {
			case elemStudent:
				m.elemStudentCount--
			case hsBoy:
				m.hsBoyCount--
			case hsGirl:
				m.hsGirlCount--
			case parent:
				m.parentCount--
			case other:
				m.otherCount--
			default:
				m.message = "unknown kind"
			}
		case "ctrl+s":
			return m, save(m.history)
		// KIND
		case "g":
			return m.AddCount(elemStudent), nil
		case "h":
			return m.AddCount(hsBoy), nil
		case "j":
			return m.AddCount(hsGirl), nil
		case "k":
			return m.AddCount(parent), nil
		case "l":
			return m.AddCount(other), nil
		}
	}
	return m, nil
}

func (m counterModel) View() string {
	if m.showAskSaveDialog {
		return "セーブしますか？(Y/n/esc)"
	} else {
		// KIND
		table := fmt.Sprintf(
			"小学生: %d\n中高生男子: %d\n中高生女子: %d\n親: %d\nその他: %d\n\n",
			m.elemStudentCount, m.hsBoyCount, m.hsGirlCount, m.parentCount, m.otherCount,
		)
		// KIND
		usage := "[g] 小学生, [h] 中高生男子, [j] 中高生女子, [k] 親, [l] その他, [q] quit\n"
		return table + usage + m.message
	}
}

func main() {
	p := tea.NewProgram(counterModel{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
