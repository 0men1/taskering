package bubble


import (
	"tsk/bubble/list"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"tsk/bubble/styles"
	tsktasks "tsk/tasks"
)

type model struct{
	model *list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd)  {
	var cmd tea.Cmd
	m.model, cmd = m.model.Update(msg)
	return m, cmd
}


func (m model) View() string {
	return styles.Style.Render(m.model.View())	
}


func Run() {
	c := tsktasks.Init()
	m := model{
		model: list.New(c),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()
	
	if err != nil {
		fmt.Println("Error running program: %v", err)
	}
}
