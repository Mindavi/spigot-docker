// This application wraps a minecraft server to run in docker
// It catches the SIGINT or SIGTERM signal (docker sends sigterm)
// and passes it to the server as stop command.
// This ensures clean shutdown of the minecraft server.
package main

import (
  "bufio"
  "flag"
  "fmt"
  "io"
  "log"
  "os"
  "os/exec"
  "os/signal"
  "strings"
  "syscall"
)

// If any error occurs, stop the world
func handleError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

// Wait for the signal that tells the server to stop and send stop command to server
func stopServer(done chan bool, pipe io.WriteCloser) {
  log.Print("Waiting for stop signal")
  <-done
  fmt.Fprintf(pipe, "stop\n")
}


func main() {
  var javaPath string
  var javaFlags string
  var minecraftPath string
  flag.StringVar(&javaPath, "j", "/usr/bin/java", "java path")
  flag.StringVar(&javaFlags, "f", "", "java flags")
  flag.StringVar(&minecraftPath, "m", "/usr/bin/spigot.jar", "Minecraft server path")
  flag.Parse()

  flagsSplit := strings.Split(javaFlags, " ")
  flagsSplit = append(flagsSplit, "-jar", minecraftPath)
  cmd := exec.Command(javaPath, flagsSplit...)
  log.Print("Created command ", cmd.Args)
  out, err := cmd.StdoutPipe()
  handleError(err)
  stderr, err := cmd.StderrPipe()
  handleError(err)
  stdin, err := cmd.StdinPipe()
  err = cmd.Start()
  handleError(err)
  signals := make(chan os.Signal, 1)
  done := make(chan bool, 1)
  signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

  go func() {
    sig := <-signals
    log.Print(sig)
    done <- true
  }()
  defer cmd.Wait()
  scanner := bufio.NewScanner(out)
  errscanner := bufio.NewScanner(stderr)
  go readApp(scanner)
  go readApp(errscanner)
  go stopServer(done, stdin)
}

func readApp(scanner *bufio.Scanner) {
  for scanner.Scan() {
    log.Print(scanner.Text())
  }
  if err := scanner.Err(); err != nil {
    log.Print(err)
  }
}
