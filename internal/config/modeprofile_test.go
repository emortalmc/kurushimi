package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestModeProfileValidity(t *testing.T) {
	for _, profile := range ModeProfiles {
		t.Run(profile.Name, func(t *testing.T) {
			require.Greater(t, profile.MinPlayers, 0, "MinPlayers must be greater than 0")
			require.Greater(t, profile.MaxPlayers, 0, "MaxPlayers must be greater than 0")
			require.GreaterOrEqual(t, profile.MaxPlayers, profile.MinPlayers, "MaxPlayers must be greater than or equal to MinPlayers")
			require.NotNil(t, profile.Selector, "Selector must not be nil")
			require.NotNil(t, profile.MatchFunction, "MatchFunction must not be nil")
		})
	}
}
