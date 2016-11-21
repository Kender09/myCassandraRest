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

func chaincodeDeploy(c *gin.Context, chain Chaincode) {
  session, _ := CqlConfig.CreateSession()
  defer session.Close()

  query := "UPDATE assets SET money = money + ? WHERE campany = ?"
  init_args := chain.Params.CtorMsg.Args[1:]
  for i := 0; i < len(init_args) / 2; i++ {
    campany := chain.Params.CtorMsg.Args[i*2]
    money := chain.Params.CtorMsg.Args[i*2 + 1]
    if err := session.Query(query, money, campany).Exec(); err != nil {
      fmt.Println(err)
      c.JSON(401, gin.H{"status": "err", "message": fmt.Sprint(err)})
      return
    }
  }
  c.JSON(200, gin.H{"status": "OK"})
}

func chaincodeInvoke(c *gin.Context, chain Chaincode) {
  session, _ := CqlConfig.CreateSession()
  defer session.Close()
  if len(chain.Params.CtorMsg.Args) != 4 {
    c.JSON(401, gin.H{"status": "Incorrect number of arguments. Expecting 4"})
    return
  }

  from_campany := chain.Params.CtorMsg.Args[1]
  to_campany := chain.Params.CtorMsg.Args[2]
  money := chain.Params.CtorMsg.Args[3]

  from_query := "UPDATE assets SET money = money - ? WHERE campany = ?"
  to_query := "UPDATE assets SET money = money + ? WHERE campany = ?"

  if err := session.Query(from_query, money, from_campany).Exec(); err != nil {
    fmt.Println(err)
    c.JSON(401, gin.H{"status": "err", "message": fmt.Sprint(err)})
    return
  }

  if err := session.Query(to_query, money, to_campany).Exec(); err != nil {
    fmt.Println(err)
    c.JSON(401, gin.H{"status": "err", "message": fmt.Sprint(err)})
    return
  }

  c.JSON(200, gin.H{"status": "OK"})
}

func chaincodeQuery(c *gin.Context, chain Chaincode) {
  session, _ := CqlConfig.CreateSession()
  defer session.Close()
  if len(chain.Params.CtorMsg.Args) != 2 {
    c.JSON(401, gin.H{"status": "Incorrect number of arguments. Expecting 2"})
    return
  }

  campany := chain.Params.CtorMsg.Args[1]
  var money string

  query := "SELECT money FROM assets WHERE campany = ?"
  if err := session.Query(query, campany).Scan(&money); err != nil {
    fmt.Println(err)
    c.JSON(401, gin.H{"status": "err", "message": fmt.Sprint(err)})
    return
  }

  c.JSON(200, gin.H{"status": "OK", "message": money})
}

func postChaincode (c *gin.Context) {
  var chain Chaincode
  if c.BindJSON(&chain) != nil {
    c.JSON(401, gin.H{"status": "unauthorized"})
    return
  }

  switch chain.Method {
    case "invoke": {
      chaincodeInvoke(c, chain)
    }
    case "query": {
      chaincodeQuery(c, chain)
    }
    case "deploy": {
      chaincodeDeploy(c, chain)
    }
    default: {
      c.JSON(401, gin.H{"status": "unauthorized"})
    }
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
