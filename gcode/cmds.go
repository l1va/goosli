package gcode

import (
	"bytes"
	"math"
	"strconv"

	. "github.com/l1va/goosli/primitives"
)

type Command interface {
	ToGCode(b *bytes.Buffer)
	LayersCount() int
}

type InclineXBack struct {
}

func (r InclineXBack) ToGCode(b *bytes.Buffer) {
	b.WriteString("M42 \n")
}
func (r InclineXBack) LayersCount() int {
	return 0
}

type InclineX struct {
}

func (r InclineX) ToGCode(b *bytes.Buffer) {
	b.WriteString("M43 \n")
}
func (r InclineX) LayersCount() int {
	return 0
}

type RotateXZ struct {
	AngleX float64
	AngleZ float64
}

func (r RotateXZ) ToGCode(b *bytes.Buffer) {
	b.WriteString("G62 X" + StrF(r.AngleX) + " Z" + StrF(r.AngleZ) + "\n")
}
func (r RotateXZ) LayersCount() int {
	return 0
}

type RotateZ struct {
	Angle float64
}

func (r RotateZ) ToGCode(b *bytes.Buffer) {
	b.WriteString("G0 A" + StrF(r.Angle) + "\n")
}
func (r RotateZ) LayersCount() int {
	return 0
}

type LayersMoving struct {
	Layers []Layer
	Index  int
}

func (lm LayersMoving) ToGCode(b *bytes.Buffer) {
	for i := 0; i < len(lm.Layers); i++ {
		b.WriteString(";LAYER:" + strconv.Itoa(i+lm.Index) + "\n")
		switchFanGCode(lm.Layers[i].FanOff, b)

		printSpeedToGCode(lm.Layers[i].WallPrintSpeed, b)
		pathesToGCode(lm.Layers[i].Paths, "OUTER_PATHES", b)

		printSpeedToGCode(lm.Layers[i].PrintSpeed, b)
		pathesToGCode(lm.Layers[i].MiddlePs, "MIDDLE_PATHES", b)
		pathesToGCode(lm.Layers[i].InnerPs, "INNER_PATHES", b)
		pathesToGCode(lm.Layers[i].Fill, "FILL_PATHES", b)
	}
}
func (lm LayersMoving) LayersCount() int {
	return len(lm.Layers)
}

func switchFanGCode(fanOff bool, b *bytes.Buffer) {
	if fanOff {
		b.WriteString("M107\n")
	} else {
		b.WriteString("M106\n")
	}
}

func printSpeedToGCode(feedrate int, b *bytes.Buffer) {
	b.WriteString("F" + strconv.Itoa(feedrate) + "\n")
}

func pathesToGCode(pths []Path, comment string, b *bytes.Buffer) {
	eOff := 0.0 //TODO: fix extruder value
	b.WriteString(";" + comment + "\n")
	for _, p := range pths {
		b.WriteString("G0 " + p.Lines[0].P1.String() + "\n")
		for _, line := range p.Lines {
			eDist := math.Sqrt(math.Pow(line.P2.X-line.P1.X, 2) + math.Pow(line.P2.Y-line.P1.Y, 2) + math.Pow(line.P2.Z-line.P1.Z, 2))
			eOff += eDist
			b.WriteString("G1 " + line.P2.String() + " E" + StrF(eOff) + "\n")
		} //TODO: optimize - not write coordinate if it was not changed
	}
}
