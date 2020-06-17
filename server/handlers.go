package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/chenyu116/generator-mobile/config"
	pb "github.com/chenyu116/generator-mobile/proto"
	"github.com/chenyu116/generator-mobile/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type result struct {
	Error string `json:"error,omitempty"`
}

func jsonError(err string) result {
	return result{
		Error: err,
	}
}

var previewList = make(map[int32]chan *exec.Cmd)

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
	baseDir := "./projects/" + req.ProjectId
	var cmds []*exec.Cmd

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		cmds = append(cmds, exec.Command("mkdir", "-p", baseDir+"/src/boot"))
		cmds = append(cmds, exec.Command("mkdir", "-p", baseDir+"/src/components"))
		cmds = append(cmds, exec.Command("mkdir", "-p", baseDir+"/src/pages"))
	}
	copyFiles := []string{"src/assets", "src/css", "src/plugins", "src/router", "src/statics", "/src/store", "/src/App.vue", "/src/index.template.html", "src/boot/i18n", "src/boot/iBeacon", "src/boot/preload", "src/boot/process", "src/boot/mapUtil", "src/boot/weixinJssdk", "src/components/HeaderWithBack.vue", "src/pages/Error404.vue", ".quasar", "quasar.conf.js"}
	for _, v := range copyFiles {
		if _, err := os.Stat(baseDir + "/" + v); os.IsNotExist(err) {
			cmds = append(cmds, exec.Command("cp", "-r", quasarDir+"/"+v, baseDir+"/"+v))
		}
	}
	linkFiles := []string{"node_modules", "babel.config.js", "jsconfig.json", "package.json", "yarn.lock", ".eslintignore", ".eslintrc.js", ".postcssrc.js"}
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
	c.AbortWithStatus(http.StatusCreated)
}
func edit(c *gin.Context) {
	var req RequestPostEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	if req.ProjectFeaturesId <= 0 {
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

	configByte, err := json.Marshal(req.Version.FeatureVersionConfig)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}

	_, err = client.UpdateProjectFeature(context.Background(), &pb.UpdateProjectFeatureRequest{
		ProjectFeaturesId:     req.ProjectFeaturesId,
		ProjectFeaturesType:   req.Type,
		ProjectFeaturesConfig: string(configByte),
		FeatureVersionId:      req.Version.FeatureVersionId,
		ProjectFeaturesName:   req.ProjectFeaturesName,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

func sendMessage(str string) {
	utils.RabbitMQ.Publish("websocketServer-fanout", "", amqp.Publishing{
		AppId: "all",
		Body:  []byte(str),
	}, nil)
}
func cancelPreview(c *gin.Context) {
	var req RequestInt32UriProjectId
	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	if req.ProjectId <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("invalid params"))
		return
	}

	if ch, ok := previewList[req.ProjectId]; ok {
		cmd := <-ch
		err := cmd.Process.Kill()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}
		delete(previewList, req.ProjectId)
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.AbortWithStatus(http.StatusNotFound)
}
func preview(c *gin.Context) {
	var req RequestPutPreview
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	if req.ProjectId <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("invalid params"))
		return
	}

	projectDir := fmt.Sprintf("/home/roger/workspace/generator-mobile/projects/%d", req.ProjectId)
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("project not initialized"))
		return
	}

	debugInfo := paramsPreviewIndex{
		DebugInfo: strings.Replace(req.Settings, "\n", "", -1),
	}
	fmt.Printf("%s", debugInfo.DebugInfo)
	indexT, err := template.ParseFiles("./packages/index.html.tmpl")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	buf := new(bytes.Buffer)
	err = indexT.Execute(buf, debugInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	indexString := strings.Replace(buf.String(), `"{`, "{", 1)
	indexString = strings.Replace(indexString, `}"`, "}", 1)
	indexString = strings.Replace(indexString, `\"`, `"`, -1)
	indexString = strings.Replace(indexString, `\'`, `'`, -1)
	indexString = strings.Replace(indexString, `\t`, ``, -1)
	buf.Reset()
	err = ioutil.WriteFile(projectDir+"/src/index.template.html", []byte(indexString), os.ModePerm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}

	go func() {
		defer func() {
			if _, ok := previewList[req.ProjectId]; ok {
				delete(previewList, req.ProjectId)
			}
		}()
		ch := make(chan *exec.Cmd, 1)
		cmd := exec.Command("quasar", "dev")
		cmd.Dir = projectDir
		ch <- cmd
		previewList[req.ProjectId] = ch
		stdoutIn, _ := cmd.StdoutPipe()
		err := cmd.Start()
		if err != nil {
			return
		}
		go func() {
			buf := make([]byte, 1024, 1024)
			for {
				n, err := stdoutIn.Read(buf[:])
				if n > 0 {
					d := buf[:n]
					sendMessage(string(d))
				}
				if err != nil {
					// Read returns io.EOF at the end of file, which is not an error for us
					if err == io.EOF {
						err = nil
					}
					return
				}
			}
		}()
		err = cmd.Wait()
		if err != nil {
			return
		}
	}()
	c.AbortWithStatus(http.StatusNoContent)
}
func deploy(c *gin.Context) {
	var req RequestInt32ProjectId
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	if req.ProjectId <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("参数错误"))
		return
	}
	projectDir := fmt.Sprintf("/home/roger/workspace/generator-mobile/projects/%d", req.ProjectId)
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("项目未初始化"))
		return
	}

	cmd := exec.Command("quasar", "build")
	cmd.Dir = projectDir
	stdoutIn, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	go func() {
		buf := make([]byte, 1024, 1024)
		for {
			n, err := stdoutIn.Read(buf[:])
			if n > 0 {
				d := buf[:n]
				sendMessage(string(d))
			}
			if err != nil {
				// Read returns io.EOF at the end of file, which is not an error for us
				if err == io.EOF {
					err = nil
				}
				return
			}
		}
	}()
	err = cmd.Wait()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}

	distDir := fmt.Sprintf("/home/roger/workspace/generator-mobile/projects/%d/dist/spa/", req.ProjectId)
	if _, err := os.Stat(distDir); os.IsNotExist(err) {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("项目未编译"))
		return
	}
	if _, err := os.Stat(projectDir + "/buildTime.lock"); os.IsNotExist(err) {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("项目未编译"))
		return
	}

	buildTimeBytes, err := ioutil.ReadFile(projectDir + "/buildTime.lock")
	cf := config.GetConfig()
	client, err := oss.New(cf.Oss.EndPoint, cf.Oss.AccessKeyID, cf.Oss.AccessKeySecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}

	// Get Bucket
	bucket, err := client.Bucket(cf.Oss.BucketName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	objectPrefix := fmt.Sprintf("project/%d/%s", req.ProjectId, buildTimeBytes)
	err = filepath.Walk(distDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || info.Name() == "buildTime.lock" {
			return nil
		}
		objectName := objectPrefix + "/" + strings.Replace(path, distDir, "", 1)
		sendMessage("sending " + objectName)
		err = bucket.PutObjectFromFile(objectName, path, oss.Origin("*"), oss.ACL(oss.ACLPublicRead))
		return err
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	sendMessage("clear")
	time.Sleep(time.Second)
	indexBytes, err := ioutil.ReadFile(distDir + "index.html")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	sendMessage(string(indexBytes))
	c.AbortWithStatus(http.StatusCreated)
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
	sendMessage("check project...")
	projectDir := fmt.Sprintf("/home/roger/workspace/generator-mobile/projects/%d", req.ProjectId)
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("project not initialized"))
		return
	}

	var debugInfo paramsPreviewIndex
	indexT, err := template.ParseFiles("./packages/index.html.tmpl")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	buf := new(bytes.Buffer)
	err = indexT.Execute(buf, debugInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	indexString := strings.Replace(buf.String(), `"{`, "{", 1)
	indexString = strings.Replace(indexString, `}"`, "}", 1)
	indexString = strings.Replace(indexString, `\"`, `"`, -1)
	indexString = strings.Replace(indexString, `\'`, `'`, -1)
	indexString = strings.Replace(indexString, `\n`, ``, -1)
	indexString = strings.Replace(indexString, `&lt;`, `<`, -1)
	indexString = strings.Replace(indexString, `"'`, `'`, -1)
	indexString = strings.Replace(indexString, `'"`, `'`, -1)
	indexString = strings.Replace(indexString, `&#39;`, `'`, -1)
	indexString = strings.Replace(indexString, `&#34;`, `"`, -1)
	buf.Reset()
	err = ioutil.WriteFile(projectDir+"/src/index.template.html", []byte(indexString), os.ModePerm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}

	sendMessage("check project ok")
	sendMessage("connect db server...")
	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("GRPC connection lost"))
		return
	}
	sendMessage("connect db server ok")
	sendMessage("request client.ProjectFeaturesByProjectId...")
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
	sendMessage("request client.ProjectFeaturesByProjectId ok")
	var cmds []*exec.Cmd
	var bootString []string
	routes := make(map[string]paramsRoutesJsRoutesParam)
	buf.Reset()
	for _, feature := range projectFeatures.Features {
		cmds = cmds[:0]
		installDir := ""
		var projectConfig featureVersionConfig
		err = json.Unmarshal([]byte(feature.ProjectFeaturesConfig), &projectConfig)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}
		sendMessage(fmt.Sprintf("generating component \"%s\"", feature.ProjectFeaturesInstallName))
		for cmpK, cmp := range projectConfig.Components {
			for cmpvK, cmpv := range cmp.Values {
				for _, f := range projectFeatures.Features {
					if cmpv.ProjectFeaturesId == f.ProjectFeaturesId {
						var fProjectConfig featureVersionConfig
						err = json.Unmarshal([]byte(f.ProjectFeaturesConfig), &fProjectConfig)
						if err != nil {
							c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
							return
						}
						projectConfig.Components[cmpK].Values[cmpvK].ProjectFeaturesConfig = fProjectConfig
						break
					}
				}
				hashIndex := strings.LastIndex(cmpv.ProjectFeaturesInstallName, "-")
				projectConfig.Components[cmpK].Values[cmpvK].ComponentHash = strings.Replace(cmpv.ProjectFeaturesInstallName[hashIndex:], "-", "C", -1)
			}
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
			_, _, err = utils.Pipeline(cmds...)
			if err != nil && !strings.HasPrefix(err.Error(), "exit") {
				c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
				return
			}
		}
		cmds = cmds[:0]

		if feature.ProjectFeaturesType != "click" {

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
				cmds = append(cmds, exec.Command("cp", packageDir+"/"+f.Name(), installDir))
			}
			if len(cmds) > 0 {
				_, _, err = utils.Pipeline(cmds...)
				if err != nil && !strings.HasPrefix(err.Error(), "exit") {
					c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
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
					cmds = append(cmds, exec.Command("/home/roger/.yarn/bin/prettier-eslint", "--config", projectDir+"/package.json", "--write", installDir+"/"+file))
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
		for _, cf := range projectConfig.Data.Values {
			if cf.FormType == "upload" {
				uploadFile, ok := cf.Value.(string)
				if !ok || uploadFile == "" {
					continue
				}
				if _, err := os.Stat(installDir + "/" + uploadFile); os.IsNotExist(err) {
					cmds = append(cmds, exec.Command("mv", "./tmp/"+uploadFile, installDir+"/"+uploadFile))
				}
			}
		}
		if len(cmds) > 0 {
			_, _, err := utils.Pipeline(cmds...)
			if err != nil && !strings.HasPrefix(err.Error(), "exit") {
				c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
				return
			}
		}
	}
	cmds = cmds[:0]
	sendMessage("generating quasar.conf.js")
	quasarT, err := template.ParseFiles("./packages/quasar.conf.js.tmpl")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	buf.Reset()
	buildTime := time.Now().Format("20060102150405")
	quasarConfig := paramsQuasarConfig{
		BootString: bootString,
		StaticDir:  fmt.Sprintf("%d/%s", req.ProjectId, buildTime),
	}
	err = quasarT.Execute(buf, quasarConfig)
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
	cmds = append(cmds, exec.Command("/home/roger/.yarn/bin/prettier-eslint", "--config", projectDir+"/package.json", "--write", projectDir+"/quasar.conf.js"))
	sendMessage("generating routes.js")
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
	cmds = append(cmds, exec.Command("/home/roger/.yarn/bin/prettier-eslint", "--config", projectDir+"/package.json", "--write", projectDir+"/src/router/routes.js"))
	if len(cmds) > 0 {
		_, _, err = utils.Pipeline(cmds...)
		if err != nil && !strings.HasPrefix(err.Error(), "exit") {
			c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
			return
		}
	}

	err = ioutil.WriteFile(projectDir+"/buildTime.lock", []byte(buildTime), os.ModePerm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	c.AbortWithStatus(http.StatusCreated)
}
