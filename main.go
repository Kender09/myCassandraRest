package main

import(
  "github.com/gin-gonic/gin"
  "github.com/gocql/gocql"
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

type cqlConfig gocql.ClusterConfig

func chaincodeQuery(chainc Chaincode) bool {
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
    }
    case "query": {
      status = chaincodeQuery(chain)
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
  cqlConfig := gocql.NewCluster("127.0.0.1")
  cqlConfig.Keyspace = "fabric"
  cqlConfig.Consistency = gocql.Quorum
  cqlConfig.Port = 9042
  r := gin.Default()

  r.POST("/chaincode", postChaincode)
  r.Run(":7050")
}
