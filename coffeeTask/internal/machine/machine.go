package machine

import (
	"errors"
	"fmt"
	"main/model"
	"main/proto"
	"main/repository"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type CoffeeMachine struct { // Струтура кофемашины
	Storage     model.StorageModel // Модель расходников машины
	Balance     int                // касса машины
	Stat        model.Stat         // статистика заказов
	Choise      int                // номер выбранного элемента
	Chosen      bool               // проверка на подтверждение
	ChoiseCoffe int                // номер выьранного коффе
	ChosenCoffe bool               // подтверждение что коффе выбрался
	Input       model.Input        // модель ввода
	InputIndex  int                // номер строки ввода
	Err         error              // ошибка
}

var (
	term               = termenv.ColorProfile()
	dot                = colorFg(" • ", "236")
	subtle             = makeFgStyle("241")
	focusedStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle            = lipgloss.NewStyle()

	focusedSubmitButton = "[ " + focusedStyle.Render("Отправить") + " ]"
	blurredSubmitButton = "[ " + blurredButtonStyle.Render("Отправить") + " ]"
)

func New(init *model.InitDate) (CoffeeMachine, error) {
	machine := CoffeeMachine{}

	setMilk(&machine, *init)
	setCorn(&machine, *init)
	setBalance(&machine, *init)
	setWater(&machine, *init)
	setCup(&machine, *init)
	setChoice(&machine)
	setChosen(&machine)
	setChoiceCoffee(&machine)
	setChosenCoffee(&machine)
	setInputIndex(&machine)
	setError(&machine)
	machine.Input.InitInput()
	return machine, nil
}
func (m *CoffeeMachine) Start() {
	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m *CoffeeMachine) Init() tea.Cmd {
	return textinput.Blink
}
func (m *CoffeeMachine) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "b":
			m.Chosen = false
		}
	}

	if !m.Chosen {
		return updateChoices(msg, m)
	}

	if m.Chosen && m.Choise == 4 {
		return m, tea.Quit
	}
	if m.Chosen && m.Choise == 0 && !m.ChosenCoffe {
		return updateChoicesCoffee(msg, m)
	}
	if m.Chosen && m.Choise == 2 {
		return updateFill(msg, m)
	}
	if m.Chosen && m.Choise == 1 {
		return m, nil
	}
	return m, nil
}
func (m *CoffeeMachine) View() string {
	var s string

	if !m.Chosen {
		s = choicesView(*m)
	}
	if m.Chosen && m.Choise == 0 && !m.ChosenCoffe {
		s = choicesViewCoffee(*m)
	}
	if m.Chosen && m.Choise == 1 {
		s = choicesViewStat(*m)
	}
	if m.Chosen && m.Choise == 4 {
		s = "Хотите закрыть?\n"
		s += subtle("enter: Ok") + dot + subtle("b:go back")
	}
	if m.ChosenCoffe {
		s = createViewCoffe(m)
		chist(m)
	}
	if m.Chosen && m.Choise == 2 {
		s = addFill(m)
	}
	if m.Chosen && m.Choise == 3 {
		m.Balance = 0
		s = "Вы успешно сняли весь баланс\n"
		s += subtle("j/k, up/down: select") + dot + subtle("q, esc: quit") + dot + subtle("b:go back")
		return s
	}
	return s
}

func updateChoices(msg tea.Msg, m *CoffeeMachine) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) { // Тут мы смотрим какую клавишу нажал пользователь и уже отталкиваемся от его действий

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down": // Опускает на один элемент ниже
			m.Choise += 1
			if m.Choise > 4 {
				m.Choise = 4
			}
		case "k", "up": //поднимает на один элемент выше
			m.Choise -= 1
			if m.Choise < 0 {
				m.Choise = 0
			}
		case "enter": // ставит в известность какой пункт был выбран
			m.Chosen = true
			return m, nil
		}
	}
	return m, nil
}

func updateChoicesCoffee(msg tea.Msg, m *CoffeeMachine) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) { // Тут мы смотрим какую клавишу нажал пользователь и уже отталкиваемся от его действий

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down": // Опускает на один элемент ниже
			m.ChoiseCoffe += 1
			if m.ChoiseCoffe > 2 {
				m.ChoiseCoffe = 2
			}
		case "k", "up": //поднимает на один элемент выше
			m.ChoiseCoffe -= 1
			if m.ChoiseCoffe < 0 {
				m.ChoiseCoffe = 0
			}
		case "enter": // ставит в известность какой пункт был выбран
			m.ChosenCoffe = true
			return m, nil
		}
	}
	return m, nil
}

func updateFill(msg tea.Msg, m *CoffeeMachine) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			inputs := []textinput.Model{
				m.Input.Water,
				m.Input.Milk,
				m.Input.Corn,
				m.Input.Cup,
			}
			s := msg.String()

			if s == "enter" && m.InputIndex == len(inputs) {
				plusWater(m)
				plusCorn(m)
				plusCup(m) //Добавить проверку на ошибки и применить соответсвующие действия
				plusMilk(m)
				if m.Err != nil {
					return m, cmd
				}
				m.Err = nil
				m.Chosen = false
				return m, cmd
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.InputIndex--
			} else {
				m.InputIndex++
			}

			if m.InputIndex > len(inputs) {
				m.InputIndex = 0
			} else if m.InputIndex < 0 {
				m.InputIndex = len(inputs)
			}

			for i := 0; i <= len(inputs)-1; i++ {
				if i == m.InputIndex {
					// Set focused state
					inputs[i].Focus()
					inputs[i].PromptStyle = focusedStyle
					inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				inputs[i].Blur()
				inputs[i].PromptStyle = noStyle
				inputs[i].TextStyle = noStyle
			}

			m.Input.Water = inputs[0]
			m.Input.Milk = inputs[1]
			m.Input.Corn = inputs[2]
			m.Input.Cup = inputs[3]

			if m.InputIndex == len(inputs) {
				m.Input.SubmitButton = focusedSubmitButton
			} else {
				m.Input.SubmitButton = blurredSubmitButton
			}

			return m, nil
		}
	}

	m, cmd = updateInputs(msg, m)
	return m, cmd
}

func choicesView(m CoffeeMachine) string {
	c := m.Choise

	tpl := "Добро пожаловать!\n\n"
	tpl += "Выберете действие!\n"
	tpl += "%s\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s",
		checkbox("Купить кофе", c == 0),
		checkbox("Посмотреть статистику", c == 1),
		checkbox("Добавить расходников", c == 2),
		checkbox("Забрать деньги", c == 3),
		checkbox("Закончить сессию", c == 4),
	)

	return fmt.Sprintf(tpl, choices)
}

func choicesViewCoffee(m CoffeeMachine) string {
	c := m.ChoiseCoffe
	tpl := ""
	tpl += "Выберете кофе:\n\n"
	tpl += "%s\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit") + dot + subtle("b:go back")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n",
		checkbox("Espresso 60р", c == 0),
		checkbox("Latte 110р", c == 1),
		checkbox("Cappuchino 140р", c == 2),
	)

	return fmt.Sprintf(tpl, choices)
}

func choicesViewStat(m CoffeeMachine) string {
	s := "В машине на данный момент:\n\n"
	s += fmt.Sprintf("%d мл воды\n", m.Storage.Water)
	s += fmt.Sprintf("%v мл молока\n", m.Storage.Milk)
	s += fmt.Sprintf("%v г кофейных зерен\n", m.Storage.CoffeCorn)
	s += fmt.Sprintf("%d стаканчиков\n", m.Storage.CupCoount)
	s += fmt.Sprintf("А в кассе тем временем %d\n\n\n", m.Balance)
	s += fmt.Sprintf("%d-espresso\n", m.Stat.Espresso)
	s += fmt.Sprintf("%d-latte\n", m.Stat.Latte)
	s += fmt.Sprintf("%d-cappuchino\n", m.Stat.Cappuchino)
	sumBuy := m.Stat.Espresso + m.Stat.Latte + m.Stat.Cappuchino
	cash := m.Stat.Espresso*repository.Espresso.Price + m.Stat.Latte*repository.Latte.Price + m.Stat.Cappuchino*repository.Cappuchino.Price
	s += fmt.Sprintf("Всего напитков продано %d на %d путинских дублонов\n\n", sumBuy, cash)
	s += subtle("q, esc: quit") + dot + subtle("b:go back")

	return s
}

func createViewCoffe(m *CoffeeMachine) string {
	switch a := m.ChoiseCoffe; a {
	case 0:
		err := checkStorage(repository.Espresso, m.Storage) //Проверяет хватает ли ресурсов машины для создания напитка
		if err != nil {
			m.ChosenCoffe = false
			return fmt.Sprintf("Попробуйте выбрать что-нибудь другое\n%s\nДля выхода нажмите b", err)
		}
		m.Balance += repository.Espresso.Price        //Плюсуем к кассе сумму заказа
		m.Stat.Espresso += 1                          //Добаляем заказ в статистику определенных напитков
		minusStorage(repository.Espresso, &m.Storage) // Минусуем ресурсы которые потребовались для выполнения заказа
	case 1:
		err := checkStorage(repository.Latte, m.Storage)
		if err != nil {
			return fmt.Sprintf("Попробуйте выбрать что-нибудь другое\n%s\nДля выхода нажмите b", err)
		}
		m.Balance += repository.Latte.Price
		m.Stat.Latte += 1
		minusStorage(repository.Latte, &m.Storage)
	case 2:
		err := checkStorage(repository.Cappuchino, m.Storage)
		if err != nil {
			return fmt.Sprintf("Попробуйте выбрать что-нибудь другое\n%s\nДля выхода нажмите b", err)
		}
		m.Balance += repository.Cappuchino.Price
		m.Stat.Cappuchino += 1
		minusStorage(repository.Cappuchino, &m.Storage)
	}

	return "Нажмите b чтобы забрать ваш кофе"
}

func addFill(m *CoffeeMachine) string {

	s := "\n"
	if m.Err != nil {
		s += fmt.Sprintf("%v\n", m.Err)
	}
	inputs := []string{
		m.Input.Water.View(),
		m.Input.Milk.View(),
		m.Input.Corn.View(),
		m.Input.Cup.View(),
	}

	for i := 0; i < len(inputs); i++ {
		s += inputs[i]
		if i < len(inputs)-1 {
			s += "\n"
		}
	}

	s += "\n\n" + m.Input.SubmitButton + "\n"
	s += subtle("j/k, up/down: select") + dot + subtle("q, esc: quit") + dot + subtle("b:go back")
	return s
}

func chist(m *CoffeeMachine) tea.Model {
	m.ChoiseCoffe = 0
	m.ChosenCoffe = false
	m.Chosen = false
	return m
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("• "+label, "#FFDEAD")
	}
	return fmt.Sprintf("  %s", label)
}

func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

func updateInputs(msg tea.Msg, m *CoffeeMachine) (*CoffeeMachine, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.Input.Water, cmd = m.Input.Water.Update(msg)
	cmds = append(cmds, cmd)

	m.Input.Milk, cmd = m.Input.Milk.Update(msg)
	cmds = append(cmds, cmd)

	m.Input.Corn, cmd = m.Input.Corn.Update(msg)
	cmds = append(cmds, cmd)

	m.Input.Cup, cmd = m.Input.Cup.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

//Для функции Buy
func checkStorage(name model.CoffeeModel, m model.StorageModel) error { //Проверка на достаточность ресурсов
	switch {
	case m.CoffeCorn < name.CoffeCorn:
		err := errors.New("к сожалению в данный момент недостаточно кофе для приготовления")
		return err
	case m.Milk < name.Milk:
		err := errors.New("к сожалению в данный момент недостаточно молока для приготовления")
		return err
	case m.Water < name.Water:
		err := errors.New("к сожалению в данный момент недостаточно воды для приготовления")
		return err
	case m.CupCoount < 1:
		err := errors.New("к сожалению в данный момент нет стаканчиков для приготовления")
		return err
	}
	return nil
}

func minusStorage(name model.CoffeeModel, m *model.StorageModel) { // Минусуем затрачиваемые ресурсы
	m.CoffeCorn -= name.CoffeCorn
	m.Milk -= name.Milk
	m.CupCoount--
	m.Water -= name.Water
}

//Относятся к new
func setMilk(m *CoffeeMachine, init model.InitDate) {
	m.Storage.Milk = init.Milk
}
func setCorn(m *CoffeeMachine, init model.InitDate) {
	m.Storage.CoffeCorn = init.CoffeCorn
}
func setWater(m *CoffeeMachine, init model.InitDate) {
	m.Storage.Water = init.Water
}
func setCup(m *CoffeeMachine, init model.InitDate) {
	m.Storage.CupCoount = init.CupCount
}
func setBalance(m *CoffeeMachine, init model.InitDate) {
	m.Balance = init.Balance
}
func setChoice(m *CoffeeMachine) {
	m.Choise = 0
}
func setChosen(m *CoffeeMachine) {
	m.Chosen = false
}
func setChoiceCoffee(m *CoffeeMachine) {
	m.ChoiseCoffe = 0
}
func setChosenCoffee(m *CoffeeMachine) {
	m.ChosenCoffe = false
}
func setInputIndex(m *CoffeeMachine) {
	m.InputIndex = 0
}
func setError(m *CoffeeMachine) {
	m.Err = nil
}

//
//Относятся к fill
func plusWater(m *CoffeeMachine) {
	wt, err := strconv.Atoi(m.Input.Water.Value())
	if err != nil {
		m.Err = errors.New("нельзя вводить буквы")
		return
	}
	if wt+m.Storage.Water > proto.WATER {
		m.Err = fmt.Errorf("в поле воды значение больше возможного. Вы можете ввести %v", proto.WATER-m.Storage.Water)
		return
	}
	m.Storage.Water += wt
}

func plusCorn(m *CoffeeMachine) {
	corn, err := strconv.Atoi(m.Input.Corn.Value())
	if err != nil {
		m.Err = errors.New("нельзя вводить буквы")
		return
	}
	if corn+m.Storage.CoffeCorn > proto.CORN {
		m.Err = fmt.Errorf("в поле Зерна значение больше возможного. Вы можете ввести %v", proto.CORN-m.Storage.CoffeCorn)
		return
	}
	m.Storage.CoffeCorn += corn
}
func plusMilk(m *CoffeeMachine) {
	milk, err := strconv.Atoi(m.Input.Milk.Value())
	if err != nil {
		m.Err = errors.New("нельзя вводить буквы")
		return
	}
	if milk+m.Storage.Milk > proto.MILK {
		m.Err = fmt.Errorf("в поле молока значение больше возможного. Вы можете ввести %v", proto.MILK-m.Storage.Milk)
		return
	}
	m.Storage.Milk += milk
}
func plusCup(m *CoffeeMachine) {
	cup, err := strconv.Atoi(m.Input.Cup.Value())
	if err != nil {
		m.Err = errors.New("нельзя вводить буквы")
		return
	}
	if cup+m.Storage.CupCoount > proto.CUP {
		m.Err = fmt.Errorf("в поле стаканчиков значение больше возможного. Вы можете ввести %v", proto.CUP-m.Storage.CupCoount)
		return
	}
	m.Storage.CupCoount += cup
}
