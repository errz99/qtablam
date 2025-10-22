// Package qtablam ...
package qtablam

import (
	// "fmt"

	"github.com/mappu/miqt/qt"
)

func NewQTablam(titles []string, data [][]string) *qt.QWidget {
	arrowCursor = qt.NewQCursor2(qt.CursorShape(qt.ArrowCursor))
	pointingHandCursor = qt.NewQCursor2(qt.CursorShape(qt.PointingHandCursor))

	defFont = DefaultFont()
	fontData = newFontData()

	// Draw area
	initColumns(titles)
	centerArea = newDrawArea(areaBack, data)
	fontData.UpdateMetrics(centerArea.rowSep)

	centerArea.SetCursor(arrowCursor)
	fieldsMenu = initMenuFields()

	centerArea.Draw()
	return centerArea.QWidget
}
