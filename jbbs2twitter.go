package main

import (
  "flag"
  "log"
  "time"
  "github.com/ChimeraCoder/anaconda"
  "github.com/mogepiyo/jbbsreader"
)

var (
  category = flag.String("category", "game", "your JBBS board category")
  banchi = flag.String("banchi", "9358", "your JBBS board banchi")
  apiKey = flag.String("api_key", "", "your api key")
  apiSecret = flag.String("api_secret", "", "your api secret")
  accessToken = flag.String("access_token", "", "your access token")
  accessTokenSecret = flag.String("access_token_secret", "", "your access token secret")
  errorRetryDurationMins = flag.Int("error_retry_duration_mins", 10, "how many minutes to wait after an error")
  jbbsRPM = flag.Int("jbbs_rpm", 1, "how many request per minutes to send to JBBS")
)

func main() {
  flag.Parse()

  log.Printf("jbbsRPM: %v\n", *jbbsRPM)

  log.Println("Generating Twitter Client")
  anaconda.SetConsumerKey(*apiKey)
  anaconda.SetConsumerSecret(*apiSecret)
  api := anaconda.NewTwitterApi(*accessToken, *accessTokenSecret)
  log.Println("Generating JBBS Board Client")
  board := jbbsreader.NewBoard(*category, *banchi)
  jbbsreader.SetGlobalRateLimitRPM(*jbbsRPM, *jbbsRPM)
  log.Println("Feeding JBBS to Twitter!")
  for ;; <-time.Tick(time.Duration(*errorRetryDurationMins) * time.Minute) {
    err := feedJBBS2Twitter(board, api)
    if err != nil {
      log.Println(err.Error())
    }
  }
}

func feedJBBS2Twitter(board *jbbsreader.Board, api *anaconda.TwitterApi) error {
  respc, errc := board.FeedNewResponses()
  for resp := range respc {
    log.Printf("Fetched %#v\n", resp)
    tweet, err := api.PostTweet(resp.Content, nil)
    if err != nil {
      // TODO: Cancel FeedNewResponses
      return err
    }
    log.Printf("Sent tweet %#v\n", tweet)
  }
  return <-errc
}
