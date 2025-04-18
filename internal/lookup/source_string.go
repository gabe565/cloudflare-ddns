// Code generated by "enumer -type Source -transform snake -linecomment -output source_string.go"; DO NOT EDIT.

package lookup

import (
	"fmt"
	"strings"
)

const _SourceName = "cloudflare_tlscloudflareopendns_tlsopendnsicanhazipipinfoipify"

var _SourceIndex = [...]uint8{0, 14, 24, 35, 42, 51, 57, 62}

const _SourceLowerName = "cloudflare_tlscloudflareopendns_tlsopendnsicanhazipipinfoipify"

func (i Source) String() string {
	if i >= Source(len(_SourceIndex)-1) {
		return fmt.Sprintf("Source(%d)", i)
	}
	return _SourceName[_SourceIndex[i]:_SourceIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _SourceNoOp() {
	var x [1]struct{}
	_ = x[CloudflareTLS-(0)]
	_ = x[Cloudflare-(1)]
	_ = x[OpenDNSTLS-(2)]
	_ = x[OpenDNS-(3)]
	_ = x[ICanHazIP-(4)]
	_ = x[IPInfo-(5)]
	_ = x[IPify-(6)]
}

var _SourceValues = []Source{CloudflareTLS, Cloudflare, OpenDNSTLS, OpenDNS, ICanHazIP, IPInfo, IPify}

var _SourceNameToValueMap = map[string]Source{
	_SourceName[0:14]:       CloudflareTLS,
	_SourceLowerName[0:14]:  CloudflareTLS,
	_SourceName[14:24]:      Cloudflare,
	_SourceLowerName[14:24]: Cloudflare,
	_SourceName[24:35]:      OpenDNSTLS,
	_SourceLowerName[24:35]: OpenDNSTLS,
	_SourceName[35:42]:      OpenDNS,
	_SourceLowerName[35:42]: OpenDNS,
	_SourceName[42:51]:      ICanHazIP,
	_SourceLowerName[42:51]: ICanHazIP,
	_SourceName[51:57]:      IPInfo,
	_SourceLowerName[51:57]: IPInfo,
	_SourceName[57:62]:      IPify,
	_SourceLowerName[57:62]: IPify,
}

var _SourceNames = []string{
	_SourceName[0:14],
	_SourceName[14:24],
	_SourceName[24:35],
	_SourceName[35:42],
	_SourceName[42:51],
	_SourceName[51:57],
	_SourceName[57:62],
}

// SourceString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func SourceString(s string) (Source, error) {
	if val, ok := _SourceNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _SourceNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Source values", s)
}

// SourceValues returns all values of the enum
func SourceValues() []Source {
	return _SourceValues
}

// SourceStrings returns a slice of all String values of the enum
func SourceStrings() []string {
	strs := make([]string, len(_SourceNames))
	copy(strs, _SourceNames)
	return strs
}

// IsASource returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Source) IsASource() bool {
	for _, v := range _SourceValues {
		if i == v {
			return true
		}
	}
	return false
}
