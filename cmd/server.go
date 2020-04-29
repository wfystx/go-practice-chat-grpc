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
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"

	"github.com/go-ini/ini"
	"github.com/spf13/cobra"
	pb "go-practice-chat-gRPC/proto"
)


// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
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
		listen, err := net.Listen("tcp", Address)

		if err != nil {
			fmt.Println("failed to listen: %v", err)
		}

		s := grpc.NewServer()

		// 注册HelloService
		pb.RegisterChatRoomServer(s, chatroomService{})

		fmt.Println("Listen on " + Address)

		err = s.Serve(listen)
		if err != nil {
			fmt.Println("Serve error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type chatroomService struct {

}

var Server = chatroomService{}

var clientMap = make(map[string]chan string)

func (cS chatroomService) Join(ctx context.Context, req *pb.JoinRequest) (*pb.Response, error) {
	clientName := req.Name
	msg := clientName + " has joined"
	for _, v := range clientMap {
		v <- msg
	}
	clientMap[clientName] = make(chan string)
	return &pb.Response{}, nil
}

func (cS chatroomService) Quit(cxt context.Context, req *pb.QuitRequest) (*pb.Response, error){
	clientName := req.Name
	msg := clientName + " just left"
	delete(clientMap, clientName)
	for _, v := range clientMap {
		v <- msg
	}
	return &pb.Response{}, nil
}

func (cS chatroomService) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	clientName := req.Name
	var msg string
	select {
	case msg = <- clientMap[clientName]:
	}
	return &pb.ReadResponse{Msg:msg}, nil
}

func (cS chatroomService) Send(ctx context.Context, req *pb.SendRequest) (*pb.Response, error) {
	clientName := req.Name
	msg := clientName + ">>> " + req.Msg
	for _, v := range clientMap {
		v <- msg
	}
	return &pb.Response{}, nil
}
