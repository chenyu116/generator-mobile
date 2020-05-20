package server

import (
	"bytes"
	"context"
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
	Error  string `json:"error,omitempty"`
	Result string `json:"result,omitempty"`
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
	quasarDir := "/home/roger/workspace/generator-mobile/quasar"
	baseDir := "./projects/" + req.ProjectId + "/"
	var cmds []*exec.Cmd

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		cmds = append(cmds, exec.Command("mkdir", baseDir))
	}
	copyFiles := []string{"src", ".quasar", "quasar.conf.js"}
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
	configString := strings.Replace(req.FeatureVersionConfigString ,`\"` , `"` , -1)
	configString = strings.Replace(configString ,`"{` , `{` , -1)
	configString = strings.Replace(configString ,`}"` , `}` , -1)
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

	fileName := fmt.Sprintf("./packages/%s-%s.zip", req.FeatureName, req.Version.FeatureVersionName)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("package: "+fileName+" not found"))
		return
	}
	var cmds []*exec.Cmd
	//baseDir := "./install/" + strconv.Itoa(int(req.ProjectId))
	//
	//if _, err := os.Stat(baseDir); os.IsNotExist(err) {
	//	cmds = append(cmds, exec.Command("mkdir", baseDir))
	//}
	featureName := fmt.Sprintf("%s-%s-%s", req.FeatureName, req.Version.FeatureVersionName, req.Type)

	//cmds = append(cmds, exec.Command("unzip", "-o", fileName, "-d", featureDir))

	installDir := ""
	if req.FeatureOnBoot {
		installDir = fmt.Sprintf("%s/src/boot/%s", projectDir, featureName)
	} else {
		uuid, err := uuid.NewUUID()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}
		installDir = fmt.Sprintf("%s/src/components/%s-%s", projectDir, featureName, uuid.String()[:8])
	}
	cmds = append(cmds, exec.Command("unzip", "-o", fileName, "-d", installDir))
	_, stderr, err := utils.Pipeline(cmds...)
	if err != nil && !strings.HasPrefix(err.Error(), "exit") {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}

	if len(stderr) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(string(stderr)))
		return
	}

	for _, v := range req.Version.FeatureVersionConfig.Data {
		t, err := template.ParseFiles(installDir + "/" + v.Template + ".tmpl")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}
		buf := new(bytes.Buffer)
		err = t.Execute(buf, v.Values)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}
		b, err := ioutil.ReadFile(installDir + "/" + v.Target)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}
		newTargetFileString := strings.Replace(string(b), "__data."+v.Template+"__", buf.String(), 1)
		newTargetFileString = strings.Replace(newTargetFileString, "&#34;", `"`, -1)
		err = ioutil.WriteFile(installDir+"/"+v.Target, []byte(newTargetFileString), os.ModePerm)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}
	}
	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("GRPC connection lost"))
		return
	}
	defer utils.DbServerGrpcConn.Put(dbServerConn)
	client := pb.NewApiClient(dbServerConn)
	if req.FeatureOnBoot {
		bootString := []string{`''`}

		reply, err := client.ProjectFeaturesByProjectId(context.Background(), &pb.ProjectFeaturesByProjectIdRequest{
			ProjectId: req.ProjectId,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}
		bootString = bootString[:0]

		for _, v := range reply.Features {
			if v.FeatureOnboot {
				bootString = append(bootString, fmt.Sprintf(`'%s-%s-%s'`, v.FeatureName, v.FeatureVersionName, v.ProjectFeaturesType))
			}
		}
		bootString = append(bootString, `'`+featureName+`'`)
		quasarConfFileByte, err := ioutil.ReadFile("./packages/quasar.conf.js")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}
		newQuasarConfFileString := strings.Replace(string(quasarConfFileByte), "__data.boot__", strings.Join(bootString, ","), 1)
		err = ioutil.WriteFile(projectDir+"/quasar.conf.js", []byte(newQuasarConfFileString), os.ModePerm)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}
	}

	reply, err := client.CreateProjectFeature(context.Background(), &pb.CreateProjectFeatureRequest{
		FeatureId:            req.FeatureId,
		ProjectFeatureType:   req.Type,
		ProjectFeatureConfig: configString,
		ProjectId:            req.ProjectId,
		FeatureVersionId:     req.Version.FeatureVersionId,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}

	fmt.Printf("%+v",reply)

	c.AbortWithStatus(http.StatusNoContent)
}
