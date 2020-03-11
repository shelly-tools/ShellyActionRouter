/**
 * @package   ShellyActionRouter
 * @copyright Thorsten Eurich
 * @license   GNU Affero General Public License (https://www.gnu.org/licenses/agpl-3.0.de.html)
 * @version  0.1
 *
 * @todo lots of documentation
 *
 * Simple Action Proxy written in Golang which makes it possible to execute multiple actions
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/gorilla/mux"
	"gopkg.in/ini.v1"
)

var netClient = &http.Client{
	Timeout: time.Second * 5,
}

type Unreachable struct {
	Status string `json:"status"`
}

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	port := cfg.Section("server").Key("port").String()

	actions, err := ini.Load("actions.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/action/{title}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]

		keys := actions.Section(title).Keys()

		for _, url := range keys {
			curl := fmt.Sprintf("%s", url)

			res, err := netClient.Get(curl)
			if err != nil {
				emp := Unreachable{Status: "offline"}
				var jsonData []byte
				jsonData, _ = json.Marshal(emp)

				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonData)

				log.Println(curl + " unreachable")
			} else {
				defer res.Body.Close()
				body, _ := ioutil.ReadAll(res.Body)
				fmt.Fprintf(w, string(body))
			}
		}
	})

	http.ListenAndServe(":"+port, r)
}
