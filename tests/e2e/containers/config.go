package containers

// ImageConfig contains all images and their respective tags
// needed for running e2e tests.
type ImageConfig struct {
	InitRepository string
	InitTag        string

	EveRepository string
	EveTag        string

	RelayerRepository string
	RelayerTag        string
}

//nolint:deadcode
const (
	// Current Git branch eve repo/version. It is meant to be built locally.
	// It is used when skipping upgrade by setting EVE_E2E_SKIP_UPGRADE to true).
	// This image should be pre-built with `make docker-build-debug` either in CI or locally.
	CurrentBranchRepository = "eve"
	CurrentBranchTag        = "debug"
	// Pre-upgrade eve repo/tag to pull.
	// It should be uploaded to Docker Hub. EVE_E2E_SKIP_UPGRADE should be unset
	// for this functionality to be used.
	// previousVersionRepository = "ghcr.io/cosmoscontracts/eve"
	// previousVersionTag        = "0.0.1-e2e"
	// Pre-upgrade repo/tag for eve initialization (this should be one version below upgradeVersion)
	// previousVersionInitRepository = "ducnt87/eve"
	// previousVersionInitTag        = "0.0.1-e2e-init-chain"
	// Hermes repo/version for relayer, use osmosis dockerhub
	relayerRepository = "osmolabs/hermes"
	relayerTag        = "0.13.0"
)

// Returns ImageConfig needed for running e2e test.
// If isUpgrade is true, returns images for running the upgrade
// If isFork is true, utilizes provided fork height to initiate fork logic
func NewImageConfig(isUpgrade, isFork bool) ImageConfig {
	config := ImageConfig{
		RelayerRepository: relayerRepository,
		RelayerTag:        relayerTag,
	}

	if !isUpgrade {
		// If upgrade is not tested, we do not need InitRepository and InitTag
		// because we directly call the initialization logic without
		// the need for Docker.
		config.EveRepository = CurrentBranchRepository
		config.EveTag = CurrentBranchTag
		return config
	}

	// If upgrade is tested, we need to utilize InitRepository and InitTag
	// to initialize older state with Docker
	// config.InitRepository = previousVersionInitRepository
	// config.InitTag = previousVersionInitTag

	// if isFork {
	// 	// Forks are state compatible with earlier versions before fork height.
	// 	// Normally, validators switch the binaries pre-fork height
	// 	// Then, once the fork height is reached, the state breaking-logic
	// 	// is run.
	// 	config.EveRepository = CurrentBranchRepository
	// 	config.EveTag = CurrentBranchTag
	// } else {
	// 	// Upgrades are run at the time when upgrade height is reached
	// 	// and are submitted via a governance proposal. Thefore, we
	// 	// must start running the previous Eve version. Then, the node
	// 	// should auto-upgrade, at which point we can restart the updated
	// 	// Eve validator container.
	// 	config.EveRepository = previousVersionRepository
	// 	config.EveTag = previousVersionTag
	// }

	return config
}
