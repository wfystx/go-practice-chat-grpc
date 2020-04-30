## Chat Room

### What's this repo
Chat room implementation written in Golang with the help of **gRPC**. Learn the language and also learn some great libs!

### What I used

- Cobra (CLI)
- Go-ini (Configs)
- gRPC
- Protobuf

### How to run

```bash
go run main.go help
```

- To get specific instructions

```bash
go run main.go server
```

- To start the server fist.

```bash
go run main.go client -n Freddie
```

- Start as many client as you want and do remember to use ```-n``` to indicate a client's name.

- After getting into the chatroom, a client can:
  - "quit" — to quit the room.
  - Or just type whatever you wanna say.

### DEMO
![Image](https://i.ibb.co/1Gbr1mQ/go-chat-grpc.gif)
