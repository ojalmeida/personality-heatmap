package phase1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"personality-heatmap/phase1/data"
	"personality-heatmap/phase1/geo"
	"personality-heatmap/phase1/models"
	"personality-heatmap/phase1/proxy"
	"strconv"
	"strings"
	"time"
)

func cls() {

	fmt.Print("\033[2J")
}

func askForChoose(prompt string, onPositive, onNegative func()) {

	var choose string
	for {

		fmt.Printf("%s (y/N) ", prompt)
		_, _ = fmt.Scanln(&choose)

		switch strings.ToLower(choose) {

		case "y", "yes":

			onPositive()

		case "n", "no", "", "\n":

			onNegative()

		default:

			fmt.Println("\nInvalid answer")
			continue

		}
	}

}

func main() {

	fmt.Println("-- Welcome to phase1 --")
	fmt.Println()
	fmt.Println("# Phase needed resources: 3 tinder accounts and a city name")
	fmt.Println("# Phase side-effects: 3 tinder accounts with fake location implemented")
	fmt.Println("# Phase output objects: 3 X-Auth-Tokens")
	fmt.Println()
	askForChoose("Do you have the needed resources and want to start phase 1 ?", start, func() { os.Exit(1) })

}

func start() {

	cls()
	fmt.Println("Right, we'll go trough the following steps:")
	fmt.Println()
	fmt.Println("1.\tSet the heatmap target city")
	fmt.Println()
	fmt.Println("2.\tLogin in the 1st account and accept location tracking")
	fmt.Println()
	fmt.Println("3.\tLogin in the 2nd account and accept location tracking")
	fmt.Println()
	fmt.Println("4.\tLogin in the 3rd account and accept location tracking")
	fmt.Println()
	fmt.Println("5.\tShow results")
	fmt.Println()
	askForChoose("Can we start step 1 ?", setCity, func() { start() })

}

func setCity() {

	type apiResponse []struct {
		PlaceID     int      `json:"place_id,omitempty"`
		Licence     string   `json:"licence,omitempty"`
		OsmType     string   `json:"osm_type,omitempty"`
		OsmID       int      `json:"osm_id,omitempty"`
		Boundingbox []string `json:"boundingbox,omitempty"`
		Lat         string   `json:"lat,omitempty"`
		Lon         string   `json:"lon,omitempty"`
		DisplayName string   `json:"display_name,omitempty"`
		Class       string   `json:"class,omitempty"`
		Type        string   `json:"type,omitempty"`
		Importance  float64  `json:"importance,omitempty"`
		Icon        string   `json:"icon,omitempty"`
	}

	var cityName string
	var cityDetails apiResponse

	cls()
	fmt.Print("City name: ")
	inputReader := bufio.NewReader(os.Stdin)
	cityName, _ = inputReader.ReadString('\n')

	cityName = strings.Trim(cityName, "\n")

	fmt.Print("Processing...")

	res, err := http.Get(fmt.Sprintf("https://nominatim.openstreetmap.org/search?city=%s&country=br&format=json", html.EscapeString(cityName)))

	if err != nil {
		fmt.Println("\nOps, this was not supposed to happen:", err.Error())
	}

	rawBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Ops, this was not supposed to happen:", err.Error())
	}

	err = json.Unmarshal(rawBody, &cityDetails)

	if err != nil {
		fmt.Println("Ops, this was not supposed to happen:", err.Error())
	}

	if err == nil {

		details := cityDetails[0]

		data.Data.City.Name = strings.Split(details.DisplayName, ",")[0]

		fmt.Print("\u001b[2K\r\n")
		fmt.Println(details.DisplayName)
		fmt.Printf("Lat: %s\n", details.Lat)
		fmt.Printf("Lng: %s\n", details.Lon)
		fmt.Println()
		askForChoose("Is it correct ?",

			func() {

				data.Data.City.Coordinates.Lat, _ = strconv.ParseFloat(details.Lat, 64)
				data.Data.City.Coordinates.Lng, _ = strconv.ParseFloat(details.Lon, 64)
				doTrackAccounts()
			},
			func() {

				fmt.Println("Okay, this is embarrassing, let's try again, please type the city name in a different manner")
				time.Sleep(time.Second * 2)
				setCity()

			},
		)

	}

}

func trackAccountInstructions() {

	cls()
	fmt.Println("Well, now the most important part")
	fmt.Println("I'll open 3 firefox windows, one followed by another")
	fmt.Println("In each of them you have to: ")
	fmt.Println("1. Login in your tinder account")
	fmt.Println("2. Accept the location tracking, asked on top of the screen")
	fmt.Println("3. Close window")
	fmt.Println()
	askForChoose("Can we start ?", trackAccountInstructions, func() { os.Exit(1) })
}

func doTrackAccounts() {

	nodesCoordinates := geo.GetNodesCoordinates(data.Data.City.Coordinates)

	cls()
	fmt.Println("Right, do the procedures for the 1st account")

	data.Data.Nodes[0] = &models.Node{
		Name: data.Data.City.Name + "-1",
		Location: models.FakeLocation{

			Coordinate: nodesCoordinates[0],
			Accuracy:   10,
		},
		APIToken: "",
	}

	node1Config := proxy.Config{
		NodeName: data.Data.Nodes[0].Name,
		FakeGPS:  data.Data.Nodes[0].Location,
	}

	track(node1Config)

	cls()
	fmt.Println("Now, for the 2nd account")

	data.Data.Nodes[1] = &models.Node{
		Name: data.Data.City.Name + "-2",
		Location: models.FakeLocation{

			Coordinate: nodesCoordinates[1],
			Accuracy:   10,
		},
		APIToken: "",
	}

	node2Config := proxy.Config{
		NodeName: data.Data.Nodes[1].Name,
		FakeGPS:  data.Data.Nodes[1].Location,
	}

	track(node2Config)

	cls()
	fmt.Println("And finally, for the 3rd account")

	data.Data.Nodes[2] = &models.Node{
		Name: data.Data.City.Name + "-3",
		Location: models.FakeLocation{

			Coordinate: nodesCoordinates[2],
			Accuracy:   10,
		},
		APIToken: "",
	}

	node3Config := proxy.Config{
		NodeName: data.Data.Nodes[2].Name,
		FakeGPS:  data.Data.Nodes[2].Location,
	}

	track(node3Config)

	cls()
	fmt.Println("Well done!")

	saveData()

	fmt.Println("Seem like we are done here, see ya!")

	os.Exit(0)

}

func track(proxyConfig proxy.Config) {

	proxy.ProxyConfig = proxyConfig

	proxy.Start()

	cmd := exec.Command("firefox", "--private", "https://tinder.com")
	cmd.Env = append(cmd.Env, "https_proxy=http://localhost:8888", "DISPLAY=:0")

	err := cmd.Start()

	if err != nil {

		fmt.Println(err.Error())

	}

	fmt.Println()
	fmt.Println("Close window when done")
	cmd.Wait()

	proxy.Stop()

}

func saveData() {

	var path string

	fmt.Print("Type the path that you want that I save this stuff: ")
	_, _ = fmt.Scanln(&path)

	path = strings.TrimRight(path, "/") + "/out.yaml"

	fmt.Printf("Saving output to %s", path)
	err := data.Save(path)

	if err != nil {
		fmt.Printf("--> Error --> %s\n", err.Error())
		fmt.Println()
		fmt.Println("Lets do it again, try type a different path")
		time.Sleep(time.Second * 2)
		saveData()

	} else {

		fmt.Print("--> Success!\n")

	}
}
