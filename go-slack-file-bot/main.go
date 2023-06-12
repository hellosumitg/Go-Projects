package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		// Print command events
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	// Slack Bot Configuration
	os.Setenv("SLACK_BOT_TOKEN", "xoxb-5430116027872-5399874769606-KieRU9XnNajUMmHs2K1zNhqB")
	os.Setenv("SLACK_APP_TOKEN", "xapp-1-A05CN667KA4-5406506655986-52d483c935807b1d14d133604454c17352d324a0639b3702d3a6d82961212d4d")
	os.Setenv("CHANNEL_ID", "C05BVDXAATF")

	// Create a new Slacker bot instance
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	// Goroutine to print command events
	go printCommandEvents(bot.CommandEvents())

	// Register the "ping" command
	bot.Command("ping", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("pong")
		},
	})

	// ----------------------YOB Calculator---------------------
	// Register the "my yob is <year>" command
	bot.Command("my yob is <year>", &slacker.CommandDefinition{
		Description: "YOB calculator",
		Examples:    []string{"my yob is 2020"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				println("error")
				return
			}
			age := 2023 - yob
			r := fmt.Sprintf("age is %d", age)
			response.Reply(r)
		},
	})

	// ---------------------Slack File Uploader---------------------
	// Create a new Slack client
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	// Configure channel and file parameters
	channelArr := []string{os.Getenv("CHANNEL_ID")}
	fileArr := []string{"free.pdf", "sample.pdf", "ZIPL.pdf"}

	// Upload files to the specified channel
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
		fmt.Printf("Name: %s, URL: %s\n", file.Name, file.URL)
	}

	// Start listening for Slack events
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
