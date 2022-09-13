
# Reconmap agent

The Reconmap agent allows clients (such as the Web client) to launch commands remotely, open interactive browser terminals, and receive push notifications. 

This is a component of many in the [Reconmap's architecture](https://reconmap.org/development/architecture.html).

## Runtime requirements 

- Docker
- Make
- Linux/Macos operating system due to dependency on OS dependent syscalls

## How to run

```shell
JWT_SECRET="One long key of your own" REDIS_HOST=localhost REDIS_PORT=6379 REDIS_PASSWORD=REconDIS ./reconmapd
```

## Development requirements 

- Golang
