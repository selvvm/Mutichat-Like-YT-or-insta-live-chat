# TCP Chat Server in Go

This project is a simple TCP chat server implemented in Go. It allows multiple clients to connect and communicate with each other in real-time via a TCP connection.

## Features

- **Concurrency**: Utilizes goroutines to handle multiple client connections concurrently, ensuring optimal server performance.
- **Error Handling**: Implements robust error handling mechanisms to gracefully manage unexpected situations and maintain server stability.
- **Safe Mode**: Provides an optional safe mode to redact sensitive information and enforce bans for malicious behavior, enhancing server security.
- **Efficient Communication**: Leverages channels for efficient communication and synchronization between goroutines, ensuring smooth message transmission.

## How to Use

1. Clone the repository
2. Build the project
3. Run the server
4. Connect clients to the server using a TCP client such as Telnet or Netcat:
5. Start chatting!

## Configuration

- **Port**: By default, the server listens on port 6969. You can change this port by modifying the `Port` constant in the code.
- **Safe Mode**: Set the `SafeMode` constant to `true` to enable safe mode, which redacts sensitive information. Set it to `false` to disable safe mode.

## Contributing

Contributions are welcome! If you find any bugs or have suggestions for improvements, feel free to open an issue or submit a pull request.


