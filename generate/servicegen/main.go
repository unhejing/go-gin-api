package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/unhejing/go-gin-api/generate/servicegen/servicepkg"
)

type Data struct {
	StructName string
	TableName  string
	HumpName   string
	InputFile  string
}

var serviceData Data

func init() {
	service := flag.String("service", "", "请输入需要生成的 service 名称\n")
	flag.Parse()
	serviceData = Data{}
	serviceData.TableName = strings.ToLower(*service)
	serviceData.HumpName = servicepkg.SQLColumnToHumpStyle(*service)
	serviceData.StructName = servicepkg.Capitalize(*service)
	serviceData.InputFile = "service"

}

func main() {
	// 生成service
	serviceData.InputFile = "service"
	filepath := serviceData.InputFile + "/" + serviceData.TableName + "_service"
	err := os.MkdirAll(filepath, 0766)
	if err != nil {
		fmt.Println("Unable to create directory, ", err.Error())
		panic(err)
	}
	fmt.Println("create dir : ", filepath)

	outputFile, err := os.Create(filepath + "/service.go")
	defer outputFile.Close()
	if err != nil {
		panic(err)
	}

	if err := servicepkg.OutputTemplate.Execute(outputFile, serviceData); err != nil {
		panic(err)
	}

	// 生成handler
	serviceData.InputFile = "api"
	filepath = serviceData.InputFile + "/" + serviceData.TableName + "_handler"
	err = os.MkdirAll(filepath, 0766)
	if err != nil {
		fmt.Println("Unable to create directory, ", err.Error())
		panic(err)
	}
	fmt.Println("create dir : ", filepath)

	outputFile, err = os.Create(filepath + "/handler.go")
	defer outputFile.Close()
	if err != nil {
		panic(err)
	}
	if err := servicepkg.HandlerOutputTemplate.Execute(outputFile, serviceData); err != nil {
		panic(err)
	}

	// 生成dto
	serviceData.InputFile = "dto"
	filepath = serviceData.InputFile + "/" + serviceData.TableName + "_dto"
	err = os.MkdirAll(filepath, 0766)
	if err != nil {
		fmt.Println("Unable to create directory, ", err.Error())
		panic(err)
	}
	fmt.Println("create dir : ", filepath)

	mdName := fmt.Sprintf("%s/page_req.go", filepath)
	reqFile, err := os.OpenFile(mdName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
	if err != nil {
		fmt.Printf("markdown file error %v\n", err.Error())
		return
	}
	pageReq := fmt.Sprintf("package %s\n", serviceData.TableName+"_dto\n")
	pageReq += fmt.Sprintf("type PageReq struct {\n")
	pageReq += fmt.Sprintf("	Page int `json:\"page\"` // 第几页\n")
	pageReq += fmt.Sprintf("	Size int `json:\"size\"` // 每页显示条数\n")
	pageReq += fmt.Sprintf("}")
	reqFile.WriteString(pageReq)
	reqFile.Close()

	mdName = fmt.Sprintf("%s/page_res.go", filepath)
	resFile, err := os.OpenFile(mdName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
	if err != nil {
		fmt.Printf("pageReq file error %v\n", err.Error())
		return
	}
	pageRes := fmt.Sprintf("package %s\n", serviceData.TableName+"_dto\n")
	pageRes += fmt.Sprintf("import \"github.com/unhejing/go-gin-api/model/%s_model\"\n", serviceData.TableName)
	pageRes += fmt.Sprintf("type PageRes struct {\n")
	pageRes += fmt.Sprintf("	List       []%s_model.%s `json:\"list\"`\n", serviceData.TableName, serviceData.StructName)
	pageRes += fmt.Sprintf("	Pagination struct {\n")
	pageRes += fmt.Sprintf("		Total        int `json:\"total\"`\n")
	pageRes += fmt.Sprintf("		CurrentPage  int `json:\"current_page\"`\n")
	pageRes += fmt.Sprintf("		PerPageCount int `json:\"per_page_count\"`\n")
	pageRes += fmt.Sprintf("	} `json:\"pagination\"`\n")
	pageRes += fmt.Sprintf("}")
	resFile.WriteString(pageRes)
	resFile.Close()
}
