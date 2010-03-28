package main

import (
	"fmt"
)

// Declare an Interface to a 3d Solid
type Solid interface {
	Volume() float
	SurfaceArea() float
}

// Contains the Fields for defining a Rectangular Prism's Dimension's
type RectPrism struct {
	l, w, h float
}

// RectPrism implements the Solid Interface
func (this *RectPrism) Volume() float {
	return(this.l * this.w * this.h)
}

func (this *RectPrism) SurfaceArea() float {
	return(2 * (this.l * this.w) + 2 * (this.l * this.h) + 2 * (this.w * this.h))
}

// This Class is going to inherit from RectPrism
type CardboardBox struct {
	// An anonymous field, all fields of RectPrism are promoted into CardboardBox
	RectPrism
	isSoggy bool
}

// This CardboardBox has the top Open so we must reimplement the SurfaceArea func
// Inherits CardboardBox
type OpenCardboardBox struct {
	CardboardBox
}

// Reimplement the SurfaceArea Function for OpenCardboardBox since it doesn't have a top
func (this *OpenCardboardBox) SurfaceArea() float {
	return(this.CardboardBox.SurfaceArea() + 2 * (this.l * this.h) + 2 * (this.w * this.h))
}

func main() int {

	fmt.Printf("\n\n");


	cbox := new(CardboardBox)
	cbox.l = 2
	cbox.w = 4
	cbox.h = 2
	cbox.isSoggy = true

	obox := new(OpenCardboardBox)
	obox.l = 2
	obox.w = 4
	obox.h = 2
	obox.isSoggy = true
	// CardboardBox implements the RectPrism interface
	// through the anonymous field RectPrismStruct
	// This Aggregates the RectPrismStruct into CardboardBox
	var rprism Solid = cbox

	fmt.Printf("      Volume: %f\n", rprism.Volume())
	fmt.Printf("Surface Area: %f\n", rprism.SurfaceArea())
	rprism = obox
	fmt.Printf("      Volume: %f\n", rprism.Volume())
	fmt.Printf("Surface Area: %f\n", rprism.SurfaceArea())
	fmt.Printf("\n\n");
	return 0;
}
