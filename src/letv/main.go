package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	id         = flag.String("id", "", "video id, required")
	resolution = flag.String("res", "720p", "video resolution, optional")
	all        = flag.Bool("a", false, "show all download link, optional")
	page       string
)

type VCode struct {
	Playurl string `xml:"playurl"`
}

func init() {
	flag.Parse()
	if *id == "" {
		usage()
	}
	page = fmt.Sprintf(`http://www.letv.com/ptv/vplay/%s.html`, *id)
}

func usage() {
	fmt.Println("usage:")
	flag.PrintDefaults()
	os.Exit(1)
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
							if *all {
								val := getVal(playUrlMap, "dispatch", *resolution)
								fmt.Printf("dispatch\n%s\n", realUrl(val[0].(string)))
								val = getVal(playUrlMap, "dispatchbak", *resolution)
								fmt.Printf("dispatchbak\n%s\n", realUrl(val[0].(string)))
								val = getVal(playUrlMap, "dispatchbak1", *resolution)
								fmt.Printf("dispatchbak1\n%s\n", realUrl(val[0].(string)))
								val = getVal(playUrlMap, "dispatchbak2", *resolution)
								fmt.Printf("dispatchbak2\n%s\n", realUrl(val[0].(string)))
							} else {
								val := getVal(playUrlMap, "dispatch", *resolution)
								fmt.Printf("%s\n", realUrl(val[0].(string)))
							}
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
