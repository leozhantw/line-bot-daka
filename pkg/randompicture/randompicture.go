package randompicture

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const endpoint = "https://www.google.com/search?safe=off&tbm=isch&q=%s"

func Random(queries []string) (string, error) {
	rand.Seed(time.Now().UnixNano())

	q := queries[rand.Intn(len(queries))]
	apiURL := fmt.Sprintf(endpoint, url.QueryEscape(q))

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile("src=\"(http[^\"]+)\"")
	matches := re.FindAllStringSubmatch(string(body), -1)

	potatoes := make([]string, len(matches))

	for index, match := range matches {
		potatoes[index] = match[1]
	}

	return potatoes[rand.Intn(len(potatoes))], nil
}
