package growthforecast

import (
  "fmt"
)

type GraphEssential struct {
  GraphName string        `json:"graph_name"`
  Id int                  `json:"id"`
  SectionName string      `json:"section_name"`
  ServiceName string      `json:"service_name"`
}

type BaseGraph struct {
  GraphEssential
  CreatedAt string     `json:"created_at"`
  Description string      `json:"description"`
  Number int              `json:"number"`
  Sort int                `json:"sort"`
  UpdatedAt string     `json:"updated_at"`
}

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
  Type string             `json:"type"`
  Ulimit int              `json:"ulimit"`
  Unit string             `json:"unit"`
}

type ComplexGraphData struct {
  Gmode string            `json:"gmode"`
  Stack bool              `json:"stack"`
  Type string             `json:"type"`
  GraphId string          `json:"graph_id"`
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

