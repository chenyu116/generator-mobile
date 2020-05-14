package server

type RequestProjectId struct {
	ProjectId int32 `form:"projectId" binding:"required" json:"projectId"`
}
