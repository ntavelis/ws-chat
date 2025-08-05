# WebSocket Chat Server

A simple WebSocket-based chat server written in Go that allows multiple clients to connect and chat in real-time.

## Features

- Real-time messaging using WebSockets
- Multiple client support
- Auto-generated funny usernames if none provided
- User connection/disconnection notifications
- Command support (currently `/exit`)
- Clean, simple architecture

## Getting Started

### Prerequisites

- Go 1.24.5 or later

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ntavelis/ws-chat.git
   cd ws-chat
   ```

2. Build the application:
   ```bash
   CGO_ENABLED=0 go build -o chat ./cmd/chat/
   ```

3. Run the server:
   ```bash
   ./chat
   ```

   By default, the server runs on port 3001. You can specify a different port:
   ```bash
   ./chat -port=8080
   ```

## Usage

### Connecting to the Chat

Connect to the WebSocket server using any WebSocket client:

- **URL**: `ws://localhost:3001`
- **With custom username**: `ws://localhost:3001?user=YourUsername`

If no username is provided, the server will generate a funny random username for you.

### Available Commands

- `/exit` - Disconnect from the chat

### Example with wscat

If you have [wscat](https://github.com/websockets/wscat) installed:

```bash
# Connect with custom username
wscat -c 'ws://localhost:3001?user=davel'
```

Example session:
```
Connected (press CTRL+C to quit)
< User davel connected.
You can pass your username by providing it as a query parameter ws://host@port?user=xxx
Type /exit to exit
< There are currently 2 users in the room.
< PeppyLlama68: Connected
< PeppyLlama68: Hey man we need to talk :)
```

### Example with JavaScript

```javascript
const ws = new WebSocket('ws://localhost:3001?user=Alice');

ws.onopen = function() {
    console.log('Connected to chat');
};

ws.onmessage = function(event) {
    console.log('Message:', event.data);
};

ws.send('Hello, everyone!');
```

## License

This project is open source and available under the [MIT License](LICENSE).
