# refernce

https://golang2.eddycjy.com/posts/ch4/01-tcp/
https://github.com/go-programming-tour-book/chatroom

# structure

- cmd: Used to Store main.main
- logic: Used to store the core business logic code of the project
- server: Used to store the controller code
- template: Used to store static template files

# benchmark

## build

```
go build ./cmd/benchmark
```

## run

```
Usage of benchmark:
  -l duration
        User login time interval (default 5s)
  -m duration
        User sending message interval (default 1m0s)
  -u int
        Number of logged-in users (default 500)
```
