package cmd

import (
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	ResBody string
	ErrMsg  struct{ err error }
)

func (e ErrMsg) Error() string { return e.err.Error() }

// GetCmd returns a cmd function that performs a simple GET command
func GetCmd(url string) func() tea.Msg {

	return func() tea.Msg {
		c := &http.Client{Timeout: 10 * time.Second}
		res, err := c.Get(url)

		if err != nil {
			return ErrMsg{err}
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return ErrMsg{err}
		}
		return ResBody(string(body))
	}
}
