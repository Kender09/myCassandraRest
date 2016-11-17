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

func postChaincode (c *gin.Context) {
  var json Chaincode
  if c.BindJSON(&json) == nil {
  }
  c.JSON(200, gin.H{
  })
}

func main() {
  router := gin.Default()

  router.POST("/chaincode", postChaincode)
  router.Run(":7050")
}
