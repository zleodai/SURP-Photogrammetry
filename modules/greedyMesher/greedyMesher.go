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

	outputFaces := []Face{}

	outputFaces = append(outputFaces, combineVoxels3(zxySlices, 0, true)...)
	outputFaces = append(outputFaces, combineVoxels3(zxySlices, 0, false)...)
	outputFaces = append(outputFaces, combineVoxels3(xyzSlices, 2, true)...)
	outputFaces = append(outputFaces, combineVoxels3(xyzSlices, 2, false)...)
	outputFaces = append(outputFaces, combineVoxels3(yxzSlices, 4, true)...)
	outputFaces = append(outputFaces, combineVoxels3(yxzSlices, 4, false)...)

	return outputFaces
}

func combineVoxels3(voxels [][][]bool, orientation int, isUp bool) []Face {
	var combinedVoxels []Face = []Face{}
	var dirOffset int
	if isUp {
		dirOffset = 1
	} else {
		dirOffset = -1
	}

	for dir, dirArray := range voxels {
		currentSlice := make([][]bool, len(voxels[0]))
		for x := 0; x < len(voxels[0]); x++ {
			currentSlice[x] = make([]bool, len(voxels[0][0]))
			for y := 0; y < len(voxels[0][0]); y++ {
				currentSlice[x][y] = voxels[dir][x][y]
			}
		}

		var assumeAir bool = false
		if (!isUp && dir == 0) || (isUp && dir == len(voxels)-1) {
			assumeAir = true
		}

		for x, xArray := range dirArray {
			for y := range xArray {
				if currentSlice[x][y] && (assumeAir || !voxels[dir+dirOffset][x][y]) {
					var xIncreasePossible bool = false
					var yIncreasePossible bool = false

					var corners [][3]int = [][3]int{{dir, x, y}}

					if (x+1 < len(currentSlice) && currentSlice[x+1][y]) && (assumeAir || !voxels[dir+dirOffset][x+1][y]) {
						xIncreasePossible = true
					} else if (y+1 < len(currentSlice[0]) && currentSlice[x][y+1]) && (assumeAir || !voxels[dir+dirOffset][x][y+1]) {
						yIncreasePossible = true
					}

					// fmt.Printf("\nxIncreasePossible: %t, yIncreasePossible: %t\n", xIncreasePossible, yIncreasePossible)
					if xIncreasePossible {
						corners = append(corners, [3]int{dir, x + 1, y})
						for xIncreasePossible {
							var currX int = corners[1][1]
							if (currX+1 < len(currentSlice) && currentSlice[currX+1][y]) && (assumeAir || !voxels[dir+dirOffset][currX+1][y]) {
								// fmt.Println("X Increase Possible")
								// fmt.Printf("	Size Constraint: %t, Voxel Present: %t, Assume Air: %t, Empty Voxel Above: %t\n", currX + 1 < len(currentSlice), currentSlice[currX+1][y], assumeAir, !voxels[dir+dirOffset][currX+1][y])
								corners[1][1] += 1
							} else {
								xIncreasePossible = false
								yIncreasePossible = true
							}
						}
						for yIncreasePossible {
							var currY = corners[1][2]
							for currX := corners[0][1]; currX <= corners[1][1]; currX++ {
								if !((currY+1 < len(currentSlice[0]) && currentSlice[currX][currY+1]) && (assumeAir || !voxels[dir+dirOffset][currX][currY+1])) {
									yIncreasePossible = false
								}
							}
							if yIncreasePossible {
								corners[1][2] += 1
							}
						}
					}
					if yIncreasePossible {
						corners = append(corners, [3]int{dir, x, y + 1})
						for yIncreasePossible {
							var currY int = corners[1][2]
							if (currY+1 < len(currentSlice[0]) && currentSlice[x][currY+1]) && (assumeAir || !voxels[dir+dirOffset][x][currY+1]) {
								// fmt.Println("Y Increase Possible")
								// fmt.Printf("	Size Constraint: %t, Voxel Present: %t, Assume Air: %t, Empty Voxel Above: %t\n", currY + 1 < len(currentSlice[0]), currentSlice[x][currY+1], assumeAir, !voxels[dir+dirOffset][x][currY+1])
								corners[1][2] += 1
							} else {
								yIncreasePossible = false
								xIncreasePossible = true
							}
						}
						for xIncreasePossible {
							var currX = corners[1][1]
							for currY := corners[0][2]; currY <= corners[1][2]; currY++ {
								if !((currX+1 < len(currentSlice) && currentSlice[currX+1][currY]) && (assumeAir || !voxels[dir+dirOffset][currX+1][currY])) {
									xIncreasePossible = false
								}
							}
							if xIncreasePossible {
								corners[1][1] += 1
							}
						}
					}

					switch length := len(corners); length {
					case 1:
						currentSlice[corners[0][1]][corners[0][2]] = false
					case 2:
						for currX := corners[0][1]; currX <= corners[1][1]; currX++ {
							for currY := corners[0][2]; currY <= corners[1][2]; currY++ {
								currentSlice[currX][currY] = false
							}
						}
					default:
						fmt.Println("Error on switch case. len(corners) must be either 1 or 2")
					}

					normalizedCorners := [][3]int{}
					for _, corner := range corners {
						newX, newY, newZ := remapXYZ(corner[0], corner[1], corner[2], FaceOrientation(orientation))
						normalizedCorners = append(normalizedCorners, [3]int{newX, newY, newZ})
					}

					var newFace Face
					if isUp {
						newFace = Face{VoxelCoords: normalizedCorners, FaceIndex: FaceOrientation(orientation)}
					} else {
						newFace = Face{VoxelCoords: normalizedCorners, FaceIndex: FaceOrientation(orientation + 1)}
					}
					combinedVoxels = append(combinedVoxels, newFace)
				}
			}
		}
	}
	return combinedVoxels
}

func combineVoxels2(voxels [][][]bool, orientation FaceOrientation, isUp bool) []Face {
	// fmt.Printf("\nCombining Voxels, Orientation: %d, isUp: %t\n", orientation, isUp)
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
			for y := 0; y < len(voxels[0][0]); y++ {
				currentSlice[x][y] = voxels[z][x][y]
			}
		}

		assumeAir := false
		if (!isUp && z == 0) || (isUp && z == len(voxels)-1) {
			assumeAir = true
		}
		for x := 0; x < len(currentSlice); x++ {
			for y := 0; y < len(currentSlice[0]); y++ {
				currVoxel := currentSlice[x][y]
				if currVoxel && (assumeAir || !voxels[z+checkOffset][x][y]) {
					// fmt.Printf("\n	CurrentVoxel: [%d, %d, %d]\n", z, x, y)
					// fmt.Printf("\n	IsFace: !voxels[%d + %d][%d][%d]", z, checkOffset, x, y)
					x1, y1, z1 := remapXYZ(x, y, z, orientation)
					corners := [][3]int{{x1, y1, z1}}

					currY, actualY := y, y
					currX := x
					nextFace, faceComplete := true, false
					if y+1 != len(voxels[0][0]) {
						nextFace = currentSlice[currX][currY+1] && (assumeAir || !voxels[z+checkOffset][currX][currY])
					} else {
						nextFace = false
					}
					for !faceComplete {
						for nextFace {
							if currY+1 != len(currentSlice[0]) {
								currY += 1
								nextFace = currentSlice[currX][currY] && (assumeAir || !voxels[z+checkOffset][currX][currY])
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
						for startY := y; startY <= actualY; startY++ {
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

func remapXYZ(x, y, z int, orientation FaceOrientation) (int, int, int) {
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
		return z, y, x
	case backward:
		return z, y, x
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
						edgeOfArray = true
					}
					var faceFound, maxXReached bool = false, false
					//Creates a new array to record the corners of the face being generated. Initalizes it with the voxel (dir, x, y)
					var faceBounds [][3]int = [][3]int{{dir, x, y}}

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
						for curX := 0; curX < cornerBounds[0]-faceBounds[1][1]; curX++ {
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
