/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"os"

	"github.com/go-ini/ini"
	"github.com/spf13/cobra"
	pb "go-practice-chat-gRPC/proto"
)

var name string
// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Cfg, err := ini.Load("conf/app.ini")
		if err != nil {
			fmt.Println("Read config file failed")
			return
		}

		Address := Cfg.Section("server").Key("address").String()

		conn, err := grpc.Dial(Address, grpc.WithInsecure())
		if err != nil {
			fmt.Println("Dial server err", err)
			return
		}

		defer conn.Close()

		client := pb.NewChatRoomClient(conn)

		joinReq := &pb.JoinRequest{Name: name}

		_, err = client.Join(context.Background(), joinReq)

		if err != nil {
			fmt.Println("Join err", err)
			return
		}

		go func() {
			for {
				readReq := &pb.ReadRequest{Name: name}
				resp, err := client.Read(context.Background(), readReq)
				if err != nil {
					fmt.Println("read err", err)
					continue
				}
				fmt.Println(resp.Msg)
			}
		}()

		for {
			reader := bufio.NewReader(os.Stdin)
			ln, _, _ := reader.ReadLine()
			msg := string(ln)
			if msg == "quit" {
				req := &pb.QuitRequest{Name: name}
				_, err = client.Quit(context.Background(), req)
				if err != nil {
					fmt.Println("Quit err", err)
				}
				break
			} else {
				req := &pb.SendRequest{Name:name, Msg:msg}
				_, err = client.Send(context.Background(), req)
				if err != nil {
					fmt.Println("Send err", err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVarP(&name, "name", "n", "", "your name")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

