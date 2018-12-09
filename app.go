package main

import (
  "os"  // used to get Env var
  "github.com/joho/godotenv"
  "fmt"
  "encoding/json"
  "io/ioutil"
  "log"
  "net/http"
  "strconv"
)

type Oncall struct {
  Oncall_info []Schedule `json:"oncalls"`
}

type Schedule struct {
  Policy Policy `json:"escalation_policy"`
  Level  int    `json:"escalation_level"`
  User   User   `json:"user"`
}

type Policy struct {
  Team string `json:"summary"`
}

type User struct {
  UserName string `json:"summary"`
}

func main() {
  // err := godotenv.Load()
  err := godotenv.Load("/Users/smarman/go/src/github.com/pagerDuty/.env")
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  auth_token := os.Getenv("PAGER_DUTY_TOKEN")  // sets PagerDuty auth token

  request, _ := http.NewRequest("GET", "https://api.pagerduty.com/oncalls?limit=100", nil)
  request.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
  request.Header.Set("Authorization", "Token token=" + auth_token)

  resp, err := http.DefaultClient.Do(request)
  if err != nil {
    log.Fatal(err)
  }

  defer resp.Body.Close()

  body, _ := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
  }

  var schedules Oncall
  err = json.Unmarshal(body, &schedules)
  if err != nil {
    log.Fatal(err)
  }

  for i :=0; i < len(schedules.Oncall_info); i ++ {
    if schedules.Oncall_info[i].Level == 1 {
      fmt.Println("Team:  " + schedules.Oncall_info[i].Policy.Team)
      fmt.Println("Level: " + strconv.Itoa(schedules.Oncall_info[i].Level))
      fmt.Println("User:  " + schedules.Oncall_info[i].User.UserName)
      fmt.Println("***********************")
    }
  }

}
