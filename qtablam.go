// Package qtablam ...
package qtablam

import (
	// "fmt"

	"github.com/mappu/miqt/qt"
)

func NewQTablam(titles []string, data [][]string) *qt.QWidget {
	arrowCursor = qt.NewQCursor2(qt.CursorShape(qt.ArrowCursor))
	pointingHandCursor = qt.NewQCursor2(qt.CursorShape(qt.PointingHandCursor))

	DefFont = DefaultFont()
	FontData = NewFontData()

	// Draw area
	centerArea := newDrawArea(titles, data)
	FontData.UpdateMetrics(centerArea.rowSep)

	centerArea.SetCursor(arrowCursor)
	fieldsMenu = initMenuFields(centerArea.columns)

	centerArea.Draw()
	return centerArea.QWidget
}
