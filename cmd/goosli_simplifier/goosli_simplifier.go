package main

import (
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
)

var (
	stl       = kingpin.Flag("stl", "Source stl file to simplify.").Short('s').Required().String()
	out       = kingpin.Flag("out", "Output path for result.").Short('o').Default("out.stl").String()
	triangles = kingpin.Flag("triangles", "Count of triangles to leave.").Short('t').Default("800").Int()
)

func main() {

	kingpin.Parse()

	mesh, err := LoadSTL(*stl)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}

	mUp, err := helpers.SimplifyMesh(mesh, *triangles)
	if err != nil {
		log.Fatal("failed to cut mesh: ", err)
	}

	err = SaveSTL(*out, mUp)
	if err != nil {
		log.Fatal("failed to save simplified mesh: ", err)
	}

}
