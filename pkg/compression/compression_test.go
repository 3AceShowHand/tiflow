package compression

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKlauspostSnappyCompression(t *testing.T) {
	original, err := os.ReadFile("./klauspost-before.json")
	require.NoError(t, err)

	compressed, err := Encode(Snappy, original)
	require.NoError(t, err)

	fmt.Printf("compressed message size is: %d", len(compressed))
}
