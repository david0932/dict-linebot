 // Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	//"log"
	//"net/http"
	"os"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
	//
	//"encoding/json"
	//"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var bot *linebot.Client

func main() {
	strID := os.Getenv("ChannelID")
	numID, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		log.Fatal("Wrong environment setting about ChannelID")
	}

	bot, err = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	received, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, result := range received.Results {
		content := result.Content()
		if content != nil && content.IsMessage && content.ContentType == linebot.ContentTypeText {
			text, err := content.TextContent()
			 dict_result := dict(text.Text)
			_, err = bot.SendText([]string{content.From}, dict_result)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
func dict(input string) string {

	root := "http://cdict.net/?q="
	para := url.QueryEscape(input)
    pos := 0
	
	dictEndPoint := root + para

	resp, err := http.Get(dictEndPoint)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
     
	//fmt.Println(string(body))
	pos = strings.Index(CToGoString(body), "content=")
	//fmt.Println(pos)
	//fmt.Println("---------------------------------\n")
	//fmt.Println(string(body[pos+8:250]))
	//fmt.Println(getContent(body,pos+8))
	//fmt.Println("---------------------------------\n")
	
	return getContent(body,pos+8)
}

func CToGoString(c []byte) string {
    n := -1
    for i, b := range c {
        if b == 0 {
            break
        }
        n = i
    }
    return string(c[:n+1])
}

func getContent(c []byte, p int) string {
   var end int
    
   
   for i:=p; i<(len(c)-p); i++ {
     if c[i] == '>' {
	    end = i
		return string(c[p:end])
	 } 
   }
   return "Not FOUND!"
}