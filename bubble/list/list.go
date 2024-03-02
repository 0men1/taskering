package list 

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	tsktasks "tsk/tasks"
)


type Model struct {
	Items 	  []tsktasks.Item
	cursor    int
	selected  map[int]struct{}
	chosenKey string
}


func New(c *tsktasks.Categories) *Model {
	m := Model {
		Items: (c.CatMap[c.CurrentIndex]),
		selected: make(map[int]struct{}),
	}
	return &m
}


func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
    switch msg := msg.(type) {

    // Is it a key press?
    case tea.KeyMsg:
	m.chosenKey = msg.String()
        // Cool, what was the actual key pressed?
	switch msg.String() {

        // These keys should exit the program.
        case "ctrl+c", "q":
            return m, tea.Quit

        // The "up" and "k" keys move the cursor up
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }

        // The "down" and "j" keys move the cursor down
        case "down", "j":
            if m.cursor < len(m.Items)-1 {
                m.cursor++
            }

        // The "enter" key and the spacebar (a literal space) toggle
        // the selected state for the item that the cursor is pointing at.
        case "enter", " ":
            _, ok := m.selected[m.cursor]
            if ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }
        }
    }

    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return m, nil
}


func (m *Model) View() string {
    // The header
    s := "What should we buy at the market?\n\n"
	chosenkey := m.chosenKey
    // Iterate over our choices
    for i, choice := range m.Items {

        // Is the cursor pointing at this choice?
        cursor := " " // no cursor
        if m.cursor == i {
            cursor = ">" // cursor!
        }

        // Is this choice selected?
        checked := " " // not selected
        if _, ok := m.selected[i]; ok {
            checked = "x" // selected!
        }

        // Render the row
	s += fmt.Sprintf(("Chosen key: %s\n %s [%s] %s\n"),chosenkey, cursor, checked, choice.Title)
    }

    // The footer
    s += "\nPress q to quit.\n"

    // Send the UI for rendering
    return s
}

