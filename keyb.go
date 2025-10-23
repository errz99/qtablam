package qtablam

import (
	"github.com/mappu/miqt/qt"
)

func onKeyPressEvent(area *DrawArea, event *qt.QKeyEvent) bool {
	modifiers := event.Modifiers()
	var refresh bool

	switch event.Key() {
	case int(qt.Key_Control):
		ModifierControl = qt.ControlModifier
	case int(qt.Key_Escape):
		area.cursorPos = -1
		refresh = true
	case int(qt.Key_F5):
		if modifiers == qt.ControlModifier {
			if area.rowSep > 0 {
				area.rowSep--
				area.UpdateRows()
				refresh = true
			}
		} else {
			if FontData.UpdateSize(area.rowSep, -1) {
				area.UpdateColsWidth()
				area.UpdateRows()
				refresh = true
			}
		}
	case int(qt.Key_F6):
		if modifiers == qt.ControlModifier {
			area.rowSep++
			area.UpdateRows()
			refresh = true
		} else {
			if FontData.UpdateSize(area.rowSep, 1) {
				area.UpdateColsWidth()
				area.UpdateRows()
				refresh = true
			}
		}
	case int(qt.Key_J):
		refresh = area.IncCursor()
	case int(qt.Key_K):
		refresh = area.DecCursor()
	case int(qt.Key_F):
		if modifiers == qt.ControlModifier {
			refresh = area.IncPage()
		}
	case int(qt.Key_B):
		if modifiers == qt.ControlModifier {
			refresh = area.DecPage()
		}
	case int(qt.Key_A):
		if modifiers == qt.ControlModifier {
			if len(area.rows) > 0 {
				area.dataActive = area.cursorPos
				// updateSongLabels()
				refresh = true
			}
		}
	case int(qt.Key_L):
		if len(area.rows) > 0 {
			switch modifiers {
			case qt.ShiftModifier:
				id := area.rows[area.cursorPos].ID
				if !area.sel.Push(id) {
					area.sel.Remove(id)
				}
				// updateSongLabels()
				refresh = true
			case qt.AltModifier:
				area.sel.Clear()
				// updateSongLabels()
				refresh = true
			case qt.ControlModifier:
				// search.SetFocus()
			default:
			}
		}
	case int(qt.Key_Home), int(qt.Key_1):
		area.GoInit()
		refresh = true
	case int(qt.Key_End), int(qt.Key_0):
		area.GoEnd()
		refresh = true
	default:
	}
	return refresh
}
