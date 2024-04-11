package utility

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"math"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unsafe"
)

func PanicIfNotNil(err error) {
	if err == nil {
		return
	}
	log.Println(err)
	debug.PrintStack()
	panic(err)
}

var envVars map[string]string

func EnvironmentVariables() map[string]string {
	if envVars != nil {
		return envVars
	}
	lines := os.Environ()
	envVars = make(map[string]string, len(lines))
	for _, line := range lines {
		comps := strings.Split(line, "=")
		if len(comps) > 1 {
			envVars[comps[0]] = comps[1]
		}
	}
	return envVars
}

type StrMap = map[string]interface{}
type AnyMap = map[interface{}]interface{}

func AnyToAnyMap(value interface{}) AnyMap {
	if value == nil {
		return nil
	}
	switch val := value.(type) {
	case AnyMap:
		return val
	case StrMap:
		count := len(val)
		if count == 0 {
			return nil
		}
		m := make(AnyMap, count)
		for k, v := range val {
			m[k] = v
		}
		return m
	default:
		return nil
	}
}

func AnyToStrMap(value interface{}) StrMap {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case StrMap:
		return v
	case AnyMap:
		l := len(v)
		if l == 0 {
			return nil
		}
		m := make(StrMap, l)
		for k, v := range v {
			m[AnyToString(k)] = v
		}
		return m
	case map[string]string:
		l := len(v)
		if l == 0 {
			return nil
		}
		m := make(StrMap, l)
		for k, v := range v {
			m[AnyToString(k)] = v
		}
		return m
	default:
		return nil
	}
}

type CustomStringConvertable interface {
	String() string
}

func AnyToString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch val := value.(type) {
	case *string:
		if val == nil {
			return ""
		}
		return *val
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case error:
		return val.Error()
	case CustomStringConvertable:
		return val.String()
	default:
		return fmt.Sprint(value)
	}
}

func AnyToInt64(value interface{}) int64 {
	if value == nil {
		return 0
	}
	switch val := value.(type) {
	case int:
		return int64(val)
	case int8:
		return int64(val)
	case int16:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return val
	case uint:
		return int64(val)
	case uint8:
		return int64(val)
	case uint16:
		return int64(val)
	case uint32:
		return int64(val)
	case uint64:
		return int64(val)
	case *string:
		if val == nil {
			return 0
		}
		if i, err := StringToInt64(*val); err == nil {
			return i
		}
	case string:
		if i, err := StringToInt64(val); err == nil {
			return i
		}
	case float32:
		return int64(val)
	case float64:
		return int64(val)
	case bool:
		if val {
			return 1
		}
		return 0
	case json.Number:
		v, _ := val.Int64()
		return v
	}
	return 0
}

func IsNumeric(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		json.Number:
		return true
	}
	return false
}

func AnyToFloat64(value interface{}) float64 {
	if value == nil {
		return 0
	}
	switch val := value.(type) {
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return val
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case *string:
		if val == nil {
			return 0
		}
		if v, err := strconv.ParseFloat(*val, 64); err == nil {
			return v
		}
	case string:
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			return v
		}
	case bool:
		if val {
			return 1
		}
		return 0
	case json.Number:
		v, _ := val.Float64()
		return v
	}
	return 0
}

func AnyToBool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch v := v.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case int8:
		return v != 0
	case int16:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case uint:
		return v != 0
	case uint8:
		return v != 0
	case uint16:
		return v != 0
	case uint32:
		return v != 0
	case uint64:
		return v != 0
	case float32:
		return v != 0
	case float64:
		return v != 0
	case string:
		if len(v) == 0 {
			return false
		}
		c := strings.ToLower(v[0:1])
		return c == "y" || c == "t" || c == "1"
	case *string:
		return v != nil && AnyToBool(*v)
	default:
		return false
	}
}

func AnyToBoolArray(any interface{}) []bool {
	if any == nil {
		return nil
	}
	switch v := any.(type) {
	case []bool:
		return v
	case []interface{}:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []int:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []int8:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []int16:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []int32:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []int64:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []uint:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []uint8:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []uint16:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []uint32:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []uint64:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []float32:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []float64:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []string:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	case []*string:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			result = append(result, AnyToBool(item))
		}
		return result
	default:
		return nil
	}
}

func AnyToInt64Array(any interface{}) []int64 {
	if any == nil {
		return nil
	}
	switch v := any.(type) {
	case []int64:
		return v
	case []interface{}:
		return AnyArrayToInt64Array(v)
	default:
		return nil
	}
}

func AnyArrayToInt64Array(arrInterface []interface{}) []int64 {
	elementArray := make([]int64, len(arrInterface))
	for i, v := range arrInterface {
		elementArray[i] = AnyToInt64(v)
	}
	return elementArray
}

func AnyToFloat64Array(any interface{}) []float64 {
	if any == nil {
		return nil
	}
	switch v := any.(type) {
	case []float64:
		return v
	case []interface{}:
		return AnyArrayToFloat64Array(v)
	case []bool, []string, []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []float32:
		return stringOrBoolOrNumberArrayToFloat64Array(v)
	default:
		return nil
	}
}

func AnyArrayToFloat64Array(arrInterface []interface{}) []float64 {
	elementArray := make([]float64, len(arrInterface))
	for i, v := range arrInterface {
		elementArray[i] = AnyToFloat64(v)
	}
	return elementArray
}

func stringOrBoolOrNumberArrayToFloat64Array(any interface{}) []float64 {
	floatArray := make([]float64, 0)
	switch array := any.(type) {
	case []string:
		for _, v := range array {
			floatArray = append(floatArray, AnyToFloat64(v))
		}
	case []bool:
		for _, v := range array {
			floatArray = append(floatArray, AnyToFloat64(v))
		}
	case []int:
		for _, v := range array {
			floatArray = append(floatArray, float64(v))
		}
	case []int8:
		for _, v := range array {
			floatArray = append(floatArray, float64(v))
		}
	case []int16:
		for _, v := range array {
			floatArray = append(floatArray, float64(v))
		}
	case []int32:
		for _, v := range array {
			floatArray = append(floatArray, float64(v))
		}
	case []int64:
		for _, v := range array {
			floatArray = append(floatArray, float64(v))
		}
	case []uint:
		for _, v := range array {
			floatArray = append(floatArray, float64(v))
		}
	case []uint8:
		for _, v := range array {
			floatArray = append(floatArray, float64(v))
		}
	case []uint16:
		for _, v := range array {
			floatArray = append(floatArray, float64(v))
		}
	case []uint32:
		for _, v := range array {
			floatArray = append(floatArray, float64(v))
		}
	case []uint64:
		for _, v := range array {
			floatArray = append(floatArray, float64(v))
		}
	case []float32:
		for _, v := range array {
			floatArray = append(floatArray, float64(v))
		}
	case []float64:
		floatArray = array
	}
	return floatArray
}

func AnyArrayToStrMap(mapInterface []interface{}) StrMap {
	if len(mapInterface)/2 < 1 {
		return nil
	}
	elementMap := make(StrMap)
	for i := 0; i < len(mapInterface)/2; i += 1 {
		key := AnyToString(mapInterface[i*2])
		elementMap[key] = mapInterface[i*2+1]
	}
	return elementMap
}

func AnyToStringArray(any interface{}) []string {
	if any == nil {
		return nil
	}
	switch v := any.(type) {
	case []string:
		return v
	case []interface{}:
		return AnyArrayToStringArray(v)
	default:
		return nil
	}
}

func AnyArrayToStringArray(arrInterface []interface{}) []string {
	elementArray := make([]string, len(arrInterface))
	for i, v := range arrInterface {
		elementArray[i] = AnyToString(v)
	}
	return elementArray
}

func StringToInt64(value string) (int64, error) {
	if index := strings.Index(value, "."); index > 0 {
		value = value[:index]
	}
	return strconv.ParseInt(value, 10, 64)
}

// BytesToString 按string的底层结构，转换[]byte
func BytesToString(b []byte) string {
	if b == nil {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes 按[]byte的底层结构，转换字符串，len与cap皆为字符串的len
func StringToBytes(s string) []byte {
	return StringPToBytes(&s)
}

func StringPToBytes(s *string) []byte {
	if s == nil {
		return nil
	}
	// 获取s的起始地址开始后的两个 uintptr 指针
	x := (*[2]uintptr)(unsafe.Pointer(s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// FindInStrings would return first index of `slice` if found `item`, or -1 for not.
func FindInStrings(slice []string, item string) int {
	return FindInSlice(slice, item)
}

func FindInInt64Array(slice []int64, item int64) int {
	return FindInSlice(slice, item)
}

func FindInFloat64Array(slice []float64, item float64) int {
	return FindInSlice(slice, item)
}

func FindInSyncMap(m *sync.Map, keys ...interface{}) interface{} {
	return FindInSyncMapWithKeys(m, keys)
}

func FindInSyncMapWithKeys(m *sync.Map, keys []interface{}) interface{} {
	if m == nil {
		return nil
	}
	l := len(keys)
	if l == 0 {
		return nil
	}
	v0, ok := m.Load(keys[0])
	if !ok || l == 1 {
		return v0
	}
	switch v := v0.(type) {
	case StrMap:
		return FindInStrMapWithKeys(v, keys[1:])
	case AnyMap:
		return FindInAnyMapWithKeys(v, keys[1:])
	default:
		return nil
	}
}

func FindInAnyMap(m AnyMap, keys ...interface{}) interface{} {
	return FindInAnyMapWithKeys(m, keys)
}

func FindInAnyMapWithKeys(m AnyMap, keys []interface{}) interface{} {
	if m == nil {
		return nil
	}
	l := len(keys)
	if l == 0 {
		return nil
	}
	value := m[keys[0]]
	if l == 1 {
		return value
	}
	switch v := value.(type) {
	case AnyMap:
		return FindInAnyMapWithKeys(v, keys[1:])
	case StrMap:
		return FindInStrMapWithKeys(v, keys[1:])
	default:
		return nil
	}
}

func FindInStrMap(m StrMap, keys ...interface{}) interface{} {
	return FindInStrMapWithKeys(m, keys)
}

func FindInStrMapWithKeys(m StrMap, keys []interface{}) interface{} {
	if m == nil {
		return nil
	}
	l := len(keys)
	if l == 0 {
		return nil
	}
	value := m[AnyToString(keys[0])]
	if l == 1 {
		return value
	}
	switch v := value.(type) {
	case AnyMap:
		return FindInAnyMapWithKeys(v, keys[1:])
	case StrMap:
		return FindInStrMapWithKeys(v, keys[1:])
	default:
		return nil
	}
}

// flatten map ,e.g
// map A
//
//	{
//		"foo":{
//			"bar":1
//		}
//	}
//
// map B = FlattenMap("",".",A)
//
//	{
//		"foo.bar":1
//	}
func FlattenMap(rootKey, delimiter string, originData StrMap) StrMap {
	result := make(StrMap, len(originData))
	for key, value := range originData {
		if value == nil {
			continue
		}
		tmpKey := key
		if rootKey != "" {
			tmpKey = rootKey + delimiter + key
		}
		if reflect.ValueOf(value).Kind() == reflect.Map {
			v := AnyToStrMap(value)
			if len(v) == 0 {
				result[key] = v
				continue
			}
			tmpMap := FlattenMap(tmpKey, delimiter, v)
			for k, v := range tmpMap {
				result[k] = v
			}
		} else {
			result[tmpKey] = value
		}
	}
	return result
}

func CanConvertToFloat32Loselessly(v float64) bool {
	absV := math.Abs(v)
	if absV < math.MaxFloat32 && absV > math.SmallestNonzeroFloat32 {
		return true
	}
	return false
}

func CanConvertToInt64Loselessly(v float64) bool {
	return v == math.Trunc(v) && v == float64(int64(v))
}

func CanConvertToInt32Loselessly(v float64) bool {
	return v == math.Trunc(v) && v < math.MaxInt32 && v > math.MinInt32
}

// StringToChunks split a string into string slices with element's size <= chunkSize
// Examples:
// StringToChunks("abcd", 1) => []string{"a", "b", "c", "d"}
// StringToChunks("abcd", 2) => []string{"ab", "cd"}
// StringToChunks("abcd", 3) => []string{"abc", "d"}
// stringToChunks("abcd", 4) => []string{"abcd"}
// stringToChunks("abcd", 5) => []string{"abcd"}
func StringToChunks(s string, chunkSize int) []string {
	var chunks []string
	strLength := len(s)
	index := 0
	for index < strLength {
		endIndex := Min(index+chunkSize, strLength)
		chunk := s[index:endIndex]
		chunks = append(chunks, chunk)
		index = endIndex
	}
	return chunks
}


func Hash(bs []byte) uint32 {
	h := fnv.New32a()
	h.Write(bs)
	return h.Sum32()
}

func IsNumString(s string) bool {
	for _, r := range []rune(s) {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func IsAlphaString(s string) bool {
	for _, r := range []rune(s) {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsAlnumString(s string) bool {
	for _, r := range []rune(s) {
		if !(unicode.IsLetter(r) || unicode.IsNumber(r)) {
			return false
		}
	}
	return true
}

var cookieReg = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func CheckEmail(str string) bool {
	r := cookieReg.MatchString(str)
	return r
}

// AnyToMapStringAny 将 Content 转化为 Map<String>Any，使其序列化时按字典序排序
func AnyToMapStringAny(content any) map[string]any {
	mContent := map[string]any{}
	sContent, err := json.Marshal(content)
	if err != nil {
		return map[string]any{}
	}
	err = json.Unmarshal(sContent, &mContent)
	if err != nil {
		return map[string]any{}
	}
	return mContent
}

// RemovePrefixKeys 删除 Map 中指定前缀的 Key
func RemovePrefixKeys(data map[string]any, prefix string) {
	// 遍历键
	for key := range data {
		// 如果键以 _ 开头
		if strings.HasPrefix(key, prefix) {
			// 删除该键
			delete(data, key)
		}

		// 如果值是 map 类型
		if m, ok := data[key].(map[string]any); ok {
			// 递归调用本函数
			RemovePrefixKeys(m, prefix)
		}
	}
}
