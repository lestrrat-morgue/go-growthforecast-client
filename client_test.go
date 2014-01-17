package growthforecast

import (
  "errors"
  "fmt"
  "os/exec"
  "testing"
)

func CreateClient() (*Client) {
  return NewClient("http://127.0.0.1")
}

func TestSanity(t *testing.T) {
  // Placeholder test to check that we can at least compile this library
  c := CreateClient()
  if c == nil {
    t.Errorf("Could not create client")
  }
}

func TestBasic(t *testing.T) {
  guard, err := startGF()
  if err != nil {
    t.Logf("We were unable to start GrowthForecast. Will not run live tests")
    t.Logf("Errors: %s", err)
    return
  }

  defer guard()
}

func startGF() (func(), error) {
  path, err := exec.LookPath("growthforecast.pl")
  if err != nil {
    return nil, errors.New(
      fmt.Sprintf(
        "Failed to lookup growthforecast.pl: %s",
        err,
      ),
    )
  }

  cmd := exec.Command(path)
  err  = cmd.Start()
  if err != nil {
    return nil, errors.New(
      fmt.Sprintf(
        "Failed to start growthforecast.pl: %s",
        err,
      ),
    )
  }
  return func() { cmd.Process.Kill() }, nil
}