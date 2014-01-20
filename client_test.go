package growthforecast

import (
  "bufio"
  "errors"
  "fmt"
  "io/ioutil"
  "math/rand"
  "net"
  "os"
  "os/exec"
  "reflect"
  "syscall"
  "testing"
  "time"
)

var GF_PORT int

func init () {
  rand.Seed(time.Now().UnixNano())
}

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

func isGraphEqual(t *testing.T, l *Graph, r *Graph) bool {
  return reflect.DeepEqual(l, r)
}

func isGraphValid(t *testing.T, g *Graph) bool {
  if g.Complex {
    t.Logf("Fetched graph should not be marked as complex")
    return false
  }

  return true
}

func assertGraphFetch(t *testing.T, g *Graph, c *Client) {
  fetched, err := c.GetGraph(g.Id)
  if err != nil {
    t.Fatalf("Failed to fetch graph by id: %s", err)
  }
  if ! isGraphValid(t, fetched) {
    t.Fatalf("Graph is not valid")
  }
  // TODO why is MD5 not present in fetched?
  fetched.MD5 = g.MD5
  if ! isGraphEqual(t, g, fetched) {
    t.Logf("Expected %+v", g)
    t.Logf("Got %+v", fetched)
    t.Fatalf("Graph fetched by id does not match")
  }

  fetched, err = c.GetGraphByPath(g.GetPath())
  if err != nil {
    t.Fatalf("Failed to fetch graph by path: %s", err)
  }
  if ! isGraphValid(t, fetched) {
    t.Fatalf("Graph is not valid")
  }
  // TODO why is MD5 not present in fetched?
  fetched.MD5 = g.MD5
  if ! isGraphEqual(t, g, fetched) {
    t.Fatalf("Graph fetched by path does not match")
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

  graphs := []*Graph {}

  c := CreateClient()
  g := NewGraph()
  g.ServiceName = "acme"
  g.SectionName = "motor"
  g.GraphName   = "oil"
  ret, err := c.CreateGraph(g)
  if err != nil {
    t.Fatalf("Failed to create graph: %s", err)
  }

  if g.GetPath() != ret.GetPath() {
    t.Fatalf(
      "Paths do not match. Expected '%s', got '%s'",
      g.GetPath(),
      ret.GetPath(),
    )
  }

  assertGraphFetch(t, ret, c)
  graphs = append(graphs, ret)

  err = c.Post(ret.GetPath(), &GraphData{ Number: 1 })
  if err != nil {
    t.Errorf(
      "Failed to post data to %s",
      ret.GetPath(),
    )
  }

  g = NewGraph()
  g.ServiceName = "acme"
  g.SectionName = "motor"
  g.GraphName = "colored"
  g.Color = "#ABCDEF"
  ret, err = c.CreateGraph(g)
  if err != nil {
    t.Fatalf("Failed to create graph with color: %s", err)
  }

  assertGraphFetch(t, ret, c)
  graphs = append(graphs, ret)

  if ret.Color != "#ABCDEF" {
    t.Fatalf("Expected color to be the #ABCDEF, but got %s", ret.Color)
  }

  list, err := c.GetGraphList()
  if err != nil {
    t.Errorf("Failed to fetch graph list: %s", err)
  } else {
    if len(list) != 2 {
      t.Errorf("Expected to receive 2 graphs, got %d", len(list))
    }

    for _, v := range list {
      path := v.GetPath()
      switch (path) {
      case "acme/motor/oil", "acme/motor/colored":
        // no op
      default:
        t.Errorf("Unexpected graph %s retrieved from GetGraphList", path)
      }
    }
  }

  cg := NewComplexGraph()
  cg.ServiceName  = "acme"
  cg.SectionName  = "motor"
  cg.GraphName    = "complex"
  for _, g = range graphs {
    cg.Data = append(cg.Data, ComplexGraphData {
      GraphId:  g.Id,
      Gmode:    g.Gmode,
      Type:     g.Type,
    })
  }
  _, err = c.CreateComplex(cg)
  if err != nil {
    t.Fatalf("Failed to create complex graph: %s", err)
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

  for p := 50000 + rand.Intn(10000); p < 65535; p++ {
    l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p))
    if err == nil {
      l.Close()
      GF_PORT = p
      break
    }
  }

  if GF_PORT == 0 {
    return nil, errors.New("Could not find an empty port")
  }

  dataDir, err := ioutil.TempDir("", "go-growthforecast-client-test")
  if err != nil {
    return nil, errors.New(
      fmt.Sprintf(
        "Could not create temporary directory: %s",
        err,
      ),
    )
  }

  cmd := exec.Command(
    path,
    fmt.Sprintf("--port=%d", GF_PORT),
    "--host=127.0.0.1",
    fmt.Sprintf("--data-dir=%s", dataDir),
  )

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

  t.Logf("Starting growthforecast.pl as: %v", cmd.Args)
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
    cmd.Process.Signal(syscall.SIGTERM)
    os.RemoveAll(dataDir)
  }
  defer func() {
    if err := recover(); err != nil {
      killproc()
      panic(err)
    }
  }()

  started := false
  addr    := fmt.Sprintf("127.0.0.1:%d", GF_PORT)
  time.Sleep(1 * time.Second)
  for timeout := time.Now().Add(10 * time.Second); timeout.After(time.Now()); {
    _, err := net.Dial("tcp", addr)
    if err == nil {
      t.Logf("Successfully connected to %s", addr)
      started = true
      break
    }
    t.Logf("Failed to connect to %s: %s", addr, err)
    time.Sleep(1 * time.Second)
  }

  if ! started {
    killproc()
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

  t.Logf("growthforecast.pl started on port %d", GF_PORT)
  go func() {
    cmd.Wait()
    t.Logf("Done wait")
  }()

  return killproc, nil
}