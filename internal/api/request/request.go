package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

func ReadJSON(r io.Reader, dst interface{}) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var typeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed json (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly-formed json")
		case errors.As(err, &typeError):
			if typeError.Field != "" {
				return fmt.Errorf("body contains bad type for json field %s (at character %d)", typeError.Field, typeError.Offset)
			}
			return fmt.Errorf("body contains bad type for json field (at character %d)", typeError.Offset)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")
		default:
			return err
		}
	}

	return nil
}
