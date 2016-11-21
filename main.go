package main

import(
  "github.com/gin-gonic/gin"
  "github.com/gocql/gocql"
  "fmt"
)


type Chaincode struct{
  Jsonrpc string `json:"jsonrpc"`
  Method  string `json:"method"`
  Params  Param  `json:"params"`
  ID      int    `json:"id"`
}

type Param struct{
  Type        int         `json:"type"`
  ChaincodeID ChaincodeID `json:"chaincodeID"`
  CtorMsg     CtorMsg     `json:"ctorMsg"`
  SecureContext string    `json:"secureContext"`
}

type ChaincodeID struct{
  Name string `json:"name"`
}

type CtorMsg struct{
  Args []string `json:"args"`
}

var CqlConfig *gocql.ClusterConfig

func chaincodeInvoke(chain Chaincode) bool {
  session, _ := CqlConfig.CreateSession()
  defer session.Close()
  if len(chain.Params.CtorMsg.Args) != 4 {
    return false
  }

  from_campany := chain.Params.CtorMsg.Args[1]
  to_campany := chain.Params.CtorMsg.Args[2]
  money := chain.Params.CtorMsg.Args[3]

  from_query := "UPDATE assets SET money = money - ? WHERE campany = ?"
  to_query := "UPDATE assets SET money = money + ? WHERE campany = ?"

  if err := session.Query(from_query, money, from_campany).Exec(); err != nil {
    fmt.Println(err)
    return false
  }

  if err := session.Query(to_query, money, to_campany).Exec(); err != nil {
    fmt.Println(err)
    return false
  }

  return true
}

func postChaincode (c *gin.Context) {
  var chain Chaincode
  if c.BindJSON(&chain) != nil {
    c.JSON(401, gin.H{"status": "unauthorized"})
    return
  }

  var status bool

  switch chain.Method {
    case "invoke": {
      status = chaincodeInvoke(chain)
    }
    case "query": {
    }
    case "deploy": {
    }
    default: {
      c.JSON(401, gin.H{"status": "unauthorized"})
      return
    }
  }

  if status {
    c.JSON(200, gin.H{"status": "OK"})
  } else {
    c.JSON(401, gin.H{"status": "unauthorized"})
  }
}

func main() {
  CqlConfig = gocql.NewCluster("127.0.0.1")
  CqlConfig.Keyspace = "fabric"
  CqlConfig.Consistency = gocql.Quorum
  CqlConfig.Port = 9042
  r := gin.Default()

  r.POST("/chaincode", postChaincode)
  r.Run(":7050")
}
