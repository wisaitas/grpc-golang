package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/wisaitas/grpc-golang/internal/grpcservice/protogenerate/pushmessage"
)

type PushMessageService interface {
	SubscribePushMessages(req *pushmessage.PushRequest, stream pushmessage.PushMessageService_SubscribePushMessagesServer) error
}

type pushMessageService struct {
	pushmessage.UnimplementedPushMessageServiceServer
}

func NewPushMessageService() *pushMessageService {
	return &pushMessageService{}
}

func (s *pushMessageService) SubscribePushMessages(req *pushmessage.PushRequest, stream pushmessage.PushMessageService_SubscribePushMessagesServer) error {
	log.Printf("Client %s subscribed to topic: %s", req.GetClientId(), req.GetTopic())

	// สร้าง context สำหรับ cancel การ streaming
	ctx := stream.Context()

	// Channel สำหรับจัดการ message queue
	messageChan := make(chan *pushmessage.PushMessage, 100)

	// เริ่ม goroutine สำหรับ generate messages
	go s.generateMessages(ctx, req, messageChan)

	// ส่ง messages ไปยัง client
	for {
		select {
		case <-ctx.Done():
			log.Printf("Client %s disconnected from topic: %s", req.GetClientId(), req.GetTopic())
			return nil
		case msg := <-messageChan:
			if err := stream.Send(msg); err != nil {
				log.Printf("Error sending message to client %s: %v", req.GetClientId(), err)
				return status.Errorf(codes.Internal, "Failed to send message: %v", err)
			}
			log.Printf("Sent message to client %s: %s", req.GetClientId(), msg.GetContent())
		}
	}
}

func (s *pushMessageService) generateMessages(ctx context.Context, req *pushmessage.PushRequest, messageChan chan<- *pushmessage.PushMessage) {
	defer close(messageChan)

	ticker := time.NewTicker(5 * time.Second) // ส่งทุก 5 วินาที
	defer ticker.Stop()

	messageCount := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			messageCount++

			// สร้าง message ตามประเภทต่างๆ
			var msgType pushmessage.MessageType
			var content string

			switch messageCount % 4 {
			case 0:
				msgType = pushmessage.MessageType_INFO
				content = fmt.Sprintf("📢 Information message #%d for %s", messageCount, req.GetClientId())
			case 1:
				msgType = pushmessage.MessageType_SUCCESS
				content = fmt.Sprintf("✅ Success notification #%d - Operation completed", messageCount)
			case 2:
				msgType = pushmessage.MessageType_WARNING
				content = fmt.Sprintf("⚠️ Warning alert #%d - Please check your settings", messageCount)
			case 3:
				msgType = pushmessage.MessageType_ERROR
				content = fmt.Sprintf("❌ Error report #%d - System issue detected", messageCount)
			}

			message := &pushmessage.PushMessage{
				Id:        fmt.Sprintf("msg_%d_%d", time.Now().Unix(), messageCount),
				Content:   content,
				Timestamp: time.Now().Format("2006-01-02 15:04:05"),
				Topic:     req.GetTopic(),
				Type:      msgType,
			}

			select {
			case messageChan <- message:
			case <-ctx.Done():
				return
			default:
				// ถ้า channel เต็ม ให้ skip message นี้
				log.Printf("Message channel full, skipping message for client %s", req.GetClientId())
			}
		}
	}
}
