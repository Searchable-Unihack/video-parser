package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	for i, f := range os.Args {
		started := time.Now()

		// Only parse files with .png
		if !strings.HasSuffix(f, ".png") {
			log.Printf("%s does not have .png extension.", f)
			continue
		}

		var output = fmt.Sprintf("%s.json", strings.TrimRight(f, ".png"))

		// Check if the output json already exists:
		if _, err := os.Stat(output); !os.IsNotExist(err) {
			log.Printf("File %s already exists", output)
			continue
		}

		text, err := fetchText(f)

		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(output, text, 0644)

		if err != nil {
			panic(err)
		}

		fmt.Printf("File %d: %s.  Output: %s\n", i, f, output)

		// Only 20 requests per minute, so throttle our requests:
		for time.Since(started) < (time.Millisecond * 3200) {
			time.Sleep(time.Millisecond * 200)
		}
	}
}

func fetchText(filename string) ([]byte, error) {
	var text []byte
	client := &http.Client{}

	subscriptionKey := os.Getenv("SUBSCRIPTION_KEY")

	reader, err := os.Open(filename)

	if err != nil {
		return text, err
	}

	url := "https://southeastasia.api.cognitive.microsoft.com/vision/v1.0/ocr?language=en&detectOrientation=true"
	//url :=  "https://southeastasia.api.cognitive.microsoft.com/vision/v1.0/analyze?visualFeatures=Categories&language=en"

	req, err := http.NewRequest("POST", url, reader)

	if err != nil {
		return text, err
	}

	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Ocp-Apim-Subscription-Key", subscriptionKey)

	resp, err := client.Do(req)

	if err != nil {
		return text, err
	}

	defer resp.Body.Close()

	text, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return text, err
	}

	return text, nil
}
