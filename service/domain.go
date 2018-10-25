/**
  @Author : hanxiaodong
*/
package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"fmt"
	"time"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type Education struct {
	ObjectType	string	`json:"docType"`
	Name	string	`json:"Name"`		// 姓名
	Gender	string	`json:"Gender"`		// 性别
	Nation	string	`json:"Nation"`		// 民族
	EntityID	string	`json:"EntityID"`		// 身份证号
	Place	string	`json:"Place"`		// 籍贯
	BirthDay	string	`json:"BirthDay"`		// 出生日期
	EnrollDate	string	`json:"EnrollDate"`		// 入学日期
	GraduationDate	string	`json:"GraduationDate"`	// 毕（结）业日期
	SchoolName	string	`json:"SchoolName"`	// 学校名称
	Major	string	`json:"Major"`	// 专业
	QuaType	string	`json:"QuaType"`	// 学历类别
	Length	string	`json:"Length"`	// 学制
	Mode	string	`json:"Mode"`	// 学习形式
	Level	string	`json:"Level"`	// 层次
	Graduation	string	`json:"Graduation"`	// 毕（结）业
	CertNo	string	`json:"CertNo"`	// 证书编号

	Photo	string	`json:"Photo"`	// 照片

	Historys	[]HistoryItem	// 当前edu的历史记录
}

type HistoryItem struct {
	TxId	string
	Education	Education
}

type ServiceSetup struct {
	ChaincodeID	string
	Client	*channel.Client
}

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}
