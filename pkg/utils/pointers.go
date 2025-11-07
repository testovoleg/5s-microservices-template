package utils

func DerefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func StrPtr(in string) *string {
	res := in
	return &res
}

func DerefBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

func BoolPtr(in bool) *bool {
	res := in
	return &res
}

func DerefInt(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}

func DerefFloat64(i *float64) float64 {
	if i != nil {
		return *i
	}
	return 0
}

func DerefInt32(i *int32) int32 {
	if i != nil {
		return *i
	}
	return 0
}

func IntPtr(in int) *int {
	res := in
	return &res
}

func Int32Ptr(in int32) *int32 {
	res := in
	return &res
}

func Float64Ptr(in float64) *float64 {
	res := in
	return &res
}

func Stringer(str []*string) []string {
	var strs []string
	for _, s := range str {
		if s == nil {
			strs = append(strs, "")
			continue
		}
		strs = append(strs, *s)
	}

	return strs
}

func StrArrPtr(str []string) *[]string {
	var res []string
	res = append(res, str...)
	return &res
}

func StrPtrArray(str []string) []*string {
	var res []*string
	for _, s := range str {
		res = append(res, &s)
	}
	return res
}
