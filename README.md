# enumnames

A small Go package for efficient mapping between integer enum values and string names.

Run `go get hermannm.dev/enumnames` to add it to your project!

## Usage

In the example below, we have a `MessageType` enum, which we want represented as a `uint8` to take
up as little space as possible. However, we also want to map each value to a name, for debugging and
marshaling/unmarshaling JSON. Here we use `enumnames`, creating a map that we can then use in our
`MessageType` methods:

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
	return msgNames.ContainsKey(msgType)
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
BenchmarkGetName-16           	1000000000	         0.7659 ns/op
BenchmarkGetNameWithMap-16    	1000000000	         7.797 ns/op
BenchmarkGetKey-16            	1000000000	         5.988 ns/op
BenchmarkGetKeyWithMap-16     	1000000000	         9.355 ns/op
PASS
ok  	hermannm.dev/enumnames	26.329s
```

The benchmarks use a `uint8` as the enum type, with 255 variants. We see that `BenchmarkGetName`,
which uses `enumnames.Map`, is 10x faster than `BenchmarkGetNameWithMap`, which uses a
`map[uint8]string`. The reverse lookup, `BenchmarkGetKey`, is 50% faster than
`BenchmarkGetKeyWithMap`, which uses a `map[string]uint8`.
