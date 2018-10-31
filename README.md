[TOC]

# 1. 需求分析与架构设计

## 1.1 需求分析

现在是一个信息化的高科技时代，许许多多的企业必须紧跟时代步伐，不断创新，才能发展壮大；而企业的发展必然离不开人才队伍的建设，也可以说创新是企业发展的动力，而人才却是企业发展的根本，所以现在各企业对于人才队伍建设十分看重，而对于人才的素质及受教育情况的要求更是重中之重。

对学历信息的查询，要么成本较高，要么比较麻烦，甚至还有一些假冒网站让人防不胜防；传统应用是将数据保存在数据库中来实现，但是现在出现的数据库由于故障或者被删、被黑造成的数据丢失的情况更是屡见不鲜，所以传统数据库并不能真正意义上确保数据的完整性及安全性。

基于这些情况，我们设计并开发了一个 `基于区块链技术的实现的学历信息征信系统`，实现了在线对学历信息的查询功能，由于区块链技术本身的特点，无须考虑数据被破坏的问题，而且杜绝了对于信息造假的情况，保证了学历信息的真实性。由于篇幅原因，我们对学历信息征信系统的应用场景进行修改及简化，实现的业务逻辑包括添加信息、修改信息、查询信息、查询详情信息等操作，实际情况下的的业务逻辑需要根据实际需求场景做出相应的调整。

由于系统需要保证人才受教育情况真实性，所以对于系统的用户而言，不可能由用户自己添加相应的学历信息，而是由具有一定权限的用户来完成添加或修改的功能。但普通用户可以通过系统溯源功能来确定信息的真伪。所以我们将系统用户的使用角色分为两种：

1. 普通用户
2. 管理员用户

普通用户具有对数据的查询功能 ，但实现查询之前必须经过登录认证：

- 用户登录：系统只针对合法用户进行授权使用，所以用户必须先进行登录才能完成相应的功能。
- 查询实现：查询分为两种方式实现
  - 根据证书编号与姓名查询：根据用户输入的证书编号与姓名进行查询。
  - 根据身份证号码查询：根据用户输入指定的身份证号码进行查询，此功能可以实现溯源。

管理员用户除具有普通用户的功能之外，额外添加了两个功能：

- 添加信息：可以向系统中添加新的学历信息。
- 修改信息：针对已存在的学历信息进行修改。

## 1.2 架构设计

我们在 [从零到壹构建基于 Fabric-SDK-Go 的Web项目实战](https://github.com/kevin-hf/kongyixueyuan) 中已经完成了一个完整的基于 `fabric-sdk-go` 的应用示例，所以我们现在使用之前的应用架构，不同的是在此应用中需要编写实现完整的链码并通过业务层调用链码中的各个函数，以实现对数据状态的操作。界面为了方便用户操作使用，仍然使用Web浏览器的方式实现。而且在此应用中我们将 `Hyperledger Fabric` 默认的状态数据库由 `LevelDB` 替换为 `CouchDB` 来实现



![架构](./img/projectArch.png)

对于 `Fabric Network`结构如下图所示：

![networkArch](./img/networkArch.png)

## 1.3 数据模型设计

由于需要向分类账本中保存数据，所以必须设计相关的结构体用于声明要保存的数据结构，用于方便的在应用中处理数据。

`Education` 结构体设计如下表所示：

| 名称           | 数据类型      | 说明                             |
| -------------- | ------------- | -------------------------------- |
| ObjectType     | string        |                                  |
| Name           | string        | 姓名                             |
| Gender         | string        | 性别                             |
| Nation         | string        | 民族                             |
| EntityID       | string        | 身份证号（记录的Key）            |
| Place          | string        | 籍贯                             |
| BirthDay       | string        | 出生日期                         |
| Photo          | string        | 照片                             |
| EnrollDate     | string        | 入学日期                         |
| GraduationDate | string        | 毕（结）业日期                   |
| SchoolName     | string        | 所读学校名称                     |
| Major          | string        | 所读专业                         |
| QuaType        | string        | 学历类别（普通、成考等）         |
| Length         | string        | 学制（两年、三年、四年、五年）   |
| Mode           | string        | 学习形式（普通全日制）           |
| Level          | string        | 层次（专科、本科、研究生、博士） |
| Graduation     | string        | 毕（结）业（毕业、结业）         |
| CertNo         | string        | 证书编号                         |
| Historys       | []HistoryItem | 当前edu的详细历史记录            |

为了能够从当前的分类状态中查询出详细的历史操作记录，我们在 `Education` 中设计了一个类型为`HistoryItem` 数组的 `Historys` 成员，表示当前状态的历史记录集。

`HistoryItem` 结构体设计如下表所示：

| 名称      | 数据类型  | 说明                   |
| --------- | --------- | ---------------------- |
| TxId      | string    | 交易编号               |
| Education | Education | 本次历史记录的详细信息 |



## 1.4 网络环境

### 1.4.1 设置环境

我们在 [从零到壹构建基于 Fabric-SDK-Go 的Web项目实战](https://github.com/kevin-hf/kongyixueyuan) （如果没有看过《从零到壹构建基于 Fabric-SDK-Go 的Web项目实战》，建议先将其学习实现，然后再学习此案例）说明了如何构建fabric网络环境，现在我们要重新完成一个新的应用，所以网络环境可以使用之前的内容，但是因为**状态数据库使用 `CouchDB` 来实现**，所以需要做出部分修改，新增与 `CouchDB` 相关的内容。为了方便起见，我们重新搭建一个应用所需的网络环境。

在`GOPATH`的`src`文件夹中新建一个目录如下：

```shell
$ mkdir -p $GOPATH/src/github.com/kongyixueyuan.com/education 
$ cd $GOPATH/src/github.com/kongyixueyuan.com/education
```

使用 `git` 命令克隆 hf-fixtures 目录当前路径

```shell
$ git clone https://github.com/kevin-hf/hf-fixtures.git
```

将 hf-fixtures 文件夹重命名为 fixtures

```shell
$ mv hf-fixtures/fixtures
```

修改`fixtures`  文件夹的所属关系为当前用户

```shell
$ sudo chown -R kevin:kevin ./fixtures
```

> 提示： kevin 为安装 Ubuntu 16.04 系统时创建的用户

进入 `fixtures` 目录

```shell
$ cd fixtures
```

为了构建区块链网络，使用 `docker` 构建处理不同角色的虚拟计算机。 在这里我们将尽可能保持简单。如果确定您的系统中已经存在相关的所需容器，或可以使用其它方式获取，则无需执行如下命令。否则请将 `fixtures` 目录下的 `pull_images.sh` 文件添加可执行权限后直接执行。

```shell
$ chmod 777 ./pull_images.sh
$ ./pull_images.sh 
```

> 提示：`pull_images.sh` 文件是下载 Fabric 环境所需容器的一个可执行脚本，下载过程需要一段时间（视网速情况而定），请耐心等待。另：请确定您的系统支持虚拟技术。

### 1.4.2 配置docker-compose.yml文件

在 `fixtures`  目录下创建一个 `docker-compose.yml`  文件并编辑

```shell
$ vim docker-compose.yml
```

1. 将  `network下的basic`   修改为  `default`

   ```yaml
   version: '2'
   
   networks:
     default:
   
   services:
   ```

2. 编辑  orderer 部分

   ```yaml
     orderer.kevin.kongyixueyuan.com:
       image: hyperledger/fabric-orderer
       container_name: orderer.kevin.kongyixueyuan.com
       environment:
         - ORDERER_GENERAL_LOGLEVEL=debug
         - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
         - ORDERER_GENERAL_LISTENPORT=7050
         - ORDERER_GENERAL_GENESISPROFILE=kongyixueyuan
         - ORDERER_GENERAL_GENESISMETHOD=file
         - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/genesis.block
         - ORDERER_GENERAL_LOCALMSPID=kevin.kongyixueyuan.com
         - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
         - ORDERER_GENERAL_TLS_ENABLED=true
         - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
         - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
         - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
       working_dir: /opt/gopath/src/github.com/hyperledger/fabric
       command: orderer
       volumes:
         - ./artifacts/genesis.block:/var/hyperledger/orderer/genesis.block
         - ./crypto-config/ordererOrganizations/kevin.kongyixueyuan.com/orderers/orderer.kevin.kongyixueyuan.com/msp:/var/hyperledger/orderer/msp
         - ./crypto-config/ordererOrganizations/kevin.kongyixueyuan.com/orderers/orderer.kevin.kongyixueyuan.com/tls:/var/hyperledger/orderer/tls
       ports:
         - 7050:7050
       networks:
         default:
           aliases:
             - orderer.kevin.kongyixueyuan.com
   ```

3. 编辑 ca 部分

   ```yaml
     ca.org1.kevin.kongyixueyuan.com:
       image: hyperledger/fabric-ca
       container_name: ca.org1.kevin.kongyixueyuan.com
       environment:
         - FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server
         - FABRIC_CA_SERVER_CA_NAME=ca.org1.kevin.kongyixueyuan.com
         - FABRIC_CA_SERVER_CA_CERTFILE=/etc/hyperledger/fabric-ca-server-config/ca.org1.kevin.kongyixueyuan.com-cert.pem
         - FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/727e69ed4a01a204cd53bf4a97c2c1cb947419504f82851f6ae563c3c96dea3a_sk
         - FABRIC_CA_SERVER_TLS_ENABLED=true
         - FABRIC_CA_SERVER_TLS_CERTFILE=/etc/hyperledger/fabric-ca-server-config/ca.org1.kevin.kongyixueyuan.com-cert.pem
         - FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server-config/727e69ed4a01a204cd53bf4a97c2c1cb947419504f82851f6ae563c3c96dea3a_sk
       ports:
         - 7054:7054
       command: sh -c 'fabric-ca-server start -b admin:adminpw -d'
       volumes:
         - ./crypto-config/peerOrganizations/org1.kevin.kongyixueyuan.com/ca/:/etc/hyperledger/fabric-ca-server-config
       networks:
         default:
           aliases:
             - ca.org1.kevin.kongyixueyuan.com
   ```

4. 声明 CouchDB 部分：

   ```yaml
     couchdb:
       container_name: couchdb
       image: hyperledger/fabric-couchdb
       # Populate the COUCHDB_USER and COUCHDB_PASSWORD to set an admin user and password
       # for CouchDB.  This will prevent CouchDB from operating in an "Admin Party" mode.
       environment:
         - COUCHDB_USER=
         - COUCHDB_PASSWORD=
       # Comment/Uncomment the port mapping if you want to hide/expose the CouchDB service,
       # for example map it to utilize Fauxton User Interface in dev environments.
       ports:
         - "5984:5984"
   ```

5. 编辑Peer部分

   1. `peer0.org1.example.com`  内容如下

      ```yaml
        peer0.org1.kevin.kongyixueyuan.com:
          image: hyperledger/fabric-peer
          container_name: peer0.org1.kevin.kongyixueyuan.com
          environment:
            - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
            - CORE_VM_DOCKER_ATTACHSTDOUT=true
            - CORE_LOGGING_LEVEL=DEBUG
            - CORE_PEER_NETWORKID=kongyixueyuan
            - CORE_PEER_PROFILE_ENABLED=true
            - CORE_PEER_TLS_ENABLED=true
            - CORE_PEER_TLS_CERT_FILE=/var/hyperledger/tls/server.crt
            - CORE_PEER_TLS_KEY_FILE=/var/hyperledger/tls/server.key
            - CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/tls/ca.crt
            - CORE_PEER_ID=peer0.org1.kevin.kongyixueyuan.com
            - CORE_PEER_ADDRESSAUTODETECT=true
            - CORE_PEER_ADDRESS=peer0.org1.kevin.kongyixueyuan.com:7051
            - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org1.kevin.kongyixueyuan.com:7051
            - CORE_PEER_GOSSIP_USELEADERELECTION=true
            - CORE_PEER_GOSSIP_ORGLEADER=false
            - CORE_PEER_GOSSIP_SKIPHANDSHAKE=true
            - CORE_PEER_LOCALMSPID=org1.kevin.kongyixueyuan.com
            - CORE_PEER_MSPCONFIGPATH=/var/hyperledger/msp
            - CORE_PEER_TLS_SERVERHOSTOVERRIDE=peer0.org1.kevin.kongyixueyuan.com
            - CORE_LEDGER_STATE_STATEDATABASE=CouchDB
            - CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb:5984
            - CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=
            - CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=
          working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
          command: peer node start
          volumes:
            - /var/run/:/host/var/run/
            - ./crypto-config/peerOrganizations/org1.kevin.kongyixueyuan.com/peers/peer0.org1.kevin.kongyixueyuan.com/msp:/var/hyperledger/msp
            - ./crypto-config/peerOrganizations/org1.kevin.kongyixueyuan.com/peers/peer0.org1.kevin.kongyixueyuan.com/tls:/var/hyperledger/tls
          ports:
            - 7051:7051
            - 7053:7053
          depends_on:
            - orderer.kevin.kongyixueyuan.com
            - couchdb
          networks:
            default:
              aliases:
                - peer0.org1.kevin.kongyixueyuan.com
      ```

   2. peer1.org1.example.com 内容如下

      ```yaml
        peer1.org1.kevin.kongyixueyuan.com:
          image: hyperledger/fabric-peer
          container_name: peer1.org1.kevin.kongyixueyuan.com
          environment:
            - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
            - CORE_VM_DOCKER_ATTACHSTDOUT=true
            - CORE_LOGGING_LEVEL=DEBUG
            - CORE_PEER_NETWORKID=kongyixueyuan
            - CORE_PEER_PROFILE_ENABLED=true
            - CORE_PEER_TLS_ENABLED=true
            - CORE_PEER_TLS_CERT_FILE=/var/hyperledger/tls/server.crt
            - CORE_PEER_TLS_KEY_FILE=/var/hyperledger/tls/server.key
            - CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/tls/ca.crt
            - CORE_PEER_ID=peer1.org1.kevin.kongyixueyuan.com
            - CORE_PEER_ADDRESSAUTODETECT=true
            - CORE_PEER_ADDRESS=peer1.org1.kevin.kongyixueyuan.com:7051
            - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer1.org1.kevin.kongyixueyuan.com:7051
            - CORE_PEER_GOSSIP_USELEADERELECTION=true
            - CORE_PEER_GOSSIP_ORGLEADER=false
            - CORE_PEER_GOSSIP_SKIPHANDSHAKE=true
            - CORE_PEER_LOCALMSPID=org1.kevin.kongyixueyuan.com
            - CORE_PEER_MSPCONFIGPATH=/var/hyperledger/msp
            - CORE_PEER_TLS_SERVERHOSTOVERRIDE=peer1.org1.kevin.kongyixueyuan.com
            - CORE_LEDGER_STATE_STATEDATABASE=CouchDB
            - CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb:5984
            - CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=
            - CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=
          working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
          command: peer node start
          volumes:
            - /var/run/:/host/var/run/
            - ./crypto-config/peerOrganizations/org1.kevin.kongyixueyuan.com/peers/peer1.org1.kevin.kongyixueyuan.com/msp:/var/hyperledger/msp
            - ./crypto-config/peerOrganizations/org1.kevin.kongyixueyuan.com/peers/peer1.org1.kevin.kongyixueyuan.com/tls:/var/hyperledger/tls
          ports:
            - 7151:7051
            - 7153:7053
          depends_on:
            - orderer.kevin.kongyixueyuan.com
            - couchdb
          networks:
            default:
              aliases:
                - peer1.org1.kevin.kongyixueyuan.com
      ```

## 1.5 测试网络环境

为了检查网络是否正常工作，使用`docker-compose`同时启动或停止所有容器。 进入`fixtures`文件夹，运行：

```shell
$ cd $GOPATH/src/github.com/kongyixueyuan.com/education/fixtures
$ docker-compose up
```

控制台会输出很多不同颜色的日志（红色不等于错误）

![启动网络](./img/dockercompose_up.png)

打开一个新终端并运行：

```shell
 $ docker ps 
```

![查看容器](./img/docker_ps.png)

将看到：两个peer，一个orderer和一个CA容器，还有一个 CouchDB 容器。 代表已成功创建了一个新的网络，可以随SDK一起使用。 要停止网络，请返回到上一个终端，按`Ctrl+C`并等待所有容器都停止。 

> **提示** ：当网络成功启动后，所有处于活动中的容器都可以访问。 也可以查看指定容器的详细日志内容。  如果想删除这些容器，需要使用`docker rm $(docker ps -aq)`将其删除 ，但在删除容器之前需要确定其在网络环境中已不再使用。
>
> 如果在网络环境启动过程中不想看到大量的日志信息，请在该启动命令中添加参数 `-d` ，如下所示： `docker-compose up -d` 。 如果要停止网络，请务必在 `docker-compose.yaml` 所在的文件夹中运行命令： `docker-compose stop` （或 使用`docker-compose down` 进行清理停止所有容器）。

最后在终端2中执行如下命令关闭网络：

```shell
$ cd $GOPATH/src/github.com/kongyixueyuan.com/education/fixtures
$ docker-compose down
```

![关闭网络](./img/dockercompose_down.png)



终端1窗口中输出如下：

![关闭网络过程](./img/dockercompose_down2.png)



# 2 SDK与链码实现

## 2.1 创建 config.yaml 文件

确认 Hyperledger Fabric 基础网络环境运行没有问题后，现在我们通过创建一个新的 config.yaml 配置文件给应用程序所使用的 Fabric-SDK-Go 配置相关参数及 Fabric 组件的通信地址

进入项目的根目录中创建一个 `config.yaml` 文件并编辑

```shell
$ cd $GOPATH/src/github.com/kongyixueyuan.com/education
$ vim config.yaml
```

config.yaml 配置文件完整内容如下:

```yaml
name: "kongyixueyuan-network"
#
# Schema version of the content. Used by the SDK to apply the corresponding parsing rules.
#
version: 1.0.0

#
# The client section used by GO SDK.
#
client:

  # Which organization does this application instance belong to? The value must be the name of an org
  # defined under "organizations"
  organization: Org1

  logging:
    level: info

  # Global configuration for peer, event service and orderer timeouts
  # if this this section is omitted, then default values will be used (same values as below)
#  peer:
#    timeout:
#      connection: 10s
#      response: 180s
#      discovery:
#        # Expiry period for discovery service greylist filter
#        # The channel client will greylist peers that are found to be offline
#        # to prevent re-selecting them in subsequent retries.
#        # This interval will define how long a peer is greylisted
#        greylistExpiry: 10s
#  eventService:
#    # Event service type (optional). If not specified then the type is automatically
#    # determined from channel capabilities.
#    type: (deliver|eventhub)
    # the below timeouts are commented out to use the default values that are found in
    # "pkg/fab/endpointconfig.go"
    # the client is free to override the default values by uncommenting and resetting
    # the values as they see fit in their config file
#    timeout:
#      connection: 15s
#      registrationResponse: 15s
#  orderer:
#    timeout:
#      connection: 15s
#      response: 15s
#  global:
#    timeout:
#      query: 180s
#      execute: 180s
#      resmgmt: 180s
#    cache:
#      connectionIdle: 30s
#      eventServiceIdle: 2m
#      channelConfig: 30m
#      channelMembership: 30s
#      discovery: 10s
#      selection: 10m

  # Root of the MSP directories with keys and certs.
  cryptoconfig:
    path: ${GOPATH}/src/github.com/kongyixueyuan.com/education/fixtures/crypto-config

  # Some SDKs support pluggable KV stores, the properties under "credentialStore"
  # are implementation specific
  credentialStore:
    path: /tmp/kongyixueyuan-store

    # [Optional]. Specific to the CryptoSuite implementation used by GO SDK. Software-based implementations
    # requiring a key store. PKCS#11 based implementations does not.
    cryptoStore:
      path: /tmp/kongyixueyuan-msp

   # BCCSP config for the client. Used by GO SDK.
  BCCSP:
    security:
     enabled: true
     default:
      provider: "SW"
     hashAlgorithm: "SHA2"
     softVerify: true
     level: 256

  tlsCerts:
    # [Optional]. Use system certificate pool when connecting to peers, orderers (for negotiating TLS) Default: false
    systemCertPool: false

    # [Optional]. Client key and cert for TLS handshake with peers and orderers
    client:
      key:
        path:
      cert:
        path:

#
# [Optional]. But most apps would have this section so that channel objects can be constructed
# based on the content below. If an app is creating channels, then it likely will not need this
# section.
#
channels:
  # name of the channel
  kevinkongyixueyuan:
    # Required. list of orderers designated by the application to use for transactions on this
    # channel. This list can be a result of access control ("org1" can only access "ordererA"), or
    # operational decisions to share loads from applications among the orderers.  The values must
    # be "names" of orgs defined under "organizations/peers"
    # deprecated: not recommended, to override any orderer configuration items, entity matchers should be used.
    # orderers:
    #  - orderer.kevin.kongyixueyuan.com

    # Required. list of peers from participating orgs
    peers:
      peer0.org1.kevin.kongyixueyuan.com:
        # [Optional]. will this peer be sent transaction proposals for endorsement? The peer must
        # have the chaincode installed. The app can also use this property to decide which peers
        # to send the chaincode install request. Default: true
        endorsingPeer: true

        # [Optional]. will this peer be sent query proposals? The peer must have the chaincode
        # installed. The app can also use this property to decide which peers to send the
        # chaincode install request. Default: true
        chaincodeQuery: true

        # [Optional]. will this peer be sent query proposals that do not require chaincodes, like
        # queryBlock(), queryTransaction(), etc. Default: true
        ledgerQuery: true

        # [Optional]. will this peer be the target of the SDK's listener registration? All peers can
        # produce events but the app typically only needs to connect to one to listen to events.
        # Default: true
        eventSource: true

      peer1.org1.kevin.kongyixueyuan.com:
        endorsingPeer: true
        chaincodeQuery: true
        ledgerQuery: true
        eventSource: true

    policies:
      #[Optional] options for retrieving channel configuration blocks
      queryChannelConfig:
        #[Optional] min number of success responses (from targets/peers)
        minResponses: 1
        #[Optional] channel config will be retrieved for these number of random targets
        maxTargets: 1
        #[Optional] retry options for query config block
        retryOpts:
          #[Optional] number of retry attempts
          attempts: 5
          #[Optional] the back off interval for the first retry attempt
          initialBackoff: 500ms
          #[Optional] the maximum back off interval for any retry attempt
          maxBackoff: 5s
          #[Optional] he factor by which the initial back off period is exponentially incremented
          backoffFactor: 2.0
      #[Optional] options for retrieving discovery info
      discovery:
        #[Optional] discovery info will be retrieved for these number of random targets
        maxTargets: 2
        #[Optional] retry options for retrieving discovery info
        retryOpts:
          #[Optional] number of retry attempts
          attempts: 4
          #[Optional] the back off interval for the first retry attempt
          initialBackoff: 500ms
          #[Optional] the maximum back off interval for any retry attempt
          maxBackoff: 5s
          #[Optional] he factor by which the initial back off period is exponentially incremented
          backoffFactor: 2.0
      #[Optional] options for the event service
      eventService:
        # [Optional] resolverStrategy specifies the peer resolver strategy to use when connecting to a peer
        # Possible values: [PreferOrg (default), MinBlockHeight, Balanced]
        #
        # PreferOrg:
        #   Determines which peers are suitable based on block height lag threshold, although will prefer the peers in the
        #   current org (as long as their block height is above a configured threshold). If none of the peers from the current org
        #   are suitable then a peer from another org is chosen.
        # MinBlockHeight:
        #   Chooses the best peer according to a block height lag threshold. The maximum block height of all peers is
        #   determined and the peers whose block heights are under the maximum height but above a provided "lag" threshold are load
        #   balanced. The other peers are not considered.
        # Balanced:
        #   Chooses peers using the configured balancer.
        resolverStrategy: PreferOrg
        # [Optional] balancer is the balancer to use when choosing a peer to connect to
        # Possible values: [Random (default), RoundRobin]
        balancer: Random
        # [Optional] blockHeightLagThreshold sets the block height lag threshold. This value is used for choosing a peer
        # to connect to. If a peer is lagging behind the most up-to-date peer by more than the given number of
        # blocks then it will be excluded from selection.
        # If set to 0 then only the most up-to-date peers are considered.
        # If set to -1 then all peers (regardless of block height) are considered for selection.
        # Default: 5
        blockHeightLagThreshold: 5
        # [Optional] reconnectBlockHeightLagThreshold - if >0 then the event client will disconnect from the peer if the peer's
        # block height falls behind the specified number of blocks and will reconnect to a better performing peer.
        # If set to 0 then this feature is disabled.
        # Default: 10
        # NOTES:
        #   - peerMonitorPeriod must be >0 to enable this feature
        #   - Setting this value too low may cause the event client to disconnect/reconnect too frequently, thereby
        #     affecting performance.
        reconnectBlockHeightLagThreshold: 10
        # [Optional] peerMonitorPeriod is the period in which the connected peer is monitored to see if
        # the event client should disconnect from it and reconnect to another peer.
        # Default: 0 (disabled)
        peerMonitorPeriod: 5s

#
# list of participating organizations in this network
#
organizations:
  Org1:
    mspid: org1.kevin.kongyixueyuan.com
    cryptoPath: peerOrganizations/org1.kevin.kongyixueyuan.com/users/{userName}@org1.kevin.kongyixueyuan.com/msp
    peers:
      - peer0.org1.kevin.kongyixueyuan.com
      - peer1.org1.kevin.kongyixueyuan.com

    # [Optional]. Certificate Authorities issue certificates for identification purposes in a Fabric based
    # network. Typically certificates provisioning is done in a separate process outside of the
    # runtime network. Fabric-CA is a special certificate authority that provides a REST APIs for
    # dynamic certificate management (enroll, revoke, re-enroll). The following section is only for
    # Fabric-CA servers.
    certificateAuthorities:
      - ca.org1.kevin.kongyixueyuan.com

#
# List of orderers to send transaction and channel create/update requests to. For the time
# being only one orderer is needed. If more than one is defined, which one get used by the
# SDK is implementation specific. Consult each SDK's documentation for its handling of orderers.
#
orderers:
  orderer.kevin.kongyixueyuan.com:
    url: localhost:7050

    # these are standard properties defined by the gRPC library
    # they will be passed in as-is to gRPC client constructor
    grpcOptions:
      ssl-target-name-override: orderer.kevin.kongyixueyuan.com
      # These parameters should be set in coordination with the keepalive policy on the server,
      # as incompatible settings can result in closing of connection.
      # When duration of the 'keep-alive-time' is set to 0 or less the keep alive client parameters are disabled
      keep-alive-time: 0s
      keep-alive-timeout: 20s
      keep-alive-permit: false
      fail-fast: false
      # allow-insecure will be taken into consideration if address has no protocol defined, if true then grpc or else grpcs
      allow-insecure: false

    tlsCACerts:
      # Certificate location absolute path
      path: ${GOPATH}/src/github.com/kongyixueyuan.com/education/fixtures/crypto-config/ordererOrganizations/kevin.kongyixueyuan.com/tlsca/tlsca.kevin.kongyixueyuan.com-cert.pem

#
# List of peers to send various requests to, including endorsement, query
# and event listener registration.
#
peers:
  peer0.org1.kevin.kongyixueyuan.com:
    # this URL is used to send endorsement and query requests
    url: localhost:7051
    # eventUrl is only needed when using eventhub (default is delivery service)
    eventUrl: localhost:7053

    grpcOptions:
      ssl-target-name-override: peer0.org1.kevin.kongyixueyuan.com
      # These parameters should be set in coordination with the keepalive policy on the server,
      # as incompatible settings can result in closing of connection.
      # When duration of the 'keep-alive-time' is set to 0 or less the keep alive client parameters are disabled
      keep-alive-time: 0s
      keep-alive-timeout: 20s
      keep-alive-permit: false
      fail-fast: false
      # allow-insecure will be taken into consideration if address has no protocol defined, if true then grpc or else grpcs
      allow-insecure: false

    tlsCACerts:
      # Certificate location absolute path
      path: ${GOPATH}/src/github.com/kongyixueyuan.com/education/fixtures/crypto-config/peerOrganizations/org1.kevin.kongyixueyuan.com/tlsca/tlsca.org1.kevin.kongyixueyuan.com-cert.pem

  peer1.org1.kevin.kongyixueyuan.com:
    # this URL is used to send endorsement and query requests
    url: localhost:7151
    # eventUrl is only needed when using eventhub (default is delivery service)
    eventUrl: localhost:7153

    grpcOptions:
      ssl-target-name-override: peer1.org1.kevin.kongyixueyuan.com
      # These parameters should be set in coordination with the keepalive policy on the server,
      # as incompatible settings can result in closing of connection.
      # When duration of the 'keep-alive-time' is set to 0 or less the keep alive client parameters are disabled
      keep-alive-time: 0s
      keep-alive-timeout: 20s
      keep-alive-permit: false
      fail-fast: false
      # allow-insecure will be taken into consideration if address has no protocol defined, if true then grpc or else grpcs
      allow-insecure: false

    tlsCACerts:
      # Certificate location absolute path
      path: ${GOPATH}/src/github.com/kongyixueyuan.com/education/fixtures/crypto-config/peerOrganizations/org1.kevin.kongyixueyuan.com/tlsca/tlsca.org1.kevin.kongyixueyuan.com-cert.pem

#
# Fabric-CA is a special kind of Certificate Authority provided by Hyperledger Fabric which allows
# certificate management to be done via REST APIs. Application may choose to use a standard
# Certificate Authority instead of Fabric-CA, in which case this section would not be specified.
#
certificateAuthorities:
  ca.org1.kevin.kongyixueyuan.com:
    url: http://localhost:7054
    tlsCACerts:
      # Certificate location absolute path
      path: ${GOPATH}/src/github.com/kongyixueyuan.com/education/fixtures/crypto-config/peerOrganizations/org1.kevin.kongyixueyuan.com/ca/ca.org1.kevin.kongyixueyuan.com-cert.pem

    # Fabric-CA supports dynamic user enrollment via REST APIs. A "root" user, a.k.a registrar, is
    # needed to enroll and invoke new users.
    registrar:
      enrollId: admin
      enrollSecret: adminpw
    # [Optional] The optional name of the CA.
    caName: ca.org1.kevin.kongyixueyuan.com

entityMatchers:
  peer:
    - pattern: (\w*)peer0.org1.kevin.kongyixueyuan.com(\w*)
      urlSubstitutionExp: localhost:7051
      eventUrlSubstitutionExp: localhost:7053
      sslTargetOverrideUrlSubstitutionExp: peer0.org1.kevin.kongyixueyuan.com
      mappedHost: peer0.org1.kevin.kongyixueyuan.com

    - pattern: (\w*)peer1.org1.kevin.kongyixueyuan.com(\w*)
      urlSubstitutionExp: localhost:7151
      eventUrlSubstitutionExp: localhost:7153
      sslTargetOverrideUrlSubstitutionExp: peer1.org1.kevin.kongyixueyuan.com
      mappedHost: peer1.org1.kevin.kongyixueyuan.com

  orderer:
    - pattern: (\w*)orderer.kevin.kongyixueyuan.com(\w*)
      urlSubstitutionExp: localhost:7050
      sslTargetOverrideUrlSubstitutionExp: orderer.kevin.kongyixueyuan.com
      mappedHost: orderer.kevin.kongyixueyuan.com

  certificateAuthorities:
    - pattern: (\w*)ca.org1.kevin.kongyixueyuan.com(\w*)
      urlSubstitutionExp: http://localhost:7054
      mappedHost: ca.org1.kevin.kongyixueyuan.com
```



## 2.2 声明结构体

在当前项目根目录中创建一个存放链码文件的 `chaincode` 目录，然后在该目录下创建一个 `eduStruct.go` 的文件并对其进行编辑

```shell
$ mkdir chaincode
$ vim chaincode/eduStruct.go
```

`eduStruct.go` 文件主要声明一个结构体，用于将多个数据包装成为一个对象，然后进行进一步的处理。该文件完整代码如下：

```go
/**
  @Author : hanxiaodong
*/
package main

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
```



## 2.3 编写链码

在 `chaincode` 目录下创建一个 `main.go` 的文件并对其进行编辑

```shell
$ vim chaincode/main.go
```

`main.go` 文件作为链码的主文件，主要声明 `Init(stub shim.ChaincodeStubInterface)、Invoke(stub shim.ChaincodeStubInterface)  `  函数，完成对链码初始化及调用的相关实现，完整代码如下：

```go
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
```



创建 `eduCC.go` 文件，该文件实现了使用链码相关的API对分类账本状态进行具体操作的各个函数：

- **PutEdu：**实现将指定的对象序列化后保存至分类账本中
- **GetEduInfo：**根据指定的Key（身份证号码）查询对应的状态，反序列后将对象返回
- **getEduByQueryString：**根据指定的查询字符串从 `CouchDB` 中查询状态
- **addEdu：**接收对象并调用 `PutEdu` 函数实现保存状态的功能
- **queryEduByCertNoAndName：**根据指定的证书编号与姓名查询状态
- **queryEduInfoByEntityID：**根据指定的身份证号码（Key）查询状态
- **updateEdu：**实现对状态进行编辑功能
- **delEdu：**从分类账本中删除状态，此功能暂不提供

`eduCC.go` 文件完整内容如下：

```go
/**
  @Author : hanxiaodong
*/

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"bytes"
)

const DOC_TYPE = "eduObj"

// 保存edu
// args: education
func PutEdu(stub shim.ChaincodeStubInterface, edu Education) ([]byte, bool) {

	edu.ObjectType = DOC_TYPE

	b, err := json.Marshal(edu)
	if err != nil {
		return nil, false
	}

	// 保存edu状态
	err = stub.PutState(edu.EntityID, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

// 根据身份证号码查询信息状态
// args: entityID
func GetEduInfo(stub shim.ChaincodeStubInterface, entityID string) (Education, bool)  {
	var edu Education
	// 根据身份证号码查询信息状态
	b, err := stub.GetState(entityID)
	if err != nil {
		return edu, false
	}

	if b == nil {
		return edu, false
	}

	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &edu)
	if err != nil {
		return edu, false
	}

	// 返回结果
	return edu, true
}

// 根据指定的查询字符串实现富查询
func getEduByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer  resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil

}

// 添加信息
// args: educationObject
// 身份证号为 key, Education 为 value
func (t *EducationChaincode) addEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}

	var edu Education
	err := json.Unmarshal([]byte(args[0]), &edu)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	// 查重: 身份证号码必须唯一
	_, exist := GetEduInfo(stub, edu.EntityID)
	if exist {
		return shim.Error("要添加的身份证号码已存在")
	}

	_, bl := PutEdu(stub, edu)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))
}

// 根据证书编号及姓名查询信息
// args: CertNo, name
func (t *EducationChaincode) queryEduByCertNoAndName(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	CertNo := args[0]
	name := args[1]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"eduObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"CertNo\":\"%s\", \"Name\":\"%s\"}}", DOC_TYPE, CertNo, name)

	// 查询数据
	result, err := getEduByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据证书编号及姓名查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的证书编号及姓名没有查询到相关的信息")
	}
	return shim.Success(result)
}

// 根据身份证号码查询详情（溯源）
// args: entityID
func (t *EducationChaincode) queryEduInfoByEntityID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}

	// 根据身份证号码查询edu状态
	b, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据身份证号码查询信息失败")
	}

	if b == nil {
		return shim.Error("根据身份证号码没有查询到相关的信息")
	}

	// 对查询到的状态进行反序列化
	var edu Education
	err = json.Unmarshal(b, &edu)
	if err != nil {
		return  shim.Error("反序列化edu信息失败")
	}

	// 获取历史变更数据
	iterator, err := stub.GetHistoryForKey(edu.EntityID)
	if err != nil {
		return shim.Error("根据指定的身份证号码查询对应的历史变更数据失败")
	}
	defer iterator.Close()

	// 迭代处理
	var historys []HistoryItem
	var hisEdu Education
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("获取edu的历史变更数据失败")
		}

		var historyItem HistoryItem
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisEdu)

		if hisData.Value == nil {
			var empty Education
			historyItem.Education = empty
		}else {
			historyItem.Education = hisEdu
		}

		historys = append(historys, historyItem)

	}

	edu.Historys = historys

	// 返回
	result, err := json.Marshal(edu)
	if err != nil {
		return shim.Error("序列化edu信息时发生错误")
	}
	return shim.Success(result)
}

// 根据身份证号更新信息
// args: educationObject
func (t *EducationChaincode) updateEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}

	var info Education
	err := json.Unmarshal([]byte(args[0]), &info)
	if err != nil {
		return  shim.Error("反序列化edu信息失败")
	}

	// 根据身份证号码查询信息
	result, bl := GetEduInfo(stub, info.EntityID)
	if !bl{
		return shim.Error("根据身份证号码查询信息时发生错误")
	}

	result.EnrollDate = info.EnrollDate
	result.GraduationDate = info.GraduationDate
	result.SchoolName = info.SchoolName
	result.Major = info.Major
	result.QuaType = info.QuaType
	result.Length = info.Length
	result.Mode = info.Mode
	result.Level = info.Level
	result.Graduation = info.Graduation
	result.CertNo = info.CertNo;

	_, bl = PutEdu(stub, result)
	if !bl {
		return shim.Error("保存信息信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息更新成功"))
}

// 根据身份证号删除信息（暂不对外提供）
// args: entityID
func (t *EducationChaincode) delEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}

	/*var edu Education
	result, bl := GetEduInfo(stub, info.EntityID)
	err := json.Unmarshal(result, &edu)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}*/

	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("删除信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息删除成功"))
}
```



## 2.4 SDK实现

链码编写好以后，我们需要使用 Fabric-SDK-Go 提供的相关 API 来实现对链码的安装及实例化操作，而无需在命令提示符中输入烦锁的相关操作命令。接下来依次完成如下步骤：

### 2.4.1 创建SDK

在 `sdkInit` 目录下新创建一个名为 `start.go` 的go文件利用 vim 编辑器进行编辑：

```shell
$ cd $GOPATH/src/github.com/kongyixueyuan.com/education
$ vim sdkInit/start.go 
```

```go
/**
  author: hanxiaodong
 */
package sdkInit

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"

)

const ChaincodeVersion  = "1.0"

func SetupSDK(ConfigFile string, initialized bool) (*fabsdk.FabricSDK, error) {

	if initialized {
		return nil, fmt.Errorf("Fabric SDK已被实例化")
	}

	sdk, err := fabsdk.New(config.FromFile(ConfigFile))
	if err != nil {
		return nil, fmt.Errorf("实例化Fabric SDK失败: %v", err)
	}

	fmt.Println("Fabric SDK初始化成功")
	return sdk, nil
}

func CreateChannel(sdk *fabsdk.FabricSDK, info *InitInfo) error {

	clientContext := sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName))
	if clientContext == nil {
		return fmt.Errorf("根据指定的组织名称与管理员创建资源管理客户端Context失败")
	}

	// New returns a resource management client instance.
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return fmt.Errorf("根据指定的资源管理客户端Context创建通道管理客户端失败: %v", err)
	}

	// New creates a new Client instance
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(info.OrgName))
	if err != nil {
		return fmt.Errorf("根据指定的 OrgName 创建 Org MSP 客户端实例失败: %v", err)
	}

	//  Returns: signing identity
	adminIdentity, err := mspClient.GetSigningIdentity(info.OrgAdmin)
	if err != nil {
		return fmt.Errorf("获取指定id的签名标识失败: %v", err)
	}

	// SaveChannelRequest holds parameters for save channel request
	channelReq := resmgmt.SaveChannelRequest{ChannelID:info.ChannelID, ChannelConfigPath:info.ChannelConfig, SigningIdentities:[]msp.SigningIdentity{adminIdentity}}
	// save channel response with transaction ID
	 _, err = resMgmtClient.SaveChannel(channelReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("创建应用通道失败: %v", err)
	}

	fmt.Println("通道已成功创建，")

	info.OrgResMgmt = resMgmtClient

	// allows for peers to join existing channel with optional custom options (specific peers, filtered peers). If peer(s) are not specified in options it will default to all peers that belong to client's MSP.
	err = info.OrgResMgmt.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("Peers加入通道失败: %v", err)
	}

	fmt.Println("peers 已成功加入通道.")
	return nil
}
```

在这个阶段，我们只初始化一个客户端，它将与 peer，CA 和 orderer进行通信。 还创建了一个指定的应用通道, 并将 Peer 节点加入到此通道中

### 2.4.2 编写测试代码

为了确保客户端能够初始化所有组件，将在启动网络的情况下进行简单的测试。 为了做到这一点，我们需要编写 Go 代码，在项目根目录下新创建一个 `main.go` 的主文件并编辑内容

```shell
$ cd $GOPATH/src/github.com/kongyixueyuan.com/education
$ vim main.go
```

`main.go` 文件完整源代码如下：

```go
/**
  author: hanxiaodong
 */
package main

import (
	"os"
	"fmt"
	"github.com/kongyixueyuan.com/education/sdkInit"
)

const (
	configFile = "config.yaml"
	initialized = false
	SimpleCC = "educc"
)

func main() {

	initInfo := &sdkInit.InitInfo{

		ChannelID: "kevinkongyixueyuan",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/kongyixueyuan.com/education/fixtures/artifacts/channel.tx",

		OrgAdmin:"Admin",
		OrgName:"Org1",
		OrdererOrgName: "orderer.kevin.kongyixueyuan.com",

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

}
```



### 2.5 安装dep工具

在运行应用程序之前，需要将 Go 源代码时行编译，但在开始编译之前，我们需要使用一个 `vendor` 目录来包含应用中所需的所有的依赖关系。 在我们的GOPATH中，我们有Fabric SDK Go和其他项目。 在尝试编译应用程序时，Golang 会在 GOPATH 中搜索依赖项，但首先会检查项目中是否存在`vendor` 文件夹。 如果依赖性得到满足，那么 Golang 就不会去检查 GOPATH 或 GOROOT。 这在使用几个不同版本的依赖关系时非常有用（可能会发生一些冲突，比如在例子中有多个BCCSP定义，通过使用像[`dep`](https://translate.googleusercontent.com/translate_c?depth=1&hl=zh-CN&rurl=translate.google.com&sl=en&sp=nmt4&tl=zh-CN&u=https://github.com/golang/dep&xid=25657,15700002,15700019,15700124,15700149,15700168,15700186,15700201&usg=ALkJrhgelyRl7D3pIJRpuA8cynagkWYHXg)这样的工具在`vendor`目录中来处理这些依赖关系。

将如下环境变量设置到用户的环境文件中(.bashrc)中

```shell
$ vim ~/.bashrc

export PATH=$PATH:$GOPATH/bin
```

执行 `source` 命令

```shell
$ source ~/.bashrc
```

安装 dep 工具

```shell
$ go get -u github.com/golang/dep/cmd/dep
```

### 2.5.1 下载所需依赖

`dep` 工具安装好之后我们来安装应用所需要的依赖

使用 `dep` 命令需要一个名为 `Gopkg.toml` 的配置文件指定依赖信息

创建一个名为`Gopkg.toml`的文件并将其复制到里面：

```shell
$ vim Gopkg.toml
```

```toml
ignored = ["github.com/kongyixueyuan.com/education/chaincode"]

[[constraint]]
  # Release v1.0.0-alpha4
  name = "github.com/hyperledger/fabric-sdk-go"
  revision = "a906355f73d060d7bf95874a9e90dc17589edbb3"

```

使用`dep`限制在 vendor 中指定希望SDK的特定版本。

> 因为在使用SDK执行安装与实例化链码部分时，使用 dep ensure 命令会报依赖错误，所以必须添加 ignored 选项

保存该文件，然后执行 `dep ensure` 命令，该命令会自动将项目所需的依赖下载至当前的 `vendor` 目录中（下载依赖可能需要一段时间）：

```shell
$ dep ensure
```

提醒：`dep ensure` 命令执行由于时间比较长，所以执行一次后即可，在后面的Makefile中可注释`@dep ensure`命令。

### 2.5.2 测试Fabric-SDK

所在依赖下载安装完成后，我们就可以进行测试

首先启动网络：

```shell
$ cd fixtures
$ docker-compose up -d
```

然后编译并运行：

```shell
$ cd ..
$ go build
$ ./education
```

命令执行后输出结果如下图所示：

### 2.5.3 利用Makefile

Fabric SDK生成一些文件，如证书，二进制文件和临时文件。 关闭网络不会完全清理环境，当需要重新启动时，这些文件将被使用以避免重复的构建过程。 对于开发，可以快速测试，但对于真正的测试环境，需要清理所有内容并从头开始。

*如何清理环境*

- 关闭你的网络： 

  ```shell
  $ cd $GOPATH/src/github.com/kongyixueyuan.com/education/fixtures 
  $ docker-compose down
  ```

- 删除证书存储（在配置文件中，`client.credentialStore`中定义）：

  ```shell
  $ rm -rf /tmp/kongyixueyuan-*
  ```

- 删除一些不是由`docker-compose`命令生成的docker容器和docker镜像： 

  ```shell
  $ docker rm -f -v `docker ps -a --no-trunc | grep "kongyixueyuan" | cut -d ' ' -f 1` 2>/dev/null
  和 
  $ docker rmi `docker images --no-trunc | grep "kongyixueyuan" | cut -d ' ' -f 1` 2>/dev/null
  ```

*如何更有效率？*

可以使用一个命令在一个步骤中自动完成所有的这些任务。 构建和启动过程也可以设置自动化。 为此，将创建一个`Makefile` 文件。 首先，确保您的系统中在  `make` 工具：

```shell
$ make --version
```

> 如果没有 `make` 工具（Ubuntu），则需要先安装 `make` ：
>
> ```shell
> $ sudo apt install make 
> ```

然后使用以下内容在项目的根目录下创建一个名为`Makefile`的文件并进行编辑：

```shell
$ cd $GOPATH/src/github.com/kongyixueyuan.com/education
$ vim Makefile
```

`Makefile` 文件完整内容如下：

```makefile
.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@dep ensure
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@cd fixtures && docker-compose up --force-recreate -d
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd fixtures && docker-compose down
	@echo "Environment down"

##### RUN
run:
	@echo "Start app ..."
	@./education

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/kongyixueyuan-* kongyixueyuan
	@docker rm -f -v `docker ps -a --no-trunc | grep "kongyixueyuan" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "kongyixueyuan" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"
```

现在完成任务：

1. 整个环境将被清理干净，
2. go程序将被编译，
3. 之后将部署网络
4. 最后该应用程序将启动并运行。

要使用它，请进入项目的根目录并使用`make`命令：

- 任务`all` ： `make`或`make all`
- 任务`clean` ：清理一切并释放网络（ `make clean` ）
- 任务`build` ：只需构建应用程序（ `make build` ）
- 任务`env-up` ：只需建立网络（ `make env-up` ）
- ...

![createchannel](./img/createchannel.png)

如果出现上图的输出结果，则说明执行成功，否则需要根据出现的错误提示进行相应的处理。

## 2.6 链码安装及实例化

### 2.6.1 使用Fabric-SDK安装及实例化链码

编辑 `sdkInit/start.go` 文件，利用Fabric-SDK提供的接口，对链码进行安装及实例化

```shell
$ vim sdkInit/start.go
```

在 `start.go` 文件中添加如下内容

```go
import (
    [......]
    
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// hanxiaodong
// QQ群（专业Fabric交流群）：862733552
func InstallAndInstantiateCC(sdk *fabsdk.FabricSDK, info *InitInfo) (*channel.Client, error) {
	fmt.Println("开始安装链码......")
	// creates new go lang chaincode package
	ccPkg, err := gopackager.NewCCPackage(info.ChaincodePath, info.ChaincodeGoPath)
	if err != nil {
		return nil, fmt.Errorf("创建链码包失败: %v", err)
	}

	// contains install chaincode request parameters
	installCCReq := resmgmt.InstallCCRequest{Name: info.ChaincodeID, Path: info.ChaincodePath, Version: ChaincodeVersion, Package: ccPkg}
	// allows administrators to install chaincode onto the filesystem of a peer
	_, err = info.OrgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return nil, fmt.Errorf("安装链码失败: %v", err)
	}

	fmt.Println("指定的链码安装成功")
	fmt.Println("开始实例化链码......")

	//  returns a policy that requires one valid
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.kevin.kongyixueyuan.com"})

	instantiateCCReq := resmgmt.InstantiateCCRequest{Name: info.ChaincodeID, Path: info.ChaincodePath, Version: ChaincodeVersion, Args: [][]byte{[]byte("init")}, Policy: ccPolicy}
	// instantiates chaincode with optional custom options (specific peers, filtered peers, timeout). If peer(s) are not specified
	_, err = info.OrgResMgmt.InstantiateCC(info.ChannelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return nil, fmt.Errorf("实例化链码失败: %v", err)
	}

	fmt.Println("链码实例化成功")

	clientChannelContext := sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.UserName), fabsdk.WithOrg(info.OrgName))
	// returns a Client instance. Channel client can query chaincode, execute chaincode and register/unregister for chaincode events on specific channel.
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("创建应用通道客户端失败: %v", err)
	}

	fmt.Println("通道客户端创建成功，可以利用此客户端调用链码进行查询或执行事务.")

	return channelClient, nil
}
```

### 2.6.2 在main中调用

编辑 `main.go` 文件

```shell
$ vim main.go
```

`main.go` 完整内容如下：

```go
/**
  author: hanxiaodong
 */
package main

import (
	"os"
	"fmt"
	"github.com/kongyixueyuan.com/education/sdkInit"
)

const (
	configFile = "config.yaml"
	initialized = false
	SimpleCC = "educc"
)

func main() {

	initInfo := &sdkInit.InitInfo{

		ChannelID: "kevinkongyixueyuan",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/kongyixueyuan.com/education/fixtures/artifacts/channel.tx",

		OrgAdmin:"Admin",
		OrgName:"Org1",
		OrdererOrgName: "orderer.kevin.kongyixueyuan.com",

		ChaincodeID: SimpleCC,
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

}
```

### 2.6.3 测试

执行 `make` 命令

```shell
$ make
```

输出如下：

![cc实例化成功](./img/installcc.png)





#3 业务层实现

##3.1 事件处理

在项目根目录下创建一个 `service` 目录作为业务层，在业务层中，我们使用 `Fabric-SDK-Go` 提供的接口对象调用相应的 API 以实现对链码的访问，最终实现对分类账本中的状态进行操作。 

```shell
$ cd $GOPATH/src/github.com/kongyixueyuan.com/education
$ mkdir service
```

在 `service` 目录下创建 `domain.go` 文件并进行编辑， 声明一个结构体及对事件相关而封装的源代码

```shell
$ vim service/domain.go
```

`domain.go` 文件完整内容如下：

```go
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
```



## 3.2 业务层调用链码实现添加状态

在 `service` 目录下创建 `eduService.go` 文件

```shell
$ vim service/eduService.go
```

在 `eduService.go` 文件中编写内容如下，通过一个 `SaveEdu` 函数实现链码的调用，向分类账本中添加状态的功能：

```go
/**
  author: hanxiaodong
  QQ群（专业Fabric交流群）：862733552
 */
package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    "encoding/json"
)

func (t *ServiceSetup) SaveEdu(edu Education) (string, error) {

	eventID := "eventAddEdu"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将edu对象序列化成为字节数组
	b, err := json.Marshal(edu)
	if err != nil {
		return "", fmt.Errorf("指定的edu对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "addEdu", Args: [][]byte{b, []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}
```



**测试添加状态**

编辑 `main.go` 文件

```shell
$ vim main.go
```

`main.go` 中创建两个 `edu` 个对象，并调用 `SaveEdu` 函数，内容如下：

```go
/**
  author: hanxiaodong
 */

package main

import (
	[......]
	"github.com/kongyixueyuan.com/education/service"
)

[......]
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
		Photo: "/static/phone/11.png",
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
		Photo: "/static/phone/22.png",
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

	//===========================================//

}
```



执行 `make` 命令运行应用程序

```shell
$ make
```

执行后如下图所示：

![测试添加状态](./img/saveedu.png)



## 3.3 调用链码实现根据证书编号与名称查询状态

通过上面的 `SaveEdu(edu Education)` 函数，实现了向分类账本中添加状态，那么我们还需要实现从该分类账本中根据指定的条件查询出相应的状态，编辑 `service/eduService.go` 文件，向该文件中添加实现根据证书编号与姓名查询状态的相应代码。

```shell
$ vim service/eduService.go
```

定义一个 `FindEduByCertNoAndName` 函数，接收两个字符串类型的参数，分别代表证书编号与姓名，该函数实现通过调用链码而实现查询状态的功能，该函数完整代码如下：

```go
[......]

func (t *ServiceSetup) FindEduByCertNoAndName(certNo, name string) ([]byte, error){

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryEduByCertNoAndName", Args: [][]byte{[]byte(certNo), []byte(name)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

```

**测试根据证书编号与名称查询状态**

编辑 `main.go` 文件

```shell
$ vim main.go
```

在 `main.go` 文件中添加调用代码如下内容：

```go
[......]
	
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

	//===========================================//

}
```

执行 `make` 命令运行应用程序

```shell
$ make
```

执行后如下图所示：

![证书编号查询测试](./img/findEduByCertNoAndName.png)



## 3.4 调用链码实现根据身份证号码查询状态

通过上面的 `FindEduByCertNoAndName(certNo, name string)` 函数，实现从该分类账本中根据指定的证书编号与姓名查询出相应的状态，下面我们来实现根据身份证号码查询状态的功能，编辑 `service/eduService.go` 文件，向该文件中添加实现根据 key 查询状态的相应代码。

```shell
$ vim service/eduService.go
```

定义一个 `FindEduInfoByEntityID` 函数，接收一个字符串类型的参数，代表身份证号码（key），该函数实现通过调用链码而实现查询状态的功能，该函数完整代码如下：

```go
[......]

func (t *ServiceSetup) FindEduInfoByEntityID(entityID string) ([]byte, error){

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryEduInfoByEntityID", Args: [][]byte{[]byte(entityID)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}
```

**测试根据身份证号码查询状态**

编辑 `main.go` 文件

```shell
$ vim main.go
```

在 `main.go` 文件中添加调用代码如下内容：

```go
[......]
	
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

	//===========================================//

}
```

执行 `make` 命令运行应用程序

```shell
$ make
```

执行后如下图所示：

![身份证号码查询测试](./img/findEduInfoByEntityID.png)



## 3.5 调用链码实现修改/添加信息状态

在一些情况下，有些人才会利用工作的业余时间进修，从而提升学历层次，我们必须要考虑到这种情况，所以需要应用程序实现对已有人员的信息进行编辑的功能；但是编辑并不能将之前的学历信息删除，而是在保留之前状态的基础之上添加新的状态，区块链技术很好的帮我们解决了这个问题。编辑 `service/eduService.go` 文件，向该文件中添加修改已有状态的相关代码。

```shell
$ vim service/eduService.go
```

定义一个 `ModifyEdu` 函数，接收一个 `Education` 类型的对象，该函数实现通过调用链码而实现对已存在的状态进行修改（添加新信息）的功能，该函数完整代码如下：

```go
[......]

func (t *ServiceSetup) ModifyEdu(edu Education) (string, error) {

	eventID := "eventModifyEdu"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将edu对象序列化成为字节数组
	b, err := json.Marshal(edu)
	if err != nil {
		return "", fmt.Errorf("指定的edu对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "updateEdu", Args: [][]byte{b, []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}
```

**测试修改状态**

编辑 `main.go` 文件

```shell
$ vim main.go
```

在 `main.go` 文件中添加调用代码如下内容：

```go
[......]
	
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
		Photo: "/static/phone/11.png",
	}
	msg, err = serviceSetup.ModifyEdu(info)
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("信息操作成功, 交易编号为: " + msg)
	}

	//===========================================//

}
```



执行 `make` 命令运行应用程序

```shell
$ make
```

执行后如下图所示：

![修改信息测试](./img/modifyEdu.png)



**查看修改之后的状态（根据身份证号码）**

状态被修改之后，我们为了确认是否真正修改成功，所以需要调用已经编写好的 `FindEduInfoByEntityID(entityID string)` 函数实现查询详情的功能。

编辑 `main.go` 文件

```shell
$ vim main.go
```

在 `main.go` 文件中添加调用代码如下内容：

```go
[......]
	
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

	//===========================================//

}
```

执行 `make` 命令运行应用程序

```shell
$ make
```

执行后如下图所示：

![修改信息测试](./img/update_findEduInfoByEntityID.png)

从终端的输出结果中可以看到详情信息已从分类账本中被成功查询，接下来我们使用根据证书编号与姓名查询修改之后的信息，看看是否正确

**查看修改之后的最新状态（根据证书编号与姓名）**

状态被修改之后，我们为了确认是否真正修改成功，所以需要调用已经编写好的 `FindEduInfoByEntityID(entityID string)` 函数实现查询详情的功能。

编辑 `main.go` 文件

```shell
$ vim main.go
```

在 `main.go` 文件中添加调用代码如下内容：

```go
[......]
	
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

	//===========================================//

}
```

执行 `make` 命令运行应用程序

```shell
$ make
```

执行后如下图所示：

![修改信息测试](./img/update_findEduByCertNoAndName.png)



# 4 控制层实现

## 4.1 设置系统用户

通过业务层已经实现了利用 `fabric-sdk-go` 调用链码查询或操作分类账本状态，接下来，我们开始实现Web应用层，应用层将其分为两个部分，

- **控制层**
- **视图层**

在项目根目录下新创建一个名为 `web`  的目录，用来存放Web应用层的所有内容

```shell
$ cd $GOPATH/src/github.com/kongyixueyuan.com/education
$ mkdir -p web/controller
```

在 `web` 目录下创建  `controller` 子目录，在该目录下创建 `userInfo.go` 、 `controllerResponse.go` 与 `controllerHandler.go` 三个文件

```shell
$ vim web/controller/userInfo.go
```

`userInfo.go` 用来模拟RDB，保存系统用户信息，作为用户登录时核对用户信息，当然，这部分大家可以使用 `MySQL` 或其它数据库来实现。

`userInfo.go` 完整代码如下：

```go
/**
  @Author : hanxiaodong
*/

package controller

import "github.com/kongyixueyuan.com/education/service"

type Application struct {
	Setup *service.ServiceSetup
}

type User struct {
	LoginName	string
	Password	string
	IsAdmin	string
}


var users []User

func init() {

	admin := User{LoginName:"Hanxiaodong", Password:"123456", IsAdmin:"T"}
	alice := User{LoginName:"ChainDesk", Password:"123456", IsAdmin:"T"}
	bob := User{LoginName:"alice", Password:"123456", IsAdmin:"F"}
	jack := User{LoginName:"bob", Password:"123456", IsAdmin:"F"}

	users = append(users, admin)
	users = append(users, alice)
	users = append(users, bob)
	users = append(users, jack)

}

func isAdmin(cuser User) bool {
	if cuser.IsAdmin == "T"{
		return true
	}
	return false
}
```



## 4.2 处理响应

创建 `controllerResponse.go`  文件

```shell
$ vim web/controller/controllerResponse.go
```

`controllerResponse.go` 主要实现对用户请求的响应，将响应结果返回给客户端浏览器。文件完整代码如下：

```go
/**
  @Author : hanxiaodong
*/

package controller

import (
	"net/http"
	"path/filepath"
	"html/template"
	"fmt"
)

func ShowView(w http.ResponseWriter, r *http.Request, templateName string, data interface{})  {

	// 指定视图所在路径
	pagePath := filepath.Join("web", "tpl", templateName)

	resultTemplate, err := template.ParseFiles(pagePath)
	if err != nil {
		fmt.Printf("创建模板实例错误: %v", err)
		return
	}

	err = resultTemplate.Execute(w, data)
	if err != nil {
		fmt.Printf("在模板中融合数据时发生错误: %v", err)
		//fmt.Fprintf(w, "显示在客户端浏览器中的错误信息")
		return
	}

}
```



## 4.3 处理请求

创建 `controllerHandler.go` 文件

```shell
$ vim web/controller/controllerHandler.go
```

`controllerHandler.go` 文件主要实现接收用户请求，并根据不同的用户请求调用业务层不同的函数，实现对分类账本的访问。其中需要声明并实现的函数：

文件完整内容如下：

```go
/**
  @Author : hanxiaodong
*/

package controller

import (
	"net/http"
	"encoding/json"
	"github.com/kongyixueyuan.com/education/service"
	"fmt"
)

var cuser User

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request)  {

	ShowView(w, r, "login.html", nil)
}

func (app *Application) Index(w http.ResponseWriter, r *http.Request)  {
	ShowView(w, r, "index.html", nil)
}

func (app *Application) Help(w http.ResponseWriter, r *http.Request)  {
	data := &struct {
		CurrentUser User
	}{
		CurrentUser:cuser,
	}
	ShowView(w, r, "help.html", data)
}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")

	var flag bool
	for _, user := range users {
		if user.LoginName == loginName && user.Password == password {
			cuser = user
			flag = true
			break
		}
	}

	data := &struct {
		CurrentUser User
		Flag bool
	}{
		CurrentUser:cuser,
		Flag:false,
	}

	if flag {
		// 登录成功
		ShowView(w, r, "index.html", data)
	}else{
		// 登录失败
		data.Flag = true
		data.CurrentUser.LoginName = loginName
		ShowView(w, r, "login.html", data)
	}
}

// 用户登出
func (app *Application) LoginOut(w http.ResponseWriter, r *http.Request)  {
	cuser = User{}
	ShowView(w, r, "login.html", nil)
}

// 显示添加信息页面
func (app *Application) AddEduShow(w http.ResponseWriter, r *http.Request)  {
	data := &struct {
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Msg:"",
		Flag:false,
	}
	ShowView(w, r, "addEdu.html", data)
}

// 添加信息
func (app *Application) AddEdu(w http.ResponseWriter, r *http.Request)  {

	edu := service.Education{
		Name:r.FormValue("name"),
		Gender:r.FormValue("gender"),
		Nation:r.FormValue("nation"),
		EntityID:r.FormValue("entityID"),
		Place:r.FormValue("place"),
		BirthDay:r.FormValue("birthDay"),
		EnrollDate:r.FormValue("enrollDate"),
		GraduationDate:r.FormValue("graduationDate"),
		SchoolName:r.FormValue("schoolName"),
		Major:r.FormValue("major"),
		QuaType:r.FormValue("quaType"),
		Length:r.FormValue("length"),
		Mode:r.FormValue("mode"),
		Level:r.FormValue("level"),
		Graduation:r.FormValue("graduation"),
		CertNo:r.FormValue("certNo"),
		Photo:r.FormValue("photo"),
	}

	app.Setup.SaveEdu(edu)
	
	r.Form.Set("certNo", edu.CertNo)
	r.Form.Set("name", edu.Name)
	app.FindCertByNoAndName(w, r)
}

func (app *Application) QueryPage(w http.ResponseWriter, r *http.Request)  {
	data := &struct {
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Msg:"",
		Flag:false,
	}
	ShowView(w, r, "query.html", data)
}

// 根据证书编号与姓名查询信息
func (app *Application) FindCertByNoAndName(w http.ResponseWriter, r *http.Request)  {
	certNo := r.FormValue("certNo")
	name := r.FormValue("name")
	result, err := app.Setup.FindEduByCertNoAndName(certNo, name)
	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	fmt.Println("根据证书编号与姓名查询信息成功：")
	fmt.Println(edu)

	data := &struct {
		Edu service.Education
		CurrentUser User
		Msg string
		Flag bool
		History bool
	}{
		Edu:edu,
		CurrentUser:cuser,
		Msg:"",
		Flag:false,
		History:false,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "queryResult.html", data)
}

func (app *Application) QueryPage2(w http.ResponseWriter, r *http.Request)  {
	data := &struct {
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Msg:"",
		Flag:false,
	}
	ShowView(w, r, "query2.html", data)
}

// 根据身份证号码查询信息
func (app *Application) FindByID(w http.ResponseWriter, r *http.Request)  {
	entityID := r.FormValue("entityID")
	result, err := app.Setup.FindEduInfoByEntityID(entityID)
	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	data := &struct {
		Edu service.Education
		CurrentUser User
		Msg string
		Flag bool
		History bool
	}{
		Edu:edu,
		CurrentUser:cuser,
		Msg:"",
		Flag:false,
		History:true,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "queryResult.html", data)
}

// 修改/添加新信息
func (app *Application) ModifyShow(w http.ResponseWriter, r *http.Request)  {
	// 根据证书编号与姓名查询信息
	certNo := r.FormValue("certNo")
	name := r.FormValue("name")
	result, err := app.Setup.FindEduByCertNoAndName(certNo, name)

	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	data := &struct {
		Edu service.Education
		CurrentUser User
		Msg string
		Flag bool
	}{
		Edu:edu,
		CurrentUser:cuser,
		Flag:true,
		Msg:"",
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "modify.html", data)
}

// 修改/添加新信息
func (app *Application) Modify(w http.ResponseWriter, r *http.Request) {
	edu := service.Education{
		Name:r.FormValue("name"),
		Gender:r.FormValue("gender"),
		Nation:r.FormValue("nation"),
		EntityID:r.FormValue("entityID"),
		Place:r.FormValue("place"),
		BirthDay:r.FormValue("birthDay"),
		EnrollDate:r.FormValue("enrollDate"),
		GraduationDate:r.FormValue("graduationDate"),
		SchoolName:r.FormValue("schoolName"),
		Major:r.FormValue("major"),
		QuaType:r.FormValue("quaType"),
		Length:r.FormValue("length"),
		Mode:r.FormValue("mode"),
		Level:r.FormValue("level"),
		Graduation:r.FormValue("graduation"),
		CertNo:r.FormValue("certNo"),
		Photo:r.FormValue("photo"),
	}

	app.Setup.ModifyEdu(edu)

	r.Form.Set("entityID", edu.EntityID)
	app.FindByID(w, r)
}
```



> 提示：用户在做一些管理操作时需要验证其它是否有相应的操作权限，需要另外进行设计。



## 4.4 指定路由

在 `web` 目录下创建一个 `webServer.go` 文件

```shell
$ vim web/webServer.go
```

该文件主要声明用户请求的路由信息，并且指定 Web 服务的启动信息。文件完整内容如下：

```go
/**
  @Author : hanxiaodong
*/

package web

import (
	"net/http"
	"fmt"
	"github.com/kongyixueyuan.com/education/web/controller"
)


// 启动Web服务并指定路由信息
func WebStart(app controller.Application)  {

	fs:= http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 指定路由信息(匹配请求)
	http.HandleFunc("/", app.LoginView)
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/loginout", app.LoginOut)

	http.HandleFunc("/index", app.Index)
	http.HandleFunc("/help", app.Help)

	http.HandleFunc("/addEduInfo", app.AddEduShow)	// 显示添加信息页面
	http.HandleFunc("/addEdu", app.AddEdu)	// 提交信息请求

	http.HandleFunc("/queryPage", app.QueryPage)	// 转至根据证书编号与姓名查询信息页面
	http.HandleFunc("/query", app.FindCertByNoAndName)	// 根据证书编号与姓名查询信息

	http.HandleFunc("/queryPage2", app.QueryPage2)	// 转至根据身份证号码查询信息页面
	http.HandleFunc("/query2", app.FindByID)	// 根据身份证号码查询信息


	http.HandleFunc("/modifyPage", app.ModifyShow)	// 修改信息页面
	http.HandleFunc("/modify", app.Modify)	//  修改信息

	http.HandleFunc("/upload", app.UploadFile)

	fmt.Println("启动Web服务, 监听端口号为: 9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("Web服务启动失败: %v", err)
	}

}
```



# 5 视图层实现

## 5.1 目录结构

在项目的web目录下新创建一个名为 `static`  的目录，用来存放Web应用视图层的所有静态内容

```shell
$ cd $GOPATH/src/github.com/kongyixueyuan.com/education
$ mkdir web/static
```

**`web/static`目录下包括四个子目录，分别为：**

- `web/static/css` ：用于存放控制页面布局及显示样式所需的 `CSS` 文件
- `web/static/js` ：用于存放编写的与用户交互的 `JavaScript` 源代码文件
- `web/static/images`：用户存放页面显示所需的所有图片文件
- `web/static/photo`：用于存储添加信息时上传的图片文件

```shell
$ mkdir -p web/static/css
$ mkdir -p web/static/images
$ mkdir -p web/static/js
$ mkdir -p web/static/photo
```



在项目的web目录下新创建一个名为 `tpl`  的目录，用来存放Web应用响应客户端的模板页面

```shell
$ mkdir web/tpl
```

在 `web/tpl` 目录下主要有如下页面：

- `login.html`：用户登录页面
- `index.html`：用户登录成功之后进入的首页面
- `help.html`： 显示帮助信息及相关操作的链接页面
- `query.html`：根据证书编号与姓名查询的页面
- `query2.html`：根据身份证号码查询的页面
- `queryResult.html`：根据不同的查询请求显示查询结果的页面
- `addEdu.html`：添加信息的页面
- `modify.html`：修改信息的页面



## 5.2 相关源码实现

相关源代码请参考：

CSS 部分：

[web/static/css/addEdu.css](https://github.com/kevin-hf/education/blob/master/web/static/css/addEdu.css)

[web/static/css/bootstrap.min.css](https://github.com/kevin-hf/education/blob/master/web/static/css/bootstrap.min.css)

[web/static/css/help.css](https://github.com/kevin-hf/education/blob/master/web/static/css/help.css)

[web/static/css/index.css](https://github.com/kevin-hf/education/blob/master/web/static/css/index.css)

[web/static/css/login.css](https://github.com/kevin-hf/education/blob/master/web/static/css/login.css)

[web/static/css/query.css](https://github.com/kevin-hf/education/blob/master/web/static/css/query.css)

[web/static/css/queryResult.css](https://github.com/kevin-hf/education/blob/master/web/static/css/queryResult.css)

[web/static/css/reset.css](https://github.com/kevin-hf/education/blob/master/web/static/css/reset.css)

JavaScript 部分

[web/static/js/bootstrap.min.js](https://github.com/kevin-hf/education/blob/master/web/static/js/bootstrap.min.js)

[web/static/js/jquery.min.js](https://github.com/kevin-hf/education/blob/master/web/static/js/jquery.min.js)

HTML 页面模板部分：

[web/tpl/addEdu.html](https://github.com/kevin-hf/education/blob/master/web/tpl/addEdu.html)

[web/tpl/help.html](https://github.com/kevin-hf/education/blob/master/web/tpl/help.html)

[web/tpl/index.html](https://github.com/kevin-hf/education/blob/master/web/tpl/index.html)

[web/tpl/login.html](https://github.com/kevin-hf/education/blob/master/web/tpl/login.html)

[web/tpl/modify.html](https://github.com/kevin-hf/education/blob/master/web/tpl/modify.html)

[web/tpl/query.html](https://github.com/kevin-hf/education/blob/master/web/tpl/query.html)

[web/tpl/query2.html](https://github.com/kevin-hf/education/blob/master/web/tpl/query2.html)

[web/tpl/queryResult.html](https://github.com/kevin-hf/education/blob/master/web/tpl/queryResult.html)



## 5.3 照片上传

在添加信息时需要额外实现一个功能－添加照片

使用jQuery Ajax功能实现

HTML代码如下：

```html
            <div class="headImg">
                <div class="uploadImg">
                    <input type="file" name="" value="上传照片" id="file">
                    +
                    <!-- <img src="./images/head.jpg" alt=""> -->
                    <img src="" alt="">
                </div>
                <p>请上传照片(120*160px)</p>
            </div>
          </div>
```

JavaScript代码如下:

```javascript
 		// 上传图片
        $('#file').unbind('change').bind('change',function() {
            event.stopPropagation();
            uploadFile('img');
            return;
        });
        // 头像图片
        var artImg;
        function uploadFile(type) {
            event.stopPropagation();
            let formData = new FormData();
            if( type == "img"){
                formData.append('file', $('#file')[0].files[0]);
            }
            $.ajax({
                url: '/upload',
                type: 'POST',
                cache: false,
                data: formData,
                processData: false,
                dataType: "json",
                contentType: false
            }).done(function (res) {
                if (res.error == "0") {
                    if( type == "img"){
                        $('.uploadImg img').attr('src',res.result.path);
                        $('#photo').val(res.result.path)
                        return artImg = res.result.path;
                    }
                } else {
                    alert("上传失败！" + res.result.msg)
                }
            }).fail(function (res) { });
        }
```



在 `web/controller` 目录下创建一个 `upload.go` 文件

```shell
$ vim web/controller/upload.go
```

`upload.go` 文件主要利用 Ajax完成 照片传功能的，完整代码如下：

```go
/**
  @Author : hanxiaodong
*/

package controller

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"crypto/rand"
	"path/filepath"
	"os"
	"mime"
	"log"
)

func (app *Application) UploadFile(w http.ResponseWriter, r *http.Request)  {

	start := "{"
	content := ""
	end := "}"

	file, _, err := r.FormFile("file")
	if err != nil {
		content = "\"error\":1,\"result\":{\"msg\":\"指定了无效的文件\",\"path\":\"\"}"
		w.Write([]byte(start + content + end))
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		content = "\"error\":1,\"result\":{\"msg\":\"无法读取文件内容\",\"path\":\"\"}"
		w.Write([]byte(start + content + end))
		return
	}

	filetype := http.DetectContentType(fileBytes)
	//log.Println("filetype = " + filetype)
	switch filetype {
	case "image/jpeg", "image/jpg":
	case "image/gif", "image/png":
	case "application/pdf":
		break
	default:
		content = "\"error\":1,\"result\":{\"msg\":\"文件类型错误\",\"path\":\"\"}"
		w.Write([]byte(start + content + end))
		return
	}

	fileName := randToken(12)	// 指定文件名
	fileEndings, err := mime.ExtensionsByType(filetype)	// 获取文件扩展名
	//log.Println("fileEndings = " + fileEndings[0])
	// 指定文件存储路径
	newPath := filepath.Join("web", "static", "photo", fileName + fileEndings[0])
	//fmt.Printf("FileType: %s, File: %s\n", filetype, newPath)

	newFile, err := os.Create(newPath)
	if err != nil {
		log.Println("创建文件失败：" + err.Error())
		content = "\"error\":1,\"result\":{\"msg\":\"创建文件失败\",\"path\":\"\"}"
		w.Write([]byte(start + content + end))
		return
	}
	defer newFile.Close()

	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		log.Println("写入文件失败：" + err.Error())
		content = "\"error\":1,\"result\":{\"msg\":\"保存文件内容失败\",\"path\":\"\"}"
		w.Write([]byte(start + content + end))
		return
	}

	path := "/static/photo/" + fileName + fileEndings[0]
	content = "\"error\":0,\"result\":{\"fileType\":\"image/png\",\"path\":\"" + path + "\",\"fileName\":\"ce73ac68d0d93de80d925b5a.png\"}"
	w.Write([]byte(start + content + end))
	return
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
```



# 6 交互演示

## 6.1 启动Web服务

最后编辑   `main.go` ，以便启动Web界面实现Web应用程序

```shell
$ vim main.go
```

添加如下内容:

```go
import(
	[......]
	"github.com/kongyixueyuan.com/education/web/controller"
	"github.com/kongyixueyuan.com/education/web"
)

func main(){}
	[......]
	
	app := controller.Application{
		Setup: &serviceSetup,
	}
	web.WebStart(app)
}
```

应用项目开发完成后，可以直接启动用来查看效果。在命令提示符中输入 `make` 命令：

```shell
$ make
```

## 6.2 访问页面

项目启动成功之后，打开浏览器访问:   [htt://localhost:9000/](http://localhost:9000/)

根据访问的URL地址系统自动响应登录页面

![login](./img/html_login.png)



输入管理员账号及密码登录验证成功，则进入系统首页面

![login](./img/html_index.png)

在首页面中点击 `查询范围`链接，进入 `help`页面，

![login](./img/html_help.png)

点击添加学历信息链接进入，添加学历信息页面

![login](./img/html_addEdu.png)



根据学历证书编号与姓名查询页面

![login](./img/html_queryResultbycert.png)

根据学历证书编号与姓名查询结果页面

![certNoAndName](./img/html_query.png)

根据身份证号码查询页面

![login](./img/html_query2.png)

根据身份证号码查询页面查询结果页面

![login](./img/html_queryResultbyentityid.png)

编辑页面

![login](./img/html_modify.png)

编辑成功自动跳转到根据身份证号码查询结果页面

![login](./img/html_modifyResult.png)



项目完整源代码，请 [点击此处](https://github.com/kevin-hf/education)

