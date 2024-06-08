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
	right
	left
	forward
	backward
)

type Face struct {
	VoxelCoords [][3]int
	FaceIndex   faceOrientation
}

func GreedyMesh(voxels [][][]uint8, threshold uint8) []Face {
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
	outputFaces = append(outputFaces, combineVoxels(zxySlices, 0, 1)...)
	outputFaces = append(outputFaces, combineVoxels(zxySlices, 0, -1)...)
	outputFaces = append(outputFaces, combineVoxels(xyzSlices, 2, 1)...)
	outputFaces = append(outputFaces, combineVoxels(xyzSlices, 2, -1)...)
	outputFaces = append(outputFaces, combineVoxels(yxzSlices, 4, 1)...)
	outputFaces = append(outputFaces, combineVoxels(yxzSlices, 4, -1)...)
	return outputFaces
}

func combineVoxels(refSlice [][][]bool, orientationOffset, orientationDirection int) []Face {
	//orientationDirection is essentially a int saying whether its an up pass or down pass,
	//orientationDirection = 1 means up pass
	//orientationDirection = -1 means down pass

	//make a copy of refSlice
	slice := make([][][]bool, len(refSlice))
	for x := 0; x < len(refSlice); x++ {
		slice[x] = make([][]bool, len(refSlice[0]))
		for y := 0; y < len(refSlice[0]); y++ {
			slice[x][y] = make([]bool, len(refSlice[0][0]))
		}
	}

	copy(slice, refSlice)

	outputFaces := []Face{}
	for dir, dirArray := range slice {
		for x, xArray := range dirArray {
			for y, voxelPresent := range xArray {
				if voxelPresent {
					// fmt.Printf("\nVoxel Present at %d, %d, %d for %d at %d direction\n", dir, x, y, orientationOffset, orientationDirection)
					var edgeOfArray bool = false
					if (dir == 0 && orientationDirection == -1) || (dir == len(dirArray)-1 && orientationDirection == 1) {
						edgeOfArray = true
					}
					var faceFound, maxXReached bool = false, false
					//Creates a new array to record the corners of the face being generated. Initalizes it with the voxel (dir, x, y)
					var faceBounds [][3]int = [][3]int{voxelCordsOffsetter([3]int{dir, x, y}, orientationOffset)}

					//Initalizes bounds for the x and y expansion of the face.
					cornerBounds := [2]int{x, y}
					for !faceFound {
						if !maxXReached {
							if slice[dir][cornerBounds[0]+1][y] && (edgeOfArray || slice[dir+orientationDirection][cornerBounds[0]+1][y]) {
								cornerBounds[0] += 1
								if len(faceBounds) == 1 {
									faceBounds = append(faceBounds, [3]int{dir, cornerBounds[0], y})
								} else {
									faceBounds[1][1] = cornerBounds[0]
								}
							} else {
								maxXReached = true
							}
						} else {
							var yIncreasePossible bool = true
							for i := 0; i < cornerBounds[0]-faceBounds[0][1]+1; i++ {
								if !(slice[dir][faceBounds[0][1]+i][cornerBounds[1]+1] && (edgeOfArray || slice[dir+orientationDirection][faceBounds[0][1]+i][cornerBounds[1]+1])) {
									yIncreasePossible = false
								}
							}
							// fmt.Printf("\n        yIncreasePossible: %t,  edgeOfArray: %t\n", yIncreasePossible, edgeOfArray)
							if yIncreasePossible {
								cornerBounds[1] += 1
								if len(faceBounds) == 1 {
									faceBounds = append(faceBounds, [3]int{dir, faceBounds[0][1], cornerBounds[1]})
								} else if faceBounds[0][1] == faceBounds[1][1] {
									faceBounds[1][2] = cornerBounds[1]
								} else if len(faceBounds) == 2 {
									faceBounds = append(faceBounds, [3]int{dir, cornerBounds[0], cornerBounds[1]})
									faceBounds = append(faceBounds, [3]int{dir, faceBounds[0][1], cornerBounds[1]})
								} else {
									faceBounds[2][2] = cornerBounds[1]
									faceBounds[3][2] = cornerBounds[1]
								}
							} else {
								faceFound = true
							}
						}
						// fmt.Printf("\n    faceFound: %t, maxXReached: %t,  cornerBounds: [%d, %d]\n", faceFound, maxXReached, cornerBounds[0], cornerBounds[1])
						// for _, face := range faceBounds {
						// 	fmt.Printf("\n    [%d, %d, %d]\n", face[0], face[1], face[2])
						// }
					}
					var newFace Face
					if orientationDirection == 1 {
						newFace = Face{VoxelCoords: faceBounds, FaceIndex: faceOrientation(orientationOffset)}
					} else {
						newFace = Face{VoxelCoords: faceBounds, FaceIndex: faceOrientation(orientationOffset + 1)}
					}
					outputFaces = append(outputFaces, newFace)

					//Lastly now u have to set the voxelPresent bool to false for the voxels already considered into outputFaces
					for curX := faceBounds[0][1]; curX < cornerBounds[0]; curX++ {
						for curY := faceBounds[0][2]; curY < cornerBounds[1]; curY++ {
							slice[dir][curX][curY] = false
						}
					}
				}
			}
		}
	}

	// fmt.Printf("\nCombined Voxels for %d in %d direction\n", orientationOffset, orientationDirection)
	return outputFaces
}

func voxelCordsOffsetter(coords [3]int, orientationOffset int) [3]int {
	//Used to change cords from zxy into the normalized xyz. Assumes that you are putting the right orientationOffset
	//Assumes the following:
	//orientationOffset = 0, when zxy
	//orientationOffset = 2, when xyz
	//orientationOffset = 4, when yxz
	var voxelCoords [3]int
	switch x := orientationOffset; x {
	case 0:
		voxelCoords = [3]int{coords[1], coords[2], coords[0]}
	case 2:
		voxelCoords = [3]int{coords[0], coords[1], coords[2]}
	case 4:
		voxelCoords = [3]int{coords[1], coords[0], coords[2]}
	default:
		fmt.Println("Error, orientationOffset param in combineVoxels() must be either 0, 2, or 4 to signify zxy, xyz, and yxz respectively")
		voxelCoords = [3]int{coords[0], coords[1], coords[2]}
	}
	return voxelCoords
}

func MeshToObj(mesh Mesh) {
	file, errs := os.Create("VoxelMatrix.JSON")
	if errs != nil {
		panic("Failed to write to file:" + errs.Error())
	}
	fmt.Println(file)
}
