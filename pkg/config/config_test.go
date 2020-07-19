package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yolossn/query2metric/pkg/config"
)

func TestFromFile(t *testing.T) {
	t.Parallel()

	// File exists
	conf, err := config.FromFile("../../test/test_config.yaml")
	require.NoError(t, err)

	require.Equal(t, len(conf.Connections), 2)

	conf, err = config.FromFile("../../test/not_exists.yaml")
	require.Error(t, err)
	require.Nil(t, conf)

	conf, err = config.FromFile("../../test/invalid_config.yaml")
	require.Error(t, err)
	require.Nil(t, conf)
}
