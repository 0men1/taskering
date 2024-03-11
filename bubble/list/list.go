package list 

import (
	tea "github.com/charmbracelet/bubbletea"
	tsktasks "tsk/tasks"
	"tsk/bubble/styles"
	"tsk/utils"
	"time"
	"google.golang.org/api/tasks/v1"
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbles/textinput"
)


type Model struct {
	Categories 	tsktasks.Categories

	entryMode  	bool	   
	deleteMode 	bool
	inputs	   	[]textinput.Model

	cursor     	int
	catcursor	int
	inputcursor 	int


	selected   	map[int]map[int]struct{}	
	viewport 	viewport.Model
}


func New(c *tsktasks.Categories) *Model {
	m := Model {
		Categories: *c,

		cursor: 0,
		catcursor: 0,
		inputcursor: 0,

		selected: make(map[int]map[int]struct{}),
		inputs: make([]textinput.Model, 6),
	}


	//Initialize the map of selected ints
	for i := range m.Categories.CatList {
		m.selected[i] = make(map[int]struct{})	
	}

	for i := range m.inputs {
		var t textinput.Model
		t = textinput.New()
		switch i {
		case 0:
			t.Placeholder = "Title"
			t.CharLimit = 30
		case 1:
			t.Placeholder = "dd (# only)"
			t.CharLimit = 2
		case 2:
			t.Placeholder = "mm (# only)"
			t.CharLimit = 2
		case 3:
			t.Placeholder = "yyyy (# only)"
			t.CharLimit = 4
		case 4:
			t.Placeholder = "Notes..."
			t.CharLimit = 100

		case 5:
			t.Placeholder = "(y/n)"
			t.CharLimit = 1
		}

		m.inputs[i] = t
	}

	return &m
}


func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	if m.entryMode {
		for i := range m.inputs {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		}
	}
	if m.deleteMode {
		for i := range m.inputs {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		}
	}


    	switch msg := msg.(type) {
    	case tea.KeyMsg:

		switch msg.Type {

		case tea.KeyShiftTab:
			if m.entryMode {
				m.prevInput()
				m.refocus()
			}
		case tea.KeyTab:
			if m.entryMode {
				m.nextInput()
				m.refocus()
			}
		}

		switch msg.String() {
        	case "ctrl+c":
            		return m, tea.Quit

		case "q":
			if(!m.entryMode && !m.deleteMode) {
				return m, tea.Quit
			}

		case "left", "h" :
			if (m.catcursor > 0 && !m.entryMode) {
				m.catcursor--
				m.cursor = 0
			}

		case "right", "l":
			if (m.catcursor < len(m.Categories.CatList)-1 && !m.entryMode) {
				m.catcursor++
				m.cursor = 0
			}

		case "up", "k":
		    if m.cursor > 0 && !m.entryMode {
			m.cursor--
		    }

		case "down", "j":
		    if m.cursor < len(m.visibleItems())-1 && !m.entryMode {
			m.cursor++
		    }

		case "enter", " ":
			if !m.entryMode {
				_, ok := m.selected[m.catcursor][m.cursor]

				if ok {
					delete(m.selected[m.catcursor], m.cursor)
				} else {
					m.selected[m.catcursor][m.cursor] = struct{}{}
				}
			}

			if m.entryMode && msg.String() != " " && !m.deleteMode{
				m.insertItemIntoCalendar()
				m.entryMode = false
				m.refocus()
			}
			if m.deleteMode && msg.String() != " " && !m.entryMode{
				m.deleteItemFromCalendar(m.cursor)
				m.deleteMode = false
			}


		case "+":
			if !m.deleteMode {
				m.entryMode = true
				m.inputs[0].Focus()
			}


		case "D":
			if !m.entryMode {
				m.deleteMode = true
				m.inputs[5].Focus()
			}


		case "esc":
			if m.entryMode {
				m.entryMode = false
				m.refocus()
			}
			if m.deleteMode {
				m.deleteMode = false
			}
		}
	}

	return m, tea.Batch(cmds...)
}


func (m *Model) View() string {
	s := "--------------------TODO LIST--------------------\n\n"

	s += m.SetCategoryString() + "\n\n"


	for i, choice := range m.visibleItems() {
		cursor := " " // no cursor
		if m.cursor == i {
		    cursor = ">" // cursor!
		}
		checked := " " // not selected
		if _, ok := m.selected[m.catcursor][i]; ok {
		    checked = "x" // selected!
		}

		s += fmt.Sprintf((" %s [%s] %s\n"), cursor, checked, choice.Title)
		t, err := time.Parse(time.RFC3339, choice.Due)
		if err != nil {
			s += fmt.Sprintf("	    Could not get the due date! \n\n\n")
		} else {
			s += fmt.Sprintf("          Due: %s, %d %d \n\n\n", 
			t.Month(),
			t.Day(),
			t.Year())
		}
	}




	if (m.deleteMode) {
		s += fmt.Sprintf("Are you sure you want to delete this task? %s \n", m.inputs[5].View())
	}


	if (m.entryMode) {
		s += fmt.Sprintf(`
		%s

		%s %s %s 

		%s
		`, 			       
			m.inputs[0].View(), 
			m.inputs[1].View(),
			m.inputs[2].View(), 
			m.inputs[3].View(), 
			m.inputs[4].View())
	}

	s += "                     ----------------------CONTROLS----------------------\n"
	s += "   (+) - Add a new task  | (D) - Delete a task    | (U) - Update a task\n" +
	     "   (h) - Left 1 Category | (l) - Right 1 Category | (j) - Down 1 task | (k) - Up 1 Task\n"

	s += "\nPress q to quit.\n"
	return s
}


func (m *Model) insertItemIntoCalendar() {
	title := m.inputs[0].Value()
	due := m.inputs[3].Value() + "-" + m.inputs[2].Value() + "-" + m.inputs[1].Value() + "T00:00:00Z"
	notes := m.inputs[4].Value()
	myTask := tsktasks.MakeTask(title, due, notes)

	newTask, err := tsktasks.InsertTask(m.Categories.CatList[m.catcursor].Id, myTask)
	if err != nil {
		fmt.Println("There was an error inserting the task: %v", err)
	}
	m.insertItemIntoModel(*newTask)
}


func (m *Model) insertItemIntoModel(task tasks.Task) {
	m.Categories.CatList[m.catcursor].Items = append(m.Categories.CatList[m.catcursor].Items, &task)
}


func (m *Model) deleteItemFromCalendar(index int) {
	if m.inputs[5].Value() == "y" {
		taskId := m.Categories.CatList[m.catcursor].Items[index].Id
		taskListId := m.Categories.CatList[m.catcursor].Id
		err := tsktasks.DeleteTask(taskListId, taskId)
		if err != nil {
			fmt.Printf("There was an error deleting the node: %v", err)
		}
		m.deleteItemFromModel(index)

	} else {
		m.deleteMode = false
	}
}


func (m *Model) deleteItemFromModel(index int) {
	m.Categories.CatList[m.catcursor].Items = utils.DeleteElement((m.Categories.CatList[m.catcursor].Items), index)
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


func (m *Model) nextInput() {
	m.inputcursor = (m.inputcursor + 1) % len(m.inputs)
}


func (m *Model) prevInput() {
	m.inputcursor--
	if m.inputcursor < 0 {
		m.inputcursor = len(m.inputs) - 1
	}
}


func (m *Model) refocus() {
	for i := range m.inputs {
		if i < 5 {
			if !m.entryMode {
				m.inputs[i].Reset()
			} else if i != m.inputcursor && m.entryMode {
				m.inputs[i].Blur()
			}
		}
	}
	if m.entryMode {
		m.inputs[m.inputcursor].Focus()
	} else {
		m.inputs[m.inputcursor].Blur()
		m.inputcursor = 0
	}
}


func (m *Model) visibleItems() []*tasks.Task {
	return m.Categories.CatList[m.catcursor].Items
}
