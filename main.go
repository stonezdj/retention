package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	host := flag.String("host", "127.0.0.1", "host")
	username := flag.String("username", "admin", "username")
	password := flag.String("password", "Harbor12345", "password")
	retentionIDList := flag.String("list", "", "retention policy id list, for example: 1,2,3,4")
	flag.Parse()
	idList := make([]int, 0)

	ids := strings.Split(*retentionIDList, ",")
	for _, id := range ids {
		idInt, err := strconv.Atoi(strings.TrimSpace(id))
		if err != nil {
			continue
		}
		idList = append(idList, idInt)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	for _, retId := range idList {
		url := fmt.Sprintf("https://%v/api/v2.0/retentions/%v", *host, retId)
		method := "GET"

		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
			return
		}
		auth := *username + ":" + *password
		authorization := base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Add("Authorization", fmt.Sprintf("Basic %v", authorization))

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		if res.StatusCode != http.StatusOK {
			fmt.Printf("Failed to get the retention id %v\n", retId)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		var plc Metadata
		fmt.Printf("Get retention success, id= %v\n", retId)
		err = json.Unmarshal(body, &plc)
		if err != nil {
			fmt.Println(err)
			return
		}
		orgCron := plc.Trigger.Settings["cron"]
		if orgCron == nil {
			fmt.Println("current job is not a cron job schedule")
			return
		}
		if orgCronStr, ok := orgCron.(string); ok {
			err = updateRetention(err, url, authorization, client, retId, &plc, "")
			if err != nil {
				fmt.Println(err)
			}
			err = updateRetention(err, url, authorization, client, retId, &plc, orgCronStr)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Updated retention policy, id=%v\n", retId)
		}
		fmt.Println("===================================")
	}

}

func updateRetention(err error, url string, authorization string, client *http.Client, retentionID int, plc *Metadata, cronStr string) error {
	orgCron := plc.Trigger.Settings["cron"]
	if orgCron == nil {
		fmt.Println("current job is not a cron job schedule")
		return nil
	}
	plc.Trigger.Settings["cron"] = cronStr
	body, err := json.Marshal(plc)
	if err != nil {
		return err
	}
	putReq, err := http.NewRequest("PUT", url, strings.NewReader(string(body)))
	putReq.Header.Add("Content-Type", "application/json")
	putReq.Header.Add("Authorization", fmt.Sprintf("Basic %v", authorization))
	if err != nil {
		return err
	}
	putRes, err := client.Do(putReq)
	if err != nil {
		return err
	}
	if putRes.StatusCode == http.StatusOK {
		fmt.Println("Updated retention schedule success.")
	} else {
		fmt.Printf("Failed to update the tag retention schedule with id %v, error code %v\n", retentionID, putRes.StatusCode)
	}
	defer putRes.Body.Close()

	body, err = ioutil.ReadAll(putRes.Body)
	if err != nil {
		return err
	}
	return nil
}
