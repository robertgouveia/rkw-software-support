package menu

import (
	"fmt"
	"log"

	"github.com/robertgouveia/do-my-job/database"
	"github.com/robertgouveia/do-my-job/tea"
)

type Param struct {
	Title string
	Value any
	Name  string
}

type Select struct {
	Title      string
	Values     []string
	Selected   any
	Name       string
	UseIndex   bool
	ValueMap   []any
	DefaultIdx int
}

type Script struct {
	Title      string
	Params     []Param
	Select     []Select
	ServerName string
	Statement  string
}

func ScriptMenu(mainMenu *tea.TeaModel) *tea.TeaModel {
	scriptMenu := tea.Create("Scripts")
	mainMenu.AddSubmenu("Scripts", scriptMenu)

	scripts := []Script{
		{
			Title: "Dispute Status Change",
			Params: []Param{
				{
					Title: "Dispute ID",
					Name:  "IssueID",
				},
			},
			Select: []Select{
				{
					Title:    "Status",
					Name:     "Status",
					UseIndex: false,
					Values: []string{
						"Logged",
						"Investigating",
						"Investigation Approved",
						"Investigation Rejected",
						"Sent To Accounts Approved",
						"Sent To Accounts Rejected",
						"Accounts Credit Created",
						"Accounts Credit Rejected",
						"Cancelled",
						"Pre-Investigation",
						"Awaiting Debit Note",
						"Awaiting RAN",
						"Disputed",
					},
				},
			},
			ServerName: "RKW Data Warehouse",
			Statement:  database.DisputeChange,
		},
		{
			Title: "Shipping Agent Service Change",
			Params: []Param{
				{
					Title: "Sales Order Number",
					Name:  "OrderNo",
				},
			},
			ServerName: "RKW Level 1",
			Statement:  database.ShippingChange,
		},
	}

	for _, script := range scripts {
		scriptMenu.AddSubmenu(script.Title, scriptTemplate(script.Title, script.Params, script.Select, script.Statement, script.ServerName))
	}

	return scriptMenu
}

func scriptTemplate(title string, params []Param, s []Select, statement, sname string) *tea.TeaModel {
	rkwScriptMenu := tea.Create(title)

	for i := range params {
		param := &params[i]

		rkwScriptMenu.AddTextInput(
			fmt.Sprintf("Set %s", param.Title),
			"Enter Param:",
			fmt.Sprintf("The current value is %v", param.Value),
			func(input string) {
				param.Value = input
				fmt.Printf("%s set to: %s and saved\n", param.Title, input)
			},
		)
	}

	for i := range s {
		rkwScriptMenu.AddSubmenu(s[i].Title, selectTemplate(&s[i]))
	}

	rkwScriptMenu.AddMenuItem("Execute", func() string {
		str := "Executing: "

		namedParams := make(map[string]interface{})

		for _, param := range params {
			if param.Name != "" && param.Value != nil {
				namedParams[param.Name] = param.Value
				str += fmt.Sprintf(" [%s:%v] ", param.Title, param.Value)
			}
		}

		for _, option := range s {
			if option.Name != "" && option.Selected != nil {
				var paramValue any
				var displayValue any

				if index, ok := option.Selected.(int); ok {
					if index >= 0 && index < len(option.Values) {
						displayValue = option.Values[index]

						if option.UseIndex {
							if len(option.ValueMap) > index {
								paramValue = option.ValueMap[index]
							} else {
								paramValue = index + 1
							}
						} else {
							paramValue = option.Values[index]
						}
					}
				} else if strValue, ok := option.Selected.(string); ok {
					displayValue = strValue
					paramValue = strValue

					if option.UseIndex {
						for i, v := range option.Values {
							if v == strValue {
								if len(option.ValueMap) > i {
									paramValue = option.ValueMap[i]
								} else {
									paramValue = i + 1
								}
								break
							}
						}
					}
				} else {
					displayValue = option.Selected
					paramValue = option.Selected
				}

				namedParams[option.Name] = paramValue
				str += fmt.Sprintf(" [%s:%v] ", option.Title, displayValue)
			}
		}

		db, err := database.Connect(sname)
		if err != nil {
			log.Fatalf("Error connecting to DB: %s", err.Error())
		}

		res, debugInfo, err := database.ExecuteWithNamedParams(db, statement, namedParams)
		if err != nil {
			log.Fatalf("Error executing DB statement: %s, Variables: %s", err.Error(), debugInfo)
		}

		rows, err := res.RowsAffected()
		if err != nil {
			log.Fatalf("Error fetching rows affected: %s", err.Error())
		}

		return str + fmt.Sprintf(" Rows Affected: %d Params: %s", rows, debugInfo)
	})

	return rkwScriptMenu
}

func selectTemplate(s *Select) *tea.TeaModel {
	rkwSelectMenu := tea.Create(s.Title)

	for i, option := range s.Values {
		value := option
		index := i

		rkwSelectMenu.AddMenuItem(value, func() string {
			s.Selected = index
			return "back"
		})
	}

	return rkwSelectMenu
}
