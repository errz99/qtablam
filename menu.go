package qtablam

import (
	"github.com/mappu/miqt/qt"
)

func initMenuFields(columns []Column) *qt.QMenu {
	fieldsMenu := qt.NewQMenu3("F&ields")

	titles := make([]string, 0, 16)
	shortcuts := make([]string, 0, 16)
	icons := make([]string, 0, 16)
	checkables := make([]bool, 0, 16)
	checks := make([]bool, 0, 16)
	cbs := make([]func(), 0, 16)

	for _, col := range columns {
		titles = append(titles, col.title)
		shortcuts = append(shortcuts, "")
		icons = append(icons, "")
		checkables = append(checkables, true)
		checks = append(checks, col.visible)
		cbs = append(cbs, onFieldsMenu)
	}

	prepareMenu(fieldsMenu, titles, shortcuts, icons, cbs)
	addChecksToMenu(fieldsMenu, checkables, checks)

	return fieldsMenu
}

func onFieldsMenu() {
	actions := fieldsMenu.Actions()
	noVisibles := true
	for i, action := range actions {
		(*pColumns)[i].visible = action.IsChecked()
		if (*pColumns)[i].visible {
			noVisibles = false
		}
	}
	if noVisibles {
		actions[0].SetChecked(true)
		(*pColumns)[0].visible = true
	}
	(*pArea).UpdateColsWidth()
	(*pArea).Draw()
}

func prepareMenu(menu *qt.QMenu, titles, shorts, icons []string, cbs []func()) {
	for i, title := range titles {
		if len(title) == 0 {
			menu.AddSeparator()
			continue
		}
		action := menu.AddAction(titles[i])
		if len(shorts[i]) > 0 {
			action.SetShortcut(qt.NewQKeySequence2(shorts[i]))
		}
		if len(icons[i]) > 0 {
			action.SetIcon(qt.QIcon_FromTheme(icons[i]))
		}
		if cbs[i] != nil {
			action.OnTriggered(cbs[i])
		}
	}
}

func addChecksToMenu(menu *qt.QMenu, checkables []bool, checks []bool) {
	actions := menu.Actions()
	var j int
	for i := range actions {
		if checkables[i] {
			actions[i].SetCheckable(true)
			actions[i].SetChecked(checks[j])
			j++
		}
	}
}
