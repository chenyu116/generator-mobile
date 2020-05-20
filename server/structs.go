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

type RequestPostInstall struct {
	RequestInt32FeatureId
	RequestInt32ProjectId
	FeatureOnBoot bool `form:"featureOnBoot" json:"featureOnBoot"`
	FeatureName string `form:"featureName" binding:"required" json:"featureName"`
	Version featureVersion `form:"version" binding:"required" json:"version"`
	FeatureVersionConfigString string `form:"featureVersionConfigString" binding:"required" json:"featureVersionConfigString"`
	Type    string `form:"type" binding:"required" json:"type"`
}

type featureVersion struct {
	FeatureVersionConfig featureVersionConfig `form:"feature_version_config" binding:"required" json:"feature_version_config"`
	FeatureVersionName string `form:"feature_version_name" binding:"required" json:"feature_version_name"`
	FeatureVersionId int32 `form:"feature_version_id" binding:"required" json:"feature_version_id"`
}
type featureVersionConfigDataValue struct {
	Key string `form:"key" binding:"required" json:"key"`
	Value interface{} `form:"value" binding:"required" json:"value"`
}
type featureVersionConfigData struct {
	Name string `form:"name" binding:"required" json:"name"`
	Template string `form:"template" binding:"required" json:"template"`
	Target string `form:"target" binding:"required" json:"target"`
	Values []featureVersionConfigDataValue `form:"values" binding:"required" json:"values"`
}
type featureVersionConfig struct {
	Data []featureVersionConfigData `form:"data" binding:"required" json:"data"`
}