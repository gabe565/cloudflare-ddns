// Code generated by "enumer -type Source -trimprefix Source -transform snake -linecomment -output source_string.go"; DO NOT EDIT.

package config

import (
	"fmt"
	"strings"
)

const _SourceName = "cloudflare_tlscloudflareopendns_tlsopendnsipinfoipify"

var _SourceIndex = [...]uint8{0, 14, 24, 35, 42, 48, 53}

const _SourceLowerName = "cloudflare_tlscloudflareopendns_tlsopendnsipinfoipify"

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
	_ = x[SourceCloudflareTLS-(0)]
	_ = x[SourceCloudflare-(1)]
	_ = x[SourceOpenDNSTLS-(2)]
	_ = x[SourceOpenDNS-(3)]
	_ = x[SourceIPInfo-(4)]
	_ = x[SourceIPify-(5)]
}

var _SourceValues = []Source{SourceCloudflareTLS, SourceCloudflare, SourceOpenDNSTLS, SourceOpenDNS, SourceIPInfo, SourceIPify}

var _SourceNameToValueMap = map[string]Source{
	_SourceName[0:14]:       SourceCloudflareTLS,
	_SourceLowerName[0:14]:  SourceCloudflareTLS,
	_SourceName[14:24]:      SourceCloudflare,
	_SourceLowerName[14:24]: SourceCloudflare,
	_SourceName[24:35]:      SourceOpenDNSTLS,
	_SourceLowerName[24:35]: SourceOpenDNSTLS,
	_SourceName[35:42]:      SourceOpenDNS,
	_SourceLowerName[35:42]: SourceOpenDNS,
	_SourceName[42:48]:      SourceIPInfo,
	_SourceLowerName[42:48]: SourceIPInfo,
	_SourceName[48:53]:      SourceIPify,
	_SourceLowerName[48:53]: SourceIPify,
}

var _SourceNames = []string{
	_SourceName[0:14],
	_SourceName[14:24],
	_SourceName[24:35],
	_SourceName[35:42],
	_SourceName[42:48],
	_SourceName[48:53],
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
