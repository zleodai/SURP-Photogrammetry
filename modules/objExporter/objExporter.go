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

func Test() byte {
	return 0
}

type PointData struct {
	Points []Point
	VoxelSize float32
}

type Point struct {
	X float32
	Y float32
	Z float32
	Value float32
}

func GetVerticesFromFaces(faces []greedyMesher.Face) [][3]float32 {
	vertSet := hashset.New()

	for _, face := range faces {
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
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 1:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) - 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 2:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) + 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 3:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[1][1]) + 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 4:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) + 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			case 5:
				newVertex := [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[0][2]) - 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[0][0]) - 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
				newVertex = [3]float32{float32(face.VoxelCoords[1][0]) + 0.5, float32(face.VoxelCoords[0][1]) - 0.5, float32(face.VoxelCoords[1][2]) + 0.5}
				vertSet.Add(getStringFromVertex(newVertex))
			default:
				fmt.Printf("\nError in faceOrientation in GetVerticesFromFaces(), Expected 0-5, Got: %d", faceOrientation)
			}
		default:
			fmt.Printf("\nError in corner amount in GetVerticesFromFaces(), Expected 1 or 2, Got: %d\n", cornerAmount)
		}
	}
	vertArray := make([][3]float32, vertSet.Size())
	for index, value := range vertSet.Values() {
		stringVertex := fmt.Sprint(value)
		vertArray[index] = getVertexFromString(stringVertex)
	}
	return vertArray
}

func getStringFromVertex(vertex [3]float32) string {
	return fmt.Sprintf("%f,%f,%f", vertex[0], vertex[1], vertex[2])
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

func PointsToJson(points [][3]float32) {
	toJSON := []Point{}
	for _, point := range points {
		toJSON = append(toJSON, Point{X: point[0], Y: point[1], Z: point[2], Value: 1})
	}

	file, errs := os.Create("NewPointCloud.JSON")
	if errs != nil {
		panic("Failed to write to file:" + errs.Error())
	}

	enc := json.NewEncoder(file)
	enc.Encode(PointData{Points: toJSON, VoxelSize: 1})
}
