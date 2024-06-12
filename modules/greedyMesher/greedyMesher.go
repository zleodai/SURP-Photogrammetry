package greedyMesher

import (
	"encoding/json"
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

type FaceOrientation int

const (
	up FaceOrientation = iota
	down
	right
	left
	forward
	backward
)

type Face struct {
	VoxelCoords [][3]int
	FaceIndex   FaceOrientation
}

type Faces struct {
	FaceArray []JSONFace
}

type JSONFace struct {
	VoxelCoords []Point
	FaceIndex   int
}

type Point struct {
	X int
	Y int
	Z int
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
	outputFaces = append(outputFaces, combineVoxels2(zxySlices, 0, true)...)
	fmt.Print(outputFaces)
	fmt.Println()
	outputFaces = append(outputFaces, combineVoxels2(zxySlices, 0, false)...)
	fmt.Print(outputFaces)
	fmt.Println()
	outputFaces = append(outputFaces, combineVoxels2(xyzSlices, 2, true)...)
	fmt.Print(outputFaces)
	fmt.Println()
	outputFaces = append(outputFaces, combineVoxels2(xyzSlices, 2, false)...)
	fmt.Print(outputFaces)
	fmt.Println()
	outputFaces = append(outputFaces, combineVoxels2(yxzSlices, 4, true)...)
	fmt.Print(outputFaces)
	fmt.Println()
	outputFaces = append(outputFaces, combineVoxels2(yxzSlices, 4, false)...)
	fmt.Print(outputFaces)
	fmt.Println()

	return outputFaces
}

func combineVoxels2(voxels [][][]bool, orientation FaceOrientation, isUp bool) []Face {
	greedyMeshed := []Face{}

	checkOffset := 0
	if isUp {
		checkOffset = 1
	} else {
		checkOffset = -1
	}
	
	for z := 0; z < len(voxels); z++ {
		currentSlice := make([][]bool, len(voxels[0]))
		for x := 0; x < len(voxels[0]); x++ {
			currentSlice[x] = make([]bool, len(voxels[0][0]))
		}

		b, _ := json.Marshal(voxels[z])

		json.Unmarshal(b, &currentSlice)

		assumeAir := false
		if (!isUp && z == 0) || (isUp && z == len(voxels) -1){
			assumeAir = true
		} 
		for x := 0; x < len(currentSlice); x++ {
			for y := 0; y < len(currentSlice[0]); y++ {
				currVoxel := currentSlice[x][y]
				if currVoxel && (assumeAir || !voxels[z + checkOffset][x][y]){
					x1, y1, z1 := remapXYZ(x,y,z, orientation)
					corners := [][3]int{{x1, y1, z1}}
					
					currY, actualY := y, y
					currX := x
					nextFace, faceComplete := true, false
					if y+1 != len(voxels[0][0]) {
						nextFace = currentSlice[currX][currY+1] && (assumeAir || !voxels[z + checkOffset][currX][currY])
					} else {
						nextFace = false
					}
					for !faceComplete {
						for nextFace {
							if currY + 1 != len(currentSlice[0]) {
								currY += 1
								nextFace = currentSlice[currX][currY] && (assumeAir || !voxels[z + checkOffset][currX][currY])
							} else {
								nextFace = false
							}
						}
						if currX == x && currY != y {
							actualY = currY
							x2, y2, z2 := remapXYZ(x, currY, z, orientation)
							corners = append(corners, [3]int{x2, y2, z2})
						} else {
							break
						}
						if currY != actualY {
							faceComplete = true
						} else {
							currX += 1
							currY = 0
						}
					}
					if len(corners) == 2 && currY == actualY && currX > x {
						x3, y3, z3 := remapXYZ(currX, y, z, orientation)
						corners = append(corners, [3]int{x3, y3, z3})
						x4, y4, z4 := remapXYZ(currX, currY, z, orientation)
						corners = append(corners, [3]int{x4, y4, z4})
					}


					for startX := x; startX <= currX; startX++ {
						for startY := y; startY <= actualY; startY++{
							currentSlice[startX][startY] = false
						}
					}

					FI := orientation
					if !isUp {
						FI += 1
					}
					greedyMeshed = append(greedyMeshed, Face{VoxelCoords: corners, FaceIndex: FI})
				}
			}
		}
	}
	return greedyMeshed
}

func remapXYZ(x,y,z int, orientation FaceOrientation) (int, int, int){
	switch orientation {
	case up:
		return x, y, z
	case down:
		return x, y, z
	case right:
		return z, x, y
	case left:
		return z, x, y
	case forward:
		return x, z, y
	case backward:
		return x, z, y
	default:
		return x, y, z
	}
}

func combineVoxels(refSlice [][][]bool, orientationOffset, orientationDirection int) []Face {
	//orientationDirection is essentially a int saying whether its an up pass or down pass,
	//orientationDirection = 1 means up pass
	//orientationDirection = -1 means down pass

	//make a copy of refSlice
	// slice := make([][][]bool, len(refSlice))
	// for x := 0; x < len(refSlice); x++ {
	// 	slice[x] = make([][]bool, len(refSlice[0]))
	// 	for y := 0; y < len(refSlice[0]); y++ {
	// 		slice[x][y] = make([]bool, len(refSlice[0][0]))
	// 	}
	// }

	// copy(slice, refSlice)

	outputFaces := []Face{}

	// fmt.Printf("\nLengths: %d, %d, %d\n", len(refSlice), len(refSlice[0]), len(refSlice[0][0]))
	
	var edgeOfArrayCount int = 0

	for dir := 0; dir < len(refSlice); dir++ {
		currentSlice := make([][]bool, len(refSlice[0]))
		for x := 0; x < len(refSlice[0]); x++ {
			currentSlice[x] = make([]bool, len(refSlice[0][0]))
		}

		copy(currentSlice, refSlice[dir])

		for x := 0; x < len(currentSlice); x++ {
			for y := 0; y < len(currentSlice[0]); y++ {
				var voxelPresent bool = currentSlice[x][y]
				if voxelPresent {
					// fmt.Printf("\nVoxel Present at %d, %d, %d for %d at %d direction\n", dir, x, y, orientationOffset, orientationDirection)
					var edgeOfArray bool = false
					if (dir == 0 && orientationDirection == -1) || (dir == len(refSlice)-1 && orientationDirection == 1) {
						edgeOfArrayCount += 1
						edgeOfArray = true
					}
					var faceFound, maxXReached bool = false, false
					//Creates a new array to record the corners of the face being generated. Initalizes it with the voxel (dir, x, y)
					var faceBounds [][3]int = [][3]int{[3]int{dir, x, y}}

					//Initalizes bounds for the x and y expansion of the face.
					cornerBounds := [2]int{x, y}

					for !faceFound {
						if !maxXReached {
							if currentSlice[cornerBounds[0]+1][y] && !(!edgeOfArray || refSlice[dir+orientationDirection][cornerBounds[0]+1][y]) {
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
								if !(currentSlice[faceBounds[0][1]+i][cornerBounds[1]+1] && !(!edgeOfArray || refSlice[dir+orientationDirection][faceBounds[0][1]+i][cornerBounds[1]+1])) {
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
						newFace = Face{VoxelCoords: faceBounds, FaceIndex: FaceOrientation(orientationOffset)}
					} else {
						newFace = Face{VoxelCoords: faceBounds, FaceIndex: FaceOrientation(orientationOffset + 1)}
					}
					outputFaces = append(outputFaces, newFace)

					//faceBounds[0] = (min, min)
					//faceBounds[1] = (max, min)
					//faceBounds[2] = (max, max)
					//faceBounds[3] = (min, max)

					//cornerBounds = {maxX, maxY}

					//Lastly now u have to set the voxelPresent bool to false for the voxels already considered into outputFaces
					// fmt.Println(currentSlice)

					if len(faceBounds) == 1 {
						currentSlice[faceBounds[0][1]][faceBounds[0][2]] = false
					} else if len(faceBounds) == 2 {
						for curX := 0; curX < cornerBounds[0] - faceBounds[1][1]; curX++ {
							currentSlice[curX][faceBounds[0][2]] = false
						}
					}

					for curX := faceBounds[0][1]; curX < cornerBounds[0]; curX++ {
						for curY := faceBounds[0][2]; curY < cornerBounds[1]; curY++ {
							currentSlice[curX][curY] = false
						}
					}
				}
			}
		}
	}

	normalizedFaces := make([]Face, len(outputFaces))

	for faceIndex, face := range outputFaces {
		var normalizedCoords [][3]int = make([][3]int, len(face.VoxelCoords))
		for index, cord := range face.VoxelCoords {
			normalizedCoord := voxelCordsOffsetter(cord, orientationOffset)
			normalizedCoords[index] = normalizedCoord
		}

		var normalizedFace Face = Face{VoxelCoords: normalizedCoords, FaceIndex: face.FaceIndex}
		normalizedFaces[faceIndex] = normalizedFace
	}

	fmt.Printf("\nEdgeOfArrayCount: %d\n", edgeOfArrayCount)

	// fmt.Printf("\nCombined Voxels for %d in %d direction\n", orientationOffset, orientationDirection)
	return normalizedFaces
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

func GenerateFaceJson(faces []Face) {
	file, errs := os.Create("Faces.JSON")
	if errs != nil {
		panic("Failed to write to file:" + errs.Error())
	}

	var jsonFaces []JSONFace = []JSONFace{}

	for _, face := range faces {
		var points []Point = []Point{}
		for _, point := range face.VoxelCoords {
			points = append(points, Point{X: point[0], Y: point[1], Z: point[2]})
		}
		jsonFaces = append(jsonFaces, JSONFace{VoxelCoords: points, FaceIndex: int(face.FaceIndex)})
	}

	enc := json.NewEncoder(file)
	enc.Encode(Faces{FaceArray: jsonFaces})
}
