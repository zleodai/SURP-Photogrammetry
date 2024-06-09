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
	Value uint8
}

type pointList struct {
	Points    []pointVal
	VoxelSize float64
}

func filterMesh(xMinMax, yMinMax, zMinMax [2]float64, points []pointCloudDecoder.Point, voxelSize float64, voxels [][][]uint8) {
	noiseMaskSize := voxelSize * 4

	voxelMask := MinMaxMesh(xMinMax, yMinMax, zMinMax, points, noiseMaskSize, false)

	subVoxelsAmt := int(noiseMaskSize / voxelSize)

	// Filters Mesh by iterating through a bigger voxel cubed Voxel and using the empty space in the bigger voxel to filter out noise.
	go func() {
		for xIndex, xArray := range voxelMask {
			for yIndex, yArray := range xArray {
				for zIndex, Value := range yArray {
					if Value != 0 && Value <= 2 {
						for x := 0; x < subVoxelsAmt; x++ {
							for y := 0; y < subVoxelsAmt; y++ {
								for z := 0; z < subVoxelsAmt; z++ {
									voxels[x+(xIndex*subVoxelsAmt)][y+(yIndex*subVoxelsAmt)][z+(zIndex*subVoxelsAmt)] = 0
								}
							}
						}
					}
				}
			}
		}
	}()

}

func MinMaxMesh(xMinMax, yMinMax, zMinMax [2]float64, points []pointCloudDecoder.Point, voxelSize float64, filterNoise bool) [][][]uint8 {
	var xSize int = int(math.Floor(math.Abs(xMinMax[0]-xMinMax[1]) / voxelSize))
	var ySize int = int(math.Floor(math.Abs(yMinMax[0]-yMinMax[1]) / voxelSize))
	var zSize int = int(math.Floor(math.Abs(zMinMax[0]-zMinMax[1]) / voxelSize))
	fmt.Println("Voxel Grid X, Y, Z:", xSize, ySize, zSize)
	fmt.Println("Total Voxel Amount:", xSize*ySize*zSize)

	voxels := make([][][]uint8, xSize+4)
	for i := 0; i < len(voxels); i++ {
		voxels[i] = make([][]uint8, ySize+4)
		for j := 0; j < len(voxels[i]); j++ {
			voxels[i][j] = make([]uint8, zSize+4)
		}
	}

	go func() {
		for _, point := range points {
			xIndex := int(math.Floor(math.Abs(xMinMax[0]-point.X) / voxelSize))
			yIndex := int(math.Floor(math.Abs(yMinMax[0]-point.Y) / voxelSize))
			zIndex := int(math.Floor(math.Abs(zMinMax[0]-point.Z) / voxelSize))
			if voxels[xIndex][yIndex][zIndex] != ^uint8(0) {
				voxels[xIndex][yIndex][zIndex] += 1
			}
		}
	}()

	//Add threshold filter
	if filterNoise {
		filterMesh(xMinMax, yMinMax, zMinMax, points, voxelSize, voxels)
	}

	runtime.GC()
	return voxels
}

func GenerateVoxelJson(voxels [][][]uint8, voxelSize float64) {
	file, errs := os.Create("VoxelMatrix.JSON")
	if errs != nil {
		panic("Failed to write to file:" + errs.Error())
	}

	const voxelValueThreshold = 0

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

	enc.Encode(pointList{Points: cleanedPoints, VoxelSize: voxelSize})
}
