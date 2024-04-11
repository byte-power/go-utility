package utility

import (
	"math/rand"
	mRand "math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	uuidTimestampLength = 40
	b32alphabet         = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	b64characterSet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	b32WordLength       = 5
	b32Chopper          = 31
	b32TimestampLength  = uuidTimestampLength / b32WordLength

	MaxSnowflakeNodeIndex = 1023
)

var serverStartTime = time.Date(2019, time.July, 12, 0, 0, 0, 0, time.UTC)
var randomIDNode *snowflake.Node

// GenerateFixedLengthUUID generates base32 uuid with the given length.
// the generated uuid consists of two parts：`timestampBase32String` and `randomBase32String`
// `timestampBase32String` is generated by base32-encoded millisecond timestamp, length is 8.
// `randomBase32String` is generated by base32-encoded random int64 which is generated by `math/rand` library.
// parameter `length` should be at least 8 if uniqueness need to be guaranteed in one millisecond.
func GenerateFixedLengthUUID(length uint8) string {
	now := time.Now()
	milliseconds := now.Sub(serverStartTime).Milliseconds()
	encodedTimestamp := uuidEncode(milliseconds, uuidTimestampLength)
	if length <= b32TimestampLength {
		return encodedTimestamp[0:length]
	}
	encodedRandomString := GenerateFixedLengthRandomString(length - b32TimestampLength)
	return encodedTimestamp + encodedRandomString
}

func uuidEncode(number int64, length int) string {
	encodeLength := int(length / b32WordLength)
	result := make([]byte, encodeLength)
	for i, j := 0, encodeLength-1; i < encodeLength; i, j = i+1, j-1 {
		x := number & b32Chopper
		result[j] = b32alphabet[x]
		number >>= b32WordLength
	}
	return BytesToString(result)
}

func generateB32EncodedSnowflakeID() string {
	snowflakeID := randomIDNode.Generate()
	f := snowflakeID.Int64()
	if f < 32 {
		return string(b32alphabet[f])
	}
	b := make([]byte, 0, 12)
	for f >= 32 {
		b = append(b, b32alphabet[f%32])
		f /= 32
	}
	b = append(b, b32alphabet[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}
	return BytesToString(b)
}

func GenerateFixedLengthRandomString(length uint8) string {
	charSetLength := int64(len(b32alphabet))
	b := make([]byte, length)
	for i := range b {
		randInt64 := rand.Int63()
		index := randInt64 % charSetLength
		b[i] = b32alphabet[index]
	}
	return BytesToString(b)
}

// InitSnowflakeNode 使用index初始化随机ID生成器，index需在[0,1023]之间.
func InitSnowflakeNode(index int64) {
	if index > MaxSnowflakeNodeIndex {
		index = index % MaxSnowflakeNodeIndex
	}
	snowflake.Epoch = serverStartTime.Unix() * 1000
	randomIDNode, _ = snowflake.NewNode(index)
	return
}

// GenerateSnowflakeID 产生随机ID.
func GenerateSnowflakeID() snowflake.ID {
	return randomIDNode.Generate()
}

func TimeInSnowflakeInt(id int64) time.Time {
	return TimeInSnowflakeID(snowflake.ID(id))
}

func TimeInSnowflakeID(id snowflake.ID) time.Time {
	// (id >> (snowflake.NodeBits + snowflake.StepBits)) + snowflake.Epoch
	ts := id.Time()
	sec := ts / 1000
	return time.Unix(sec, (ts-sec*1000)*int64(time.Millisecond))
}

func GenUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func StringWithCharset(length int, charset string) string {
	seededRand := mRand.New(
		mRand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenID(prefix string) string {
	length := 13
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	return prefix + StringWithCharset(length, charset)
}

// GenPassword 生成随机密码，32位长度，只包含小写英文及数字，开头不是数字
func GenPassword() string {
	rand.Seed(time.Now().UnixNano())

	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	password := make([]rune, 32)

	password[0] = letters[rand.Intn(25)]
	for i := 1; i < 32; i++ {
		password[i] = letters[rand.Intn(len(letters))]
	}
	return string(password)
}
