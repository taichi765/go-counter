package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type personKind int

const (
	elemStudent personKind = iota
	hsBoy
	hsGirl
	parent
	other
)

func (k personKind) String() string {
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

type counterModel struct {
	elemStudentCount int
	hsBoyCount       int
	hsGirlCount      int
	parentCount      int
	otherCount       int

	showSaveDialog bool

	history history
	message string
}

func (m counterModel) Init() tea.Cmd {
	return nil
}

func (m counterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.showSaveDialog {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y":
				msg, err := m.history.Save()
				if err != nil {
					log.Fatal(err)
				}
				m.message = msg
				return m, tea.Quit
			case "n":
				return m, tea.Quit
			case "esc":
				m.showSaveDialog = false
				return m, nil
			}
		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.showSaveDialog = true
		case "backspace":
			last := m.history.Pop()
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
				log.Fatal("[Update] unknown kind")
			}
		case "ctrl+s":
			msg, err := m.history.Save()
			if err != nil {
				log.Fatal(err)
			}
			m.message = msg
		case "g":
			m.elemStudentCount++
			m.history.Add(historyEntry{time: time.Now(), kind: elemStudent})
			m.message = ""
		case "h":
			m.hsBoyCount++
			m.history.Add(historyEntry{time: time.Now(), kind: hsBoy})
			m.message = ""
		case "j":
			m.hsGirlCount++
			m.history.Add(historyEntry{time: time.Now(), kind: hsGirl})
			m.message = ""
		case "k":
			m.parentCount++
			m.history.Add(historyEntry{time: time.Now(), kind: parent})
			m.message = ""
		case "l":
			m.otherCount++
			m.history.Add(historyEntry{time: time.Now(), kind: other})
			m.message = ""
		}
	}
	return m, nil
}

func (m counterModel) View() string {
	if m.showSaveDialog {
		return "セーブしますか？(Y/n/esc)"
	} else {
		table := fmt.Sprintf(
			"小学生: %d\n中高生男子: %d\n中高生女子: %d\n親: %d\nその他: %d\n\n",
			m.elemStudentCount, m.hsBoyCount, m.hsGirlCount, m.parentCount, m.otherCount,
		)
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
