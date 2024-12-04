package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := config.New(configFile)
	ctx := context.Background()
	defer ctx.Done()

	log := logger.New(config.Logger.Level, config.Logger.Output)
	defer log.Close()

	addr := fmt.Sprintf("%v:%v", config.Endpoint.Host, config.Endpoint.GRPCPort)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(err.Error())
		os.Exit(1) //nolint:gocritic
	}

	client := pb.NewEventsClient(conn)
	log.Info("Sending CREATE request")
	createReq := &pb.NewEventRequest{
		Event: &pb.TransientEvent{
			Title:       "event 1",
			Description: "some description",
			StartTime:   time.Date(2024, 11, 25, 12, 0, 0, 0, time.UTC).Unix(),
			EndTime:     time.Date(2024, 11, 25, 12, 30, 0, 0, time.UTC).Unix(),
			OwnerEmail:  "user1@example.com",
			Notify:      "",
		},
	}
	scalarResponse, err := client.CreateEvent(ctx, createReq)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	event := scalarResponse.GetEvent()
	if event == nil {
		log.Error("Received response but payload is empty")
		os.Exit(1)
	}

	log.Info("Success",
		"Id", event.Id,
		"Title", event.Title,
		"StartTime", event.StartTime,
		"EndTime", event.EndTime,
		"Description", event.Description,
		"OwnerEmail", event.OwnerEmail)

	log.Info("Sending GET request")
	getRequest := &pb.EventIdRequest{Id: event.Id}
	scalarResponse, err = client.GetEvent(ctx, getRequest)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	event = scalarResponse.GetEvent()
	if event == nil {
		log.Error("Received response but payload is empty")
		os.Exit(1)
	}
	log.Info("Success")

	log.Info("Sending UPDATE request")
	updateReq := &pb.UpdateEventRequest{
		Id: event.Id,
		Event: &pb.TransientEvent{
			Title:       event.Title,
			Description: "new description",
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
			OwnerEmail:  event.OwnerEmail,
			Notify:      event.Notify,
		},
	}

	scalarResponse, err = client.UpdateEvent(ctx, updateReq)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	event = scalarResponse.GetEvent()
	if event == nil {
		log.Error("Received response but payload is empty")
		os.Exit(1)
	}

	log.Info("Success",
		"Id", event.Id,
		"Title", event.Title,
		"StartTime", event.StartTime,
		"EndTime", event.EndTime,
		"Description", event.Description,
		"OwnerEmail", event.OwnerEmail)

	log.Info("Sending DELETE request")
	getRequest = &pb.EventIdRequest{Id: event.Id}
	_, err = client.DeleteEvent(ctx, getRequest)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	log.Info("Success")
}
