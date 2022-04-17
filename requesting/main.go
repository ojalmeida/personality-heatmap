package requesting

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gosuri/uilive"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"personality-heatmap/data"
	"personality-heatmap/data/database"
	"personality-heatmap/models"
	"strings"
	"time"
)

var phase1Path string

func init() {

	file, err := os.OpenFile("/tmp/requester.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 660)

	if err != nil {

		panic(err.Error())

	}

	log.SetOutput(file)

}

func Start() {

	fmt.Println("-- Welcome to phase 2 --")
	fmt.Println()
	fmt.Println("# Phase needed resources: 3 tinder APITokens with properly location configured")
	fmt.Println("# Phase side-effects: ---")
	fmt.Println("# Phase output objects: 3 profile databases")
	fmt.Println()
	askForChoose("Do you have the needed resources and want to start phase 2 ?", start, func() { os.Exit(1) })

}

func cls() {

	fmt.Print("\033[H\033[2J")
}

func askForChoose(prompt string, onPositive, onNegative func()) {

	var choose string
	for {

		fmt.Printf("%s (Y/n) ", prompt)
		_, _ = fmt.Scanln(&choose)

		switch strings.ToLower(choose) {

		case "y", "yes", "", "\n":

			onPositive()

		case "n", "no":

			onNegative()

		default:

			fmt.Println("\nInvalid answer")
			continue

		}
	}

}

func start() {

	cls()
	fmt.Print("Type the .yaml file path exported in phase 1: ")

	_, _ = fmt.Scanln(&phase1Path)

	err := data.LoadData(phase1Path)

	if err == nil {
		doPopulateDB()
	}

}

func doPopulateDB() {

	pathSlice := strings.Split(phase1Path, "/")
	dir := "/" + strings.Join(pathSlice[1:len(pathSlice)-1], "/")

	database.CreateDatabases(dir)
	populateDatabases()

}

func populateDatabases() {

	cls()
	var numberOfProfiles int
	fmt.Print("How many profiles do you want to populate (e.g 2000) ? ")
	_, _ = fmt.Scanln(&numberOfProfiles)

	doneChans := [3]chan bool{

		make(chan bool),
		make(chan bool),
		make(chan bool),
	}
	numberOfProfilesChans := [3]chan int{

		make(chan int),
		make(chan int),
		make(chan int),
	}
	timeoutChans := [3]chan bool{

		make(chan bool),
		make(chan bool),
		make(chan bool),
	}

	for n, node := range data.Data.Nodes {

		go saveProfiles(node, numberOfProfiles, database.Databases[n], doneChans[n], numberOfProfilesChans[n], timeoutChans[n])

	}

	allDone := 0.0
	insertedProfiles := [3]int{}
	done := [3]bool{}
	timeout := [3]bool{}

	cls()

	writer1 := uilive.New()
	writer2 := uilive.New()
	writer3 := uilive.New()

	writer1.Start()
	writer2.Start()
	writer3.Start()

	for {

		time.Sleep(time.Millisecond * 100)

		if done[0] {

			_, _ = fmt.Fprintf(writer1, "# %s\n\t%d/%d inserted profiles -> Done!\n", data.Data.Nodes[0].Name, insertedProfiles[0], numberOfProfiles)

		} else if timeout[0] {

			_, _ = fmt.Fprintf(writer1, "# %s\n\t%d/%d inserted profiles -> Timeout received!\n", data.Data.Nodes[0].Name, insertedProfiles[0], numberOfProfiles)

		} else {

			_, _ = fmt.Fprintf(writer1, "# %s\n\t%d/%d inserted profiles\n", data.Data.Nodes[0].Name, insertedProfiles[0], numberOfProfiles)

		}

		if done[1] {

			_, _ = fmt.Fprintf(writer1, "# %s\n\t%d/%d inserted profiles -> Done!\n", data.Data.Nodes[1].Name, insertedProfiles[1], numberOfProfiles)

		} else if timeout[1] {

			_, _ = fmt.Fprintf(writer1, "# %s\n\t%d/%d inserted profiles -> Timeout received!\n", data.Data.Nodes[1].Name, insertedProfiles[1], numberOfProfiles)

		} else {

			_, _ = fmt.Fprintf(writer1, "# %s\n\t%d/%d inserted profiles\n", data.Data.Nodes[1].Name, insertedProfiles[1], numberOfProfiles)

		}

		if done[2] {

			_, _ = fmt.Fprintf(writer1, "# %s\n\t%d/%d inserted profiles-> Done!\n", data.Data.Nodes[2].Name, insertedProfiles[2], numberOfProfiles)

		} else if timeout[2] {

			_, _ = fmt.Fprintf(writer1, "# %s\n\t%d/%d inserted profiles -> Timeout received!\n", data.Data.Nodes[2].Name, insertedProfiles[2], numberOfProfiles)

		} else {

			_, _ = fmt.Fprintf(writer1, "# %s\n\t%d/%d inserted profiles\n", data.Data.Nodes[2].Name, insertedProfiles[2], numberOfProfiles)

		}

		if allDone == 1.0 {

			time.Sleep(time.Second * 2)
			break

		}

		select {

		case <-timeoutChans[0]:

			timeout[0] = true

		case <-timeoutChans[1]:

			timeout[1] = true

		case <-timeoutChans[2]:

			timeout[2] = true

		case <-doneChans[0]:

			allDone += 1 / 3
			done[0] = true

		case <-doneChans[1]:

			allDone += 1 / 3
			done[1] = true

		case <-doneChans[2]:

			allDone += 1 / 3
			done[2] = true

		case i := <-numberOfProfilesChans[0]:

			insertedProfiles[0] = i

		case i := <-numberOfProfilesChans[1]:

			insertedProfiles[1] = i

		case i := <-numberOfProfilesChans[2]:

			insertedProfiles[2] = i

		default:
			continue
		}

	}

	writer1.Stop()
	writer2.Stop()
	writer3.Stop()

}

func saveProfiles(node models.Node, limit int, db *sql.DB, done chan bool, numberOfInserts chan int, timeout chan bool) {

	insertedProfiles := 0

outer:
	for {

		log.Println("Getting profiles from", node.Name)

		profiles := getProfiles(node.APIToken, timeout)

		select {

		case <-timeout:

			break outer

		default:
			break

		}

		for _, profile := range profiles {

			if insertedProfiles >= limit {
				break outer
			}

			err := database.Insert(profile, db)

			if err != nil {

				log.Printf("error when inserting %s in database %s : %s\n", profile.User.Name, node.Name, err.Error())
				continue
			}

			log.Printf("%s inserted in database %s", profile.User.Name, node.Name)

			insertedProfiles += 1

			numberOfInserts <- insertedProfiles

		}

		time.Sleep(time.Second)

	}

	done <- true

}

func getProfiles(apiToken string, timeout chan bool) []models.Profile {

	profiles := []models.Profile{}

	client := http.Client{Timeout: time.Second * 10}
	recURL, _ := url.Parse("https://api.gotinder.com/v2/recs/core?locale=pt")

	req := &http.Request{

		URL:    recURL,
		Header: map[string][]string{"X-Auth-Token": {apiToken}},
	}

	res, err := client.Do(req)

	if err != nil || res.StatusCode != 200 {

		return []models.Profile{}

	}

	rawBody, _ := ioutil.ReadAll(res.Body)

	if string(rawBody) == `{"meta":{"status":200},"data":{"timeout":1800000}}` {

		timeout <- true

	}

	var apiResponse models.APIResponse

	err = json.Unmarshal(rawBody, &apiResponse)

	if err != nil {
		return []models.Profile{}
	}

	for _, result := range apiResponse.Data.Results {

		profile := models.Profile{
			User:     result.User,
			Distance: result.DistanceMi,
		}

		if len(result.ExperimentInfo.UserInterests.SelectedInterests) != 0 {

			interests := []string{}

			for _, interest := range result.ExperimentInfo.UserInterests.SelectedInterests {

				interests = append(interests, interest.Name)

			}

			profile.Interests = interests

		}

		passProfile(profile.User.ID, apiToken)

		profiles = append(profiles, profile)

	}

	return profiles
}

func passProfile(id, apiToken string) {

	passURL, _ := url.Parse(fmt.Sprintf("https://api.gotinder.com/pass/%s", id))

	client := http.Client{Timeout: time.Second * 10}

	req := &http.Request{

		URL:    passURL,
		Header: map[string][]string{"X-Auth-Token": {apiToken}},
	}

	_, _ = client.Do(req)

}
