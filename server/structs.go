package server

type RequestInt32ProjectId struct {
	ProjectId int32 `form:"projectId" binding:"required" json:"projectId"`
}

type RequestStringProjectId struct {
	ProjectId string `form:"projectId" binding:"required" json:"projectId"`
}

type RequestInt32FeatureId struct {
	FeatureId int32 `form:"featureId" binding:"required" json:"featureId"`
}