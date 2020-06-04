package server

import pb "github.com/chenyu116/generator-mobile/proto"

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
	Version             featureVersion `form:"version" binding:"required" json:"version"`
	Type                string         `form:"type" binding:"required" json:"type"`
	ProjectFeaturesName string         `form:"projectFeaturesName" json:"projectFeaturesName"`
}

type featureVersion struct {
	FeatureVersionConfig featureVersionConfig `form:"feature_version_config" binding:"required" json:"feature_version_config"`
	FeatureVersionName   string               `form:"feature_version_name" binding:"required" json:"feature_version_name"`
	FeatureVersionId     int32                `form:"feature_version_id" binding:"required" json:"feature_version_id"`
}
type featureVersionConfigDataValue struct {
	Key      string      `form:"key" binding:"required" json:"key"`
	FormType string      `form:"formType" binding:"required" json:"formType"`
	Value    interface{} `form:"value" binding:"required" json:"value"`
	Name     string      `form:"name" binding:"required" json:"name"`
}
type featureVersionConfigData struct {
	Name     string                          `form:"name" json:"name"`
	Template string                          `form:"template" json:"template"`
	Values   []featureVersionConfigDataValue `form:"values" json:"values"`
}

type featureVersionConfigComponent struct {
	Name          string                               `form:"name" json:"name"`
	Template      string                               `form:"template" binding:"required" json:"template"`
	Key           string                               `form:"key" binding:"required" json:"key"`
	Limit         int                                  `form:"limit" json:"limit"`
	Accept        []string                             `form:"accept"  json:"accept"`
	Values        []featureVersionConfigComponentValue `form:"values" binding:"required" json:"values"`
}

type featureVersionConfigComponentValue struct {
	ProjectFeaturesId          int32                `form:"project_features_id" binding:"required" json:"project_features_id"`
	ProjectFeaturesConfig      featureVersionConfig `form:"project_features_config" binding:"required" json:"project_features_config"`
	ComponentHash              string               `form:"componentHash" json:"componentHash"`
	ProjectFeaturesInstallName string               `form:"project_features_install_name" json:"project_features_install_name"`
}

type featureVersionConfig struct {
	Data         featureVersionConfigData        `form:"data" json:"data"`
	Dependencies []string                        `form:"dependencies" json:"dependencies"`
	Components   []featureVersionConfigComponent `form:"components" json:"components"`
}

type paramsRoutesJsRoutesParam struct {
	Path string
	Page string
}
type paramsRoutesJs struct {
	EntranceName string
	Routes       map[string]paramsRoutesJsRoutesParam
}

type paramsTemplateParse struct {
	InstallDir    string
	Config        featureVersionConfig
	DataValues    map[string]interface{}
}

type paramsUploadFile struct {
	Dst  string
	File string
}

type projectFeature struct {
	pb.ProjectFeatureAll
	featureVersionConfig
}
