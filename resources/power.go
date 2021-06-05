package resources

import (
	"io/ioutil"
	"log"
	"math"
	"strconv"
)

const dir = "/sys/class/power_supply/BAT0/"

var (
	PowerOn bool
)

func GetMicroWatts() int {
	vString, err := ioutil.ReadFile(dir + "voltage_now")
	if err != nil {
		log.Printf("Couldn't read voltage: %s\n", err.Error())
		return 0
	}
	cString, err := ioutil.ReadFile(dir + "current_now")
	if err != nil {
		log.Printf("Couldn't read current: %s\n", err.Error())
		return 0
	}

	// we can ignore errors here because the reads will always be numbers
	v, _ := strconv.Atoi(string(vString))
	c, _ := strconv.Atoi(string(cString))
	v /= 1000
	c /= 1000
	return v * c
}

func GetWatts() float64 {
	return float64(GetMicroWatts()) / math.Pow(10, 6)

}
