package qtablam

import (
	// "fmt"
	"slices"
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

	areaBack    = [4]int{240, 240, 240, 255}
	areaColor   = [4]int{250, 250, 250, 255}
	lineColor   = [4]int{220, 220, 220, 255}
	textColor   = [4]int{0, 0, 0, 255}
	cursorColor = [4]int{255, 251, 100, 127}
	selColor    = [4]int{200, 200, 200, 127}

	DefFont  string
	FontData MyFontData

	columns         = make([]Column, 0, 8)
	fieldsMenu      *qt.QMenu
	ModifierControl = qt.NoModifier
	showDataTable   = true
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

type Row struct {
	ID    int
	texts []string
}

type DrawArea struct {
	*qt.QGraphicsView
	width      int
	height     int
	colSep     int
	rowSep     int
	offx       int
	rows       int
	rowOff     int
	offInc     int
	cursorPos  int
	data       []Row
	sel        Selection
	dataActive int
}

func newDrawArea(backColor [4]int, data [][]string) DrawArea {
	var brush = qt.NewQBrush3(myQColor(backColor))

	rows := make([]Row, 0, 16)
	for i, row := range data {
		rows = append(rows, Row{ID: i, texts: row})
	}

	var area = DrawArea{qt.NewQGraphicsView2(), 0, 0, 0, 0, 0, len(data), 0, 0, 0, rows, newSelection(), -1}
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

	area.OnResizeEvent(func(super func(event *qt.QResizeEvent), event *qt.QResizeEvent) {
		area.ResizeUpdate(event.Size().Width(), event.Size().Height())
		area.Draw()
	})
	area.OnMouseMoveEvent(func(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
		onAreaMouseMoveEvent(area, event)
	})
	area.OnWheelEvent(func(super func(event *qt.QWheelEvent), event *qt.QWheelEvent) {
		if onAreaWheelEvent(&area, event) {
			area.Draw()
		}
	})
	area.OnMousePressEvent(func(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
		if onAreaPressEvent(&area, event) {
			area.Draw()
		}
	})
	area.OnMouseDoubleClickEvent(func(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
		if onAreaDoubleClickEvent(&area, event) {
			area.Draw()
		}
	})

	area.OnKeyPressEvent(func(super func(event *qt.QKeyEvent), event *qt.QKeyEvent) {
		if onKeyPressEvent(&centerArea, event) {
			centerArea.Draw()
		}
	})

	area.OnKeyReleaseEvent(func(super func(event *qt.QKeyEvent), event *qt.QKeyEvent) {
		switch event.Key() {
		case int(qt.Key_Control):
			ModifierControl = qt.NoModifier
		default:
		}
	})

	return area
}

func (da *DrawArea) UpdateColsWidth() {
	var colsWidth int
	for _, col := range columns {
		if col.visible {
			colsWidth += col.width*FontData.W + centerArea.colSep
		}
	}
	if colsWidth < da.width {
		da.offx = (da.width - colsWidth) / 2
	} else {
		da.offx = 0
	}
}

func (da *DrawArea) UpdateRows() {
	da.rows = da.height / (FontData.H + da.rowSep)
}

func (da *DrawArea) SetCursorPosition(pos int) {
	da.cursorPos = pos
	da.rowOff = pos - da.rows/2
	if da.rowOff < 0 {
		da.rowOff = 0
	}
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
	if !showDataTable {
		return
	}
	w := da.width
	h := da.height
	fw := float64(w)
	fh := float64(h)

	daPen := qt.NewQPen3(myQColor(areaColor))
	textPen := qt.NewQPen3(myQColor(textColor))
	linePen := qt.NewQPen3(myQColor(lineColor))
	cursorPen := qt.NewQPen3(myQColor(cursorColor))
	selPen := qt.NewQPen3(myQColor(selColor))

	font := FontData.Font
	fontw := float64(FontData.W)
	fonth := float64(FontData.H)

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
		h := float64(FontData.H+da.rowSep) - recPadding*2
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
		if slices.Contains(da.sel.Elems, da.data[i].ID) {
			font.SetBold(true)
			drawRec(selPen, vpos)
		}

		for j, col := range columns {
			if !col.visible {
				continue
			}
			putItem(col.width, da.data[i].texts[j])
			xpos += float64(col.width)*fontw + xsep
		}
		font.SetBold(false)
		scene.AddLine4(offx, vpos, float64(da.width)-offx, vpos, linePen)
	}
	scene.AddLine4(offx, vpos+fonth+ysep, float64(da.width)-offx, vpos+fonth+ysep, linePen)
}
