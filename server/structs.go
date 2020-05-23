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
	FeatureOnBoot              bool           `form:"featureOnBoot" json:"featureOnBoot"`
	FeatureName                string         `form:"featureName" binding:"required" json:"featureName"`
	Version                    featureVersion `form:"version" binding:"required" json:"version"`
	FeatureVersionConfigString string         `form:"featureVersionConfigString" binding:"required" json:"featureVersionConfigString"`
	Type                       string         `form:"type" binding:"required" json:"type"`
	RoutePath                  string         `form:"routePath" json:"routePath"`
	ProjectFeaturesName        string         `form:"projectFeaturesName" json:"projectFeaturesName"`
}

type featureVersion struct {
	FeatureVersionConfig featureVersionConfig `form:"feature_version_config" binding:"required" json:"feature_version_config"`
	FeatureVersionName   string               `form:"feature_version_name" binding:"required" json:"feature_version_name"`
	FeatureVersionId     int32                `form:"feature_version_id" binding:"required" json:"feature_version_id"`
}
type featureVersionConfigDataValue struct {
	Key                   string      `form:"key" binding:"required" json:"key"`
	Type                  string      `form:"type" binding:"required" json:"type"`
	Value                 interface{} `form:"value" binding:"required" json:"value"`
	ProjectFeaturesConfig featureVersionConfig      `form:"project_features_config" json:"project_features_config"`
}
type featureVersionConfigData struct {
	Name     string                          `form:"name" binding:"required" json:"name"`
	Template string                          `form:"template" binding:"required" json:"template"`
	Target   string                          `form:"target" binding:"required" json:"target"`
	Values   []featureVersionConfigDataValue `form:"values" binding:"required" json:"values"`
}

type featureVersionConfigFeature struct {
	featureVersionConfigData
	Limit int      `form:"limit" binding:"required" json:"limit"`
	Type  []string `form:"type"  json:"type"`
}

type featureVersionConfig struct {
	Data     []featureVersionConfigData    `form:"data" json:"data"`
	Features []featureVersionConfigFeature `form:"features" json:"features"`
}

type paramsRoutesJsRoutesParam struct {
	Path string
	Page string
}
type paramsRoutesJs struct {
	EntranceName string
	Routes       []paramsRoutesJsRoutesParam
}

type paramsTemplateParse struct {
	InstallDir string
	Values     []featureVersionConfigDataValue
}

type paramsUploadFile struct {
	Dst  string
	File string
}
