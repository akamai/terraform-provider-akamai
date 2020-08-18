package gtm

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patrickmn/go-cache"
)

// KeepLogCorrelationID Used for logging footer
var KeepLogCorrelationID *string

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

func makeCache() *cache.Cache {
	profilecache := cache.New(5*time.Minute, 10*time.Minute)
	return profilecache
}

func CreateNonce() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("[DEBUG] Generate Uuid failed %s", err)
		return ""
	}
	return uuid.String()
}
