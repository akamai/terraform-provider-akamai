package akamai

import (
	"crypto/sha1"
	"encoding/hex"
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
