package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	id         string
	resolution string
	page       string
)

type VCode struct {
	Playurl string `xml:"playurl"`
}

func init() {
	flag.StringVar(&id, "id", "", "video id")
	flag.StringVar(&resolution, "res", "720p", "video resolution")
	flag.Parse()
	if id == "" {
		panic(errors.New("video should not be empty!"))
	}
	page = fmt.Sprintf(`http://www.letv.com/ptv/vplay/%s.html`, id)
}

func main() {
	var err error
	if resp, err := http.Get(page); err == nil {
		defer resp.Body.Close()
		if data, err := ioutil.ReadAll(resp.Body); err == nil {
			html := string(data)
			idx := strings.Index(html, "v_code=")
			vCode := html[idx+7:]
			vCode = vCode[:strings.Index(vCode, "'")]
			//url decode
			if vCode, err = url.QueryUnescape(vCode); err == nil {
				//decode base64
				if data, err = base64.StdEncoding.DecodeString(vCode); err == nil {
					var vCodeObj VCode
					if err = xml.Unmarshal(data, &vCodeObj); err == nil {
						var playUrlMap map[string]interface{}
						if err = json.Unmarshal([]byte(vCodeObj.Playurl), &playUrlMap); err == nil {
							val := getVal(playUrlMap, "dispatch", resolution)
							fmt.Printf("%s\n", realUrl(val[0].(string)))
						}
					}
				}
			}

		}
	}
	if err != nil {
		panic(err)
	}
}

func getVal(playUrlMap map[string]interface{}, dispatch, resolution string) []interface{} {
	val := playUrlMap[dispatch].(map[string]interface{})[resolution]
	return val.([]interface{})
}

func realUrl(fakeUrl string) string {
	idx := strings.Index(fakeUrl, "?")
	idx1 := strings.LastIndex(fakeUrl[:idx], "/")
	idx2 := strings.Index(fakeUrl[7:], "/")
	fakeCode := fakeUrl[idx1+1 : idx]
	realCode, _ := base64.StdEncoding.DecodeString(fakeCode)
	return fakeUrl[:idx2+8] + string(realCode) + fakeUrl[idx:]
}
