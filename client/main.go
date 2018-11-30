/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"../remote_terminal_api"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := remote_terminal_api.NewRemoteTerminalClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	lineInput := ""
	for {
		fmt.Printf("\nrt> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan() // use `for scanner.Scan()` to keep reading
		lineInput = scanner.Text()
		if lineInput == "exit" {
			break
		}
		r, err := c.ExecuteCommand(ctx, &remote_terminal_api.CommandRequest{Cmd: lineInput})
		if err != nil {
			fmt.Printf("Failed to execute: %v", err)
		} else {
			fmt.Printf(r.Response)
		}
	}
}
