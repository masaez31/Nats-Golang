package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"strings"
	"time"
)

func main() {
	
	natsAddr := flag.String("nats", "nats://localhost:4222", "Address of the NATS server")
	channel := flag.String("channel", "chat", "Name of the chat channel")
	username := flag.String("name", "anonymous", "Your name in the chat")
	flag.Parse()

	
	nc, err := nats.Connect(*natsAddr)
	if err != nil {
		fmt.Printf("Error connecting to NATS server: %v\n", err)
		os.Exit(1)
	}
	defer nc.Close()

	fmt.Printf("Connected to NATS server at %s\n", *natsAddr)
	fmt.Printf("Joined channel: %s as %s\n", *channel, *username)

	
	js, err := nc.JetStream()
	if err != nil {
		fmt.Printf("Error enabling JetStream: %v\n", err)
		os.Exit(1)
	}

	
	streamName := *channel
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{streamName},
	})
	if err != nil && !strings.Contains(err.Error(), "stream already exists") {
		fmt.Printf("Error creating stream: %v\n", err)
		os.Exit(1)
	}

	
	fmt.Println("Fetching messages from the last hour...")
	subCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	subOpts := []nats.SubOpt{
		nats.StartTime(now.Add(-1 * time.Hour)),
	}

	sub, err := js.SubscribeSync(streamName, subOpts...)
	if err == nil {
		for {
			msg, err := sub.NextMsgWithContext(subCtx)
			if err != nil {
				break
			}
			fmt.Println(string(msg.Data))
		}
	}

	
	liveSub, err := js.Subscribe(streamName, func(msg *nats.Msg) {
		
		if strings.Contains(string(msg.Data), fmt.Sprintf("%s:", *username)) {
			return 
		}
		fmt.Println(string(msg.Data))
	})
	if err != nil {
		fmt.Printf("Error subscribing to channel: %v\n", err)
		os.Exit(1)
	}
	defer liveSub.Unsubscribe()

	
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		if !scanner.Scan() {
			break
		}
		message := scanner.Text()
		if strings.ToLower(message) == "exit" {
			break
		}

		fullMessage := fmt.Sprintf("%s: %s", *username, message)
		if _, err := js.Publish(streamName, []byte(fullMessage)); err != nil {
			fmt.Printf("Error publishing message: %v\n", err)
		}
	}

	fmt.Println("Exiting chat...")
}
