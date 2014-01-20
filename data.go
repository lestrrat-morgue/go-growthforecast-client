package growthforecast

import (
  "fmt"
)

type GraphEssential struct {
  GraphName string        `json:"graph_name"`
  Id int                  `json:"id,omitempty"`
  SectionName string      `json:"section_name"`
  ServiceName string      `json:"service_name"`
}

type GraphData struct {
  Number int              `json:"number"`
  Color string            `json:"color"`
  Mode string             `json:"mode"`
}

type BaseGraph struct {
  GraphEssential
  Complex bool            `json:"complex"`
  CreatedAt string        `json:"created_at,omitempty"`
  Description string      `json:"description,omitempty"`
  Number int              `json:"number"`
  Sort int                `json:"sort"`
  Type string             `json:"type"`
  UpdatedAt string        `json:"updated_at,omitempty"`
}

const (
  MODE_COUNT string = "count"
  MODE_GAUGE        = "gauge"
  MODE_MODIFIED     = "modified"
  GMODE_GAUGE       = "gauge"
  GMODE_SUBTRACT    = "subtract"
  TYPE_AREA         = "AREA"
  TYPE_LINE1        = "LINE1"
  TYPE_LINE2        = "LINE2"
)
const (
  MODE_DEFAULT      = MODE_GAUGE
  GMODE_DEFAULT     = GMODE_GAUGE
  TYPE_DEFAULT      = TYPE_AREA
)

type Graph struct {
  BaseGraph
  Adjust string           `json:"adjust"`
  AdjustVal string        `json:"adjustval"`
  Color string            `json:"color"`
  Gmode string            `json:"gmode"`
  Llimit int              `json:"llmit"`
  MD5 string              `json:"md5"`
  Meta string             `json:"meta"`
  Mode string             `json:"mode"`
  Stype string            `json:"stype"`
  Sllimit int             `json:"sllimit"`
  Sulimit int             `json:"sulimit"`
  Ulimit int              `json:"ulimit"`
  Unit string             `json:"unit"`
}

type ComplexGraphData struct {
  Gmode string            `json:"gmode"`
  Stack bool              `json:"stack"`
  Type string             `json:"type"`
  GraphId int             `json:"graph_id"`
}

type ComplexGraph struct {
  BaseGraph
  Data []ComplexGraphData `json:"data"`
  Sumup bool              `json:"sumup"`
}

type GraphList []GraphEssential
type ComplexList []GraphEssential

func (self *GraphEssential) GetPath() string {
  return fmt.Sprintf(
    "%s/%s/%s",
    self.ServiceName,
    self.SectionName,
    self.GraphName,
  )
}

func NewGraph() (*Graph) {
  g := &Graph {}
  g.Gmode = GMODE_DEFAULT
  g.Mode  = MODE_DEFAULT
  g.Type  = TYPE_DEFAULT
  return g
}

func NewComplexGraph() (*ComplexGraph) {
  g := &ComplexGraph {}
  g.Complex = true
  return g
}

