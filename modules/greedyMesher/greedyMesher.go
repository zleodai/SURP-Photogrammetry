package greedyMesher

import (
	"fmt"
	"os"
)

func Test() byte {
	return 0
}

type Mesh struct {
	Vertices []Vertix
	Faces    []int
}

type Vertix struct {
	X float64
	Y float64
	Z float64
}

func GreedyMesh(voxels [][][]uint8, threshold uint8) {
	zxySlices := make([][][]bool, len(voxels[0][0]))
	for z := 0; z < len(voxels[0][0]); z++ {
		zxySlices[z] = make([][]bool, len(voxels))
		for x := 0; x < len(voxels); x++ {
			zxySlices[z][x] = make([]bool, len(voxels[0]))
		}
	}

	xyzSlices := make([][][]bool, len(voxels))
	for x := 0; x < len(voxels); x++ {
		xyzSlices[x] = make([][]bool, len(voxels[0]))
		for y := 0; y < len(voxels[0]); y++ {
			xyzSlices[x][y] = make([]bool, len(voxels[0][0]))
		}
	}

	yxzSlices := make([][][]bool, len(voxels[0]))
	for y := 0; y < len(voxels[0]); y++ {
		yxzSlices[y] = make([][]bool, len(voxels))
		for x := 0; x < len(voxels); x++ {
			yxzSlices[y][x] = make([]bool, len(voxels[0][0]))
		}
	}

	go func() {
		for x, xArray := range voxels {
			for y, yArray := range xArray {
				for z, Value := range yArray {
					if Value > threshold {
						zxySlices[z][x][y] = true
						xyzSlices[x][y][z] = true
						yxzSlices[y][x][z] = true
					}
				}
			}
		}
	}()

	// for z, xyArray := range zxySlices {

	// 	for x, yArray := range xyArray {
	// 		for y, value := range yArray {

	// 		}
	// 	}
	// }
}

func MeshToObj(mesh Mesh) {
	file, errs := os.Create("VoxelMatrix.JSON")
	if errs != nil {
		panic("Failed to write to file:" + errs.Error())
	}
	fmt.Println(file)
}
