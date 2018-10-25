/**
  @Author : hanxiaodong
*/

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"github.com/hyperledger/fabric/protos/peer"
)

type EducationChaincode struct {

}

func (t *EducationChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response{

	return shim.Success(nil)
}

func (t *EducationChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response{
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()

	if fun == "addEdu"{
		return t.addEdu(stub, args)		// 添加信息
	}else if fun == "queryEduByCertNoAndName" {
		return t.queryEduByCertNoAndName(stub, args)		// 根据证书编号及姓名查询信息
	}else if fun == "queryEduInfoByEntityID" {
		return t.queryEduInfoByEntityID(stub, args)	// 根据身份证号码及姓名查询详情
	}else if fun == "updateEdu" {
		return t.updateEdu(stub, args)		// 根据证书编号更新信息
	}else if fun == "delEdu"{
		return t.delEdu(stub, args)	// 根据证书编号删除信息
	}

	return shim.Error("指定的函数名称错误")

}

func main(){
	err := shim.Start(new(EducationChaincode))
	if err != nil{
		fmt.Printf("启动EducationChaincode时发生错误: %s", err)
	}
}
