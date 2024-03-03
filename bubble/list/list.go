package list 

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	tsktasks "tsk/tasks"
	"tsk/bubble/styles"
)


type Model struct {
	Categories tsktasks.Categories
	Items 	   []tsktasks.Item
	cursor     int
	catcursor  int
	selected   map[int]map[int]struct{}	
	chosenKey  string
}


func New(c *tsktasks.Categories) *Model {
	m := Model {
		Categories: *c,
		Items: c.CatList[c.CurrentIndex].Items,
		cursor: 0,
		catcursor: 0,
		selected: make(map[int]map[int]struct{}),
	}

	return &m
}


func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.KeyMsg:
	m.chosenKey = msg.String()
	switch msg.String() {

        case "ctrl+c", "q":
            return m, tea.Quit

    	case "left", "h":
		if (m.catcursor > 0) {
			m.catcursor--
			m.cursor = 0
			m.SetItems()
		}

	case "right", "l":
		if (m.catcursor < len(m.Categories.CatList)-1) {
			m.catcursor++
			m.cursor = 0
			m.SetItems()
		}

        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }

        case "down", "j":
            if m.cursor < len(m.Items)-1 {
                m.cursor++
            }

        case "enter", " ":

	    if _, ok := m.selected[m.catcursor]; !ok { // Initialize the map of int structs if it doesnt exist
		m.selected[m.catcursor] = make(map[int]struct{})
	    }

            _, ok := m.selected[m.catcursor][m.cursor]
            if ok {
                delete(m.selected[m.catcursor], m.cursor)
            } else {
                m.selected[m.catcursor][m.cursor] = struct{}{}
            }
        }
    }
    return m, nil
}



func (m *Model) View() string {
    s := "What should we buy at the market?\n\n"
	chosenkey := m.chosenKey

    s += m.SetCategoryString() + "\n\n"

    for i, choice := range m.Items {

        cursor := " " // no cursor
        if m.cursor == i {
            cursor = ">" // cursor!
        }

        checked := " " // not selected
        if _, ok := m.selected[m.catcursor][i]; ok {
            checked = "x" // selected!
        }

	s += fmt.Sprintf(("Chosen key: %s\n %s [%s] %s\n"),chosenkey, cursor, checked, choice.Title)
    }
    s += "\nPress q to quit.\n"
    return s
}

func (m *Model) SetItems() {
	m.Items = m.Categories.CatList[m.catcursor].Items	
}


func (m *Model) SetCategoryString() (string) {
	s := "Categories:\n"
	for i, cat := range m.Categories.CatList {
		if m.Categories.CatList[m.catcursor].Id == cat.Id {
			s += fmt.Sprintf(("  %d: %s  "), i,  styles.CatStyle.Render(cat.Title))	
		} else {
			s += fmt.Sprintf(("  %d: %s  "), i,  cat.Title)	
		}
	}
	return s
}


