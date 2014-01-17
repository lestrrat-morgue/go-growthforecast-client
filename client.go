package growthforecast

import (
// "io/ioutil"
  "encoding/json"
  "errors"
  "fmt"
  "net/http"
  "net/url"
)

type Client struct {
  BaseURL string
}

//    client := growthforecast.NewClient("http://gf.mycompany.com")
//    g, err := client.GetGraph("service/section/graph")
//    if err != nil {
//      log.Fatalf("Error while fetching graph: %s", err)
//    }
func NewClient(base string) *Client {
  return &Client { base }
}

func (self *Client) createURL(path string) string {
  return fmt.Sprintf("%s/%s", self.BaseURL, path)
}

func (self *Client) getJSON(path string) (*http.Response, error) {
  url := self.createURL(path)
  res, err := http.Get(url)
  if err != nil {
    return nil, err
  }

  if res.StatusCode != 200 {
    return nil, errors.New(
      fmt.Sprintf(
        "HTTP request to %s failed",
        url,
      ),
    )
  }

  return res, nil
}

func (self *Client) getDecoder(path string) (*json.Decoder, error) {
  res, err := self.getJSON(path)
  if err != nil {
    return nil, err
  }

  return json.NewDecoder(res.Body), nil
}

// Note that when you create a graph, you can only specify the basic 
// parameter
func (self *Client) CreateGraph(graph *Graph) error {
  values := url.Values{
    "number": {fmt.Sprintf("%d",graph.Number)},
  }
  if graph.Mode != "" {
    values.Add("mode", graph.Mode)
  }
  if graph.Color != "" {
    values.Add("color", graph.Color)
  }
  url := self.createURL(graph.GetPath())
  res, err := http.PostForm(url, values)
  if err != nil {
    return err
  }

  if res.StatusCode != 200 {
    return errors.New(
      fmt.Sprintf(
        "HTTP request to %s failed",
        url,
      ),
    )
  }

  return nil
}

func (self *Client) GetGraph(id string) (*Graph, error) {
  dec, err := self.getDecoder(fmt.Sprintf("/json/graph/%s", id))
  if err != nil {
    return nil, err
  }

  var e Graph
  err = dec.Decode(&e)
  if err != nil {
    return nil, err
  }

  return &e, nil
}

func (self *Client) GetComplex(id string) (*ComplexGraph, error){
  dec, err := self.getDecoder(fmt.Sprintf("/json/complex/%s", id))
  if err != nil {
    return nil, err
  }

  var e ComplexGraph
  err = dec.Decode(&e)
  if err != nil {
    return nil, err
  }

  return &e, nil
}

func (self *Client) GetGraphList() (GraphList, error) {
  dec, err := self.getDecoder("/json/list/graph")
  if err != nil {
    return nil, err
  }

  var e GraphList
  err = dec.Decode(&e)
  if err != nil {
    return nil, err
  }

  return e, nil
}

func (self *Client) GetComplexList() (ComplexList, error) {
  dec, err := self.getDecoder("/json/list/complex")
  if err != nil {
    return nil, err
  }

  var e ComplexList
  err = dec.Decode(&e)
  if err != nil {
    return nil, err
  }

  return e, nil
}

