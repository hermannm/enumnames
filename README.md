# enumnames

A small Go package for efficient mapping between integer enum values and string names.

Run `go get hermannm.dev/enumnames` to add it to your project!

**Docs:** [pkg.go.dev/hermannm.dev/enumnames](https://pkg.go.dev/hermannm.dev/enumnames)

**Contents:**

- [Usage](#usage)
- [Benchmarks](#benchmarks)
- [Maintainer's guide](#maintainers-guide)

## Usage

In the example below, we have a `MessageType` enum, which we want represented as a `uint8` to take
up as little space as possible. However, we also want to map each value to a name, for debugging and
marshaling/unmarshaling JSON. Here we use `enumnames`, creating a map that we can then use in our
`MessageType` methods:

<!-- @formatter:off -->
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
<!-- @formatter:on -->

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

## Maintainer's guide

### Publishing a new release

- Run tests and linter ([`golangci-lint`](https://golangci-lint.run/)):
  ```
  go test ./... && golangci-lint run
  ```
- Add an entry to `CHANGELOG.md` (with the current date)
    - Remember to update the link section, and bump the version for the `[Unreleased]` link
- Create commit and tag for the release (update `TAG` variable in below command):
  ```
  TAG=vX.Y.Z && git commit -m "Release ${TAG}" && git tag -a "${TAG}" -m "Release ${TAG}" && git log --oneline -2
  ```
- Push the commit and tag:
  ```
  git push && git push --tags
  ```
    - Our release workflow will then create a GitHub release with the pushed tag's changelog entry
