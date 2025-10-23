// Package qtablam ...
package qtablam

import (
	// "fmt"

	"github.com/mappu/miqt/qt"
)

func NewQTablam(titles []string, data [][]string) *DrawArea {
	arrowCursor = qt.NewQCursor2(qt.CursorShape(qt.ArrowCursor))
	pointingHandCursor = qt.NewQCursor2(qt.CursorShape(qt.PointingHandCursor))

	DefFont = DefaultFont()
	FontData = NewFontData()

	// Draw area
	area := newDrawArea(titles, data)
	FontData.UpdateMetrics(area.rowSep)

	area.SetCursor(arrowCursor)
	fieldsMenu = initMenuFields(area.columns)

	return area
}
