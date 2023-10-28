package main

import (
    "flag"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "httpdump/pkg/log"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    port := flag.String("p", "8080", "http port")
    respType := flag.String("t", "html", "response Content-Type: html, json, txt")
    logPath := flag.String("l", "", "log path")

    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "Options:\n")
        flag.PrintDefaults()
    }
    flag.Parse()

    // init logger
    log.Init(*logPath, "trace")

    app := fiber.New(fiber.Config{
        DisableStartupMessage:     true,
        DisableDefaultDate:        true,
        DisableHeaderNormalizing:  true,
        DisableDefaultContentType: true,
    })
    // Match any request
    app.Use(func(c *fiber.Ctx) error {
        log.Trace(c.IP(), c.Request().String())
        err := c.Next()
        // print response
        // log.Debug(c.Response().String())
        return err
    })
    app.Get("/", func(c *fiber.Ctx) error {
        switch *respType {
        case "json":
            return c.JSON(fiber.Map{"hello": "world"})
        case "txt":
            c.Set("Content-Type", "text/plain; charset=utf-8")
            return c.SendString("hello world")
        default:
            c.Set("Content-Type", "text/html; charset=utf-8")
            return c.SendString("hello world")
        }
    })
    // fmt.Println(*port)
    go log.Fatal(app.Listen(fmt.Sprintf(":%s", *port)))

    watchSignal()
}

func watchSignal() {
    log.Info("Server pid:", os.Getpid())

    sigs := make(chan os.Signal, 1)
    // https://pkg.go.dev/os/signal
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
    for {
        // 没有信号就阻塞，从而避免主协程退出
        sig := <-sigs
        log.Info("Get signal:", sig)
        switch sig {
        default:
            log.Info("Stop")
            return
        }
    }
}
