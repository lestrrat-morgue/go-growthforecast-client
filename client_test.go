package growthforecast

import (
  "testing"
)

func CreateClient() (*Client) {
  return NewClient("gf.ops.dev.livedoor.net", 80)
}

func TestSanity(t *testing.T) {
  // Placeholder test to check that we can at least compile this library
  c := NewClient("dummy", 80)
  if c == nil {
    t.Errorf("Could not create client")
  }
}