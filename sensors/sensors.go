package sensors

import (
	"encoding/json"
	"os/exec"
	"strings"

	. "github.com/xpetit/x/v2"
)

// MaxTemperature returns the maximum temperature reported by lm-sensors
func MaxTemperature() (max float64) {
	var data map[string]map[string]json.RawMessage
	C(json.Unmarshal(C2(exec.Command("sensors", "-j").CombinedOutput()), &data))
	for _, v := range data {
		for k, v := range v {
			if k == "Adapter" {
				continue
			}
			var data map[string]float64
			C(json.Unmarshal(v, &data))
			for k, v := range data {
				if strings.HasSuffix(k, "_input") && v > max {
					max = v
				}
			}
		}
	}
	return
}
