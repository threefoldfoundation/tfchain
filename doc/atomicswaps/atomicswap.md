# Cross chain-atomic swaps

## Theory

A cross-chain swap is a trade between two users of different cryptocurrencies. For example, one party may send Threefold Tokens to a second party's Threefold Address, while the second party would send Bitcoin to the first party's Bitcoin address. However, as the blockchains are unrelated and transactions cannot be reversed, this does not protect against one of the parties not honoring their end of the deal. One common solution to this problem is to introduce a mutually-trusted third party for escrow. An atomic cross-chain swap solves this problem without the need for a third party. On top of that, it achieves waterproof validation without introducing the problems and complexities introduced by a escrow-based validation system.

Atomic swaps involve each party paying into a contract transaction, one contract for each blockchain. The contracts contain an output that is spendable by either party, but the rules required for redemption are different for each party involved. 

## Example

Let's assume Bob wants to buy 567 TFT from Alice for 0.1234BTC

Bob creates a bitcoin address and Alice creates a Threefold Address.

Bob initiates the swap, he generates a 32-byte secret and hashes it
using the SHA256 algorithm, resulting in a 32-byte hashed secret.

Bob now creates a swap transaction, as a smart contract, and publishes it on the Bitcoin chain, it has 0.1234BTC as an output and the output can be redeemed (used as input) using either 1 of the following conditions:
- a timeout has passed (48hours) and claimed by Bob's refund address;
- the money is claimed by Alice's registered address and the secret is given that hashes to the hashed secret created by Bob 

This means Alice can claim the Bitcoin if she has the secret and if the atomic swap process fails, Bob can always reclaim his BTC after the timeout.

Bob sends this contract and the transaction id of this transaction on the Bitcoin chain to Alice, making sure he does not share the secret of course.

 Now Alice validates if everything is as agreed (=audit)after which She creates a similar transaction on the Rivine chain but with a timeout for a refund of only 24 hours and she uses the same hash secret as the first contract for Bob to claim the tokens.
 This transaction has 9876 tokens as an output and the output can be redeemed( used as input) using either 1 of the following conditions:
- a timeout has passed ( 24hours) and claimed by the seller's refund address
- the secret is given that hashes to the hash secret Bob created (= same one as used in the bitcoin swap transaction) and claimed by the buyer's address

For Bob to claim the Threefold Tokens, he has to use and as such disclose the secret.

The magic of the atomic swap lies in the fact that the same secret is used to claim the tokens in both swap transactions but it is not disclosed in the contracts because only the hash of the secret is used there. The moment Bob claims the Threefold Tokens, he discloses the secret and Alice has enough time left to claim the bitcoin because the timeout of the first contract is longer than the one of the second contract.
Of course, either Bob or Alice can be the initiator or the participant.

## Walkthroughs

Walkthroughs of the above Bitcoin example:
- [using Btc Core (full node)](defaultbtcatomicswap.md)
- [using Electrum (thin client)](electrumbtcatomicswap.md)

Ethereum walkthroughs:
- [using Go Ethereum (full or light node)](defaultethatomicswap.md)

## References

Threefold atomic swaps are default [Rivine](https://github.com/threefoldtech/rivine) atomic swaps.
