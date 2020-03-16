package fixtures

const (
	TestFilename = "ports.json"
)

var (
	TestJSONDataSmall = struct {
		Data string
		Len  uint64
	}{
		Data: `
		{
			"AEAJM": {
			"name": "Ajman",
			"city": "Ajman",
			"country": "United Arab Emirates",
			"alias": [],
			"regions": [],
			"coordinates": [
				55.5136433,
				25.4052165
			],
			"province": "Ajman",
			"timezone": "Asia/Dubai",
			"unlocs": [
				"AEAJM"
			],
			"code": "52000"
			},
			"AEAUH": {
			"name": "Abu Dhabi",
			"coordinates": [
				54.37,
				24.47
			],
			"city": "Abu Dhabi",
			"province": "Abu Z¸aby [Abu Dhabi]",
			"country": "United Arab Emirates",
			"alias": [],
			"regions": [],
			"timezone": "Asia/Dubai",
			"unlocs": [
				"AEAUH"
			],
			"code": "52001"
			}
		}
		`,
		Len: 2,
	}
)
