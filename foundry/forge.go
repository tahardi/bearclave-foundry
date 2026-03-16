package foundry

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"os/exec"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tahardi/bearclave-foundry/internal"
)

const (
	ForgeCommand  = "forge"
	ForgeVersion  = "1.5.1-stable"
	ScriptCommand = "script"

	BroadcastFlag  = "--broadcast"
	PrivateKeyFlag = "--private-key"
	RPCFlag        = "--rpc-url"

	// BroadcastPath is where the `forge script` stores the resulting broadcast file.
	// It should be: <broadcast_dir>/<script_name>/<chain_id>/run-latest.json
	//
	// Example: ../../contracts/broadcast/HelloWorld.s.sol/31337/run-latest.json
	BroadcastPath = "%s/%s/%d/run-latest.json"

	// ScriptName is the name of the file containing the scrip to deploy a contract.
	// It should be: <contract_name>.s.sol
	//
	// Example: HelloWorld.s.sol
	ScriptName = "%s.s.sol"

	// ScriptPath is the script path to use with the `forge script` command.
	// It should be: <script_dir>/<script_name>:<contract_name>Script
	//
	// Example: ../../contracts/scripts/HelloWorld.s.sol:HelloWorldScript
	ScriptPath = "%s/%s:%sScript"
)

var (
	ErrForge = errors.New("forge")
)

type Forge struct {
	broadcastDir string
	scriptDir    string
	url          string
	chainID      *big.Int
}

func NewForge(
	ctx context.Context,
	broadcastDir string,
	scriptDir string,
	url string,
	chainID *big.Int,
) (*Forge, error) {
	_, err := exec.LookPath(ForgeCommand)
	if err != nil {
		msg := ForgeCommand + " command not found in system path"
		return nil, anvilError(msg, err)
	}

	cmd := exec.CommandContext(ctx, ForgeCommand, "--version")
	out, err := cmd.Output()
	if err != nil {
		msg := "running " + ForgeCommand + " --version: "
		return nil, anvilError(msg, err)
	}

	version := strings.TrimSpace(string(bytes.TrimSpace(out)))
	if !strings.Contains(version, ForgeVersion) {
		slog.Warn(
			"version mismatch",
			"command", ForgeCommand,
			"expected", ForgeVersion,
			"got", version,
		)
	}

	return &Forge{
		broadcastDir: broadcastDir,
		scriptDir:    scriptDir,
		url:          url,
		chainID:      chainID,
	}, nil
}

// DeployContract deploys a smart contract via the `forge script` command.
// The command format is:
// forge script <script_path> --rpc-url <rpc_url> --private-key <private_key> --broadcast
//
// Example:
//
//	forge script ../../contracts/scripts/HelloWorld.s.sol:HelloWorldScript \
//	    --rpc-url http://127.0.0.1:8545 \
//	    --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
//	    --broadcast
func (f *Forge) DeployContract(
	ctx context.Context,
	contractName string,
	owner *Account,
) (*common.Address, error) {
	scriptName := fmt.Sprintf(ScriptName, contractName)
	scriptPath := fmt.Sprintf(ScriptPath, f.scriptDir, scriptName, contractName)
	args := []string{
		ScriptCommand, scriptPath,
		RPCFlag, f.url,
		PrivateKeyFlag, owner.PrivateKeyHex(),
		BroadcastFlag,
	}

	deploy := exec.CommandContext(ctx, ForgeCommand, args...)
	out, err := deploy.CombinedOutput()
	if err != nil {
		msg := "deploying contract: " + string(out)
		return nil, forgeError(msg, err)
	}

	broadcastPath := fmt.Sprintf(BroadcastPath, f.broadcastDir, scriptName, f.chainID)
	bytes, err := os.ReadFile(broadcastPath)
	if err != nil {
		msg := "reading broadcast file: " + broadcastPath
		return nil, forgeError(msg, err)
	}

	broadcast := &internal.Broadcast{}
	err = json.Unmarshal(bytes, broadcast)
	if err != nil {
		return nil, forgeError("unmarshaling broadcast", err)
	}

	address, err := broadcast.GetContractAddress(contractName)
	if err != nil {
		return nil, forgeError("getting contract address", err)
	}
	return address, nil
}
