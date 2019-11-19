package dnsv2

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type Recordset struct {
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	TTL   int      `json:"ttl"`
	Rdata []string `json:"rdata"`
} //`json:"recordsets"`

type MetadataH struct {
	ShowAll       bool `json:"showAll"`
	TotalElements int  `json:"totalElements"`
} //`json:"metadata"`

type RecordSetResponse struct {
	Metadata   MetadataH   `json:"metadata"`
	Recordsets []Recordset `json:"recordsets"`
}

/*
type RecordSetResponse struct {
	Metadata struct {
		ShowAll       bool `json:"showAll"`
		TotalElements int  `json:"totalElements"`
	} `json:"metadata"`
	Recordsets []struct {
		Name  string   `json:"name"`
		Type  string   `json:"type"`
		TTL   int      `json:"ttl"`
		Rdata []string `json:"rdata"`
	} `json:"recordsets"`
}
*/

/*
{
  "metadata": {
    "zone": "example.com",
    "types": [
      "A"
    ],
    "page": 1,
    "pageSize": 25,
    "totalElements": 2
  },
  "recordsets": [
    {
      "name": "www.example.com",
      "type": "A",
      "ttl": 300,
      "rdata": [
        "10.0.0.2",
        "10.0.0.3"
      ]
    },
    {
      "name": "mail.example.com",
      "type": "A",
      "ttl": 300,
      "rdata": [
        "192.168.0.1",
        "192.168.0.2"
      ]
    }
  ]
}

*/

func FullIPv6(ip net.IP) string {
	dst := make([]byte, hex.EncodedLen(len(ip)))
	_ = hex.Encode(dst, ip)
	return string(dst[0:4]) + ":" +
		string(dst[4:8]) + ":" +
		string(dst[8:12]) + ":" +
		string(dst[12:16]) + ":" +
		string(dst[16:20]) + ":" +
		string(dst[20:24]) + ":" +
		string(dst[24:28]) + ":" +
		string(dst[28:])
}

func padvalue(str string) string {
	v_str := strings.Replace(str, "m", "", -1)
	v_float, err := strconv.ParseFloat(v_str, 32)
	if err != nil {
		return "FAIL"
	}
	v_result := fmt.Sprintf("%.2f", v_float)

	return v_result
}

// Used to pad coordinates to x.xxm format
func padCoordinates(str string) string {

	s := strings.Split(str, " ")
	lat_d, lat_m, lat_s, lat_dir, long_d, long_m, long_s, long_dir, altitude, size, horiz_precision, vert_precision := s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], s[8], s[9], s[10], s[11]
	return lat_d + " " + lat_m + " " + lat_s + " " + lat_dir + " " + long_d + " " + long_m + " " + long_s + " " + long_dir + " " + padvalue(altitude) + "m " + padvalue(size) + "m " + padvalue(horiz_precision) + "m " + padvalue(vert_precision) + "m"

}

func NewRecordSetResponse(name string) *RecordSetResponse {
	recordset := &RecordSetResponse{}
	return recordset
}

func GetRecordList(zone string, name string, record_type string) (*RecordSetResponse, error) {
	records := NewRecordSetResponse(name)

	req, err := client.NewRequest(
		Config,
		"GET",
		"/config-dns/v2/zones/"+zone+"/recordsets?types="+record_type+"&showAll=true",
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}
	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, &ZoneError{zoneName: name}
	} else {
		err = client.BodyJSON(res, records)
		if err != nil {
			return nil, err
		}
		return records, nil
	}
}

func GetRdata(zone string, name string, record_type string) ([]string, error) {
	records, err := GetRecordList(zone, name, record_type)
	if err != nil {
		return nil, err
	}

	var arrLength int
	for _, c := range records.Recordsets {
		if c.Name == name {
			arrLength = len(c.Rdata)
		}
	}

	rdata := make([]string, 0, arrLength)

	for _, r := range records.Recordsets {
		if r.Name == name {
			for _, i := range r.Rdata {
				str := i

				if record_type == "AAAA" {
					addr := net.ParseIP(str)
					result := FullIPv6(addr)
					str = result
				} else if record_type == "LOC" {
					str = padCoordinates(str)
				}
				rdata = append(rdata, str)
			}
		}
	}
	return rdata, nil
}
