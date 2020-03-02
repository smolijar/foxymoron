# Go Cool

## `flags`: args parsing
Having built in param parsing is just cool
```go
// go run . -port=8080
port := flag.Int("port", 8000, "Local server port")
flag.Parse()
```

## `flags`: print usage
Print usage generated from flags with description and defaults is cool
```go
flag.Parse()
if *help == true {
    // COOL: print flag usage
    flag.PrintDefaults()
    os.Exit(0)
}
```
```
foxymoron

  -help
        Print help
  -port int
        Local server port (default 8000)
  -token string
        GitLab API token
  -url string
        GitLab URL (default "http://gitlab.com")
```

## Ad-hod go functions
```go
res := make(...)
for _, t := range tasks {
    res = append(res, execute(t))
}
```

Create ad-hoc blocking-to-async functions

```go
res := make([]X)
resChan := make(chan X)
// spawn tasks
for _, t := range tasks {
    go func() {
        resChan <- execute(t)
    }
}
// consume results as they are ready
for i := 0; i < len(tasks); ++i {
    // COOL: use `<-channel` without expicit assignment
    res = append(res, <-res)
}
```

## Loging with `log`
```go
// COOL: you can use default logger from `log` and it outputs by default `2020/01/11 17:35:28 Retireved ...`
// COOL: you can use %v for default formatting
log.Printf("Hello %v", world)
```

## Non-blocking channel read/write
Select takes first ready communication operator.
Optional default is noop ready communication operator.

E.g. run 10 goroutines and let them all behave uniformly, even though you just need (and are going to collect) only the first result. Otherwise you'd need to let write only the first one, have buffered channel or effectively consume all produced messages concurrently not to block the tasks.
```go
// COOL: non-blocking write
select {
case resultChannel <- result:
default:
}
```

## Go VS code intellisense
Enable `go.useLanguageServer` in VS Code to enable renaming via `gopls`

## Uninstalling packages
Run `go mod tidy` to clean unimported packages [ref](https://stackoverflow.com/questions/57186705/how-to-remove-an-installed-package-using-go-modules)

## Decoupled methods

With go's syntax for methods, you can easily code-split classes.

Decoupled methods - controller example https://github.com/swaggo/swag/tree/master/example/celler/controller

## Why is Go succesfull
About simplicity, stability and core principles
[Why is Go succesfull](https://www.youtube.com/watch?v=cQ7STILAS0M)