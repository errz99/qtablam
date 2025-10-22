// Package qtablam ...
package qtablam

import (
	// "fmt"
	// "slices"
	// "strconv"
	"unicode/utf8"

	"github.com/mappu/miqt/qt"
)

const (
	verPadding = 3
	recPadding = 2
)

var (
	pointedColumn *Column
	centerArea    DrawArea

	arrowCursor        *qt.QCursor
	pointingHandCursor *qt.QCursor
	activeCursor       = arrowCursor

	// areaBack    = [4]int{240, 240, 240, 255}
	areaBack    = [4]int{24, 240, 240, 255}
	areaColor   = [4]int{250, 250, 250, 255}
	lineColor   = [4]int{220, 220, 220, 255}
	textColor   = [4]int{0, 0, 0, 255}
	cursorColor = [4]int{255, 251, 100, 127}
	selColor    = [4]int{200, 200, 200, 127}

	defFont  string
	fontData MyFontData

	columns    = make([]Column, 0, 8)
	fieldsMenu *qt.QMenu
)

type Column struct {
	title     string
	width     int
	visible   bool
	resizable bool
	texts     []string
}

func newColumn(title string, width int) Column {
	texts := make([]string, 0)
	return Column{
		title:     title,
		width:     width,
		visible:   true,
		resizable: true,
		texts:     texts,
	}
}

func initColumns(titles []string) {
	for _, title := range titles {
		columns = append(columns, newColumn(title, len(title)+5))
	}
}

type DrawArea struct {
	*qt.QGraphicsView
	width     int
	height    int
	colSep    int
	rowSep    int
	offx      int
	rows      int
	rowOff    int
	offInc    int
	cursorPos int
	data      [][]string
}

func newDrawArea(backColor [4]int, data [][]string) DrawArea {
	var brush = qt.NewQBrush3(myQColor(backColor))

	var area = DrawArea{qt.NewQGraphicsView2(), 0, 0, 0, 0, 0, len(data), 0, 0, 0, data}
	area.colSep = 6
	area.rowSep = 2
	area.offInc = 2
	area.cursorPos = -1

	var scene = qt.NewQGraphicsScene()
	area.SetScene(scene)

	area.SetSizePolicy2(qt.QSizePolicy__Expanding, qt.QSizePolicy__Expanding)
	area.SetTransformationAnchor(qt.QGraphicsView__AnchorUnderMouse)
	area.SetViewportUpdateMode(qt.QGraphicsView__MinimalViewportUpdate)

	area.SetFrameStyle(0)
	area.SetBackgroundBrush(brush)

	return area
}

func (da *DrawArea) UpdateColsWidth() {
	var colsWidth int
	for _, col := range columns {
		if col.visible {
			colsWidth += col.width*fontData.w + centerArea.colSep
		}
	}
	if colsWidth < da.width {
		da.offx = (da.width - colsWidth) / 2
	} else {
		da.offx = 0
	}
}

func (da *DrawArea) UpdateRows() {
	da.rows = da.height / (fontData.h + da.rowSep)
}

func (da *DrawArea) ResizeUpdate(newWidth, newHeight int) {
	da.width = newWidth
	da.UpdateColsWidth()
	if newHeight != da.height {
		da.height = newHeight
		da.UpdateRows()
	}
}

func (da *DrawArea) IncRowOff() {
	for i := da.offInc; i > 0; i-- {
		if da.rowOff < len(da.data)-(i-1) {
			da.rowOff++
		}
	}
}

func (da *DrawArea) DecRowOff() {
	for i := da.offInc; i > 0; i-- {
		if da.rowOff > i-1 {
			da.rowOff--
		}
	}
}

func (da *DrawArea) IncCursor() bool {
	if da.cursorPos < da.rowOff || da.cursorPos > da.rowOff+da.rows {
		da.cursorPos = da.rowOff
		return true
	} else if da.cursorPos < len(da.data)-1 {
		da.cursorPos++
		if da.cursorPos > da.rows-2 && da.rowOff < len(da.data)-da.rows*3/4 {
			da.rowOff++
		}
		return true
	}
	return false
}

func (da *DrawArea) DecCursor() bool {
	if da.cursorPos < da.rowOff || da.cursorPos > da.rowOff+da.rows {
		da.cursorPos = da.rowOff + da.rows - 1
		return true
	} else if da.cursorPos > 0 {
		da.cursorPos--
		if da.cursorPos < da.rowOff {
			if da.rowOff > 0 {
				da.rowOff--
			}
		}
		return true
	}
	return false
}

func (da *DrawArea) IncPage() bool {
	if da.rowOff <= len(da.data)-da.rows {
		da.rowOff += da.rows
		return true
	} else if da.cursorPos < len(da.data)-1 {
		da.cursorPos = len(da.data) - 1
		return true
	}
	return false
}

func (da *DrawArea) DecPage() bool {
	if da.rowOff >= da.rows {
		da.rowOff -= da.rows
		return true
	} else if da.cursorPos > 0 {
		da.cursorPos = 0
		return true
	}
	return false
}

func (da *DrawArea) GoInit() {
	da.rowOff = 0
	da.cursorPos = 0
}

func (da *DrawArea) GoEnd() {
	da.rowOff = (len(da.data) + 1) - da.rows
	da.cursorPos = len(da.data) - 1
}

func myQColor(color [4]int) *qt.QColor {
	return qt.NewQColor11(color[0], color[1], color[2], color[3])
}

func (da *DrawArea) Draw() {
	// if !showSongsList {
	// 	return
	// }
	w := da.width
	h := da.height
	fw := float64(w)
	fh := float64(h)

	daPen := qt.NewQPen3(myQColor(areaColor))
	textPen := qt.NewQPen3(myQColor(textColor))
	linePen := qt.NewQPen3(myQColor(lineColor))
	cursorPen := qt.NewQPen3(myQColor(cursorColor))
	// selPen := qt.NewQPen3(myQColor(selColor))

	font := fontData.font
	fontw := float64(fontData.w)
	fonth := float64(fontData.h)

	scene := da.Scene()
	scene.Clear()
	da.SetSceneRect2(0, 0, fw, fh)

	offx := float64(da.offx)
	xpos := offx
	vpos := 0.0
	xsep := float64(da.colSep)
	ysep := float64(da.rowSep)

	vpadding := verPadding
	if da.rowSep > 1 && vpadding > 0 {
		vpadding--
	}

	scene.AddRect6(offx, 0, fw-offx*2, fh, daPen, daPen.Brush())

	drawText := func(text string, xpos, ypos float64) {
		textItem := scene.AddText2(text, font)
		textItem.SetDefaultTextColor(textPen.Color())
		y := ypos - float64(vpadding) + float64(da.rowSep)
		textItem.SetPos2(xpos, y)
	}

	// Header
	font.SetBold(true)

	for i, col := range columns {
		if col.visible {
			if i > 0 {
				scene.AddLine4(xpos, 0, xpos, fh, linePen)
			}
			drawText(col.title, xpos, 0)
			xpos += float64(col.width)*fontw + xsep
		}
	}

	// Rows
	font.SetBold(false)

	tempStr := func(text string, width int) string {
		txtWidth := utf8.RuneCountInString(text)
		if txtWidth <= width {
			return text
		}
		newWidth := txtWidth - (txtWidth - width)
		return text[:newWidth]
	}

	putItem := func(colWidth int, text string) {
		str := tempStr(text, colWidth)
		drawText(str, xpos, vpos)
	}

	drawRec := func(pen *qt.QPen, vpos float64) {
		y := vpos + ysep - float64(da.rowSep) + recPadding
		h := float64(fontData.h+da.rowSep) - recPadding*2
		scene.AddRect6(offx, y, float64(da.width)-offx*2, h, pen, pen.Brush())
	}

	for i := da.rowOff; i < len(da.data); i++ {
		xpos = offx
		vpos += fonth + ysep
		if i == da.rows+da.rowOff {
			break
		}

		if i == da.cursorPos {
			drawRec(cursorPen, vpos)
		}
		// if slices.Contains(data.MySel.Elems, data.MyList.Songs[i].ID) {
		// 	font.SetBold(true)
		// 	drawRec(selPen, vpos)
		// }

		for j, col := range columns {
			if !col.visible {
				continue
			}
			putItem(col.width, da.data[i][j])
			xpos += float64(col.width)*fontw + xsep
		}
		font.SetBold(false)
		scene.AddLine4(offx, vpos, float64(da.width)-offx, vpos, linePen)
	}
	scene.AddLine4(offx, vpos+fonth+ysep, float64(da.width)-offx, vpos+fonth+ysep, linePen)
}

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
