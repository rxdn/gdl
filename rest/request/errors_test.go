package request

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestObjectError(t *testing.T) {
	raw := `
{
  "code": 50035,
  "errors": {
    "access_token": {
      "_errors": [
        {
          "code": "BASE_TYPE_REQUIRED",
          "message": "This field is required"
        }
      ]
    }
  },
  "message": "Invalid Form Body"
}`

	var e ApiV8Error
	require.NoError(t, e.UnmarshalJSON([]byte(raw)))

	require.Equal(t, 50035, e.Code)
	require.Equal(t, "Invalid Form Body", e.Message)

	require.Len(t, e.Errors, 1)
	require.Equal(t, "access_token", e.Errors[0].FieldName)
	require.Equal(t, "BASE_TYPE_REQUIRED", e.Errors[0].Code)
	require.Equal(t, "This field is required", e.Errors[0].Message)
}

func TestArrayError(t *testing.T) {
	raw := `
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
}`

	var e ApiV8Error
	require.NoError(t, e.UnmarshalJSON([]byte(raw)))

	require.Equal(t, 50035, e.Code)
	require.Equal(t, "Invalid Form Body", e.Message)

	require.Len(t, e.Errors, 2)

	fst := e.Errors[0]
	require.Equal(t, "activities[0].platform", fst.FieldName)
	require.Equal(t, "BASE_TYPE_CHOICES", fst.Code)
	require.Equal(t, "Value must be one of ('desktop', 'android', 'ios').", fst.Message)

	snd := e.Errors[1]
	require.Equal(t, "activities[0].type", snd.FieldName)
	require.Equal(t, "BASE_TYPE_CHOICES", snd.Code)
	require.Equal(t, "Value must be one of (0, 1, 2, 3, 4, 5).", snd.Message)
}

func TestNestedArray(t *testing.T) {
	raw := `
{
  "message": "Invalid Form Body",
  "code": 50035,
  "errors": {
    "components": {
      "0": {
        "components": {
          "1": {
            "emoji": {
              "id": {
                "_errors": [
                  {
                    "code": "BUTTON_COMPONENT_INVALID_EMOJI",
                    "message": "Invalid emoji"
                  }
                ]
              }
            }
          }
        }
      }
    }
  }
}`

	var e ApiV8Error
	require.NoError(t, e.UnmarshalJSON([]byte(raw)))

	require.Equal(t, 50035, e.Code)
	require.Equal(t, "Invalid Form Body", e.Message)

	require.Len(t, e.Errors, 1)
	require.Equal(t, "components[0].components[1].emoji.id", e.Errors[0].FieldName)
}
