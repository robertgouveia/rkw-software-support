package tea

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	bubble "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuItemType int

const (
	ContentItem MenuItemType = iota
	SubmenuItem
	TextInputItem
)

type MenuItem struct {
	Title     string
	Content   func() string
	SubMenu   *TeaModel
	ItemType  MenuItemType
	OnSubmit  func(string)
	Prompt    string
	InputDesc string
}

type TextInputModel struct {
	Parent      *TeaModel
	TextInput   textinput.Model
	Title       string
	Prompt      string
	Description string
	OnSubmit    func(string)
}

func NewTextInputModel(parent *TeaModel, title, prompt, description string, onSubmit func(string)) *TextInputModel {
	ti := textinput.New()
	ti.Placeholder = "Enter text..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	return &TextInputModel{
		Parent:      parent,
		TextInput:   ti,
		Title:       title,
		Prompt:      prompt,
		Description: description,
		OnSubmit:    onSubmit,
	}
}

func (m TextInputModel) Init() bubble.Cmd {
	return textinput.Blink
}

func (m *TextInputModel) Update(msg bubble.Msg) (bubble.Model, bubble.Cmd) {
	var cmd bubble.Cmd

	switch msg := msg.(type) {
	case bubble.KeyMsg:
		switch msg.Type {
		case bubble.KeyEnter:
			if m.OnSubmit != nil {
				m.OnSubmit(m.TextInput.Value())
			}
			return m.Parent, nil
		case bubble.KeyEsc:
			return m.Parent, nil
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m TextInputModel) View() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s\n\n", m.Title))

	if m.Description != "" {
		b.WriteString(fmt.Sprintf("%s\n\n", m.Description))
	}

	b.WriteString(fmt.Sprintf("%s\n\n", m.Prompt))
	b.WriteString(m.TextInput.View())
	b.WriteString("\n\nPress Enter to submit, Esc to cancel")

	return b.String()
}

type TeaModel struct {
	MenuItems []MenuItem
	Title     string
	Parent    *TeaModel

	SelectedMenu string

	Cursor   int
	Selected int
	Quitting bool
	Back     bool

	TitleStyle  lipgloss.Style
	ItemStyle   lipgloss.Style
	CursorStyle lipgloss.Style
}

func (m TeaModel) Init() bubble.Cmd {
	return nil
}

func Create(title string) *TeaModel {
	return &TeaModel{
		MenuItems:   []MenuItem{},
		Title:       title,
		TitleStyle:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")),
		ItemStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD")),
		CursorStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF875F")),
	}
}

func (m *TeaModel) AddMenuItem(title string, contentFunc func() string) {
	m.MenuItems = append(m.MenuItems, MenuItem{
		Title:    title,
		Content:  contentFunc,
		ItemType: ContentItem,
	})
}

func (m *TeaModel) AddSubmenu(title string, submenu *TeaModel) {
	submenu.Parent = m
	m.MenuItems = append(m.MenuItems, MenuItem{
		Title:    title,
		SubMenu:  submenu,
		ItemType: SubmenuItem,
	})
}

func (m *TeaModel) AddTextInput(title, prompt, description string, onSubmit func(string)) {
	m.MenuItems = append(m.MenuItems, MenuItem{
		Title:     title,
		ItemType:  TextInputItem,
		OnSubmit:  onSubmit,
		Prompt:    prompt,
		InputDesc: description,
	})
}

func (m *TeaModel) Update(msg bubble.Msg) (bubble.Model, bubble.Cmd) {
	if m.Back {
		m.Back = false
		return m.Parent, nil
	}

	switch msg := msg.(type) {
	case bubble.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, quitAfterDelay()
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.MenuItems)-1 {
				m.Cursor++
			}
		case "enter", " ":
			m.Selected = m.Cursor
			selectedItem := m.MenuItems[m.Selected]

			switch selectedItem.ItemType {
			case SubmenuItem:
				return selectedItem.SubMenu, nil
			case ContentItem:
				if selectedItem.Content != nil {
					result := selectedItem.Content()

					if result == "back" && m.Parent != nil {
						return m.Parent, nil
					}

					m.SelectedMenu = result
					m.Quitting = true
					return m, quitAfterDelay()
				}
			case TextInputItem:
				inputModel := NewTextInputModel(
					m,
					selectedItem.Title,
					selectedItem.Prompt,
					selectedItem.InputDesc,
					selectedItem.OnSubmit,
				)
				m.Parent.Cursor = 0
				m.Parent.Selected = 0
				return inputModel, textinput.Blink
			}
		case "backspace", "esc", "left", "h":
			if m.Parent != nil {
				return m.Parent, nil
			}
		}
	}

	return m, nil
}

func (m TeaModel) View() string {
	if m.Quitting {
		return "Exiting...\n"
	}

	s := fmt.Sprintf("%s\n\n", m.TitleStyle.Render(m.Title))

	for i, item := range m.MenuItems {
		cursor := " "
		if m.Cursor == i {
			cursor = m.CursorStyle.Render(">")
		}

		indicator := ""
		switch item.ItemType {
		case SubmenuItem:
			indicator = " ▶"
		case TextInputItem:
			indicator = " ✎"
		}

		s += fmt.Sprintf("%s [%s]%s\n", cursor, m.ItemStyle.Render(item.Title), indicator)
	}

	s += "\n(↑/↓) Navigate   (Enter) Select   "
	if m.Parent != nil {
		s += "(Esc) Back   "
	}
	s += "(q) Quit\n"

	return s
}

func (m *TeaModel) Run() (bubble.Model, error) {
	p := bubble.NewProgram(m)
	model, err := p.Run()

	if teaModel, ok := model.(*TeaModel); ok && teaModel.SelectedMenu != "" {
		fmt.Println(teaModel.SelectedMenu)
	}

	return model, err
}

func quitAfterDelay() bubble.Cmd {
	return func() bubble.Msg {
		time.Sleep(500 * time.Millisecond)
		return bubble.Quit()
	}
}
