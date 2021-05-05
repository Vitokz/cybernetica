package model

import(
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	blurredSubmitButton = "[ " + blurredButtonStyle.Render("Отправить") + " ]"
	blurredButtonStyle = lipgloss.NewStyle()
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

type Input struct{
	Water textinput.Model
	Corn textinput.Model
	Milk textinput.Model
	Cup textinput.Model
	SubmitButton  string
}

func (i *Input) InitInput() {
	water := textinput.NewModel()
	water.Placeholder = "Вода в мл (макс: 5000)"
	water.Focus()
	water.PromptStyle = focusedStyle
	water.TextStyle = focusedStyle
	water.CharLimit = 4
    i.Water=water

	corn := textinput.NewModel()
	corn.Placeholder = "Зерно (макс: 900)"
	corn.CharLimit = 64
	corn.CharLimit = 3
    i.Corn=corn

	milk := textinput.NewModel()
	milk.Placeholder = "Молоко в мл (макс: 1000)"
	milk.CharLimit = 32
	milk.CharLimit = 4
    i.Milk=milk

	cup := textinput.NewModel()
	cup.Placeholder = "Стаканы (макс 50)"
	cup.CharLimit = 32
	cup.CharLimit = 2
    i.Cup=cup

	i.SubmitButton=blurredSubmitButton
}