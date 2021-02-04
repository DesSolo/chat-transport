### Chat Transport
Simple daemon for forwarding messages between chats

#### Usage
```bash
$ ./chat-transport -h
Usage of ./chat-transport:
  -c string
        config file path (default "config.toml")
  -v    show current version
```

#### Supported transports
|Messenger|Chats|Channels|
|--|---|---|
|Telegram|:heavy_check_mark:|:heavy_multiplication_x:|
|VK|:clock130:|:clock130:|

#### Configuration
See [example.conf](https://github.com/DesSolo/chat-transport/blob/master/example/config.conf)