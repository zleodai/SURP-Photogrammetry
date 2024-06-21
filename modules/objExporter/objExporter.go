package objExporter

import (
	"encoding/json"
	"fmt"
	"math"
	"modules/greedyMesher"
	"os"
	"reflect"
	"strconv"
	"strings"

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

func GetMeshFacesFromVertices(faces []greedyMesher.Face, vertices [][3]int, vertexMatrix [][][]bool, vertexMap map[string]int) [][3]int {
	meshFaces := make([][3]int, 0)

	for _, face := range faces {
		switch cornerAmount := len(face.VoxelCoords); cornerAmount {
		case 1:
			var customEdgeColliders [][2][3]float32 = [][2][3]float32{}
			meshFaces = append(meshFaces, triangulateVertices(face.VoxelCoords, vertexMap, customEdgeColliders)...)
		case 2:

		default:
			fmt.Printf("\nError in face cornerAmount in GetMeshFacesFromVertices. Expected: 1 or 2 Got: %d\n", cornerAmount)
		}
	}

	return meshFaces
}

func triangulateVertices(vertices [][3]int, vertexMap map[string]int, customEdgeColliders [][2][3]float32) [][3]int {
	triangulatedFaces := make([][3]int, 0)
	createdEdges := make([][2][3]int, 0)
	createdEdgeColliders := make([][2][3]float32, 0)
	//For the rest of the function we assume createdEdges and createdEdgeColliders will contain same elements with the same index locations. Only difference is createdEdges contains the int versions of the edge

	for _, edgeCollider := range customEdgeColliders {
		createdEdgeColliders = append(createdEdgeColliders, edgeCollider)
		var intEdge [2][3]int = [2][3]int{[3]int{int(edgeCollider[0][0] + 0.1), int(edgeCollider[0][1] + 0.1), int(edgeCollider[0][2] + 0.1)}, [3]int{int(edgeCollider[1][0] + 0.1), int(edgeCollider[1][1] + 0.1), int(edgeCollider[1][2] + 0.1)}}
		createdEdges = append(createdEdges, intEdge)
	}

	//We first iterate through all the vertices to create all the triangles on the sides
	for index, vertex := range vertices {
		//targetVertex is the vertex that the vertex will attempt to create a triangle with
		//for now we assume that each vertex automatically will attempt to form a triangle with its proceeding vertex
		var targetVertex [3]int
		var targetVertexIndex int
		if index < len(vertices)-1 {
			targetVertex = vertices[index+1]
			targetVertexIndex = index+1
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
			//If an edge was not created we know that a triangle with the targetVertex has not been created
			var midPoint = getMidPoint(vertex, targetVertex)

			var midPointOnEdge bool = false
			var midPointEdgeIndex int = 0

			//First check if middlepoint is on an existing edge
			for intEdgeIndex, intEdge := range createdEdges {
				var float32Edge [2][3]float32 = [2][3]float32{{float32(intEdge[0][0]), float32(intEdge[0][1]), float32(intEdge[0][2])}, {float32(intEdge[1][0]), float32(intEdge[1][1]), float32(intEdge[1][2])}}

				if determineOrientation(float32Edge[0], float32Edge[1], midPoint) == 0 && determineEdgePointIntersection(float32Edge, midPoint) {
					midPointEdgeIndex = intEdgeIndex
				}
			}
			
			//Iterate through otherVertexs to see which ones can make a triangle with targetVertex
			var triangleCreated bool = false
			for otherIndex, otherVertex := range vertices {
				if !triangleCreated && (otherIndex != index && otherIndex != targetVertexIndex) {
					var collisionDetected bool = false
					var newEdge [2][3]float32 = [2][3]float32{midPoint, {float32(otherVertex[0]), float32(otherVertex[1]), float32(otherVertex[2])}}

					//Iterate through all edgeColliders (except for any edgeCollider midPoint is already on) to detect if creating a new edge will cause any collisions
					for edgeIndex, edge := range createdEdgeColliders {
						if !triangleCreated && !(midPointOnEdge && edgeIndex == midPointEdgeIndex) {
							if determineCollision(newEdge, edge) {
								collisionDetected = true
							}
						}
					}

					//Check and see if the newEdge will collide with vertex or targetVertex
					var float32Vertex [3]float32 = [3]float32{float32(vertex[0]), float32(vertex[1]), float32(vertex[2])}
					if determineOrientation(newEdge[0], newEdge[1], float32Vertex) == 0 && determineEdgePointIntersection(newEdge, float32Vertex) {
						collisionDetected = true
					}
					var float32TargetVertex [3]float32 = [3]float32{float32(targetVertex[0]), float32(targetVertex[1]), float32(targetVertex[2])}
					if determineOrientation(newEdge[0], newEdge[1], float32TargetVertex) == 0 && determineEdgePointIntersection(newEdge, float32TargetVertex) {
						collisionDetected = true
					}

					//When no collision is detected we make the new triangle
					if !collisionDetected{
						var newFace [3]int = [3]int{vertexMap[getStringFromIntVertex(vertex)], vertexMap[getStringFromIntVertex(targetVertex)], vertexMap[getStringFromIntVertex(otherVertex)]}

						//First try to detect if the triangle we are making is has already been created
						var triangleAlreadyCreated bool = false

						for _, face := range triangulatedFaces {
							if reflect.DeepEqual(newFace, face) {
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

		for !allPossibleTrianglesMade {
			var triangleMade bool = false
			var oldTargetVertexIndex int = targetVertexIndex

			//Check all vertices to see if a new targetVertex can be found
			for otherIndex, otherVertex := range vertices {
				if otherIndex != index && otherIndex != oldTargetVertexIndex{
					var newEdge [2][3]float32 = [2][3]float32 {{float32(vertex[0]), float32(vertex[1]), float32(vertex[1])}, {float32(otherVertex[0]), float32(otherVertex[1]), float32(otherVertex[1])}}
					
					//Check to see if newEdge created collides with any existing colliders
					var collisionFound bool = false
					for _, edge := range createdEdgeColliders {
						if determineCollision(newEdge, edge) {
							collisionFound = true

							if determineOrientation(newEdge[0], newEdge[1], edge[0]) == 0 && determineOrientation(newEdge[0], newEdge[1], edge[1]) == 0 {
								collisionFound = false
							}
						}
					}

					if !collisionFound {
						targetVertex = otherVertex
						targetVertexIndex = otherIndex
					}
				}
			}

			//Create new triangle if targetVertexIndex is different from oldTargetVertexIndex
			if oldTargetVertexIndex != targetVertexIndex {
				//Repeat code from original triangle creation on sides
				var midPoint = getMidPoint(vertex, targetVertex)

				var midPointOnEdge bool = false
				var midPointEdgeIndex int = 0

				//First check if middlepoint is on an existing edge
				for intEdgeIndex, intEdge := range createdEdges {
					var float32Edge [2][3]float32 = [2][3]float32{{float32(intEdge[0][0]), float32(intEdge[0][1]), float32(intEdge[0][2])}, {float32(intEdge[1][0]), float32(intEdge[1][1]), float32(intEdge[1][2])}}

					if determineOrientation(float32Edge[0], float32Edge[1], midPoint) == 0 && determineEdgePointIntersection(float32Edge, midPoint) {
						midPointEdgeIndex = intEdgeIndex
					}
				}
				
				//Iterate through otherVertexs to see which ones can make a triangle with targetVertex
				var triangleCreated bool = false
				for otherIndex, otherVertex := range vertices {
					if !triangleCreated && (otherIndex != index && otherIndex != targetVertexIndex) {
						var collisionDetected bool = false
						var newEdge [2][3]float32 = [2][3]float32{midPoint, {float32(otherVertex[0]), float32(otherVertex[1]), float32(otherVertex[2])}}

						//Iterate through all edgeColliders (except for any edgeCollider midPoint is already on) to detect if creating a new edge will cause any collisions
						for edgeIndex, edge := range createdEdgeColliders {
							if !triangleCreated && !(midPointOnEdge && edgeIndex == midPointEdgeIndex) {
								if determineCollision(newEdge, edge) {
									collisionDetected = true
								}
							}
						}

						//Check and see if the newEdge will collide with vertex or targetVertex
						var float32Vertex [3]float32 = [3]float32{float32(vertex[0]), float32(vertex[1]), float32(vertex[2])}
						if determineOrientation(newEdge[0], newEdge[1], float32Vertex) == 0 && determineEdgePointIntersection(newEdge, float32Vertex) {
							collisionDetected = true
						}
						var float32TargetVertex [3]float32 = [3]float32{float32(targetVertex[0]), float32(targetVertex[1]), float32(targetVertex[2])}
						if determineOrientation(newEdge[0], newEdge[1], float32TargetVertex) == 0 && determineEdgePointIntersection(newEdge, float32TargetVertex) {
							collisionDetected = true
						}

						//When no collision is detected we make the new triangle
						if !collisionDetected{
							var newFace [3]int = [3]int{vertexMap[getStringFromIntVertex(vertex)], vertexMap[getStringFromIntVertex(targetVertex)], vertexMap[getStringFromIntVertex(otherVertex)]}

							//First try to detect if the triangle we are making is has already been created
							var triangleAlreadyCreated bool = false

							for _, face := range triangulatedFaces {
								if reflect.DeepEqual(newFace, face) {
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
			
			//Last check to see if no new triangles were made
			if !triangleMade {
				allPossibleTrianglesMade = true
			}
		}
	}

	return triangulatedFaces
}

func TriangulateVerticesTester() {
	fmt.Println("\nTesting triangulateVertices()")
	p1 := [3]int{1, 1, 0}
	p2 := [3]int{2, 1, 0}
	p3 := [3]int{2, 2, 0}
	p4 := [3]int{1, 2, 0}
	vertices := [][3]int{p1, p2, p3, p4}
	customEdgeColliders := [][2][3]float32{}
	vertexMap := map[string]int{}
	vertexMap[getStringFromIntVertex(p1)] = 0
	vertexMap[getStringFromIntVertex(p2)] = 1
	vertexMap[getStringFromIntVertex(p3)] = 2
	vertexMap[getStringFromIntVertex(p4)] = 3

	fmt.Print("\nTest Case 1\n	Input: ")
	fmt.Print(vertices)
	fmt.Print("\n	Result:\n")
	result := triangulateVertices(vertices, vertexMap, customEdgeColliders)
	fmt.Print(result)
	fmt.Println()

	p1 = [3]int{1, 1, 0}
	p2 = [3]int{1, 2, 0}
	p3 = [3]int{2, 2, 0}
	p4 = [3]int{2, 1, 0}
	vertices = [][3]int{p1, p2, p3, p4}
	customEdgeColliders = [][2][3]float32{}
	vertexMap = map[string]int{}
	vertexMap[getStringFromIntVertex(p1)] = 0
	vertexMap[getStringFromIntVertex(p2)] = 1
	vertexMap[getStringFromIntVertex(p3)] = 2
	vertexMap[getStringFromIntVertex(p4)] = 3

	fmt.Print("\nTest Case 2\n	Input: ")
	fmt.Print(vertices)
	fmt.Print("\n	Result:\n")
	result = triangulateVertices(vertices, vertexMap, customEdgeColliders)
	fmt.Print(result)
	fmt.Println()

	p1 = [3]int{1, 1, 0}
	p2 = [3]int{1, 2, 0}
	p3 = [3]int{1, 5, 0}
	p4 = [3]int{5, 5, 0}
	p5 := [3]int{5, 4, 0}
	p6 := [3]int{5, 1, 0}
	vertices = [][3]int{p1, p2, p3, p4, p5, p6}
	customEdgeColliders = [][2][3]float32{}
	vertexMap = map[string]int{}
	vertexMap[getStringFromIntVertex(p1)] = 0
	vertexMap[getStringFromIntVertex(p2)] = 1
	vertexMap[getStringFromIntVertex(p3)] = 2
	vertexMap[getStringFromIntVertex(p4)] = 3
	vertexMap[getStringFromIntVertex(p5)] = 4
	vertexMap[getStringFromIntVertex(p6)] = 5

	fmt.Print("\nTest Case 3\n	Input: ")
	fmt.Print(vertices)
	fmt.Print("\n	Result:\n")
	result = triangulateVertices(vertices, vertexMap, customEdgeColliders)
	fmt.Print(result)
	fmt.Println()

	p1 = [3]int{1, 1, 0}
	p2 = [3]int{1, 4, 0}
	p3 = [3]int{2, 4, 0}
	p4 = [3]int{2, 5, 0}
	p5 = [3]int{1, 5, 0}
	p6 = [3]int{1, 8, 0}
	p7 := [3]int{8, 8, 0}
	p8 := [3]int{8, 1, 0}
	vertices = [][3]int{p1, p2, p3, p4, p5, p6, p7, p8}
	vertexMap = map[string]int{}
	customEdgeColliders = [][2][3]float32{}
	customEdgeColliders = append(customEdgeColliders, [2][3]float32{getMidPoint([3]int{0, 4, 0}, [3]int{0, 5, 0}), getMidPoint(p3, p4)})
	vertexMap[getStringFromIntVertex(p1)] = 0
	vertexMap[getStringFromIntVertex(p2)] = 1
	vertexMap[getStringFromIntVertex(p3)] = 2
	vertexMap[getStringFromIntVertex(p4)] = 3
	vertexMap[getStringFromIntVertex(p5)] = 4
	vertexMap[getStringFromIntVertex(p6)] = 5
	vertexMap[getStringFromIntVertex(p7)] = 6
	vertexMap[getStringFromIntVertex(p8)] = 7

	fmt.Print("\nTest Case 4\n	Input: ")
	fmt.Print(vertices)
	fmt.Print("\n	Result:\n")
	result = triangulateVertices(vertices, vertexMap, customEdgeColliders)
	fmt.Print(result)
	fmt.Println()
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

	var yOffsetted bool = false
	if delta == 0 {
		yOffsetted = true
		alpha = math.Sqrt(math.Pow(float64(pointB[0]-pointA[0]), 2) + math.Pow(float64(pointB[1]-pointA[1])+0.0001, 2))
		delta = math.Atan2(float64(pointB[1]-pointA[1])+0.0001, float64(pointB[0]-pointA[0]))
	}

	var pX float32 = float32(math.Cos(delta) * (alpha / 2)) + float32(pointA[0])
	var pY float32 = float32(math.Sin(delta) * (alpha / 2)) + float32(pointA[1])

	if yOffsetted {
		pY -= 0.0001/2
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

	p1 = [3]int{1, 1, 0}
	q1 = [3]int{1, 10, 0}

	fmt.Printf("\nTest Case 2, Input: [%d, %d], [%d, %d] Got:", p1[0], p1[1], q1[0], q1[1])
	result = getMidPoint(p1, q1)
	fmt.Printf("[%f, %f]", result[0], result[1])

	p1 = [3]int{1, 1, 0}
	q1 = [3]int{4, 4, 0}

	fmt.Printf("\nTest Case 3, Input: [%d, %d], [%d, %d] Got:", p1[0], p1[1], q1[0], q1[1])
	result = getMidPoint(p1, q1)
	fmt.Printf("[%f, %f]", result[0], result[1])

	p1 = [3]int{0, 0, 0}
	q1 = [3]int{0, 0, 0}

	fmt.Printf("\nTest Case 4, Input: [%d, %d], [%d, %d] Got:", p1[0], p1[1], q1[0], q1[1])
	result = getMidPoint(p1, q1)
	fmt.Printf("[%f, %f]", result[0], result[1])

	p1 = [3]int{1, 1, 0}
	q1 = [3]int{-4, -4, 0}

	fmt.Printf("\nTest Case 5, Input: [%d, %d], [%d, %d] Got:", p1[0], p1[1], q1[0], q1[1])
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
	var p2X float64 = math.Cos(delta) * (alpha - offsetDistance) + float64(edge[0][0])
	var p2Y float64 = math.Sin(delta) * (alpha - offsetDistance) + float64(edge[0][1])

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

func GetVerticesFromFaces(faces []greedyMesher.Face) ([][3]int, [][][]bool, map[string]int) {
	vertSet := hashset.New()

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
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 1:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 2:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 3:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 4:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 5:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			default:
				fmt.Printf("\nError in faceOrientation in GetVerticesFromFaces(), Expected 0-5, Got: %d", faceOrientation)
			}
		case 2:
			switch faceOrientation := face.FaceIndex; faceOrientation {
			case 0:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 1:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 2:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 3:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 4:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 5:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
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
	return vertArray, vertMatrix, vertMap
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
