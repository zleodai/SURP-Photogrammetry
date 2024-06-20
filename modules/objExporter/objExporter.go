package objExporter

import (
	"encoding/json"
	"fmt"
	"modules/greedyMesher"
	"os"
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

func GetMeshFacesFromVertices(faces []greedyMesher.Face, vertices [][3]int, vertexMatrix [][][]bool, vertexMap map[string]int) [][]int {
	meshFaces := make([][]int, 0)

	for _, face := range faces {
		switch cornerAmount := len(face.VoxelCoords); cornerAmount {
		case 1:
			meshFaces = append(meshFaces, triangulateVertices(face.VoxelCoords, vertexMap)...)
		case 2:
			
		default:
			fmt.Printf("\nError in face cornerAmount in GetMeshFacesFromVertices. Expected: 1 or 2 Got: %d\n", cornerAmount)
		}
	}

	return meshFaces
}

func triangulateVertices(vertices [][3]int, vertexMap map[string]int) [][]int{
	triangulatedFaces := make([][]int, 0)



	return triangulatedFaces
}

func determineCollision(edge [2][3]int, existingEdges [][2][3]int) bool {
	var collisionFound bool = false
	for _, existingEdge := range existingEdges {
		
	}
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
		intVertex := [3]int{int(vertex[0]*2), int(vertex[1]*2), int(vertex[2]*2)}
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
		toJSON = append(toJSON, Point{X: float32(point[0])/2, Y: float32(point[1])/2, Z: float32(point[2])/2, Value: 1})
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
