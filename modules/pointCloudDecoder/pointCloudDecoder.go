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

func DecodeFromPath(jsonFilePath string) PointCloud {
	jsonFile, err := os.ReadFile(jsonFilePath)
	check(err)

	var pointCloudData PointCloud

	json.Unmarshal(jsonFile, &pointCloudData)

	return pointCloudData
}
