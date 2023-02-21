package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"salt-tortilla/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var httpVerbs = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
}

// [ ] input -> []textinput.Model
// en initModel crear los inputs necesarios
// [ ] añadir backwards/forwards movement
// [ ] make request
// [ ] pass everything to utils, etc.
// [ ] cobra??
type model struct {
	textInput textinput.Model
	typing    bool
	url       string
	httpVerb  string
	headers   []string
	body      string
	stage     int
	cursor    int
}

type Stage int

const (
	Url Stage = iota
	HttpVerb
	Headers
	Body
	Fetch
)

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) SetValue(v string, s Stage) {
	switch s {
	case Url:
		m.url = v
	case HttpVerb:
		m.httpVerb = v
	case Headers:
		m.headers = append(m.headers, v)
	case Body:
		m.body = v
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.typing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				return m, tea.Quit
			case "tab":
				m.SetValue(m.textInput.Value(), Stage(m.stage))
				m.textInput.SetValue("")

			case "enter":
				m.SetValue(m.textInput.Value(), Stage(m.stage))
				m.stage++
				if Stage(m.stage) == HttpVerb || Stage(m.stage) == Fetch {
					m.typing = false
				}
				m.textInput.SetValue("")
				m.textInput.Focus()
			}
		}
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				return m, tea.Quit
			case "up", "k":
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(httpVerbs) - 1
				}
			case "down", "j":
				m.cursor++
				if m.cursor >= len(httpVerbs) {
					m.cursor = 0
				}
			case "enter":
				m.httpVerb = httpVerbs[m.cursor]
				m.stage++
				m.typing = true

			}
		}
	}
	return m, nil
}

func (m model) View() string {
	headersHead := fmt.Sprintf("Input the Headers:\n\nURL: %s\nHTTP Verb: %s\nHeaders:%s\n\n%s", m.url, m.httpVerb, utils.PrintSlice(m.headers), m.textInput.View())
	bodyHead := fmt.Sprintf("Input the Body:\n\nURL: %s\nHTTP Verb: %s\nHeaders: %s\n\n%s", m.url, m.httpVerb, m.headers, m.textInput.View())

	if m.typing {
		switch Stage(m.stage) {
		case Url:
			return fmt.Sprintf("Input the URL:\n\n%s", m.textInput.View())
			// return utils.URLString(m.textInput.View)
		case Headers:
			// utils.HeadersString(m.url, m.httpVerb, m.headers, m.textInput.View())
			return fmt.Sprintf("%s\n\n(press tab to input new Header, press enter to input Body)", headersHead)
		case Body:
			return bodyHead
		}
	} else if Stage(m.stage) == HttpVerb {
		s := strings.Builder{}
		s.WriteString(utils.HttpVerbString(m.url))

		for i := 0; i < len(httpVerbs); i++ {
			if m.cursor == i {
				s.WriteString("[*] ")
			} else {
				s.WriteString("[ ] ")
			}
			s.WriteString(httpVerbs[i])
			s.WriteString("\n")
		}
		s.WriteString("\n(press q, Ctrl+C or esc to exit)\n")
		return s.String()
	} else {
		s := fmt.Sprintf("\nURL: %s\nHTTP Verb: %s\nHeaders: %s\nBody: %s\n\n", m.url, m.httpVerb, utils.PrintSlice(m.headers), m.body)
		return s
	}
	return "Press Ctrl+C to exit"
}

func main() {
	t := textinput.New()
	t.Placeholder = "http://localhost:8000"
	t.Focus()
	initialModel := model{
		textInput: t,
		typing:    true,
	}
	m, err := tea.NewProgram(initialModel).Run()

	if err != nil {
		log.Fatalf("Oh no, something went wrong!\n%s", err)
	}

	if m, ok := m.(model); ok {
		fmt.Printf("Url: %s\nHttp Verb: %s\nHeaders: %v\nBody: %s\n\n", m.url, m.httpVerb, m.headers, m.body)
	}
}
