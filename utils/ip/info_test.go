package ip

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIP(t *testing.T) {
	ip := "77.91.123.163"
	data := IPInfo{
		IP:       ip,
		Hostname: "onorridg.tech",
		City:     "Meppel",
		Country:  "NL",
		Loc:      "52.6958,6.1944",
		Org:      "AS44477 STARK INDUSTRIES SOLUTIONS LTD",
		Timezone: "Europe/Amsterdam",
	}
	result, err := Info(ip)
	assert.Equal(t, nil, err)
	assert.Equal(t, data, *result)
}
