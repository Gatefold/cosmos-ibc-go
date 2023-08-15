package host

import (
	"regexp"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
)

// pathRegex defines a regular expression to extract values between forward slashes (/) from a path.
// [^/] is a negated set to match any character which is not a forward slash.
// + matches one or more occurrences of the preceding token (up to n slashes).
//
// Examples:
// - a/b/c matches a, b, and c
// - /a matches a
var pathRegex = regexp.MustCompile("[^/]+")

// ParseIdentifier parses the sequence from the identifier using the provided prefix. This function
// does not need to be used by counterparty chains. SDK generated connection and channel identifiers
// are required to use this format.
func ParseIdentifier(identifier, prefix string) (uint64, error) {
	if !strings.HasPrefix(identifier, prefix) {
		return 0, errorsmod.Wrapf(ErrInvalidID, "identifier doesn't contain prefix `%s`", prefix)
	}

	splitStr := strings.Split(identifier, prefix)
	if len(splitStr) != 2 {
		return 0, errorsmod.Wrapf(ErrInvalidID, "identifier must be in format: `%s{N}`", prefix)
	}

	// sanity check
	if splitStr[0] != "" {
		return 0, errorsmod.Wrapf(ErrInvalidID, "identifier must begin with prefix %s", prefix)
	}

	sequence, err := strconv.ParseUint(splitStr[1], 10, 64)
	if err != nil {
		return 0, errorsmod.Wrap(err, "failed to parse identifier sequence")
	}
	return sequence, nil
}

// MustParseClientStatePath returns the client ID from a client state path. It panics
// if the provided path is invalid or if the clientID is empty.
func MustParseClientStatePath(path string) string {
	clientID, err := parseClientStatePath(path)
	if err != nil {
		panic(err.Error())
	}

	return clientID
}

// parseClientStatePath returns the client ID from a client state path. It returns
// an error if the provided path is invalid.
func parseClientStatePath(path string) (string, error) {
	matches := pathRegex.FindAllString(path, -1)

	if len(matches) != 3 {
		return "", errorsmod.Wrapf(ErrInvalidPath, "cannot parse client state path %s, %v", path, matches)
	}

	if matches[0] != string(KeyClientStorePrefix) {
		return "", errorsmod.Wrapf(ErrInvalidPath, "path does not begin with client store prefix: expected %s, got %s", KeyClientStorePrefix, matches[0])
	}

	if matches[2] != KeyClientState {
		return "", errorsmod.Wrapf(ErrInvalidPath, "path does not end with client state key: expected %s, got %s", KeyClientState, matches[2])
	}

	if strings.TrimSpace(matches[1]) == "" {
		return "", errorsmod.Wrap(ErrInvalidPath, "clientID is empty")
	}

	clientID := matches[1]

	return clientID, nil
}

// ParseConnectionPath returns the connection ID from a full path. It returns
// an error if the provided path is invalid.
func ParseConnectionPath(path string) (string, error) {
	matches := pathRegex.FindAllString(path, -1)

	// localhost connection paths are just /path
	if len(matches) == 1 {
		return matches[0], nil
	}

	if len(matches) != 2 {
		return "", errorsmod.Wrapf(ErrInvalidPath, "cannot parse connection path %s", path)
	}

	connectionID := matches[1]

	return connectionID, nil
}

// ParseChannelPath returns the port and channel ID from a full path. It returns
// an error if the provided path is invalid.
func ParseChannelPath(path string) (string, string, error) {
	split := strings.Split(path, "/")
	if len(split) < 5 {
		return "", "", errorsmod.Wrapf(ErrInvalidPath, "cannot parse channel path %s", path)
	}

	if split[1] != KeyPortPrefix || split[3] != KeyChannelPrefix {
		return "", "", errorsmod.Wrapf(ErrInvalidPath, "cannot parse channel path %s", path)
	}

	portID := split[2]
	channelID := split[4]

	return portID, channelID, nil
}

// MustParseConnectionPath returns the connection ID from a full path. Panics
// if the provided path is invalid.
func MustParseConnectionPath(path string) string {
	connectionID, err := ParseConnectionPath(path)
	if err != nil {
		panic(err)
	}
	return connectionID
}

// MustParseChannelPath returns the port and channel ID from a full path. Panics
// if the provided path is invalid.
func MustParseChannelPath(path string) (string, string) {
	portID, channelID, err := ParseChannelPath(path)
	if err != nil {
		panic(err)
	}
	return portID, channelID
}
