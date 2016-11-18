package main

import(
  "github.com/gin-gonic/gin"
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
  r := gin.Default()

  r.POST("/chaincode", postChaincode)
  r.Run(":7050")
}
