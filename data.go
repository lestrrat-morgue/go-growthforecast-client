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

const MODE_COUNT      = "count"
const MODE_GAUGE      = "gauge"
const MODE_MODIFIED   = "modified"
const MODE_DEFAULT    = MODE_GAUGE

const GMODE_GAUGE     = "gauge"
const GMODE_SUBTRACT  = "subtract"
const GMODE_DEFAULT   = GMODE_GAUGE

const TYPE_AREA       = "AREA"
const TYPE_LINE1      = "LINE1"
const TYPE_LINE2      = "LINE2"
const TYPE_DEFAULT    = TYPE_AREA

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

