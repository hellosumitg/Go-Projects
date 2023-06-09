package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {

	os.Setenv("SLACK_BOT_TOKEN", "xoxb-5407372789825-5418504853424-yHdKGDFtdyaIMrwK6CzDWJOT")
	os.Setenv("CHANNEL_ID", "C05B65RSF6K")
	api := slack.New(os.Getenv(("SLACK_BOT_TOKEN")))
	channelArr := []string{os.Getenv("CHANNEL_ID")}           // dynamic array in golang(i.e slice)
	fileArr := []string{"free.pdf", "sample.pdf", "zipl.pdf"} // here we can put multiple files

	for i := 0; i < len(fileArr); i++ {
		params := slack.FileUploadParameters{
			Channels: channelArr,
			File:     fileArr[i],
		}
		file, err := api.UploadFile(params)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("Name:%s\n", file.Name)
	}
}
