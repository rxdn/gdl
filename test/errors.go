package main

var objectError = []byte(`
{
	"code": 50035,
	"errors": {
		"access_token": {
			"_errors": [{
				"code": "BASE_TYPE_REQUIRED",
				"message": "This field is required"
			}]
		}
	},
	"message": "Invalid Form Body"
}
`)

var arrayError = []byte(`
{
    "code": 50035,
    "errors": {
        "activities": {
            "0": {
                "platform": {
                    "_errors": [
                        {
                            "code": "BASE_TYPE_CHOICES",
                            "message": "Value must be one of ('desktop', 'android', 'ios')."
                        }
                    ]
                },
                "type": {
                    "_errors": [
                        {
                            "code": "BASE_TYPE_CHOICES",
                            "message": "Value must be one of (0, 1, 2, 3, 4, 5)."
                        }
                    ]
                }
            }
        }
    },
    "message": "Invalid Form Body"
}
`)
