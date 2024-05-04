package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
	Children  []int  `json:"children"`
	Id        int    `json:"id"`
	Points    int    `json:"points"`
}

type SearchResult struct {
	Hits []Item `json:"hits"`
}

type editorFinishedMsg struct{ err error }

func openUrl(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	exec.Command(cmd, args...).Start() //nolint:gosec
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type listItem struct {
	title, desc string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(listItem)
			if ok && len(i.desc) > 0 {
				openUrl(i.desc)
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	endpoint := "http://hn.algolia.com/api/v1/"
	stories := getFrontPage(endpoint)

	items := []list.Item{}

	for _, v := range stories {
		items = append(items, listItem{title: v.Title, desc: v.Url})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "hnt - Hacker News Terminal"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func getFrontPage(endpoint string) []Item {
	url := endpoint + "search?tags=front_page"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	// We Read the response body on the line below
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var searchResult SearchResult
	err2 := json.Unmarshal(body, &searchResult)
	if err2 != nil {
		fmt.Println(err2)
	}

	return searchResult.Hits
}

func getItem(endpoint string, id int) Item {
	url := endpoint + "items/" + fmt.Sprint(id)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var item Item
	json.Unmarshal(body, &item)
	return item
}
