package voxelMesher

import (
	"encoding/json"
	"fmt"
	"math"
	"modules/pointCloudDecoder"
	"os"
	"runtime"
	"strconv"
)

type pointVal struct {
	X     float64
	Y     float64
	Z     float64
	Value float64
}

type pointList struct {
	Points []pointVal
}

func Test() byte {
	return 0
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Mesh(xArray, yArray, zArray []pointCloudDecoder.Point, voxelSize float64) [][][]int {
	var xSize int = int(math.Floor(math.Abs(xArray[0].X-xArray[len(xArray)-1].X) / voxelSize))
	var ySize int = int(math.Floor(math.Abs(yArray[0].Y-yArray[len(yArray)-1].Y) / voxelSize))
	var zSize int = int(math.Floor(math.Abs(zArray[0].Z-zArray[len(zArray)-1].Z) / voxelSize))
	fmt.Println("\nVoxel Grid X, Y, Z:", xSize, ySize, zSize)
	fmt.Println("Total Voxel Amount:", xSize*ySize*zSize)

	voxels := make([][][]int, xSize+1)
	for i := 0; i < len(voxels); i++ {
		voxels[i] = make([][]int, ySize+1)
		for j := 0; j < len(voxels[i]); j++ {
			voxels[i][j] = make([]int, zSize+1)
		}
	}

	for _, point := range xArray {
		xIndex := int(math.Floor(math.Abs(xArray[0].X-point.X) / voxelSize))
		yIndex := int(math.Floor(math.Abs(yArray[0].Y-point.Y) / voxelSize))
		zIndex := int(math.Floor(math.Abs(zArray[0].Z-point.Z) / voxelSize))
		voxels[xIndex][yIndex][zIndex] += 1
	}

	//GenerateVoxelJson(voxels)

	runtime.GC()
	return voxels
}

func MinMaxMesh(xMinMax, yMinMax, zMinMax [2]float64, points []pointCloudDecoder.Point, voxelSize float64) [][][]int {
	var xSize int = int(math.Floor(math.Abs(xMinMax[0]-xMinMax[1]) / voxelSize))
	var ySize int = int(math.Floor(math.Abs(yMinMax[0]-yMinMax[1]) / voxelSize))
	var zSize int = int(math.Floor(math.Abs(zMinMax[0]-zMinMax[1]) / voxelSize))
	fmt.Println("Voxel Grid X, Y, Z:", xSize, ySize, zSize)
	fmt.Println("Total Voxel Amount:", xSize*ySize*zSize)

	voxels := make([][][]int, xSize+1)
	for i := 0; i < len(voxels); i++ {
		voxels[i] = make([][]int, ySize+1)
		for j := 0; j < len(voxels[i]); j++ {
			voxels[i][j] = make([]int, zSize+1)
		}
	}

	for _, point := range points {
		xIndex := int(math.Floor(math.Abs(xMinMax[0]-point.X) / voxelSize))
		yIndex := int(math.Floor(math.Abs(yMinMax[0]-point.Y) / voxelSize))
		zIndex := int(math.Floor(math.Abs(zMinMax[0]-point.Z) / voxelSize))
		voxels[xIndex][yIndex][zIndex] += 1
	}

	//GenerateVoxelJson(voxels)

	runtime.GC()
	return voxels
}

func IterativeMesh(xMinMax, yMinMax, zMinMax [2]float64, points []pointCloudDecoder.Point, voxelEndSize float64, Iterations, scaleFactor int) {
	var xSize int = int(math.Floor(math.Abs(xMinMax[0]-xMinMax[1]) / voxelEndSize))
	var ySize int = int(math.Floor(math.Abs(yMinMax[0]-yMinMax[1]) / voxelEndSize))
	var zSize int = int(math.Floor(math.Abs(zMinMax[0]-zMinMax[1]) / voxelEndSize))

	const minimumVoxelAmount = 8

	masterVoxels := make([][][]float64, xSize+1)
	for i := 0; i < len(masterVoxels); i++ {
		masterVoxels[i] = make([][]float64, ySize+1)
		for j := 0; j < len(masterVoxels[i]); j++ {
			masterVoxels[i][j] = make([]float64, zSize+1)
		}
	}

	fmt.Printf("\nX Max/Min: %s/%s, Y Max/Min: %s/%s, Z Max/Min: %s/%s\n", strconv.FormatFloat(xMinMax[0], 'f', -1, 64), strconv.FormatFloat(xMinMax[1], 'f', -1, 64), strconv.FormatFloat(yMinMax[0], 'f', -1, 64), strconv.FormatFloat(yMinMax[1], 'f', -1, 64), strconv.FormatFloat(zMinMax[0], 'f', -1, 64), strconv.FormatFloat(zMinMax[1], 'f', -1, 64))

	fmt.Printf("\nx bounds: %s, y bounds: %s, z bounds: %s\n", strconv.FormatFloat(math.Abs(xMinMax[0]-xMinMax[1])/minimumVoxelAmount, 'f', -1, 64), strconv.FormatFloat(math.Abs(yMinMax[0]-yMinMax[1])/minimumVoxelAmount, 'f', -1, 64), strconv.FormatFloat(math.Abs(zMinMax[0]-zMinMax[1])/minimumVoxelAmount, 'f', -1, 64))

	//voxelSize := voxelEndSize * math.Pow(float64(scaleFactor), float64(Iterations-1))
	for _, point := range points {
		for i := 0; i < Iterations; i++ {
			size := math.Pow(float64(scaleFactor), float64(i-1))
			xIndex := int(math.Floor(math.Abs(xMinMax[0]-point.X) / float64(size)))
			yIndex := int(math.Floor(math.Abs(yMinMax[0]-point.Y) / float64(size)))
			zIndex := int(math.Floor(math.Abs(zMinMax[0]-point.Z) / float64(size)))

			subVoxelsAmt := int(size / voxelEndSize)

			if size != voxelEndSize {
				for x := xIndex * subVoxelsAmt; x < (xIndex + 1) * subVoxelsAmt; x++ {
					for y := yIndex * subVoxelsAmt; y < (yIndex + 1) * subVoxelsAmt; y++ {
						for z := zIndex * subVoxelsAmt; z < (zIndex + 1) * subVoxelsAmt; z++ {
							masterVoxels[x][y][z] += 1 / math.Pow(float64(subVoxelsAmt), 3)
						}
					}
				}
			} else {
				masterVoxels[xIndex][yIndex][zIndex] += 1
			}
		}
	}

	GenerateVoxelJson(masterVoxels, voxelEndSize)
}

func GenerateVoxelJson(voxels [][][]float64, voxelSize float64) {
	file, errs := os.Create("VoxelMatrix.JSON")
	if errs != nil {
		panic("Failed to write to file:" + errs.Error())
	}

	const voxelValueThreshold = 0.25

	enc := json.NewEncoder(file)
	cleanedPoints := []pointVal{}

	for xIndex, xArray := range voxels {
		for yIndex, yArray := range xArray {
			for zIndex, Value := range yArray {
				if Value > voxelValueThreshold {
					point := pointVal{X: float64(xIndex) * voxelSize, Y: float64(yIndex) * voxelSize, Z: float64(zIndex) * voxelSize, Value: Value}
					cleanedPoints = append(cleanedPoints, point)
				}
			}
		}
	}

	fmt.Println("Clean Points Got: " + strconv.Itoa(len(cleanedPoints)))

	enc.Encode(pointList{Points: cleanedPoints})
}
