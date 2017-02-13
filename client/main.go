package main

import (
    "context"
    "log"

    "google.golang.org/grpc"

    "github.com/vendasta/event-store/pb"
)

func main() {
    conn, err := grpc.Dial("minikube.xxx:31957", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewEventStoreClient(conn)

    r, err := c.ListEvents(context.Background(), &pb.ListEventsRequest{AggregateType:pb.AggregateType(pb.AggregateType_CAMPAIGN)})
    if err != nil {
        log.Fatalf("could not get stats: %v", err)
    }

    log.Printf("Topics %s", r)
}