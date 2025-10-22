package qtablam

import (
	"unicode/utf8"

	"github.com/mappu/miqt/qt"
)

func onAreaMouseMoveEvent(area DrawArea, event *qt.QMouseEvent) {
	x := event.X()
	y := event.Y()
	fw := FontData.W

	if x > area.offx && x < area.width-area.offx {
		if y < FontData.H+area.rowSep {
			total := area.offx
			res := -1
			for i, col := range columns {
				if col.visible {
					total += col.width*fw + area.colSep
					res = i
					if total > x {
						break
					}
				}
			}
			if activeCursor != pointingHandCursor {
				activeCursor = pointingHandCursor
				centerArea.SetCursor(pointingHandCursor)
			}
			if res >= 0 {
				pointedColumn = &columns[res]
			}
		} else {
			if activeCursor != arrowCursor {
				activeCursor = arrowCursor
				centerArea.SetCursor(arrowCursor)
				pointedColumn = nil
			}
		}
	}
}

func onAreaWheelEvent(area *DrawArea, event *qt.QWheelEvent) bool {
	delta := event.Delta()

	if delta < 0 {
		if activeCursor == pointingHandCursor {
			area.UpdateColsWidth()
			if pointedColumn != nil && pointedColumn.resizable {
				if pointedColumn.width > utf8.RuneCountInString(pointedColumn.title) {
					pointedColumn.width--
				}
			}
		} else {
			area.IncRowOff()
		}
	} else if delta > 0 {
		if activeCursor == pointingHandCursor {
			area.UpdateColsWidth()
			if pointedColumn != nil && pointedColumn.resizable {
				pointedColumn.width++
			}
		} else {
			area.DecRowOff()
		}
	}
	return true
}

func onAreaPressEvent(area *DrawArea, event *qt.QMouseEvent) bool {
	x := event.X()
	y := event.Y()
	modifiers := event.Modifiers()
	button := event.Button()

	if x > area.offx && x < area.width-area.offx {
		if y < FontData.H+area.rowSep {
			if button == 2 {
				globalPos := area.MapToGlobal(event.Pos())
				fieldsMenu.Popup(globalPos)
			} else {
				// field := pointedColumn.field
				// switch field {
				// case data.FieldDiscNr:
				// 	return false
				// default:
				// 	data.MyList.Sort(field)
				// 	return true
				// }
			}
		} else {
			ypos := y/(FontData.H+area.rowSep) + area.rowOff - 1
			if button == 2 {
				id := area.data[ypos].ID
				if !area.sel.Push(id) {
					area.sel.Remove(id)
				}
				return true
			}
			if ypos < len(area.data) {
				switch modifiers {
				case qt.NoModifier:
					area.cursorPos = ypos
					return true
				case qt.AltModifier:
					area.sel.Clear()
					return true
				case qt.ShiftModifier:
					if !area.sel.Push(ypos) {
						area.sel.Remove(ypos)
					}
					return true
				default:
				}
			}
		}
	}
	return false
}

func onAreaDoubleClickEvent(area *DrawArea, event *qt.QMouseEvent) bool {
	x := event.X()
	y := event.Y()

	if y > FontData.H+area.rowSep && (x > area.offx && x < area.width-area.offx) {
		position := y/(FontData.H+area.rowSep) + area.rowOff - 1
		if position < len(area.data) {
			area.dataActive = position
			// updateSongLabels()
		}
		return true
	}
	return false
}
