package pointCloudDecoder

import (
	"encoding/json"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type PointCloud struct {
	Views      []View      `json:"views"`
	Intrinsics []Intrinsic `json:"intrinsics"`
	Poses      []Pose      `json:"poses"`
	Structures []Structure `json:"structure"`
}

type View struct {
	ViewId      string `json:"viewId"`
	PoseId      string `json:"poseId"`
	FrameId     string `json:"frameId"`
	IntrinsicId string `json:"intrinsicId"`
	ResectionId string `json:"resectionId"`
	Width       string `json:"width"`
	Height      string `json:"height"`
}

type Intrinsic struct {
	IntrinsicId        string `json:"intrinsicId"`
	Width              string `json:"width"`
	Height             string `json:"height"`
	SensorWidth        string `json:"sensorWidth"`
	SensorHeight       string `json:"sensorHeight"`
	SerialNumber       string `json:"serialNumber"`
	Type               string `json:"type"`
	InitializationMode string `json:"initializationMode"`
	InitialFocalLength string `json:"initialFocalLength"`
	FocalLength        string `json:"focalLength"`
}

type Pose struct {
	PoseId   string   `json:"poseId"`
	PoseData PoseData `json:"pose"`
}

type PoseData struct {
	Transform TransformData `json:"transform"`
	Locked    string        `json:"locked"`
}

type TransformData struct {
	Rotation []string `json:"rotation"`
	Center   []string `json:"center"`
}

type Structure struct {
	LandmarkId   string        `json:"landmarkId"`
	DescType     string        `json:"descType"`
	Color        []string      `json:"color"`
	X            []string      `json:"X"`
	Observations []Observation `json:"observations"`
}

type Observation struct {
	ObservationId string   `json:"observationId"`
	FeatureId     string   `json:"featureId"`
	X             []string `json:"x"`
	Scale         string   `json:"scale"`
}

func DecodeFromRawJsonFromPath(jsonFilePath string) PointCloud {
	jsonFile, err := os.ReadFile(jsonFilePath)
	check(err)

	var pointCloudData PointCloud

	json.Unmarshal(jsonFile, &pointCloudData)

	return pointCloudData
}

func GenerateFloatJson(jsonFilePath string, jsonFileName string) string {
	file, errs := os.Create(jsonFileName)
	if errs != nil {
		panic("Failed to write to file:" + errs.Error())
	}
	defer file.Close()

	var data PointCloud = DecodeFromRawJsonFromPath(jsonFilePath)

	_, errs = file.WriteString("{\"points\": \n	[")
	check(errs)

	for index, structure := range data.Structures {
		_, errs = file.WriteString("\n		{")
		check(errs)
		for index, x := range structure.X {
			switch checker := index; checker {
			case 0:
				_, errs = file.WriteString("\n			\"x\": " + x + ",")
				check(errs)
			case 1:
				_, errs = file.WriteString("\n			\"y\": " + x + ",")
				check(errs)
			case 2:
				_, errs = file.WriteString("\n			\"z\": " + x)
				check(errs)
			default:
				panic("More than 3 cordinates found cannot calculate 4d or greater coords")
			}
		}
		if index != len(data.Structures)-1 {
			_, errs = file.WriteString("\n		}, ")
			check(errs)
		} else {
			_, errs = file.WriteString("\n		}")
			check(errs)
		}
	}

	_, errs = file.WriteString("\n	]\n}")
	check(errs)

	var outputPath string = "./" + jsonFileName
	return outputPath
}
