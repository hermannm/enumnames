# enumnames

A small Go package for efficient mapping between integer enum values and string names.

Run `go get hermannm.dev/enumnames` to add it to your project!

## Usage

In the example below, we have a `MessageType` enum, which we want represented as a `uint8` to take up as little space as possible. However, we also want to map each value to a name, for debugging and marshaling/unmarshaling JSON. Here we use `enumnames`, creating a map that we can then use in our `MessageType` methods:

```go
import "hermannm.dev/enumnames"

type MessageType uint8

const (
	JoinLobbyMessage MessageType = iota + 1
	ReadyMessage
	StartGameMessage
)

var msgNames = enumnames.NewMap(map[MessageType]string{
	JoinLobbyMessage: "JOIN_LOBBY",
	ReadyMessage:     "READY",
	StartGameMessage: "START_GAME",
})

func (msgType MessageType) IsValid() bool {
	return msgNames.ContainsEnumValue(msgType)
}

func (msgType MessageType) String() string {
	return msgNames.GetNameOrFallback(msgType, "INVALID_MESSAGE_TYPE")
}

func (msgType MessageType) MarshalJSON() ([]byte, error) {
	return msgNames.MarshalToNameJSON(msgType)
}

func (msgType *MessageType) UnmarshalJSON(bytes []byte) error {
	return msgNames.UnmarshalFromNameJSON(bytes, msgType)
}
```

## Benchmarks

Result of running `go test -bench=. -benchtime=10s`:

```
goos: linux
goarch: amd64
pkg: hermannm.dev/enumnames
cpu: AMD Ryzen 7 PRO 6850U with Radeon Graphics
BenchmarkGetName-16                     	1000000000	         0.7203 ns/op
BenchmarkGetNameWithMap-16              	1000000000	         7.505 ns/op
BenchmarkEnumValueFromName-16           	1000000000	         6.079 ns/op
BenchmarkEnumValueFromNameWithMap-16    	1000000000	         9.025 ns/op
PASS
ok  	hermannm.dev/enumnames	25.787s
```

The benchmarks use a `uint8` as the enum type, with 255 variants. We see that `BenchmarkGetName`, which uses `enumnames.Map`, is 10x faster than `BenchmarkGetNameWithMap`, which uses `map[uint8]string`. The reverse lookup, `BenchmarkEnumValueFromName`, is almost 50% faster than `BenchmarkEnumValueFromNameWithMap`, which uses a `map[string]uint8`.
