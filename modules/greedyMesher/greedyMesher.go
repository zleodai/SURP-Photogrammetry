package greedyMesher

import (
	"fmt"
	"os"
)

type Mesh struct {
	Vertices []Vertix
	Faces    []int
}

type Vertix struct {
	X float64
	Y float64
	Z float64
}

type faceOrientation int

const (
	up faceOrientation = iota
	down
	left
	right
	forward
	backward
)

type Face struct {
	VoxelCoords [3]int
	FaceIndex   faceOrientation
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

	outputFaces := []Face{}
	for z, zArray := range zxySlices {
		for x, xArray := range zArray {
			for y, voxelPresent := range xArray {
				if voxelPresent && z != len(zArray)-1 && z != 0 {
					if zxySlices[z+1][x][y] {
						//Relative Upface found
						voxelCoords := [3]int{x, y, z}
						var newFace Face = Face{VoxelCoords: voxelCoords, FaceIndex: 0}
						outputFaces = append(outputFaces, newFace)
					}
					if zxySlices[z-1][x][y] {
						//Relative Downface found
						voxelCoords := [3]int{x, y, z}
						var newFace Face = Face{VoxelCoords: voxelCoords, FaceIndex: 1}
						outputFaces = append(outputFaces, newFace)
					}
				}
			}
		}
	}
	for x, xArray := range xyzSlices {
		for y, yArray := range xArray {
			for z, voxelPresent := range yArray {
				if voxelPresent && x != len(xArray)-1 && x != 0 {
					if xyzSlices[x+1][y][z] {
						//Relative Upface found
						voxelCoords := [3]int{x, y, z}
						var newFace Face = Face{VoxelCoords: voxelCoords, FaceIndex: 3}
						outputFaces = append(outputFaces, newFace)
					}
					if xyzSlices[x-1][y][z] {
						//Relative Downface found
						voxelCoords := [3]int{x, y, z}
						var newFace Face = Face{VoxelCoords: voxelCoords, FaceIndex: 2}
						outputFaces = append(outputFaces, newFace)
					}
				}
			}
		}
	}
	for y, yArray := range yxzSlices {
		for x, xArray := range yArray {
			for z, voxelPresent := range xArray {
				if voxelPresent && y != len(yArray)-1 && y != 0 {
					if yxzSlices[y+1][x][z] {
						//Relative Upface found
						voxelCoords := [3]int{x, y, z}
						var newFace Face = Face{VoxelCoords: voxelCoords, FaceIndex: 4}
						outputFaces = append(outputFaces, newFace)
					}
					if yxzSlices[y-1][x][z] {
						//Relative Downface found
						voxelCoords := [3]int{x, y, z}
						var newFace Face = Face{VoxelCoords: voxelCoords, FaceIndex: 5}
						outputFaces = append(outputFaces, newFace)
					}
				}
			}
		}
	}
	fmt.Print(outputFaces)
}

// func combineVoxels(slice [][][]bool, mode int) []Face {
// 	outputFaces := []Face{}
// 	for dir, dirArray := range slice {
// 		for x, xArray := range dirArray {
// 			for y, voxelPresent := range xArray {
// 				if voxelPresent {
// 					if slice[dir+1][x][y] || dir == len(dirArray)-1{
// 						//Relative Upface found
// 						var newFace Face = Face{VoxelCoords: {}}
// 						outputFaces = append(outputFaces, )
// 					if slice[dir-1][x][y] || dir == 0 {
// 						//Relative Downface found
// 						outputFaces = append(outputFaces, )
// 					}
// 				}
// 			}
// 		}
// 	}
// }

func MeshToObj(mesh Mesh) {
	file, errs := os.Create("VoxelMatrix.JSON")
	if errs != nil {
		panic("Failed to write to file:" + errs.Error())
	}
	fmt.Println(file)
}
