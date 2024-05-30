package goravel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Define an envelope type.
type envelope map[string]interface{}

// readIDParam reads interpolated "id" from request URL and returns it and nil. If there is an error
// it returns and 0 and an error.
func (g *Goravel) ReadIDParam(r *http.Request) (int64, error) {
	idFromParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idFromParam, 10, 64)

	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// writeJSON marshals data structure to encoded JSON response. It returns an error if there are
// any issues, else error is nil.
func (g *Goravel) WriteJSON(w http.ResponseWriter, status int, data envelope,
	headers ...http.Header) error {
	// Use the json.MarshalIndent() function so that whitespace is added to the encoded JSON. Use
	// no line prefix and tab indents for each element.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// At this point, we know that we won't encounter any more errors before writing the response,
	// so it's safe to add any headers that we want to include. We loop through the header map
	// and add each header to the http.ResponseWriter header map. Note that it's OK if the
	// provided header map is nil. Go doesn't through an error if you try to range over (
	// or generally, read from) a nil map
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	// Add the "Content-Type: application/json" header, then write the status code and JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(js); err != nil {
		g.ErrorLog.Println(err)
		return err
	}

	return nil
}

// readJSON decodes request Body into corresponding Go type. It triages for any potential errors
// and returns corresponding appropriate errors.
func (g *Goravel) ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB to prevent
	// any potential nefarious DoS attacks.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize the json.Decoder, and call the DisallowUnknownFields() method on it
	// before decoding. So, if the JSON from the client includes any field which
	// cannot be mapped to the target destination, the decoder will return an error
	// instead of just ignoring the field.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the request body to the destination.
	err := dec.Decode(dst)
	if err != nil {
		// If there is an error during decoding, start the error triage...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// Use the error.As() function to check whether the error has the type *json.SyntaxError.
		// If it does, then return a plain-english error message which includes the location
		// of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON at (charcter %d)", syntaxError.Offset)

		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax error in the JSON. So, we check for this using errors.Is() and return
		// a generic error message. There is an open issue regarding this at
		// https://github.com/golang/go/issues/25956
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Likewise, catch any *json.UnmarshalTypeError errors.
		// These occur when the JSON value is the wrong type for the target destination.
		// If the error relates to a specific field, then we include that in our error message
		// to make it easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)",
				unmarshalTypeError.Offset)

		// An io.EOF error will be returned by Decode() if the request body is empty. We check
		// for this with errors.Is() and return a plain-english error message instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// If the JSON contains a field which cannot be mapped to the target destination
		// then Decode() will now return an error message in the format "json: unknown
		// field "<name>"". We check for this, extract the field name from the error,
		// and interpolate it into our custom error message.
		// Note, that there's an open issue at https://github.com/golang/go/issues/29035
		// regarding turning this into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// If the request body exceeds 1MB in size then decode will now fail with the
		// error "http: request body too large". There is an open issue about turning
		// this into a distinct error type at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// A json.InvalidUnmarshalError error will be returned if we pass a non-nil
		// pointer to Decode(). We catch this and panic, rather than returning an error
		// to our handler. At the end of this chapter we'll talk about panicking
		// versus returning, and discuss why it's an appropriate thing to do in this specific
		// situation.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// For anything else, return the error message as-is.
		default:
			return err
		}
	}

	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value then this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body, and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// url.Values:
// type Values map[string][]string

// readString is a helper method on application type that returns a string value from the URL query
// string, or the provided default value if no matching key is found.
func (g *Goravel) ReadStrings(qs url.Values, key string, defaultValue string) string {
	// Extract the value for a given key from the URL query string.
	// If no key exists this will return an empty string "".
	s := qs.Get(key)

	// If no key exists (or the value is empty) then return the default value
	if s == "" {
		return defaultValue
	}

	// Otherwise, return the string
	return s
}

// readCSV is a helper method on application type that reads a string value from the URL query
// string and then splits it into a slice on the comma character. If no matching key is found
// then it returns the provided default value.
func (g *Goravel) ReadCSV(qs url.Values, key string, defaultValue []string) []string {
	// Extract the value from the URL query string
	csv := qs.Get(key)

	// if no key exists (or the value is empty) then return the default value
	if csv == "" {
		return defaultValue
	}

	// Otherwise, parse the value into a []string slice and return it.
	return strings.Split(csv, ",")
}

// readInt is a helper method on application type that reads a string value from the URL query
// string and converts it to an integer before returning. If no matching key is found then it
// returns the provided default value. If the value couldn't be converted to an integer, then we
// record an error message in the provided Validator instance, and return the default value.
func (g *Goravel) ReadInt(qs url.Values, key string, defaultValue int) (int, error) {
	// Extract the value from the URL query string.
	s := qs.Get(key)

	// If no key exists (or the value is empty) then return the default value.
	if s == "" {
		return defaultValue, nil
	}

	// Try to convert the string value to an int. If this fails, add an error message to the
	// validator instance and return the default value.
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue, errors.New("invalid query parameter")
	}

	// Otherwise, return the converted integer value.
	return i, nil
}
