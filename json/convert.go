package json

import (
	"encoding/json"

	"github.com/csunibo/synta"
)

func Convert(synta synta.Synta) Synta {
	// TODO: convert
	return Synta{}
}

func ToJson(synta synta.Synta) (buf []byte, err error) {
	buf, err = json.Marshal(Convert(synta))
	return
}
