package growthforecast

import (
// "io/ioutil"
  "encoding/json"
  "errors"
  "fmt"
  "net/http"
)

type Client struct {
  Host string
  Port int
}

func NewClient(host string, port int) *Client {
  return &Client { host, port }
}

func (self *Client) getJSON(path string) (*http.Response, error) {
  url := fmt.Sprintf(
    "http://%s:%d%s",
    self.Host,
    self.Port,
    path,
  )
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

