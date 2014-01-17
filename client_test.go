package growthforecast

import (
  "bufio"
  "errors"
  "fmt"
  "net"
  "os"
  "os/exec"
  "testing"
  "time"
)

var GF_PORT int

func CreateClient() (*Client) {
  return NewClient(fmt.Sprintf("http://127.0.0.1:%d", GF_PORT))
}

func TestSanity(t *testing.T) {
  // Placeholder test to check that we can at least compile this library
  c := CreateClient()
  if c == nil {
    t.Errorf("Could not create client")
  }
}

func TestBasic(t *testing.T) {
  guard, err := startGF(t)
  if err != nil {
    t.Logf("We were unable to start GrowthForecast. Will not run live tests")
    t.Logf("Errors: %s", err)
    return
  }
  defer guard()

  c := CreateClient()
  g := NewGraph()
  g.ServiceName = "acme"
  g.SectionName = "motor"
  g.GraphName   = "oil"
  err = c.CreateGraph(g)

  if err != nil {
    t.Errorf("Failed to create graph: %s", err)
  }
}

func startGF(t *testing.T) (func(), error) {
  path, err := exec.LookPath("growthforecast.pl")
  if err != nil {
    return nil, errors.New(
      fmt.Sprintf(
        "Failed to lookup growthforecast.pl: %s",
        err,
      ),
    )
  }

  // XXX This needs fixing
  GF_PORT = 5125

  cmd := exec.Command(path, "--port=5125", "--host=127.0.0.1")

  stderrpipe, err := cmd.StderrPipe()
  if err != nil {
    return nil, errors.New(
      fmt.Sprintf(
        "Failed to open pipe to stderr: %s",
        err,
      ),
    )
  }
  stdoutpipe, err := cmd.StdoutPipe()
  if err != nil {
    return nil, errors.New(
      fmt.Sprintf(
        "Failed to open pipe to stdout: %s",
        err,
      ),
    )
  }
  pipes := []struct {
    Out *os.File
    Rdr *bufio.Reader
  } {
    { os.Stdout, bufio.NewReader(stdoutpipe) },
    { os.Stderr, bufio.NewReader(stderrpipe) },
  }

  err  = cmd.Start()
  if err != nil {
    return nil, errors.New(
      fmt.Sprintf(
        "Failed to start growthforecast.pl: %s",
        err,
      ),
    )
  }

  killed := false
  killproc := func() {
    killed = true
    t.Logf("Killing growthforecast.pl")
    cmd.Process.Kill()
  }
  defer func() {
    if err := recover(); err != nil {
      killproc()
      panic(err)
    }
  }()

  started := false
  addr    := fmt.Sprintf("127.0.0.1:%d", GF_PORT)
  for timeout := time.Now().Add(10 * time.Second); timeout.After(time.Now()); {
    _, err := net.Dial("tcp", addr)
    if err == nil {
      started = true
      break
    }
    time.Sleep(1 * time.Second)
  }

  if ! started {
    return nil, errors.New(
      fmt.Sprintf(
        "Failed to connect to port %d",
        GF_PORT,
      ),
    )
  }

  for _, p := range pipes {
    go func(out *os.File, in *bufio.Reader) {
      for !killed {
        str, err := in.ReadBytes('\n')
        if str != nil {
          out.Write(str)
          out.Sync()
        }

        if err != nil {
          break
        }
      }
    }(p.Out, p.Rdr)
  }

  t.Logf("growthforecast.pl started")
  go cmd.Wait()

  return killproc, nil
}