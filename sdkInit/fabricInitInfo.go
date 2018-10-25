/**
  author: kevin
 */

package sdkInit

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
)

type InitInfo struct {
	ChannelID     string
	ChannelConfig string
	OrgAdmin      string
	OrgName       string
	OrdererOrgName	string
	OrgResMgmt *resmgmt.Client

	ChaincodeID	string
	ChaincodeGoPath	string
	ChaincodePath	string
	UserName	string
}
