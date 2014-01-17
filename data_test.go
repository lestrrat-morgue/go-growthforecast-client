package growthforecast

import(
  "testing"
)

func TestComplexGraph(t *testing.T) {
  g := NewComplexGraph()
  if ! g.Complex {
    t.Errorf("g.Complex should be true! (%+v)", g)
  }

  g.ServiceName = "acme"
  g.SectionName = "motor"
  g.GraphName   = "visitors"

  path := g.GetPath()
  if path != "acme/motor/visitors" {
    t.Errorf("Path expected 'acme/motor/visitors', got '%s'", path)
  }
}

func TestGraph(t *testing.T) {
  g := NewGraph()

  if g.Complex {
    t.Errorf("g.Complex should be false! (%+v)", g)
  }

  g.ServiceName = "acme"
  g.SectionName = "motor"
  g.GraphName   = "visitors"

  path := g.GetPath()
  if path != "acme/motor/visitors" {
    t.Errorf("Path expected 'acme/motor/visitors', got '%s'", path)
  }
}
