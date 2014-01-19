package growthforecast

/*

GrowthForecast (http://kazeburo.github.io/GrowthForecast/) is a data 
visualization tool. This library gives you an easy way to interact
with the GrowthForecast server

*/

import (
  "bytes"
  "encoding/json"
  "errors"
  "fmt"
  "net/http"
  "net/url"
  "strings"
)

type Client struct {
  BaseURL string
}

/*

NewClient() creates a new growthforecast.Client struct. You just need
to pass it a base URL where all API endpoints will automatically be
generated from.

    client := growthforecast.NewClient("http://gf.mycompany.com")
    g, err := client.GetGraph("service/section/graph")
    if err != nil {
        log.Fatalf("Error while fetching graph: %s", err)
    }

*/
func NewClient(base string) *Client {
  return &Client { base }
}

func (self *Client) createURL(path string) string {
  if strings.HasPrefix(path, "/") {
    path = path[1:len(path)]
  }
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
func (self *Client) CreateGraph(graph *Graph) (*Graph, error) {
  values := url.Values{
    "number": {fmt.Sprintf("%d",graph.Number)},
  }
  if graph.Mode != "" {
    values.Add("mode", graph.Mode)
  }
  if graph.Color != "" {
    values.Add("color", graph.Color)
  }
  url := self.createURL(fmt.Sprintf("api/%s", graph.GetPath()))
  res, err := http.PostForm(url, values)
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

  var jres struct {
    Data Graph      `json:"data"`
    Error int       `json:"error"`
  }
  dec := json.NewDecoder(res.Body)
  err = dec.Decode(&jres)
  if err != nil {
    return nil, errors.New(
      fmt.Sprintf(
        "Failed to decode JSON: %s",
        err,
      ),
    )
  }

  if jres.Error != 0 {
    return nil, errors.New(
      fmt.Sprintf(
        "Error response: %s",
        jres.Error,
      ),
    )
  }

  return &jres.Data, nil
}

func (self *Client) CreateComplex(graph *ComplexGraph) (*ComplexGraph, error) {
  payload, err := json.Marshal(graph)
  if err != nil {
    return nil, errors.New(
      fmt.Sprintf(
        "Failed to encode json data: %s", err,
      ),
    )
  }

  url := self.createURL("json/create/complex")
  res, err := http.Post(
    url,
    "application/json",
    bytes.NewReader(payload),
  )
  if err != nil {
    return nil, errors.New(
      fmt.Sprintf(
        "Failed to post to %s: %s", url, err,
      ),
    )
  }

  if res.StatusCode != 200 {
    return nil, errors.New(
      fmt.Sprintf(
        "HTTP request to %s failed with %s",
        url,
        res.Status,
      ),
    )
  }

  var jres struct {
    Location string `json:"location"`
    Error int
  }
  dec := json.NewDecoder(res.Body)
  err = dec.Decode(&jres)
  if err != nil {
    return nil, errors.New(
      fmt.Sprintf(
        "Failed to decode JSON: %s",
        err,
      ),
    )
  }

  if jres.Error != 0 {
    return nil, errors.New(
      fmt.Sprintf(
        "Error response: %s",
        jres.Error,
      ),
    )
  }

  return self.GetComplexByPath(graph.GetPath())
}

func (self *Client) GetGraph(id int) (*Graph, error) {
  // It's actually exactly the same as GetGraphByPath, but we just
  // implement this conversion for easy of use
  return self.GetGraphByPath(fmt.Sprintf("%d", id))
}

func (self *Client) GetGraphByPath(path string) (*Graph, error) {
  dec, err := self.getDecoder(fmt.Sprintf("/json/graph/%s", path))
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

func (self *Client) GetComplex(id int) (*ComplexGraph, error){
  return self.GetComplexByPath(fmt.Sprintf("%d", id))
}

func (self *Client) GetComplexByPath(path string) (*ComplexGraph, error){
  dec, err := self.getDecoder(fmt.Sprintf("/json/complex/%s", path))
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

/*

Fetches the list of ComplexGraphs registered in the GrowthForecast instance.

*/
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

