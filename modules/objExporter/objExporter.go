package objExporter

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"modules/greedyMesher"
	"modules/pointSorter"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
)

type PointData struct {
	Points    []Point
	VoxelSize float32
}

type Point struct {
	X     float32
	Y     float32
	Z     float32
	Value float32
}

func GetMeshFacesFromVertices(faces []greedyMesher.Face, faceMap map[string][][3]int, vertices [][3]int, vertexMatrix [][][]bool, vertexMap map[string]int) [][3]int {
	meshFaces := make([][3]int, 0)

	for _, face := range faces {
		faceVerticies := faceMap[getStringFromFace(face)]
		involvedVertices := [][3]int{}

		var ignoreZ bool = faceVerticies[0][2] == faceVerticies[1][2]
		var ignoreX bool = faceVerticies[0][0] == faceVerticies[1][0]
		var ignoreY bool = faceVerticies[0][1] == faceVerticies[1][1]
		
		var preppOffset int
		
		//ABCD are the respective corners of the box created from two vertices, [min, min] [max, max]
		//AB will be the inital up pass
		//BC will be the following pass rigthwards
		//CD will be the downpass
		//DA lastly will be the pass leftward
		if ignoreZ {
			preppOffset = 0
			//Ignore Z here
			//AB, abIndex being yMin to yMax, x set to xMin
			for abIndex := faceVerticies[0][1]; abIndex <= faceVerticies[1][1]; abIndex++ {
				if vertexMatrix[faceVerticies[0][0]][abIndex][faceVerticies[0][2]] {
					var newVertex [3]int = [3]int{faceVerticies[0][0], abIndex, faceVerticies[0][2]}
					//First check if vertex is already in the array
					var alreadyInArray bool = false
					for _, vertex := range involvedVertices {
						if reflect.DeepEqual(newVertex, vertex) {
							alreadyInArray = true
						}
					}
					if !alreadyInArray {
						involvedVertices = append(involvedVertices, newVertex)
					}
				}
			}
			//BC, bcIndex being xMin to xMax, y set to yMax
			for bcIndex := faceVerticies[0][0]; bcIndex <= faceVerticies[1][0]; bcIndex++ {
				if vertexMatrix[bcIndex][faceVerticies[1][1]][faceVerticies[0][2]] {
					var newVertex [3]int = [3]int{bcIndex, faceVerticies[1][1], faceVerticies[0][2]}
					var alreadyInArray bool = false
					for _, vertex := range involvedVertices {
						if reflect.DeepEqual(newVertex, vertex) {
							alreadyInArray = true
						}
					}
					if !alreadyInArray {
						involvedVertices = append(involvedVertices, newVertex)
					}
				}
			}
			//CD, cdIndex being yMax to yMin, x set to xMax
			for cdIndex := faceVerticies[1][1]; cdIndex >= faceVerticies[0][1]; cdIndex-- {
				if vertexMatrix[faceVerticies[1][0]][cdIndex][faceVerticies[0][2]] {
					var newVertex [3]int = [3]int{faceVerticies[1][0], cdIndex, faceVerticies[0][2]}
					var alreadyInArray bool = false
					for _, vertex := range involvedVertices {
						if reflect.DeepEqual(newVertex, vertex) {
							alreadyInArray = true
						}
					}
					if !alreadyInArray {
						involvedVertices = append(involvedVertices, newVertex)
					}
				}
			}
			//DA, daIndex being xMax to xMin, y set to yMin
			for daIndex := faceVerticies[1][0]; daIndex >= faceVerticies[0][0]; daIndex-- {
				if vertexMatrix[daIndex][faceVerticies[0][1]][faceVerticies[0][2]] {
					var newVertex [3]int = [3]int{daIndex, faceVerticies[0][1], faceVerticies[0][2]}
					var alreadyInArray bool = false
					for _, vertex := range involvedVertices {
						if reflect.DeepEqual(newVertex, vertex) {
							alreadyInArray = true
						}
					}
					if !alreadyInArray {
						involvedVertices = append(involvedVertices, newVertex)
					}
				}
			}
		} else if ignoreX {
			preppOffset = 1
			//Ignore X here
			//AB, abIndex being yMin to yMax, z set to zMin
			for abIndex := faceVerticies[0][1]; abIndex <= faceVerticies[1][1]; abIndex++ {
				if vertexMatrix[faceVerticies[0][0]][abIndex][faceVerticies[0][2]] {
					var newVertex [3]int = [3]int{faceVerticies[0][0], abIndex, faceVerticies[0][2]}
					//First check if vertex is already in the array
					var alreadyInArray bool = false
					for _, vertex := range involvedVertices {
						if reflect.DeepEqual(newVertex, vertex) {
							alreadyInArray = true
						}
					}
					if !alreadyInArray {
						involvedVertices = append(involvedVertices, newVertex)
					}
				}
			}
			//BC, bcIndex being zMin to zMax, y set to yMax
			for bcIndex := faceVerticies[0][2]; bcIndex <= faceVerticies[1][2]; bcIndex++ {
				if vertexMatrix[faceVerticies[0][0]][faceVerticies[1][1]][bcIndex] {
					var newVertex [3]int = [3]int{faceVerticies[0][0], faceVerticies[1][1], bcIndex}
					var alreadyInArray bool = false
					for _, vertex := range involvedVertices {
						if reflect.DeepEqual(newVertex, vertex) {
							alreadyInArray = true
						}
					}
					if !alreadyInArray {
						involvedVertices = append(involvedVertices, newVertex)
					}
				}
			}
			//CD, cdIndex being yMax to yMin, z set to zMax
			for cdIndex := faceVerticies[1][1]; cdIndex >= faceVerticies[0][1]; cdIndex-- {
				if vertexMatrix[faceVerticies[0][0]][cdIndex][faceVerticies[1][2]] {
					var newVertex [3]int = [3]int{faceVerticies[1][0], cdIndex, faceVerticies[0][2]}
					var alreadyInArray bool = false
					for _, vertex := range involvedVertices {
						if reflect.DeepEqual(newVertex, vertex) {
							alreadyInArray = true
						}
					}
					if !alreadyInArray {
						involvedVertices = append(involvedVertices, newVertex)
					}
				}
			}
			//DA, daIndex being zMax to zMin, y set to yMin
			for daIndex := faceVerticies[1][2]; daIndex >= faceVerticies[0][2]; daIndex-- {
				if vertexMatrix[faceVerticies[0][0]][faceVerticies[0][1]][daIndex] {
					var newVertex [3]int = [3]int{faceVerticies[0][0], faceVerticies[0][1], daIndex}
					var alreadyInArray bool = false
					for _, vertex := range involvedVertices {
						if reflect.DeepEqual(newVertex, vertex) {
							alreadyInArray = true
						}
					}
					if !alreadyInArray {
						involvedVertices = append(involvedVertices, newVertex)
					}
				}
			}
		} else if ignoreY {
			preppOffset = 2
			//Ignore Y here
			//AB, abIndex being zMin to zMax, x set to xMin
			for abIndex := faceVerticies[0][2]; abIndex <= faceVerticies[1][2]; abIndex++ {
				if vertexMatrix[faceVerticies[0][0]][faceVerticies[0][1]][abIndex] {
					var newVertex [3]int = [3]int{faceVerticies[0][0], faceVerticies[0][1], abIndex}
					//First check if vertex is already in the array
					var alreadyInArray bool = false
					for _, vertex := range involvedVertices {
						if reflect.DeepEqual(newVertex, vertex) {
							alreadyInArray = true
						}
					}
					if !alreadyInArray {
						involvedVertices = append(involvedVertices, newVertex)
					}
				}
			}
			//BC, bcIndex being xMin to xMax, z set to zMax
			for bcIndex := faceVerticies[0][0]; bcIndex <= faceVerticies[1][0]; bcIndex++ {
				if vertexMatrix[bcIndex][faceVerticies[0][1]][faceVerticies[1][2]] {
					var newVertex [3]int = [3]int{bcIndex, faceVerticies[0][1], faceVerticies[1][2]}
					var alreadyInArray bool = false
					for _, vertex := range involvedVertices {
						if reflect.DeepEqual(newVertex, vertex) {
							alreadyInArray = true
						}
					}
					if !alreadyInArray {
						involvedVertices = append(involvedVertices, newVertex)
					}
				}
			}
			//CD, cdIndex being zMax to zMin, x set to xMax
			for cdIndex := faceVerticies[1][2]; cdIndex >= faceVerticies[0][2]; cdIndex-- {
				if vertexMatrix[faceVerticies[1][0]][faceVerticies[0][1]][cdIndex] {
					var newVertex [3]int = [3]int{faceVerticies[1][0], faceVerticies[0][1], cdIndex}
					var alreadyInArray bool = false
					for _, vertex := range involvedVertices {
						if reflect.DeepEqual(newVertex, vertex) {
							alreadyInArray = true
						}
					}
					if !alreadyInArray {
						involvedVertices = append(involvedVertices, newVertex)
					}
				}
			}
			//DA, daIndex being xMax to xMin, z set to zMin
			for daIndex := faceVerticies[1][0]; daIndex >= faceVerticies[0][0]; daIndex-- {
				if vertexMatrix[daIndex][faceVerticies[0][1]][faceVerticies[0][2]] {
					var newVertex [3]int = [3]int{daIndex, faceVerticies[0][1], faceVerticies[0][2]}
					var alreadyInArray bool = false
					for _, vertex := range involvedVertices {
						if reflect.DeepEqual(newVertex, vertex) {
							alreadyInArray = true
						}
					}
					if !alreadyInArray {
						involvedVertices = append(involvedVertices, newVertex)
					}
				}
			}
		}
		preppedVertices, preppedVertixMap := prepFaceCoordsForTriangle(involvedVertices, int(face.FaceIndex), vertexMap, preppOffset)
		result := triangulateVertices(preppedVertices, preppedVertixMap)
		meshFaces = append(meshFaces, result...)
	}

	return meshFaces
}

func prepFaceCoordsForTriangle(vertices [][3]int, faceOrientation int, vertexMap map[string]int, preppOffset int) ([][3]int, map[string]int) {
	preppedVertices := [][3]int{}
	preppedVertixMap := map[string]int{}

	switch preppOffset {
	case 0:
		//z is the same
		return vertices, vertexMap
	case 1:
		//x is the same
		for _, vertex := range vertices {
			index := vertexMap[getStringFromIntVertex(vertex)]
			preppedVertex := [3]int{vertex[1], vertex[2], vertex[0]}
			preppedVertices = append(preppedVertices, preppedVertex)
			preppedVertixMap[getStringFromIntVertex(preppedVertex)] = index
		}
	case 2:
		//y is the same
		for _, vertex := range vertices {
			index := vertexMap[getStringFromIntVertex(vertex)]
			preppedVertex := [3]int{vertex[0], vertex[2], vertex[1]}
			preppedVertices = append(preppedVertices, preppedVertex)
			preppedVertixMap[getStringFromIntVertex(preppedVertex)] = index
		}
	default:
		fmt.Printf("\nError, preppOffset in prepFaceCoordsForTriangle() must be from 0 - 2, Got: %d", preppOffset)
	}

	return preppedVertices, preppedVertixMap
}

func triangulateVertices(vertices [][3]int, vertexMap map[string]int) [][3]int {
	triangulatedFaces := make([][3]int, 0)
	createdEdges := make([][2][3]int, 0)
	createdEdgeColliders := make([][2][3]float32, 0)
	//For the rest of the function we assume createdEdges and createdEdgeColliders will contain same elements with the same index locations. Only difference is createdEdges contains the int versions of the edge
	//We first iterate through all the vertices to create all the triangles on the sides
	for index, vertex := range vertices {
		var debug bool = false

		if index == 404 {
			debug = true
		}

		if debug {
			fmt.Printf("\n		For vertex %d: [%d, %d]", index, vertex[0], vertex[1])
		}

		if debug {
			fmt.Printf("\n		EdgeColliders %f", createdEdgeColliders)
		}

		//targetVertex is the vertex that the vertex will attempt to create a triangle with
		//for now we assume that each vertex automatically will attempt to form a triangle with its proceeding vertex
		var targetVertex [3]int
		var targetVertexIndex int
		if index < len(vertices)-1 {
			targetVertex = vertices[index+1]
			targetVertexIndex = index + 1
		} else {
			targetVertex = vertices[0]
			targetVertexIndex = 0
		}

		//Check to see if an edge with the vertex and targetVertex has already been created
		var vertexConnectedToTarget bool = false
		for _, edge := range createdEdges {
			var vertexFound bool = false
			var targetVertexFound bool = false

			for _, edgeVertex := range edge {
				if reflect.DeepEqual(vertex, edgeVertex) {
					vertexFound = true
				} else if reflect.DeepEqual(targetVertex, edgeVertex) {
					targetVertexFound = true
				}
			}

			if vertexFound && targetVertexFound {
				vertexConnectedToTarget = true
			}
		}

		if !vertexConnectedToTarget {
			if debug {
				fmt.Printf("\n		Target Vertex %d: [%d, %d]", index+1, targetVertex[0], targetVertex[1])
			}
			//If an edge was not created we know that a triangle with the targetVertex has not been created
			var midPoint = getMidPoint(vertex, targetVertex)

			if debug {
				fmt.Printf("\n		Middle Vertex: [%f, %f]", midPoint[0], midPoint[1])
			}

			var midPointOnEdge bool = false
			var midPointEdgeIndex int = 0

			//First check if middlepoint is on an existing edge
			for intEdgeIndex, edge := range createdEdgeColliders {
				if determineOrientation(edge[0], midPoint, edge[1]) == 0 && determineEdgePointIntersection(edge, midPoint) {
					midPointEdgeIndex = intEdgeIndex
					midPointOnEdge = true
				}
			}

			if debug && midPointOnEdge {
				fmt.Printf("\n		midPointEdgeIndex: %d", midPointEdgeIndex)
			}

			//Iterate through otherVertexs to see which ones can make a triangle with targetVertex
			var triangleCreated bool = false
			
			var heapArray = []pointSorter.HeapItem{}


			for otherIndex, otherVertex := range vertices {
				var distance float32 = getDistance(vertex, otherVertex)
				point := pointSorter.HeapItem{Value:[4]int{otherVertex[0], otherVertex[1], otherVertex[2], otherIndex}, Priority:float64(distance), Index:otherIndex}
				heapArray = append(heapArray, point)
			}

			distanceOrderedVerticeArray, indexArray := pointSorter.SortVertices(heapArray)


			for x, otherVertex := range distanceOrderedVerticeArray {
				otherIndex := indexArray[x]
				if debug {
					fmt.Printf("\n\n			Other Vertex %d: [%d, %d]", otherIndex, otherVertex[0], otherVertex[1])
				}
				if !triangleCreated && (otherIndex != index && otherIndex != targetVertexIndex) {
					var collisionDetected bool = false
					var newEdge [2][3]float32 = [2][3]float32{midPoint, {float32(otherVertex[0]), float32(otherVertex[1]), float32(otherVertex[2])}}

					//Iterate through all edgeColliders (except for any edgeCollider midPoint is already on) to detect if creating a new edge will cause any collisions
					for edgeIndex, edge := range createdEdgeColliders {
						if !collisionDetected && (!triangleCreated && !(midPointOnEdge && edgeIndex == midPointEdgeIndex)) {
							if determineCollision(newEdge, edge) {
								if debug {
									fmt.Print("\n				Collision 1 detected")
								}
								if debug {
									fmt.Printf("\n				Collision on Edge [%f, %f], [%f, %f]", edge[0][0], edge[0][1], edge[1][0], edge[1][1])
								}
								collisionDetected = true
							}
						}
					}

					//Special Cases
					if !collisionDetected {
						var float32Vertex [3]float32 = [3]float32{float32(vertex[0]), float32(vertex[1]), float32(vertex[2])}
						var float32TargetVertex [3]float32 = [3]float32{float32(targetVertex[0]), float32(targetVertex[1]), float32(targetVertex[2])}
						//Check and see if the newEdge will collide with targetVertex
						if determineOrientation(newEdge[0], newEdge[1], float32Vertex) == 0 && determineEdgePointIntersection(newEdge, float32Vertex) {
							if debug {
								fmt.Print("\n				Collision 2 detected")
							}
							collisionDetected = true
						}
						//Check and see if the newEdge will collide with vertex
						if determineOrientation(newEdge[0], newEdge[1], float32TargetVertex) == 0 && determineEdgePointIntersection(newEdge, float32TargetVertex) {
							if debug {
								fmt.Print("\n				Collision 3 detected")
							}
							collisionDetected = true
						}

						if !collisionDetected {
							//Check to see if an edge created by vertex and otherVertex will collide with other points
							var float32OtherVertex [3]float32 = [3]float32{float32(otherVertex[0]), float32(otherVertex[1]), float32(otherVertex[2])}
							for indexX, pointX := range vertices {
								if !collisionDetected && (indexX != index && indexX != targetVertexIndex && indexX != otherIndex) {
									var float32PointX [3]float32 = [3]float32{float32(pointX[0]), float32(pointX[1]), float32(pointX[2])}
									if determineOrientation(float32Vertex, float32OtherVertex, float32PointX) == 0 && determineEdgePointIntersection([2][3]float32{float32Vertex, float32OtherVertex}, float32PointX) {
										if debug {
											fmt.Print("\n				Collision 4 detected")
										}
										collisionDetected = true
									}
								}
							}
						}

						if !collisionDetected {
							//Check to see if an edge created by targetVertex and otherVertex will collide with other points
							var float32OtherVertex [3]float32 = [3]float32{float32(otherVertex[0]), float32(otherVertex[1]), float32(otherVertex[2])}
							for indexX, pointX := range vertices {
								if !collisionDetected && (indexX != index && indexX != targetVertexIndex && indexX != otherIndex) {
									var float32PointX [3]float32 = [3]float32{float32(pointX[0]), float32(pointX[1]), float32(pointX[2])}
									if determineOrientation(float32TargetVertex, float32OtherVertex, float32PointX) == 0 && determineEdgePointIntersection([2][3]float32{float32TargetVertex, float32OtherVertex}, float32PointX) {
										if debug {
											fmt.Print("\n				Collision 5 detected")
										}
										collisionDetected = true
									}
								}
							}
						}

						if !collisionDetected {
							//Check to see if midpoint is the same point as another point
							for _, vertexX := range vertices {
								var float32VertexX [3]float32 = [3]float32{float32(vertexX[0]), float32(vertexX[1]), float32(vertexX[2])}
								if reflect.DeepEqual(midPoint, float32VertexX) {
									if debug {
										fmt.Print("\n					Collision 6 detected")
									}
									collisionDetected = true
								}
							}
						}

						if !collisionDetected {
							//Check to see if triangle created will harbor another point
							for indexX, vertexX := range vertices {
								if indexX != index && indexX != targetVertexIndex && indexX != otherIndex{
									if determineIfVertexInTriangle(vertex, targetVertex, otherVertex, vertexX) {
										if debug {
											fmt.Print("\n					Collision 7 detected")
										}
										collisionDetected = true
									}
								}
							}
						}
					}

					if debug {
						fmt.Printf("\n				collisionDetected: %t", collisionDetected)
					}

					//When no collision is detected we make the new triangle
					if !collisionDetected {
						if debug {
							fmt.Printf("\n				No Collision Detected For [%d, %d]", otherVertex[0], otherVertex[1])
						}
						var newFace [3]int = [3]int{vertexMap[getStringFromIntVertex(vertex)], vertexMap[getStringFromIntVertex(targetVertex)], vertexMap[getStringFromIntVertex(otherVertex)]}

						//First try to detect if the triangle we are making is has already been created
						var triangleAlreadyCreated bool = false

						for _, face := range triangulatedFaces {
							newFaceCopy := []int{newFace[0], newFace[1], newFace[2]}
							sort.Ints(newFaceCopy)
							faceCopy := []int{face[0], face[1], face[2]}
							sort.Ints(faceCopy)
							if reflect.DeepEqual(newFaceCopy, faceCopy) {
								triangleAlreadyCreated = true
							}
						}

						if !triangleAlreadyCreated {
							triangulatedFaces = append(triangulatedFaces, newFace)

							var abIntEdge [2][3]int = [2][3]int{vertex, targetVertex}
							var bcIntEdge [2][3]int = [2][3]int{targetVertex, otherVertex}
							var caIntEdge [2][3]int = [2][3]int{otherVertex, vertex}
							createdEdges = append(createdEdges, abIntEdge, bcIntEdge, caIntEdge)

							var abEdge [2][3]float32 = edgeOffsetter(abIntEdge)
							var bcEdge [2][3]float32 = edgeOffsetter(bcIntEdge)
							var caEdge [2][3]float32 = edgeOffsetter(caIntEdge)
							createdEdgeColliders = append(createdEdgeColliders, abEdge, bcEdge, caEdge)
							triangleCreated = true
						}
					}
				}
			}
		}
	}

	//Now we attempt to fill in any remaining triangles
	for index, vertex := range vertices {
		var targetVertex [3]int
		var targetVertexIndex int = index
		var allPossibleTrianglesMade bool = false

		var debug bool = false

		if index == 404 {
			debug = true
		}

		if debug {
			fmt.Printf("\n		For Vertex %d: [%d, %d]", index, vertex[0], vertex[1])
		}

		for !allPossibleTrianglesMade {
			var triangleMade bool = false
			var oldTargetVertexIndex int = targetVertexIndex

			//Check all vertices to see if a new targetVertex can be found
			for otherIndex, otherVertex := range vertices {
				if otherIndex != index && otherIndex != oldTargetVertexIndex {
					var newEdge [2][3]float32 = [2][3]float32{{float32(vertex[0]), float32(vertex[1]), float32(vertex[1])}, {float32(otherVertex[0]), float32(otherVertex[1]), float32(otherVertex[1])}}

					//Check to see if newEdge created collides with any existing colliders
					var collisionFound bool = false
					for _, edge := range createdEdgeColliders {
						if determineCollision(newEdge, edge) {
							//First Check to see if newEdge is parralel with existing collider, if it isnt parralell collision is found
							if !(determineOrientation(newEdge[0], newEdge[1], edge[0]) == 0 && determineOrientation(newEdge[0], newEdge[1], edge[1]) == 0) {
								collisionFound = true
							}
						}
					}

					if !collisionFound {
						if debug {
							fmt.Printf("\n\n		Found Target Vertex %d: [%d, %d]", otherIndex, otherVertex[0], otherVertex[1])
						}
						targetVertex = otherVertex
						targetVertexIndex = otherIndex

						//Repeat code from original triangle creation on sides
						if debug {
							fmt.Printf("\n			Target Vertex %d: [%d, %d]", targetVertexIndex, targetVertex[0], targetVertex[1])
						}
						//If an edge was not created we know that a triangle with the targetVertex has not been created
						var midPoint = getMidPoint(vertex, targetVertex)

						if debug {
							fmt.Printf("\n			Middle Vertex: [%f, %f]", midPoint[0], midPoint[1])
						}

						var midPointOnEdge bool = false
						var midPointEdgeIndex int = 0

						//First check if middlepoint is on an existing edge
						for intEdgeIndex, edge := range createdEdgeColliders {
							if determineOrientation(edge[0], midPoint, edge[1]) == 0 && determineEdgePointIntersection(edge, midPoint) {
								midPointEdgeIndex = intEdgeIndex
								midPointOnEdge = true
							}
						}

						if debug && midPointOnEdge {
							fmt.Printf("\n			midPointEdgeIndex: %d", midPointEdgeIndex)
						}

						//Iterate through otherVertexs to see which ones can make a triangle with targetVertex
						var triangleCreated bool = false
						var heapArray = []pointSorter.HeapItem{}

						for otherIndex, otherVertex := range vertices {
							var distance float32 = getDistance(vertex, otherVertex)
							point := pointSorter.HeapItem{Value:[4]int{otherVertex[0], otherVertex[1], otherVertex[2], otherIndex}, Priority:float64(distance), Index:otherIndex}
							heapArray = append(heapArray, point)
						}

						distanceOrderedVerticeArray, indexArray := pointSorter.SortVertices(heapArray)

						for x, otherVertex := range distanceOrderedVerticeArray {
							otherIndex := indexArray[x]
							if debug {
								fmt.Printf("\n\n				Other Vertex %d: [%d, %d]", otherIndex, otherVertex[0], otherVertex[1])
							}
							if !triangleCreated && (otherIndex != index && otherIndex != targetVertexIndex) {
								var collisionDetected bool = false
								var newEdge [2][3]float32 = [2][3]float32{midPoint, {float32(otherVertex[0]), float32(otherVertex[1]), float32(otherVertex[2])}}

								//Iterate through all edgeColliders (except for any edgeCollider midPoint is already on) to detect if creating a new edge will cause any collisions
								for edgeIndex, edge := range createdEdgeColliders {
									if !collisionDetected && (!triangleCreated && !(midPointOnEdge && edgeIndex == midPointEdgeIndex)) {
										if determineCollision(newEdge, edge) {
											if debug {
												fmt.Print("\n					Collision 1 detected")
											}
											collisionDetected = true
										}
									}
								}

								//Special Cases
								if !collisionDetected {
									var float32Vertex [3]float32 = [3]float32{float32(vertex[0]), float32(vertex[1]), float32(vertex[2])}
									var float32TargetVertex [3]float32 = [3]float32{float32(targetVertex[0]), float32(targetVertex[1]), float32(targetVertex[2])}
									//Check and see if the newEdge will collide with targetVertex
									if determineOrientation(newEdge[0], newEdge[1], float32Vertex) == 0 && determineEdgePointIntersection(newEdge, float32Vertex) {
										if debug {
											fmt.Print("\n					Collision 2 detected")
										}
										collisionDetected = true
									}
									//Check and see if the newEdge will collide with vertex
									if determineOrientation(newEdge[0], newEdge[1], float32TargetVertex) == 0 && determineEdgePointIntersection(newEdge, float32TargetVertex) {
										if debug {
											fmt.Print("\n					Collision 3 detected")
										}
										collisionDetected = true
									}

									if !collisionDetected {
										//Check to see if an edge created by vertex and otherVertex will collide with other points
										var float32OtherVertex [3]float32 = [3]float32{float32(otherVertex[0]), float32(otherVertex[1]), float32(otherVertex[2])}
										for indexX, pointX := range vertices {
											if !collisionDetected && (indexX != index && indexX != targetVertexIndex && indexX != otherIndex) {
												var float32PointX [3]float32 = [3]float32{float32(pointX[0]), float32(pointX[1]), float32(pointX[2])}
												if determineOrientation(float32Vertex, float32OtherVertex, float32PointX) == 0 && determineEdgePointIntersection([2][3]float32{float32Vertex, float32OtherVertex}, float32PointX) {
													if debug {
														fmt.Print("\n					Collision 4 detected")
													}
													collisionDetected = true
												}
											}
										}
									}

									if !collisionDetected {
										//Check to see if an edge created by targetVertex and otherVertex will collide with other points
										var float32OtherVertex [3]float32 = [3]float32{float32(otherVertex[0]), float32(otherVertex[1]), float32(otherVertex[2])}
										for indexX, pointX := range vertices {
											if !collisionDetected && (indexX != index && indexX != targetVertexIndex && indexX != otherIndex) {
												var float32PointX [3]float32 = [3]float32{float32(pointX[0]), float32(pointX[1]), float32(pointX[2])}
												if determineOrientation(float32TargetVertex, float32OtherVertex, float32PointX) == 0 && determineEdgePointIntersection([2][3]float32{float32TargetVertex, float32OtherVertex}, float32PointX) {
													if debug {
														fmt.Print("\n					Collision 5 detected")
													}
													collisionDetected = true
												}
											}
										}
									}

									if !collisionDetected {
										//Check to see if midpoint is the same point as another point
										for _, vertexX := range vertices {
											var float32VertexX [3]float32 = [3]float32{float32(vertexX[0]), float32(vertexX[1]), float32(vertexX[2])}
											if reflect.DeepEqual(midPoint, float32VertexX) {
												if debug {
													fmt.Print("\n					Collision 6 detected")
												}
												collisionDetected = true
											}
										}
									}

									if !collisionDetected {
										//Check to see if triangle created will harbor another point
										for indexX, vertexX := range vertices {
											if indexX != index && indexX != targetVertexIndex && indexX != otherIndex{
												if determineIfVertexInTriangle(vertex, targetVertex, otherVertex, vertexX) {
													if debug {
														fmt.Print("\n					Collision 7 detected")
													}
													collisionDetected = true
												}
											}
										}
									}
								}

								if debug {
									fmt.Printf("\n					collisionDetected: %t", collisionDetected)
								}

								//When no collision is detected we make the new triangle
								if !collisionDetected {
									if debug {
										fmt.Printf("\n					No Collision Detected For [%d, %d]", otherVertex[0], otherVertex[1])
									}
									var newFace [3]int = [3]int{vertexMap[getStringFromIntVertex(vertex)], vertexMap[getStringFromIntVertex(targetVertex)], vertexMap[getStringFromIntVertex(otherVertex)]}

									//First try to detect if the triangle we are making is has already been created
									var triangleAlreadyCreated bool = false

									for _, face := range triangulatedFaces {
										newFaceCopy := []int{newFace[0], newFace[1], newFace[2]}
										sort.Ints(newFaceCopy)
										faceCopy := []int{face[0], face[1], face[2]}
										sort.Ints(faceCopy)
										if reflect.DeepEqual(newFaceCopy, faceCopy) {
											triangleAlreadyCreated = true
											if debug {
												fmt.Print("\n					triangleAlreadyCreated")
											}
										}
									}

									if !triangleAlreadyCreated {
										if debug {
											fmt.Print("\n						Creating new Triangle")
										}
										triangulatedFaces = append(triangulatedFaces, newFace)

										var abIntEdge [2][3]int = [2][3]int{vertex, targetVertex}
										var bcIntEdge [2][3]int = [2][3]int{targetVertex, otherVertex}
										var caIntEdge [2][3]int = [2][3]int{otherVertex, vertex}
										createdEdges = append(createdEdges, abIntEdge, bcIntEdge, caIntEdge)

										var abEdge [2][3]float32 = edgeOffsetter(abIntEdge)
										var bcEdge [2][3]float32 = edgeOffsetter(bcIntEdge)
										var caEdge [2][3]float32 = edgeOffsetter(caIntEdge)
										createdEdgeColliders = append(createdEdgeColliders, abEdge, bcEdge, caEdge)
										triangleMade = true
									}
								}
							}
						}
					}
				}
			}

			//Last check to see if no new triangles were made
			if !triangleMade {
				allPossibleTrianglesMade = true
			}
		}
	}

	return triangulatedFaces
}

func TriangulateVerticesTester(testTime bool) {
	fmt.Println("\nTesting triangulateVertices()")

	var p1 [3]int
	var p2 [3]int
	var p3 [3]int
	var p4 [3]int
	var p5 [3]int
	var p6 [3]int
	var p7 [3]int
	var p8 [3]int
	var vertices [][3]int
	var vertexMap map[string]int
	var result [][3]int

	p1 = [3]int{0, 0, 0}
	p2 = [3]int{0, 6, 0}
	p3 = [3]int{4, 6, 0}
	p4 = [3]int{6, 6, 0}
	p5 = [3]int{6, 3, 0}
	p6 = [3]int{6, 0, 0}
	p7 = [3]int{4, 0, 0}
	p8 = [3]int{1, 0, 0}
	vertices = [][3]int{p1, p2, p3, p4, p5, p6, p7, p8}
	vertexMap = map[string]int{}
	vertexMap[getStringFromIntVertex(p1)] = 0
	vertexMap[getStringFromIntVertex(p2)] = 1
	vertexMap[getStringFromIntVertex(p3)] = 2
	vertexMap[getStringFromIntVertex(p4)] = 3
	vertexMap[getStringFromIntVertex(p5)] = 4
	vertexMap[getStringFromIntVertex(p6)] = 5
	vertexMap[getStringFromIntVertex(p7)] = 6
	vertexMap[getStringFromIntVertex(p8)] = 7

	fmt.Print("\nTest Case 5\n	Input: ")
	fmt.Print(vertices)
	result = triangulateVertices(vertices, vertexMap)
	fmt.Print("\n	Result:")
	fmt.Print(result)
	fmt.Println()

	if testTime {
		fmt.Println("\nTesting 100000 simple runs")
		start := time.Now()
		for i := 0; i < 100000; i++ {
			p1 = [3]int{1 * rand.IntN(100), 1 * rand.IntN(100), 0}
			p2 = [3]int{1 * rand.IntN(100), 2 * rand.IntN(100), 0}
			p3 = [3]int{2 * rand.IntN(100), 2 * rand.IntN(100), 0}
			p4 = [3]int{2 * rand.IntN(100), 1 * rand.IntN(100), 0}
			vertices = [][3]int{p1, p2, p3, p4}
			vertexMap = map[string]int{}
			vertexMap[getStringFromIntVertex(p1)] = 0
			vertexMap[getStringFromIntVertex(p2)] = 1
			vertexMap[getStringFromIntVertex(p3)] = 2
			vertexMap[getStringFromIntVertex(p4)] = 3
			triangulateVertices(vertices, vertexMap)
		}
		duration := time.Since(start)

		fmt.Println("Execution Took: ")
		fmt.Print(duration)
		fmt.Println()

		fmt.Println("\nTesting 60000 complex runs")
		start = time.Now()
		for i := 0; i < 20000; i++ {
			p1 = [3]int{1 * rand.IntN(100), 1 * rand.IntN(100), 0}
			p2 = [3]int{1 * rand.IntN(100), 2 * rand.IntN(100), 0}
			p3 = [3]int{1 * rand.IntN(100), 5 * rand.IntN(100), 0}
			p4 = [3]int{5 * rand.IntN(100), 5 * rand.IntN(100), 0}
			p5 = [3]int{5 * rand.IntN(100), 4 * rand.IntN(100), 0}
			p6 = [3]int{5 * rand.IntN(100), 1 * rand.IntN(100), 0}
			vertices = [][3]int{p1, p2, p3, p4, p5, p6}
			vertexMap = map[string]int{}
			vertexMap[getStringFromIntVertex(p1)] = 0
			vertexMap[getStringFromIntVertex(p2)] = 1
			vertexMap[getStringFromIntVertex(p3)] = 2
			vertexMap[getStringFromIntVertex(p4)] = 3
			vertexMap[getStringFromIntVertex(p5)] = 4
			vertexMap[getStringFromIntVertex(p6)] = 5
			triangulateVertices(vertices, vertexMap)

			p1 = [3]int{1 * rand.IntN(100), 1 * rand.IntN(100), 0}
			p2 = [3]int{1 * rand.IntN(100), 4 * rand.IntN(100), 0}
			p3 = [3]int{2 * rand.IntN(100), 4 * rand.IntN(100), 0}
			p4 = [3]int{2 * rand.IntN(100), 5 * rand.IntN(100), 0}
			p5 = [3]int{1 * rand.IntN(100), 5 * rand.IntN(100), 0}
			p6 = [3]int{1 * rand.IntN(100), 8 * rand.IntN(100), 0}
			p7 = [3]int{8 * rand.IntN(100), 8 * rand.IntN(100), 0}
			p8 = [3]int{8 * rand.IntN(100), 1 * rand.IntN(100), 0}
			vertices = [][3]int{p1, p2, p3, p4, p5, p6, p7, p8}
			vertexMap = map[string]int{}
			vertexMap[getStringFromIntVertex(p1)] = 0
			vertexMap[getStringFromIntVertex(p2)] = 1
			vertexMap[getStringFromIntVertex(p3)] = 2
			vertexMap[getStringFromIntVertex(p4)] = 3
			vertexMap[getStringFromIntVertex(p5)] = 4
			vertexMap[getStringFromIntVertex(p6)] = 5
			vertexMap[getStringFromIntVertex(p7)] = 6
			vertexMap[getStringFromIntVertex(p8)] = 7
			triangulateVertices(vertices, vertexMap)

			p1 = [3]int{0, 0, 0}
			p2 = [3]int{0, 6 * rand.IntN(100), 0}
			p3 = [3]int{4 * rand.IntN(100), 6 * rand.IntN(100), 0}
			p4 = [3]int{6 * rand.IntN(100), 6 * rand.IntN(100), 0}
			p5 = [3]int{6 * rand.IntN(100), 3 * rand.IntN(100), 0}
			p6 = [3]int{6 * rand.IntN(100), 0, 0}
			p7 = [3]int{4 * rand.IntN(100), 0, 0}
			p8 = [3]int{1 * rand.IntN(100), 0, 0}
			vertices = [][3]int{p1, p2, p3, p4, p5, p6, p7, p8}
			vertexMap = map[string]int{}
			vertexMap[getStringFromIntVertex(p1)] = 0
			vertexMap[getStringFromIntVertex(p2)] = 1
			vertexMap[getStringFromIntVertex(p3)] = 2
			vertexMap[getStringFromIntVertex(p4)] = 3
			vertexMap[getStringFromIntVertex(p5)] = 4
			vertexMap[getStringFromIntVertex(p6)] = 5
			vertexMap[getStringFromIntVertex(p7)] = 6
			vertexMap[getStringFromIntVertex(p8)] = 7
			triangulateVertices(vertices, vertexMap)
		}
		duration = time.Since(start)

		fmt.Println("Execution Took: ")
		fmt.Print(duration)
		fmt.Println()
	}
}

func getDistance(pointA [3]int, pointB[3]int) float32 {
	return float32(math.Sqrt(math.Pow(float64(pointB[0]-pointA[0]), 2) + math.Pow(float64(pointB[1]-pointA[1]), 2)))
}

func getMidPoint(pointA [3]int, pointB [3]int) [3]float32 {
	//Similar to the rest of the helper functions this function ignores the z axis. The z is only inputed and outputed to retain information about the point
	//Only calculates the midpoint in terms of the z access
	//Also assumes pointA has the same z value as pointB

	//alpha is the distance between pointA and pointB
	var alpha float64 = math.Sqrt(math.Pow(float64(pointB[0]-pointA[0]), 2) + math.Pow(float64(pointB[1]-pointA[1]), 2))
	//delta is the angle created from a triangle that is made from pointA, pointB, and [3]int{pointB[0], pointA[1]}
	var delta float64 = math.Atan2(float64(pointB[1]-pointA[1]), float64(pointB[0]-pointA[0]))

	var pX float32
	var pY float32

	if delta > 3.14159 || delta == 0 {
		pX = float32(pointA[0]+pointB[0]) / 2
		pY = float32(pointA[1])
	} else {
		pX = float32(math.Cos(delta)*(alpha/2)) + float32(pointA[0])
		pY = float32(math.Sin(delta)*(alpha/2)) + float32(pointA[1])
	}

	return [3]float32{pX, pY, float32(pointA[2])}
}

func GetMidPointTester() {
	fmt.Println("\nTesting getMidPoint()")
	p1 := [3]int{1, 1, 0}
	q1 := [3]int{10, 1, 0}

	fmt.Printf("\nTest Case 1, Input: [%d, %d], [%d, %d] Got:", p1[0], p1[1], q1[0], q1[1])
	result := getMidPoint(p1, q1)
	fmt.Printf("[%f, %f]", result[0], result[1])

	p1 = [3]int{10, 1, 0}
	q1 = [3]int{1, 1, 0}

	fmt.Printf("\nTest Case 2, Input: [%d, %d], [%d, %d] Got:", p1[0], p1[1], q1[0], q1[1])
	result = getMidPoint(p1, q1)
	fmt.Printf("[%f, %f]", result[0], result[1])

	p1 = [3]int{1, 10, 0}
	q1 = [3]int{1, 1, 0}

	fmt.Printf("\nTest Case 3, Input: [%d, %d], [%d, %d] Got:", p1[0], p1[1], q1[0], q1[1])
	result = getMidPoint(p1, q1)
	fmt.Printf("[%f, %f]", result[0], result[1])

	p1 = [3]int{1, 1, 0}
	q1 = [3]int{1, 10, 0}

	fmt.Printf("\nTest Case 4, Input: [%d, %d], [%d, %d] Got:", p1[0], p1[1], q1[0], q1[1])
	result = getMidPoint(p1, q1)
	fmt.Printf("[%f, %f]", result[0], result[1])

	fmt.Println()
}

func edgeOffsetter(edge [2][3]int) [2][3]float32 {
	//For now offsets by 0.01
	var offsetDistance float64 = 0.01

	//alpha here is the hypotnuse created from the triangle that is made from edge[0], edge[1], and [3]int{edge[1][0], edge[0][1]}
	var alpha float64 = math.Sqrt(math.Pow(float64(edge[1][0]-edge[0][0]), 2.0) + math.Pow(float64(edge[1][1]-edge[0][1]), 2))
	//delta here is the angle created from the triangle that is made from edge[0], edge[1], and [3]int{edge[1][0], edge[0][1]}
	var delta float64 = math.Atan2(float64(edge[1][1]-edge[0][1]), float64(edge[1][0]-edge[0][0]))

	var yOffsetted bool = false
	if delta == 0 {
		yOffsetted = true
		alpha = math.Sqrt(math.Pow(float64(edge[1][0]-edge[0][0]), 2.0) + math.Pow(float64(edge[1][1]-edge[0][1])+0.0001, 2))
		delta = math.Atan2(float64(edge[1][1]-edge[0][1])+0.0001, float64(edge[1][0]-edge[0][0]))
	}

	//p1 and p2 are the new edge[0] and edge[1] respectively
	var p2X float64 = math.Cos(delta)*(alpha-offsetDistance) + float64(edge[0][0])
	var p2Y float64 = math.Sin(delta)*(alpha-offsetDistance) + float64(edge[0][1])

	if yOffsetted {
		p2Y -= 0.0001
	}

	var p2 [3]float32 = [3]float32{float32(p2X), float32(p2Y), float32(edge[1][2])}

	delta = math.Atan2(float64(edge[1][0]-edge[0][0]), float64(edge[1][1]-edge[0][1]))

	var p1X float64 = float64(edge[1][0]) - (math.Sin(delta) * (alpha - offsetDistance))
	var p1Y float64 = float64(edge[1][1]) - (math.Cos(delta) * (alpha - offsetDistance))

	var p1 [3]float32 = [3]float32{float32(p1X), float32(p1Y), float32(edge[0][2])}

	return [2][3]float32{p1, p2}
}

func EdgeOffsetterTester() {
	fmt.Println("\nTesting edgeOffsetter()")
	p1 := [3]int{1, 1, 0}
	q1 := [3]int{10, 1, 0}

	fmt.Printf("\nTest Case 1, Input: [[%d, %d],[%d, %d] Got:", p1[0], p1[1], q1[0], q1[1])
	result := edgeOffsetter([2][3]int{p1, q1})
	fmt.Printf("[%f, %f]", result[0], result[1])

	p1 = [3]int{0, 0, 0}
	q1 = [3]int{1, 1, 0}

	fmt.Printf("\nTest Case 2, Input: [[%d, %d],[%d, %d] Got:", p1[0], p1[1], q1[0], q1[1])
	result = edgeOffsetter([2][3]int{p1, q1})
	fmt.Printf("[%f, %f]", result[0], result[1])

	p1 = [3]int{1, 1, 0}
	q1 = [3]int{1, 2, 0}

	fmt.Printf("\nTest Case 3, Input: [[%d, %d],[%d, %d] Got:", p1[0], p1[1], q1[0], q1[1])
	result = edgeOffsetter([2][3]int{p1, q1})
	fmt.Printf("[%f, %f]", result[0], result[1])

	fmt.Println()
}

func determineIfVertexInTriangle(p0 [3]int, p1 [3]int, p2 [3]int, p [3]int) bool {
	p0x := float32(p0[0])
	p0y := float32(p0[1])
	p1x := float32(p1[0])
	p1y := float32(p1[1])
	p2x := float32(p2[0])
	p2y := float32(p2[1])
	px := float32(p[0])
	py := float32(p[1])

	Area := 0.5 *(-p1y*p2x + p0y*(-p1x + p2x) + p0x*(p1y - p2y) + p1x*p2y);
	s := 1/(2*Area)*(p0y*p2x - p0x*p2y + (p2y - p0y)*px + (p0x - p2x)*py);
	t := 1/(2*Area)*(p0x*p1y - p0y*p1x + (p0y - p1y)*px + (p1x - p0x)*py);

	return (s >= 0 && t >= 0) && (1-s-t) >= 0
}

func determineCollision(edge [2][3]float32, otherEdge [2][3]float32) bool {
	//a = edge[0], b = edge[1], c = otherEdge[0], d = otherEdge[1]
	var abcOrientation int = determineOrientation(edge[0], edge[1], otherEdge[0])
	var abdOrientation int = determineOrientation(edge[0], edge[1], otherEdge[1])
	var cdaOrientation int = determineOrientation(otherEdge[0], otherEdge[1], edge[0])
	var cdbOrientation int = determineOrientation(otherEdge[0], otherEdge[1], edge[1])

	if abcOrientation != abdOrientation && cdaOrientation != cdbOrientation {
		return true
	}
	if abcOrientation == 0 && determineEdgePointIntersection(edge, otherEdge[0]) {
		return true
	}
	if abdOrientation == 0 && determineEdgePointIntersection(edge, otherEdge[1]) {
		return true
	}
	if cdaOrientation == 0 && determineEdgePointIntersection(otherEdge, edge[0]) {
		return true
	}
	if cdbOrientation == 0 && determineEdgePointIntersection(otherEdge, edge[1]) {
		return true
	}
	return false
}

func determineOrientation(pointA [3]float32, pointB [3]float32, pointC [3]float32) int {
	//For usage with 3d points however this only check orientation of x and y. Ignores z
	var orientationValue float32 = ((pointB[1] - pointA[1]) * (pointC[0] - pointB[0])) - ((pointB[0] - pointA[0]) * (pointC[1] - pointB[1]))

	if orientationValue == 0 {
		return 0
	}
	if orientationValue > 0 {
		return 1
	}
	return 2
}

func determineEdgePointIntersection(edge [2][3]float32, point [3]float32) bool {
	//For usage only where the edge points are collinear with the point, Also assumes you are checking on a xy plane only. Ignores z. Includes [3]int for specific use case in triangulateVertices()
	if (point[0] <= max(edge[0][0], edge[1][0]) && point[0] >= min(edge[0][0], edge[1][0])) && (point[1] <= max(edge[0][1], edge[1][1]) && point[1] >= min(edge[0][1], edge[1][1])) {
		return true
	}
	return false
}

func DetermineCollisionTester() {
	fmt.Println("Testing determineCollision()")

	p1 := [3]float32{1, 1, 0}
	q1 := [3]float32{10, 1, 0}
	p2 := [3]float32{1, 2, 0}
	q2 := [3]float32{10, 2, 0}

	fmt.Printf("\nTest Case 1, Expected: False, Got: %t\n", determineCollision([2][3]float32{p1, q1}, [2][3]float32{p2, q2}))

	p1 = [3]float32{10, 0, 0}
	q1 = [3]float32{0, 10, 0}
	p2 = [3]float32{0, 0, 0}
	q2 = [3]float32{10, 10, 0}

	fmt.Printf("\nTest Case 2, Expected: True, Got: %t\n", determineCollision([2][3]float32{p1, q1}, [2][3]float32{p2, q2}))

	p1 = [3]float32{-5, -5, 0}
	q1 = [3]float32{0, 0, 0}
	p2 = [3]float32{1, 1, 0}
	q2 = [3]float32{10, 10, 0}

	fmt.Printf("\nTest Case 3, Expected: False, Got: %t\n", determineCollision([2][3]float32{p1, q1}, [2][3]float32{p2, q2}))

	p1 = [3]float32{0, 0, 0}
	q1 = [3]float32{1, 1, 0}
	p2 = [3]float32{2, 2, 0}
	q2 = [3]float32{1.00001, 1.00001, 0}

	fmt.Printf("\nTest Case 3, Expected: False, Got: %t\n", determineCollision([2][3]float32{p1, q1}, [2][3]float32{p2, q2}))

	p1 = [3]float32{0, 0, 0}
	q1 = [3]float32{0, 1, 0}
	p2 = [3]float32{0, 1, 0}
	q2 = [3]float32{1, 1, 0}

	fmt.Printf("\nTest Case 3, Expected: True, Got: %t\n", determineCollision([2][3]float32{p1, q1}, [2][3]float32{p2, q2}))
}

func GetVerticesFromFaces(faces []greedyMesher.Face) ([][3]int, [][][]bool, map[string]int, map[string][][3]int) {
	vertSet := hashset.New()

	faceMap := map[string][][3]int{}

	var maxX int = 0
	var maxY int = 0
	var maxZ int = 0

	for _, face := range faces {
		for _, coord := range face.VoxelCoords {
			maxX = max(maxX, coord[0])
			maxY = max(maxY, coord[1])
			maxZ = max(maxZ, coord[2])
		}

		switch cornerAmount := len(face.VoxelCoords); cornerAmount {
		case 1:
			switch faceOrientation := face.FaceIndex; faceOrientation {
			case 0:
				newVerts := [][3]int{}
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				faceMap[getStringFromFace(face)] = newVerts
			case 1:
				newVerts := [][3]int{}
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				faceMap[getStringFromFace(face)] = newVerts
			case 2:
				newVerts := [][3]int{}
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				faceMap[getStringFromFace(face)] = newVerts
			case 3:
				newVerts := [][3]int{}
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				faceMap[getStringFromFace(face)] = newVerts
			case 4:
				newVerts := [][3]int{}
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				faceMap[getStringFromFace(face)] = newVerts
			case 5:
				newVerts := [][3]int{}
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				faceMap[getStringFromFace(face)] = newVerts
			default:
				fmt.Printf("\nError in faceOrientation in GetVerticesFromFaces(), Expected 0-5, Got: %d", faceOrientation)
			}
		case 2:
			switch faceOrientation := face.FaceIndex; faceOrientation {
			case 0:
				newVerts := [][3]int{}
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				faceMap[getStringFromFace(face)] = newVerts
			case 1:
				newVerts := [][3]int{}
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				faceMap[getStringFromFace(face)] = newVerts
			case 2:
				newVerts := [][3]int{}
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				faceMap[getStringFromFace(face)] = newVerts
			case 3:
				newVerts := [][3]int{}
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				faceMap[getStringFromFace(face)] = newVerts
			case 4:
				newVerts := [][3]int{}
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				faceMap[getStringFromFace(face)] = newVerts
			case 5:
				newVerts := [][3]int{}
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVerts = append(newVerts, [3]int{int(newVertex[0]*2), int(newVertex[1]*2), int(newVertex[2]*2)})
				faceMap[getStringFromFace(face)] = newVerts
			default:
				fmt.Printf("\nError in faceOrientation in GetVerticesFromFaces(), Expected 0-5, Got: %d", faceOrientation)
			}
		default:
			fmt.Printf("\nError in corner amount in GetVerticesFromFaces(), Expected 1 or 2, Got: %d\n", cornerAmount)
		}
	}
	vertMatrix := make([][][]bool, (maxX+1)*2)
	for x := 0; x < (maxX+1)*2; x++ {
		vertMatrix[x] = make([][]bool, (maxY+1)*2)
		for y := 0; y < (maxY+1)*2; y++ {
			vertMatrix[x][y] = make([]bool, (maxZ+1)*2)
		}
	}

	vertMap := make(map[string]int)

	vertArray := make([][3]int, vertSet.Size())
	for index, value := range vertSet.Values() {
		stringVertex := fmt.Sprint(value)
		vertex := getVertexFromString(stringVertex)
		intVertex := [3]int{int(vertex[0] * 2), int(vertex[1] * 2), int(vertex[2] * 2)}
		vertArray[index] = intVertex
		vertMatrix[intVertex[0]][intVertex[1]][intVertex[2]] = true
		vertMap[getStringFromIntVertex(intVertex)] = index
	}
	return vertArray, vertMatrix, vertMap, faceMap
}

func getStringFromFace(face greedyMesher.Face) string {
	switch cornerCount := len(face.VoxelCoords); cornerCount {
	case 1:
		return fmt.Sprintf("%d:%d,%d,%d", face.FaceIndex, face.VoxelCoords[0][0], face.VoxelCoords[0][1], face.VoxelCoords[0][2])
	case 2:
		return fmt.Sprintf("%d:%d,%d,%d,%d,%d,%d", face.FaceIndex, face.VoxelCoords[0][0], face.VoxelCoords[0][1], face.VoxelCoords[0][2], face.VoxelCoords[1][0], face.VoxelCoords[1][1], face.VoxelCoords[1][2])
	default: 
		fmt.Printf("\nError in getStringFromFace. Expected len(face.VoxelCoords) to be 1 or 2, Got: %d\n", cornerCount)
		return "404"
	}
}

func getStringFromVertex(vertex [3]float32) string {
	return fmt.Sprintf("%f,%f,%f", vertex[0], vertex[1], vertex[2])
}

func getStringFromIntVertex(vertex [3]int) string {
	return fmt.Sprintf("%d,%d,%d", vertex[0], vertex[1], vertex[2])
}

func getVertexFromString(vertex string) [3]float32 {
	splitStrings := strings.Split(vertex, ",")
	x, err := strconv.ParseFloat(splitStrings[0], 32)
	if err != nil {
		fmt.Printf("\nError on getVertexFromString attempted to parse %s into float.\n", splitStrings[0])
		return [3]float32{}
	}
	y, err := strconv.ParseFloat(splitStrings[1], 32)
	if err != nil {
		fmt.Printf("\nError on getVertexFromString attempted to parse %s into float.\n", splitStrings[1])
		return [3]float32{}
	}
	z, err := strconv.ParseFloat(splitStrings[2], 32)
	if err != nil {
		fmt.Printf("\nError on getVertexFromString attempted to parse %s into float.\n", splitStrings[2])
		return [3]float32{}
	}
	return [3]float32{float32(x), float32(y), float32(z)}
}

func getIntVertexFromString(vertex string) [3]int {
	splitStrings := strings.Split(vertex, ",")
	x, err := strconv.ParseFloat(splitStrings[0], 32)
	if err != nil {
		fmt.Printf("\nError on getVertexFromString attempted to parse %s into float.\n", splitStrings[0])
		return [3]int{}
	}
	y, err := strconv.ParseFloat(splitStrings[1], 32)
	if err != nil {
		fmt.Printf("\nError on getVertexFromString attempted to parse %s into float.\n", splitStrings[1])
		return [3]int{}
	}
	z, err := strconv.ParseFloat(splitStrings[2], 32)
	if err != nil {
		fmt.Printf("\nError on getVertexFromString attempted to parse %s into float.\n", splitStrings[2])
		return [3]int{}
	}
	return [3]int{int(x), int(y), int(z)}
}

func PointsToJson(points [][3]int) {
	toJSON := []Point{}
	for _, point := range points {
		toJSON = append(toJSON, Point{X: float32(point[0]) / 2, Y: float32(point[1]) / 2, Z: float32(point[2]) / 2, Value: 1})
	}

	file, errs := os.Create("FacePointData.JSON")
	if errs != nil {
		panic("Failed to write to file:" + errs.Error())
	}

	enc := json.NewEncoder(file)
	enc.Encode(PointData{Points: toJSON, VoxelSize: 0.1})
}

func VoxelsToJson(voxels [][][]uint8) {
	file, errs := os.Create("VoxelMatrix.JSON")
	if errs != nil {
		panic("Failed to write to file:" + errs.Error())
	}

	const voxelValueThreshold = 0

	enc := json.NewEncoder(file)
	cleanedPoints := []Point{}
	for xIndex, xArray := range voxels {
		for yIndex, yArray := range xArray {
			for zIndex, value := range yArray {
				if value > voxelValueThreshold {
					point := Point{X: float32(xIndex), Y: float32(yIndex), Z: float32(zIndex), Value: float32(value)}
					cleanedPoints = append(cleanedPoints, point)
				}
			}
		}
	}

	fmt.Println("Clean Points Got: " + strconv.Itoa(len(cleanedPoints)))

	enc.Encode(PointData{Points: cleanedPoints, VoxelSize: 1})
}

func ToOBJFile(filename string, vertices, faces [][3]int) {
	file, errs := os.Create(filename + ".obj")
	if errs != nil {
		panic("Failed to write to file:" + errs.Error())
	}
	defer file.Close()

	w := io.Writer(file)

	for _, vertex := range vertices {
		io.WriteString(w, "v" + " " + fmt.Sprint(vertex[0]) + " " + fmt.Sprint(vertex[1]) + " " + fmt.Sprint(vertex[2]) + "\n")
	}
	io.WriteString(w, "\n")
	for _, face := range faces {
		io.WriteString(w, "f" + " " + fmt.Sprint(face[0]+1) + " " + fmt.Sprint(face[1]+1) + " " + fmt.Sprint(face[2]+1) + "\n")
	}
}
