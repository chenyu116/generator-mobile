package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/chenyu116/generator-mobile/proto"
	"github.com/chenyu116/generator-mobile/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type result struct {
	Error string `json:"error,omitempty"`
}

func jsonError(err string) result {
	return result{
		Error: err,
	}
}

func projects(c *gin.Context) {
	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("grpc connection lost"))
		return
	}
	defer utils.DbServerGrpcConn.Put(dbServerConn)
	client := pb.NewApiClient(dbServerConn)
	reply, err := client.Projects(context.Background(), &pb.ProjectsRequest{
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, reply.Projects)
}

func projectInitialized(c *gin.Context) {

	dirs, _ := ioutil.ReadDir("./projects")
	fmt.Println(dirs)
	var projects []string
	for _, fi := range dirs {
		projects = append(projects, "'"+fi.Name()+"'")
	}
	if len(projects) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, projects)
		return
	}

	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("GRPC connection lost"))
		return
	}
	defer utils.DbServerGrpcConn.Put(dbServerConn)
	client := pb.NewApiClient(dbServerConn)
	reply, err := client.ProjectInitialized(context.Background(), &pb.ProjectInitializedRequest{
		Projects: strings.Join(projects, ","),
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, reply.Projects)
}

func features(c *gin.Context) {
	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("GRPC connection lost"))
		return
	}
	defer utils.DbServerGrpcConn.Put(dbServerConn)
	client := pb.NewApiClient(dbServerConn)
	reply, err := client.Features(context.Background(), &pb.FeaturesRequest{
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, reply.Feature)
}

func feature(c *gin.Context) {
	var req RequestInt32FeatureId
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	if req.FeatureId <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("invalid params"))
		return
	}
	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("GRPC connection lost"))
		return
	}
	defer utils.DbServerGrpcConn.Put(dbServerConn)
	client := pb.NewApiClient(dbServerConn)
	reply, err := client.Feature(context.Background(), &pb.FeatureRequest{
		FeatureId: req.FeatureId,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, reply.Feature)
}

func projectFeatures(c *gin.Context) {
	var req RequestInt32ProjectId
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	if req.ProjectId <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("miss project_id"))
		return
	}
	projectId := strconv.Itoa(int(req.ProjectId))
	baseDir := "./projects/" + projectId

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, jsonError("需要初始化项目"))
		return
	}
	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("GRPC connection lost"))
		return
	}
	defer utils.DbServerGrpcConn.Put(dbServerConn)
	client := pb.NewApiClient(dbServerConn)
	reply, err := client.ProjectFeaturesByProjectId(context.Background(), &pb.ProjectFeaturesByProjectIdRequest{
		ProjectId: req.ProjectId,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, reply.Features)
}
func projectInit(c *gin.Context) {
	var req RequestStringProjectId
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	quasarDir := "/home/roger/workspace/quasar"
	baseDir := "./projects/" + req.ProjectId + "/"
	var cmds []*exec.Cmd

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		cmds = append(cmds, exec.Command("mkdir", baseDir))
	}
	copyFiles := []string{"src", ".quasar", "quasar.conf.js.tmpl"}
	for _, v := range copyFiles {
		if _, err := os.Stat(baseDir + "/" + v); os.IsNotExist(err) {
			cmds = append(cmds, exec.Command("cp", "-r", quasarDir+"/"+v, baseDir))
		}
	}
	linkFiles := []string{"node_modules", "babel.config.js", "jsconfig.json", "package.json", "yarn.lock", ".eslintignore", ".eslintrc.js", ".gitignore", ".postcssrc.js"}
	for _, v := range linkFiles {
		if _, err := os.Stat(baseDir + "/" + v); os.IsNotExist(err) {
			cmds = append(cmds, exec.Command("ln", "-s", quasarDir+"/"+v, baseDir))
		}
	}
	// Run the pipeline
	_, _, err := utils.Pipeline(cmds...)
	if err != nil && !strings.HasPrefix(err.Error(), "exit") {
		fmt.Println("err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}

	// Print the stderr, if any
	//if len(stderr) > 0 {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(string(stderr)))
	//	return
	//}

	//cmd := exec.Command("quasar", "build")
	//cmd.Dir = baseDir
	////显示运行的命令
	//stdout, err := cmd.StdoutPipe()
	//if err != nil {
	//	fmt.Println(err)
	//	c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//	return
	//}
	//
	//cmd.Start()
	//reader := bufio.NewReader(stdout)
	//
	////实时循环读取输出流中的一行内容
	//buf := new(bytes.Buffer)
	//for {
	//	line, err2 := reader.ReadString('\n')
	//	if err2 != nil || io.EOF == err2 {
	//		break
	//	}
	//	buf.WriteString(line + "<br />")
	//}
	//
	//cmd.Wait()

	c.AbortWithStatus(http.StatusOK)
}

func upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	newUUID, err := uuid.NewUUID()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	extIndex := strings.LastIndex(file.Filename, ".")

	filename := newUUID.String() + file.Filename[extIndex:]
	if err := c.SaveUploadedFile(file, "./tmp/"+filename); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"file": filename})

}
func install(c *gin.Context) {
	var req RequestPostInstall
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	if req.FeatureId <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("invalid params"))
		return
	}
	if req.ProjectId <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("invalid params"))
		return
	}

	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("GRPC connection lost"))
		return
	}
	defer utils.DbServerGrpcConn.Put(dbServerConn)
	client := pb.NewApiClient(dbServerConn)
	featureDetails, err := client.Feature(context.Background(), &pb.FeatureRequest{
		FeatureId: req.FeatureId,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	featureName := fmt.Sprintf("%s-%s-%s-%s", featureDetails.Feature.FeatureName, req.Version.FeatureVersionName, req.Type, newUUID.String()[:8])

	configByte, err := json.Marshal(req.Version.FeatureVersionConfig)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}

	_, err = client.CreateProjectFeature(context.Background(), &pb.CreateProjectFeatureRequest{
		FeatureId:                  req.FeatureId,
		ProjectFeaturesType:        req.Type,
		ProjectFeaturesConfig:      string(configByte),
		ProjectId:                  req.ProjectId,
		FeatureVersionId:           req.Version.FeatureVersionId,
		ProjectFeaturesInstallName: featureName,
		ProjectFeaturesName:        req.ProjectFeaturesName,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}
func build(c *gin.Context) {
	var req RequestInt32ProjectId
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	if req.ProjectId <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("invalid params"))
		return
	}
	//return
	//configString := strings.Replace(req.FeatureVersionConfigString, `\"`, `"`, -1)
	//configString = strings.Replace(configString, `"{`, `{`, -1)
	//configString = strings.Replace(configString, `}"`, `}`, -1)
	//uuid , err := uuid.NewUUID()
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//	return
	//}
	projectDir := fmt.Sprintf("./projects/%d", req.ProjectId)
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("project not initialized"))
		return
	}

	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("GRPC connection lost"))
		return
	}
	client := pb.NewApiClient(dbServerConn)
	projectFeatures, err := client.ProjectFeaturesByProjectId(context.Background(), &pb.ProjectFeaturesByProjectIdRequest{
		ProjectId: req.ProjectId,
	})
	if err != nil {
		utils.DbServerGrpcConn.Put(dbServerConn)
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	utils.DbServerGrpcConn.Put(dbServerConn)

	var cmds []*exec.Cmd
	var bootString []string
	routes := make(map[string]paramsRoutesJsRoutesParam)
	buf := new(bytes.Buffer)
	for _, feature := range projectFeatures.Features {
		cmds = cmds[:0]
		installDir := ""
		var projectConfig featureVersionConfig
		err = json.Unmarshal([]byte(feature.ProjectFeaturesConfig), &projectConfig)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}

		if feature.FeatureOnboot {
			installDir = fmt.Sprintf("%s/src/boot/%s", projectDir, feature.ProjectFeaturesInstallName)
		} else {
			installDir = fmt.Sprintf("%s/src/components/%s", projectDir, feature.ProjectFeaturesInstallName)
		}
		if _, err := os.Stat(installDir); os.IsNotExist(err) {
			cmds = append(cmds, exec.Command("mkdir", installDir))
		}

		//for _, v := range projectConfig.Data.Values {
		//	if v.FormType == "upload" {
		//		if _, err := os.Stat(installDir + "/" + v.Value.(string)); os.IsNotExist(err) {
		//			cmds = append(cmds, exec.Command("cp", "./tmp/"+v.Value.(string), installDir))
		//		}
		//	}
		//}

		if len(cmds) > 0 {
			_, stderr, err := utils.Pipeline(cmds...)
			if err != nil && !strings.HasPrefix(err.Error(), "exit") {
				c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
				return
			}

			if len(stderr) > 0 {
				c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(string(stderr)))
				return
			}
		}
		cmds = cmds[:0]

		if feature.ProjectFeaturesType == "click" {
			continue
		}

		packageName := fmt.Sprintf("%s-%s", feature.FeatureName, feature.FeatureVersionName)
		packageDir := fmt.Sprintf("./packages/%s", packageName)
		if _, err := os.Stat(packageDir); os.IsNotExist(err) {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("package: \""+packageName+"\" not found"))
			return
		}

		files, err := ioutil.ReadDir(packageDir)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}

		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".tmpl") {
				continue
			}
			if _, err := os.Stat(installDir + "/" + f.Name()); os.IsNotExist(err) {
				cmds = append(cmds, exec.Command("cp", packageDir+"/"+f.Name(), installDir))
			}
		}
		if len(cmds) > 0 {
			_, stderr, err := utils.Pipeline(cmds...)
			if err != nil && !strings.HasPrefix(err.Error(), "exit") {
				c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
				return
			}

			if len(stderr) > 0 {
				c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(string(stderr)))
				return
			}
			cmds = cmds[:0]
		}

		writeFiles := make(map[string]string)
		dataValues := make(map[string]interface{})
		for _, v := range projectConfig.Data.Values {
			dataValues[v.Key] = v.Value
		}

		newParamsTemplateParse := paramsTemplateParse{
			InstallDir: strings.Replace(installDir, projectDir+"/src/", "", 1),
			Config:     projectConfig,
			DataValues: dataValues,
		}
		//fmt.Printf("%+v\n\n", newParamsTemplateParse)
		if projectConfig.Data.Template != "" {
			t, err := template.ParseFiles(packageDir + "/" + projectConfig.Data.Template)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
				return
			}
			buf := new(bytes.Buffer)
			err = t.Execute(buf, newParamsTemplateParse)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
				return
			}

			targetString := strings.Replace(buf.String(), "&#34;", `"`, -1)
			targetString = strings.Replace(targetString, "&#39;", `'`, -1)
			targetString = strings.Replace(targetString, "&lt;", `<`, -1)
			targetString = strings.Replace(targetString, `|"`, "", -1)
			targetString = strings.Replace(targetString, `"|`, "", -1)
			writeFiles[projectConfig.Data.Template] = targetString
		}

		for _, v := range projectConfig.Components {
			if _, ok := writeFiles[v.Template]; !ok {
				t, err := template.ParseFiles(packageDir + "/" + v.Template)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
					return
				}
				//fmt.Println("newParamsTemplateParse", newParamsTemplateParse)
				buf.Reset()
				err = t.Execute(buf, newParamsTemplateParse)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
					return
				}
				targetString := strings.Replace(buf.String(), "&#34;", `"`, -1)
				targetString = strings.Replace(targetString, "&#39;", `'`, -1)
				writeFiles[v.Template] = targetString
			}
		}
		if len(writeFiles) > 0 {
			for file, s := range writeFiles {
				err := ioutil.WriteFile(installDir+"/"+file, []byte(s), os.ModePerm)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
					return
				}
				cmds = append(cmds, exec.Command("prettier", "--config", projectDir+"/package.json", "--write", installDir+"/"+file))
			}
		}

		if len(cmds) > 0 {
			_, stderr, err := utils.Pipeline(cmds...)
			if err != nil && !strings.HasPrefix(err.Error(), "exit") {
				c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
				return
			}

			if len(stderr) > 0 {
				c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(string(stderr)))
				return
			}
		}
		if feature.FeatureOnboot {
			bootString = append(bootString, `'`+feature.GetProjectFeaturesInstallName()+`'`)
		}

		if feature.ProjectFeaturesType == "entrance" {
			routes["/"] = paramsRoutesJsRoutesParam{
				Path: "/",
				Page: "components/" + feature.GetProjectFeaturesInstallName() + "/Index.vue",
			}
		} else if feature.ProjectFeaturesType == "page" {
			for _, cf := range projectConfig.Data.Values {
				if cf.Key == "routePath" {
					path, ok := cf.Value.(string)
					if !ok {
						break
					}
					if _, ok = routes[path]; ok {
						break
					}
					routes[path] = paramsRoutesJsRoutesParam{
						Path: cf.Value.(string),
						Page: "components/" + feature.GetProjectFeaturesInstallName() + "/Index.vue",
					}
					break
				}
			}
		}
	}

	quasarT, err := template.ParseFiles("./packages/quasar.conf.js.tmpl")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	buf.Reset()
	err = quasarT.Execute(buf, bootString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	targetString := strings.Replace(buf.String(), "&#34;", `"`, -1)
	targetString = strings.Replace(targetString, "&#39;", `'`, -1)
	err = ioutil.WriteFile(projectDir+"/quasar.conf.js", []byte(targetString), os.ModePerm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}

	routesJsTemp, err := template.ParseFiles("./packages/routes.js.tmpl")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	buf.Reset()
	err = routesJsTemp.Execute(buf, routes)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	err = ioutil.WriteFile(projectDir+"/src/router/routes.js", buf.Bytes(), os.ModePerm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}

	//cmds = append(cmds, exec.Command("unzip", "-o", fileName, "-d", featureDir))

	//var uploadFiles []paramsUploadFile
	//writeFiles := make(map[string]string)
	//buf := new(bytes.Buffer)
	//
	//dataValues := make(map[string]interface{})
	//for _, v := range req.Version.FeatureVersionConfig.Data.Values {
	//	dataValues[v.Key] = v.Value
	//}
	//newParamsTemplateParse := paramsTemplateParse{
	//	InstallDir: strings.Replace(installDir, projectDir+"/src/", "", 1),
	//	Config:     req.Version.FeatureVersionConfig,
	//	DataValues: dataValues,
	//}
	//// parse Data
	//if req.Version.FeatureVersionConfig.Data.Template != "" {
	//	t, err := template.ParseFiles(packageDir + "/" + req.Version.FeatureVersionConfig.Data.Template)
	//	if err != nil {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//	//fmt.Println("newParamsTemplateParse", newParamsTemplateParse)
	//
	//	err = t.Execute(buf, newParamsTemplateParse)
	//	if err != nil {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//	targetString := strings.Replace(buf.String(), "&#34;", `"`, -1)
	//	targetString = strings.Replace(targetString, "&#39;", `'`, -1)
	//	targetString = strings.Replace(targetString, "&lt;", `<`, -1)
	//	targetString = strings.Replace(targetString, `|"`, "", -1)
	//	targetString = strings.Replace(targetString, `"|`, "", -1)
	//	writeFiles[req.Version.FeatureVersionConfig.Data.Template] = targetString
	//
	//	for _, value := range req.Version.FeatureVersionConfig.Data.Values {
	//		if value.FormType == "upload" {
	//			uploadPath, ok := value.Value.(string)
	//			if ok && uploadPath != "" && len(uploadPath) > 32 {
	//				fmt.Printf("upload %+v", uploadPath)
	//				uploadFiles = append(uploadFiles, paramsUploadFile{
	//					Dst:  installDir,
	//					File: uploadPath,
	//				})
	//			}
	//		}
	//	}
	//}
	//for _, v := range req.Version.FeatureVersionConfig.Components {
	//	if _, ok := writeFiles[v.Template]; !ok {
	//		t, err := template.ParseFiles(packageDir + "/" + v.Template)
	//		if err != nil {
	//			continue
	//			//c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//			//return
	//		}
	//		//fmt.Println("newParamsTemplateParse", newParamsTemplateParse)
	//		buf.Reset()
	//		err = t.Execute(buf, newParamsTemplateParse)
	//		if err != nil {
	//			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//			return
	//		}
	//		targetString := strings.Replace(buf.String(), "&#34;", `"`, -1)
	//		targetString = strings.Replace(targetString, "&#39;", `'`, -1)
	//		writeFiles[req.Version.FeatureVersionConfig.Data.Template] = targetString
	//	}
	//}
	//
	//configByte, err := json.Marshal(req.Version.FeatureVersionConfig)
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//	return
	//}
	//
	//cmds = cmds[:0]
	//if len(writeFiles) > 0 {
	//	for file, s := range writeFiles {
	//		err := ioutil.WriteFile(installDir+"/"+file, []byte(s), os.ModePerm)
	//		if err != nil {
	//			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//			return
	//		}
	//		cmds = append(cmds, exec.Command("prettier", "--config", projectDir+"/package.json", "--write", installDir+"/"+file))
	//	}
	//}
	//
	//if len(cmds) > 0 {
	//	_, stderr, err := utils.Pipeline(cmds...)
	//	if err != nil && !strings.HasPrefix(err.Error(), "exit") {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//
	//	if len(stderr) > 0 {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(string(stderr)))
	//		return
	//	}
	//}
	//
	//if featureDetails.Feature.FeatureOnboot {
	//	bootString := []string{}
	//
	//	for _, v := range projectFeatures.Features {
	//		if v.FeatureOnboot {
	//			bootString = append(bootString, v.GetProjectFeaturesInstallName())
	//		}
	//	}
	//	bootString = append(bootString, `'`+featureName+`'`)
	//	quasarT, err := template.ParseFiles("./packages/quasar.conf.js.tmpl")
	//	if err != nil {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//	buf.Reset()
	//	err = quasarT.Execute(buf, bootString)
	//	err = ioutil.WriteFile(projectDir+"/quasar.conf.js", buf.Bytes(), os.ModePerm)
	//	if err != nil {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//} else {
	//	entranceName := ""
	//	routes := make(map[string]paramsRoutesJsRoutesParam)
	//	//var routes []paramsRoutesJsRoutesParam
	//	for _, v := range projectFeatures.Features {
	//		if v.GetProjectFeaturesType() == "entrance" {
	//			entranceName = v.GetProjectFeaturesInstallName()
	//		}
	//
	//		var projectFeatureConfig featureVersionConfig
	//		err = json.Unmarshal([]byte(v.GetProjectFeaturesConfig()), &projectFeatureConfig)
	//		if err == nil {
	//			for _, cf := range projectFeatureConfig.Data.Values {
	//				if cf.Key == "routePath" {
	//					path, ok := cf.Value.(string)
	//					if !ok {
	//						break
	//					}
	//					if _, ok = routes[path]; ok {
	//						break
	//					}
	//					routes[path] = paramsRoutesJsRoutesParam{
	//						Path: cf.Value.(string),
	//						Page: "components/" + v.GetProjectFeaturesInstallName() + "/Index.vue",
	//					}
	//					break
	//				}
	//			}
	//		}
	//	}
	//	if req.Type == "entrance" {
	//		entranceName = featureName
	//	}
	//
	//	if req.Type == "page" {
	//		for _, cf := range req.Version.FeatureVersionConfig.Data.Values {
	//			if cf.Key == "routePath" {
	//				path, ok := cf.Value.(string)
	//				if !ok {
	//					break
	//				}
	//				if _, ok = routes[path]; ok {
	//					break
	//				}
	//				routes[path] = paramsRoutesJsRoutesParam{
	//					Path: cf.Value.(string),
	//					Page: "components/" + featureName + "/Index.vue",
	//				}
	//				break
	//			}
	//		}
	//	}
	//
	//	routesJs := paramsRoutesJs{
	//		EntranceName: "components/" + entranceName + "/Index.vue",
	//		Routes:       routes,
	//	}
	//	routesJsTemp, err := template.ParseFiles("./packages/routes.js.tmpl")
	//	if err != nil {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//	buf.Reset()
	//	err = routesJsTemp.Execute(buf, routesJs)
	//	if err != nil {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//	err = ioutil.WriteFile(projectDir+"/src/router/routes.js", buf.Bytes(), os.ModePerm)
	//	if err != nil {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//}
	////
	//if len(uploadFiles) > 0 {
	//	cmds = cmds[:0]
	//	for _, v := range uploadFiles {
	//		cmds = append(cmds, exec.Command("mv", "./tmp/"+v.File, v.Dst))
	//	}
	//	_, stderr, err := utils.Pipeline(cmds...)
	//	if err != nil && !strings.HasPrefix(err.Error(), "exit") {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//
	//	if len(stderr) > 0 {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(string(stderr)))
	//		return
	//	}
	//}
	//_ = configByte
	//_, err = client.CreateProjectFeature(context.Background(), &pb.CreateProjectFeatureRequest{
	//	FeatureId:                  req.FeatureId,
	//	ProjectFeaturesType:        req.Type,
	//	ProjectFeaturesConfig:      string(configByte),
	//	ProjectId:                  req.ProjectId,
	//	FeatureVersionId:           req.Version.FeatureVersionId,
	//	ProjectFeaturesInstallName: featureName,
	//	ProjectFeaturesName:        req.ProjectFeaturesName,
	//})
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//	return
	//}
	c.AbortWithStatus(http.StatusNoContent)
}
