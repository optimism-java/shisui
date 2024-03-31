package beacon

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/protolambda/zrnt/eth2/beacon/common"

	"github.com/protolambda/ztyp/codec"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/protolambda/zrnt/eth2/configs"
	"github.com/stretchr/testify/require"
)

type Entry struct {
	ContentKey   string `json:"content_key"`
	ContentValue string `json:"content_value"`
}

func TestLightClientBootstrap(t *testing.T) {
	data, err := os.ReadFile("./test_data/light_client_bootstrap.json")
	require.NoError(t, err)
	en := new(Entry)
	err = json.Unmarshal(data, en)
	require.NoError(t, err)
	value := hexutil.MustDecode(en.ContentValue)
	forkName := value[:4]

	value = value[4:]
	reader := codec.NewDecodingReader(bytes.NewReader(value), uint64(len(value)))

	spec := &SpecWithForkDigest{
		Spec:       configs.Mainnet,
		ForkDigest: common.ForkDigest(forkName),
	}

	bootStrap := &LightClientBootstrap{}
	err = bootStrap.Deserialize(spec, reader)
	require.NoError(t, err)

	var buf bytes.Buffer
	writer := codec.NewEncodingWriter(&buf)
	err = bootStrap.Serialize(spec, writer)
	require.NoError(t, err)
	require.True(t, bytes.Equal(value, buf.Bytes()))
}
