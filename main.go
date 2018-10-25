/**
  author: kevin
 */

package main

import (
	"os"
	"fmt"
	"github.com/kongyixueyuan.com/education/sdkInit"
	"github.com/kongyixueyuan.com/education/service"
	"encoding/json"
	"github.com/kongyixueyuan.com/education/web/controller"
	"github.com/kongyixueyuan.com/education/web"
)

const (
	configFile = "config.yaml"
	initialized = false
	EduCC = "educc"
)

func main() {

	initInfo := &sdkInit.InitInfo{

		ChannelID: "kevinkongyixueyuan",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/kongyixueyuan.com/education/fixtures/artifacts/channel.tx",

		OrgAdmin:"Admin",
		OrgName:"Org1",
		OrdererOrgName: "orderer.kevin.kongyixueyuan.com",

		ChaincodeID: EduCC,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath: "github.com/kongyixueyuan.com/education/chaincode/",
		UserName:"User1",
	}

	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer sdk.Close()

	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	channelClient, err := sdkInit.InstallAndInstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(channelClient)

	//===========================================//

	serviceSetup := service.ServiceSetup{
		ChaincodeID:EduCC,
		Client:channelClient,
	}

	edu := service.Education{
		Name: "张三",
		Gender: "男",
		Nation: "汉",
		EntityID: "101",
		Place: "北京",
		BirthDay: "1991年01月01日",
		EnrollDate: "2009年9月",
		GraduationDate: "2013年7月",
		SchoolName: "中国政法大学",
		Major: "社会学",
		QuaType: "普通",
		Length: "四年",
		Mode: "普通全日制",
		Level: "本科",
		Graduation: "毕业",
		CertNo: "111",
		Photo: "/static/photo/11.png",
	}

	edu2 := service.Education{
		Name: "李四",
		Gender: "男",
		Nation: "汉",
		EntityID: "102",
		Place: "上海",
		BirthDay: "1992年02月01日",
		EnrollDate: "2010年9月",
		GraduationDate: "2014年7月",
		SchoolName: "中国人民大学",
		Major: "行政管理",
		QuaType: "普通",
		Length: "四年",
		Mode: "普通全日制",
		Level: "本科",
		Graduation: "毕业",
		CertNo: "222",
		Photo: "/static/photo/22.png",
	}

	msg, err := serviceSetup.SaveEdu(edu)
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}

	msg, err = serviceSetup.SaveEdu(edu2)
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}

	// 根据证书编号与名称查询信息
	result, err := serviceSetup.FindEduByCertNoAndName("222","李四")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var edu service.Education
		json.Unmarshal(result, &edu)
		fmt.Println("根据证书编号与姓名查询信息成功：")
		fmt.Println(edu)
	}

	// 根据身份证号码查询信息
	result, err = serviceSetup.FindEduInfoByEntityID("101")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var edu service.Education
		json.Unmarshal(result, &edu)
		fmt.Println("根据身份证号码查询信息成功：")
		fmt.Println(edu)
	}

	// 修改/添加信息
	info := service.Education{
		Name: "张三",
		Gender: "男",
		Nation: "汉",
		EntityID: "101",
		Place: "北京",
		BirthDay: "1991年01月01日",
		EnrollDate: "2013年9月",
		GraduationDate: "2015年7月",
		SchoolName: "中国政法大学",
		Major: "社会学",
		QuaType: "普通",
		Length: "两年",
		Mode: "普通全日制",
		Level: "研究生",
		Graduation: "毕业",
		CertNo: "333",
		Photo: "/static/photo/11.png",
	}
	msg, err = serviceSetup.ModifyEdu(info)
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("信息操作成功, 交易编号为: " + msg)
	}

	// 根据身份证号码查询信息
	result, err = serviceSetup.FindEduInfoByEntityID("101")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var edu service.Education
		json.Unmarshal(result, &edu)
		fmt.Println("根据身份证号码查询信息成功：")
		fmt.Println(edu)
	}

	// 根据证书编号与名称查询信息
	result, err = serviceSetup.FindEduByCertNoAndName("333","张三")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var edu service.Education
		json.Unmarshal(result, &edu)
		fmt.Println("根据证书编号与姓名查询信息成功：")
		fmt.Println(edu)
	}

	/*// 删除信息
	msg, err = serviceSetup.DelEdu("101")
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("信息删除成功, 交易编号为: " + msg)
	}

	// 根据身份证号码查询信息
	result, err = serviceSetup.FindEduInfoByEntityID("101")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("根据身份证号码查询信息失败，指定身份证号码的信息不存在或已被删除...")
	} else {
		var edu service.Education
		json.Unmarshal(result, &edu)
		fmt.Println("根据身份证号码查询信息成功：")
		fmt.Println(edu)
	}*/

	//===========================================//

	app := controller.Application{
		Setup: &serviceSetup,
	}
	web.WebStart(app)

}
