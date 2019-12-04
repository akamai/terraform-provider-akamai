package akamai

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"os"
	"reflect"

	"github.com/hashicorp/terraform/helper/schema"
)

func getSingleSchemaSetItem(d interface{}) map[string]interface{} {
	ss := d.(*schema.Set)
	list := ss.List()

	if len(list) == 0 || list[0] == nil {
		return nil
	}

	return list[0].(map[string]interface{})
}

func getSetList(d interface{}) ([]interface{}, bool) {
	if ss, ok := d.(*schema.Set); ok {
		return ss.List(), ok
	}

	return nil, false
}

func unmarshalSetString(d interface{}) ([]string, bool) {
	schemaSet, ok := d.(*schema.Set)

	if !ok {
		return nil, false
	}

	schemaList := schemaSet.List()
	stringSet := make([]string, len(schemaList))

	for i, v := range schemaList {
		stringSet[i] = v.(string)
	}

	return stringSet, ok
}

func readNullableString(d interface{}) *string {
	str, ok := d.(string)

	if !ok || len(str) == 0 {
		return nil
	}

	return &str
}

func getSHAString(rdata string) string {
	h := sha1.New()
	h.Write([]byte(rdata))

	sha1hashtest := hex.EncodeToString(h.Sum(nil))
	return sha1hashtest
}

func logfile(filename string, text string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func jsonBytesEqual(b1, b2 []byte) bool {
	var o1 interface{}
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 interface{}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

/*
func suppressEquivalentJsonDiffs(k, old, new string, d *schema.ResourceData) bool {
	ob := bytes.NewBufferString("")
	if err := json.Compact(ob, []byte(old)); err != nil {
		return false
	}

	nb := bytes.NewBufferString("")
	if err := json.Compact(nb, []byte(new)); err != nil {
		return false
	}

	return jsonBytesEqual(ob.Bytes(), nb.Bytes())
}

func jsonBytesEqual(b1, b2 []byte) bool {
	var o1 interface{}
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 interface{}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

// ValidateJsonString is a SchemaValidateFunc which tests to make sure the
// supplied string is valid JSON.
func ValidateJsonString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := NormalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	return
}
func NormalizeJsonString(jsonString interface{}) (string, error) {
	var j interface{}

	if jsonString == nil || jsonString.(string) == "" {
		return "", nil
	}

	s := jsonString.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	bytes, _ := json.Marshal(j)
	return string(bytes[:]), nil
}
*/
