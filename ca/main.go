package main

import(
  "github.com/gin-gonic/gin"
  "fmt"
)

var ScUsers = []string{"jim", "lukas", "diego", "binhn"}

func getCA (c *gin.Context) {
  sc := c.Param("secureContext")
  fmt.Println(sc)
  for _, v:= range ScUsers {
    if v == sc {
      c.JSON(200, gin.H{"status": "OK"})
      return
    }
  }
  c.JSON(401, gin.H{"status": "NG"})
}

func main() {
  r := gin.Default()
  r.GET("/ca/:secureContext", getCA)

  r.Run(":7052")
}
