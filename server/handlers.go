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
	fmt.Printf("%+v", req.Version.FeatureVersionConfig.Components)
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

	//fileName := fmt.Sprintf("./packages/%s-%s.zip", req.FeatureName, req.Version.FeatureVersionName)
	//if _, err := os.Stat(fileName); os.IsNotExist(err) {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("package: "+fileName+" not found"))
	//	return
	//}
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

	packageName := fmt.Sprintf("%s-%s", featureDetails.Feature.FeatureName, req.Version.FeatureVersionName)
	packageDir := fmt.Sprintf("./packages/%s", packageName)
	if _, err := os.Stat(packageDir); os.IsNotExist(err) {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("package: \""+packageName+"\" not found"))
		return
	}

	var cmds []*exec.Cmd
	//baseDir := "./install/" + strconv.Itoa(int(req.ProjectId))
	//
	//if _, err := os.Stat(baseDir); os.IsNotExist(err) {
	//	cmds = append(cmds, exec.Command("mkdir", baseDir))
	//}
	featureName := fmt.Sprintf("%s-%s-%s", featureDetails.Feature.FeatureName, req.Version.FeatureVersionName, req.Type)

	//cmds = append(cmds, exec.Command("unzip", "-o", fileName, "-d", featureDir))

	installDir := ""
	if featureDetails.Feature.FeatureOnboot {
		installDir = fmt.Sprintf("%s/src/boot/%s", projectDir, featureName)
	} else {
		if req.Type == "entrance" || !featureDetails.Feature.FeatureReuse {
			installDir = fmt.Sprintf("%s/src/components/%s", projectDir, featureName)
		} else {
			newUUID, err := uuid.NewUUID()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
				return
			}
			installDir = fmt.Sprintf("%s/src/components/%s-%s", projectDir, featureName, newUUID.String()[:8])
			featureName = fmt.Sprintf("%s-%s", featureName, newUUID.String()[:8])
		}

	}
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		cmds = append(cmds, exec.Command("mkdir", installDir))
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
	}
	var uploadFiles []paramsUploadFile
	writeFiles := make(map[string]string)
	buf := new(bytes.Buffer)

	newParamsTemplateParse := paramsTemplateParse{
		InstallDir: strings.Replace(installDir, projectDir+"/src/", "", 1),
		Config:     req.Version.FeatureVersionConfig,
	}
	fmt.Printf("featureNameSplit %+v", newParamsTemplateParse)
	// parse Data
	if req.Version.FeatureVersionConfig.Data.Template != "" {
		t, err := template.ParseFiles(packageDir + "/" + req.Version.FeatureVersionConfig.Data.Template)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}
		//fmt.Println("newParamsTemplateParse", newParamsTemplateParse)

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
		writeFiles[req.Version.FeatureVersionConfig.Data.Template] = targetString

		for _, value := range req.Version.FeatureVersionConfig.Data.Values {
			if value.FormType == "upload" {
				uploadPath, ok := value.Value.(string)
				if ok && uploadPath != "" && len(uploadPath) > 32 {
					fmt.Printf("upload %+v", uploadPath)
					uploadFiles = append(uploadFiles, paramsUploadFile{
						Dst:  installDir,
						File: uploadPath,
					})
				}
			}
		}
	}
	for _, v := range req.Version.FeatureVersionConfig.Components {
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
			writeFiles[req.Version.FeatureVersionConfig.Data.Template] = targetString
		}
	}

	configByte, err := json.Marshal(req.Version.FeatureVersionConfig)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}

	//fmt.Printf("%+v", req.Version.FeatureVersionConfig.Data)
	//for _, v := range req.Version.FeatureVersionConfig.Data {
	//	t, err := template.ParseFiles(packageDir + "/" + v.Template + ".tmpl")
	//	if err != nil {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//	newParamsTemplateParse := paramsTemplateParse{
	//		InstallDir: strings.Replace(installDir, projectDir+"/src/", "", 1),
	//		Values:     v.Values,
	//	}
	//	//fmt.Println("newParamsTemplateParse", newParamsTemplateParse)
	//	buf := new(bytes.Buffer)
	//	err = t.Execute(buf, newParamsTemplateParse)
	//	if err != nil {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//	//fmt.Println("buf.String()", buf.String())
	//	targetString := ""
	//	if s, ok := writeFiles[v.Target]; ok {
	//		targetString = s
	//	} else {
	//		b, err := ioutil.ReadFile(packageDir + "/" + v.Target)
	//		if err != nil {
	//			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//			return
	//		}
	//		targetString = string(b)
	//
	//	}
	//	targetString = strings.Replace(targetString, "__data."+v.Template+"__", buf.String(), 1)
	//	targetString = strings.Replace(targetString, "&#34;", `"`, -1)
	//	targetString = strings.Replace(targetString, "&#39;", `'`, -1)
	//	writeFiles[v.Target] = targetString
	//
	//	for _, value := range v.Values {
	//		if value.Type == "upload" {
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
	//
	//for _, v := range req.Version.FeatureVersionConfig.Features {
	//	targetString := ""
	//	if s, ok := writeFiles[v.Target]; ok {
	//		targetString = s
	//	} else {
	//		b, err := ioutil.ReadFile(packageDir + "/" + v.Target)
	//		if err != nil {
	//			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//			return
	//		}
	//		targetString = string(b)
	//	}
	//	fmt.Printf("Values %+v",v.Values)
	//	fmt.Printf("targetString %s\n\n",targetString)
	//	buf := new(bytes.Buffer)
	//	var t *template.Template
	//	if v.Target == v.Template {
	//		t, err = template.New(".").Parse(targetString)
	//	} else {
	//		t, err = template.ParseFiles(packageDir + "/" + v.Template + ".tmpl")
	//	}
	//	if err != nil {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//
	//	if len(v.Values) > 0 {
	//		err = t.Execute(buf, v.Values)
	//		if err != nil {
	//			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//			return
	//		}
	//	}
	//	targetString = strings.Replace(targetString, "__feature."+v.Template+"__", buf.String(), 1)
	//	targetString = strings.Replace(targetString, "&#34;", `"`, -1)
	//	targetString = strings.Replace(targetString, "&#39;", `'`, -1)
	//	writeFiles[v.Target] = targetString
	//}
	//
	cmds = cmds[:0]
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

	//reply, err := client.ProjectFeaturesByProjectId(context.Background(), &pb.ProjectFeaturesByProjectIdRequest{
	//	ProjectId: req.ProjectId,
	//})
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//	return
	//}
	//if req.FeatureOnBoot {
	//	bootString := []string{}
	//
	//	for _, v := range reply.Features {
	//		if v.FeatureOnboot {
	//			bootString = append(bootString, fmt.Sprintf(`'%s-%s-%s'`, v.FeatureName, v.FeatureVersionName, v.ProjectFeaturesType))
	//		}
	//	}
	//	bootString = append(bootString, `'`+featureName+`'`)
	//	quasarConfFileByte, err := ioutil.ReadFile("./packages/quasar.conf.js.tmpl")
	//	if err != nil {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//	newQuasarConfFileString := strings.Replace(string(quasarConfFileByte), "__data.boot__", strings.Join(bootString, ","), 1)
	//	err = ioutil.WriteFile(projectDir+"/quasar.conf.js", []byte(newQuasarConfFileString), os.ModePerm)
	//	if err != nil {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
	//		return
	//	}
	//} else {
	//	entranceName := ""
	//	var routes []paramsRoutesJsRoutesParam
	//	for _, v := range reply.Features {
	//		if v.GetProjectFeaturesType() == "entrance" {
	//			entranceName = v.GetProjectFeaturesInstallName()
	//		}
	//		if v.GetProjectFeaturesRoutePath() != "" {
	//			routes = append(routes, paramsRoutesJsRoutesParam{
	//				Path: v.GetProjectFeaturesRoutePath(),
	//				Page: "components/" + v.GetProjectFeaturesInstallName() + "/Index.vue",
	//			})
	//		}
	//	}
	//	if req.Type == "entrance" {
	//		entranceName = featureName
	//	}
	//
	//	if req.Type == "page" {
	//		routes = append(routes, paramsRoutesJsRoutesParam{
	//			Path: req.RoutePath,
	//			Page: "components/" + featureName + "/Index.vue",
	//		})
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
	//	buf := new(bytes.Buffer)
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
	//
	if len(uploadFiles) > 0 {
		cmds = cmds[:0]
		for _, v := range uploadFiles {
			cmds = append(cmds, exec.Command("mv", "./tmp/"+v.File, v.Dst))
		}
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
	_ = configByte
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
