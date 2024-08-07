package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func handleStream(c *fiber.Ctx) error {
	pm := NewProviderManager()
	c.Set("content-Type", "text/plain; charset=utf-8")
	c.Set("Transfer-Encoding", "chunked")

	// Create channels for handling different events
	done := make(chan struct{})
	errCh := make(chan error, 1)

	// Goroutine for streaming data
	go func() {
		defer close(done)
		ctx, cancel := context.WithCancel(c.Context())
		defer cancel()
		for i := 0; i < 25; i++ {
			select {
			case <-c.Context().Done(): // Handle client disconnect
				fmt.Println("Client disconnected")
				return
			case err := <-errCh: // Handle errors from streaming
				fmt.Printf("Error encountered: %v\n", err)
				return
			default:
				data, err := pm.StreamData(ctx, i)
				if err != nil {
					errCh <- fmt.Errorf("%w", err)
					return
				}

				// Simulate data generation
				if _, err := c.Write(data); err != nil {
					errCh <- fmt.Errorf("write error: %w", err)
					return
				}
				c.WriteString("\n")

				time.Sleep(10 * time.Microsecond)
			}
		}
		fmt.Println("[*] streaming completed for the request.")
	}()

	// Wait for completion, disconnection, or errors
	select {
	case <-done:
		return nil
	case err := <-errCh:
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Internal server error: %v", err))
	}
}
