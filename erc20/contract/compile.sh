DIR="$(dirname "$BASH_SOURCE")"

solc --bin --overwrite -o "${DIR}/bin" "${DIR}/proxy.sol"
solc --bin --overwrite -o "${DIR}/bin" "${DIR}/tokenV0.sol"
