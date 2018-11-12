# 3Bot

A _3Bot_ —a virtual service— is a digital avatar for a person or
group of persons running in a ThreeFold Container.

The ThreeFold chain ([tfchain][tfchain]) is also used to store the information for every registered 3Bot. A registration happens by creating a 3Bot Registration Tx containing the initial information about a 3Bot on the [tfchain][tfchain].

This has several benefits:

- All 3Bot registration records are distributed and available at all times;
- It is ensured that all [names](#bot-name) (used by 3Bots) are unique and owned by a single [public key](#public-key);
- It is guaranteed that [the required fees](#fees) are paid without the need of a central service, third party financial service or even fiat currency;

Fees are paid in TFT, as part of the initial registration Tx as well as any update Tx after that. Records of a 3Bot will exist forever in the `tfchain` registry, once created, but it can be checked if a 3Bot is active by comparing its queried expiration epoch with the current time as indicated by `tfchain`.

It is up to to 3Bot DNS services to keep its own DNS registry up to date, taking into account the expiration of names as well as any manual updates applied to 3Bot records over time.

## Index

1. [Records](#records): explains what 3Bot records are;
    * 1.1 [Record Updates](#record-updates): explains how [a 3Bot record](#records) can be updated;
2. [Fees](#fees): explains the fees that have to be paid for a 3Bot transaction and how it is computed;
3. [Consensus Rules](#consensus-rules): explains the consensus rules that apply to all 3Bot transactions;
4. [Types](#types): explains types specific to 3Bot transactions/[records](#records).

## Records

The information of a 3Bot is stored in a record, a structural object, identified by the unique (32-bit unsigned integer) identifier of the 3Bot. The record is created and updated using Transaction versions specifically made available for this purpose. All records can contain the following information:

- **Unique ID**: a unique incremental/sequential (4-byte integral) identifier, assigned to every registered 3Bot. The first valid BotID is `1`, not `0`. The order of the ID is based on the total 3Bot Transaction Count. Meaning that if the `unique ID` counter is at 10, and a block is registered containing 3 transactions of which the first and third are 3Bot registration Tx's, the one of the first Tx will be assigned unique ID `10`, while the latter Tx will be assigned `11`. If a few blocks later there is another 3Bot registration Tx it will be assigned `12` and so on...;
- **List of Names**: inspired by DNS names, it are one or multiple optional [names](#bot-name) that can be assigned to a 3Bot, such that you can reach a 3Bot using one of its [names](#bot-name), rather than having to directly use its [IP address or hostname](#network-address). The [tfchain][tfchain] registry defines no link between **the list of names** and **the list of addresses**, this is a detail that has to be worked out by the services (such as 3Bot DNS services) that consume this data;
- **List of Network Addresses**: [IPv4/6 addresses or (domain) hostnames](#network-address) that can be used to reach a 3Bot on. It is optional and can be left empty (if and only if there is at least one [name](#bot-name) registered) as to be able to register a bot simply to reserve one or multiple [names](#bot-name) for it already, without the 3Bot actually being active yet);
- **Public Key**: The unique [Public Key](#public-key) (the [ed25519][ed25519] algorithm is the only supported one for the initial deployment of this feature) that is used by the 3Bot to proof that it has the authority to change its record, as to be able to make any future updates as well as the initial registration;
- **Expiration Epoch Time**: Expiration Epoch Time, defining until when the [names](#bot-name) for a given 3Bot are active/claimed. Beyond this Epoch time the [names](#bot-name) will still be stored in the record, but should be seen as inactive by the consumer of this data (e.g. 3Bot DNS services). This implies also that when a 3Bot is expired, that any 3Bot (including this 3Bot) can (re)claim the expired [names](#bot-name);
    - Note that the record of an expired 3Bot might still contain the [names](#bot-name) as defined by that 3Bot prior to expiring, even though the 3Bot no longer owns these [names](#bot-name). Therefore it is very important that any service sitting on top of a 3Bot record DB checks the expiration date prior to consumption;

Ideally a 3Bot record database stores this information as compact as possible, but this is not a strict requirement. What is required however that the database respects the limits imposed for all used types. You can read more about these limits in [the Consensus Rules chapter](#consensus-rules) chapter.

Note that a single 3Bot will get a unique ID assigned only once, at the point of registration. Once defined it isn't changed or decoupled from the [public key](#public-key) (unique to that 3Bot as well), no matter what or how many updates it receives.

> For now a 3Bot can only get to know its unique ID once its registration Tx is accepted by the consensus as part of a created block. Once that is the case, an up-to-date explorer node will be able to return the 3Bot's record (including its unique ID) given the correct (string/text encoded) public key. See [the Rest API](#rest-api) chapter for more information.
>
> In the future we can probably handle this more elegantly using the [Rivine][rivine]-developed Electrum module. This way the 3Bot doesn't have to poll an explorer node until it knows its unique ID, until than however there is no other option. See the [Rivine][rivine] issue at <https://github.com/threefoldtech/rivine/issues/408> for more information about the upcoming Electrum module.

Extra information, which is not strictly required in order to consume the data, that could be stored by an explorer for a given 3Bot record:

- A list of identifiers of all the transactions that have affected the 3Bot record, which includes all update Tx's (should those exist) as well as the initial registration Tx. If such a list is maintained, it is recommended to preserve the order of the transactions in the order defined by the blockchain;
  - The TransactionDB module (shipping with the Go reference implementation of tfchain), stores this;

## Record Updates

Once a 3Bot is registered, its record can be updated without having to register (again). This saves you in fees and saves the blockchain in the amount of bytes to store.

Please read [the Fees chapter](#fees) to know the total amount of required additional fees for any given combination of updates. Also make sure to understand [the consensus rules](#consensus-rules) as it defines the limits of an update (as well as registration).

The following updates (to a record) can happen:

- the number of months can be extended;
- an inactive 3Bot can be activated again (by extending the number of months);
    - (!) this has the effect that the time of the update block gets used as the start time of the new 3Bot activity period;
    - (!) this also has as effect that all [names](#bot-name) that were still registered in the inactive 3Bot's record up to that point get implicitly removed;
- one or multiple [name(s)](#bot-name) can be added (if the [name](#bot-name) is available and the bot has less than 5 [names](#bot-name) after applying the [names](#bot-name) to-be removed);
- one or multiple [name(s)](#bot-name) can be removed (only if the 3Bot owns these [names](#bot-name));
- one or multiple [network address(es)](#network-address) can be added (if the bot has less than 10 [addresses](#network-address) after applying the [addresses](#network-address) to-be removed);
- one or multiple [network address(es)](#network-address) can be removed (if the bot has these [addresses](#network-address) registered);

A 3Bot (record) cannot be deleted (the blockchain never forgets, unless it forks). You can however deactivate it, by ensuring all [network addresses](#network-address) are removed. No refunds are given. Should you want you can also remove all [(DNS) names](#bot-name) to free them up already (again no refunds are given), otherwise they'll expire once the record's Expiration Epoch time has been reached. Deleting data from a record requires no additional fees.

## Fees

Registering a new 3Bot as well as other actions require additional fees, That go on top of the regular required (minimum) transaction fee (of `0.1 TFT`).

The 3Bot fees, paid to the Threefold Foundation, are as follows:

- registration of a new 3Bot (static price): `100 TFT`;
- monthly fee per 3Bot (static price): `10 TFT`:
  - first month is for free at registration time;
- deletion and transfer of a [name](#bot-name): free;
- per [name](#bot-name): `50 TFT`:
  - at registration time the fee is to be paid only for each additional [name](#bot-name),
    while the first one is for free;
  - when modifying a 3Bot record the fee is applied to each added [name](#bot-name);
- [network address](#network-address) info change (static price): `20 TFT`;

The monthly fee is a static value, and ensures the 3Bot remains active. An inactive bot will still exist in the registry, but will no longer be supported by any ThreeFold Foundation service that runs on top of such registry.

Registering a 3Bot also requires a minimum of one month
(which is already paid for in the registration costs of `100 TFT`),
as a consequence it is not possible to register an inactive 3Bot.
In other words, a 3Bot can only become inactive by not paying the
required monthly fee of `10 TFT` before its expiration timestamp has been reached at least one block less than the highest block.

A 3Bot can register [one name](#bot-name) and up to 10 [network addresses](#network-address) free of charge. Modifying [addresses](#network-address) or adding names post-registration is never free however. At any given block height, a 3Bot is only allowed up to 5 [names](#bot-name) and 10 [network addresses](#network-address).

Additionally the following discounts on the monthly fees apply:

- `30%` discount when paying `[12,23]` months<sup>(1)</sup> at once;
- `50%` discount when paying `24+` months<sup>(1)</sup> at once;

> (1) one month is defined as `30 * 24 * 60 * 60 = 2592000` seconds.

### Example: a minimal 3Bot

In order to minimize the costs for a 3Bot one can therefore choose
to register only what is included in the required frees and use what is free:

- 0 or 1 [name(s)](#bot-name);
- 0 to 10 [network address)(es)](#network-address);

Which would give us the following example additional fee table for the registration for the 3Bot:

|number of months|additional fees in TFT|total discount in TFT|discount per month in TFT|
|-|-|-|-|
|1|`100`|`0`|`0`|
|3|`100`|`0`|`0`|
|12|`174`|`36`|`3`|
|24|`210`|`120`|`5`|

As you can see, the difference between 12 months and 24 months is pretty small, making it pretty attractive to sign up immediately for a 2 year period. While saving you a lot of coins, it doesn't lock you to a specific (set of) [name(s)](#bot-name), as this information (as well as [the network addresses](#network-address) used) can still be changed, without affecting the activity period of the 3Bot (or its (to be) paid fees).

## Consensus Rules

Once you understand how the [fees](#fees) work and what properties [a 3Bot record](#records) contains, you'll notice that the consensus rules are straightforward.

Here is the complete list of rules applied on all 
3Bot transactions:

- The total sum of Miner Fees has to equal at least the minimum transaction fee (`0.1 TFT`);
- The additional fees have to be exactly the amount of additional fees computed as described in [the Fees chapter](#fees), for simplicity<sup>(2)</sup> and fairness there can be no extra fees given;
- All fees (meaning the combination of miner and additional fees) should be funded with given coin inputs;
- Each coin input has to be valid according to the standard rules;
- The refund coin output is optional, and there can only be one;
- At any _resulting_ point no more than 5 [names](#network-address) can be registered for a single 3Bot (_resulting_ meaning that if you update a 3Bot that already has 4 [names](#bot-name) you can add 2 [names](#bot-name) ONLY if you also remove 1 in that same update Tx);
- At any _resulting_ point no more than 10 [network addresses](#network-address) can be registered for a single 3Bot (_resulting_ meaning that if you update a 3Bot that already has 9 addresses you can add 2 [addresses](#network-address) ONLY if you also remove 1 in that same update Tx);
- All [names](#network-address) have to be valid (more about this later);
- All [network addresses](#network-address) have to be valid, a [network address](#network-address) can be: IPv4, IPv6 or a (domain) hostname);
- At any resulting point the number of months (stored as an epoch time, defining a range between the current chain time and that epoch time) has to be in the inclusive range of `[0, 24]` (`0` implying the 3Bot is inactive);
- The signature has to be valid:
  - meaning the input data is as expected, and completely based on the given Tx data;
  - the signature is signed using the private key paired with the known/given [public key](#public-key) (only at registration the public key is given);

> (2) the 3Bot fee is implicitly defined. In other words it is not defined in the Transaction,
but instead has be computed. Computing the extra fee that is to be paid for a 3Bot transaction
can be done in two different ways:
> - (2a) By manually computing the required fee;
> - (2b) By using the following algorithm: `sum(inputs) - refundOutput - minterFee`;
>
> As a result, should a 3Bot fee be allowed to include tips,
> the internal consensus engine would be
> required to look up the values of each used (coin) input
> (as unspent output) in order to know how much is
> is to be transferred as a blockchain-generated coin output
> to the Threefold Foundation (2a). By ruling this option out,
> and only accepting the exact required amount to be paid,
> we can quickly compute the fee ourself (2b), using nothing
> other than the content of the transaction it relates to
> in order to know how the expected coin output to the
> Threefold Foundation would look like.

Each [name](#bot-name) has to be formatted according to the following rules:
- It can be maximum 63 bytes long;
- It can consists of a group of characters, separated by the `.` (dot) character;
- Each group has to have at least 5 characters and can contain only numerical and alphabetical ASCII characters (both lowercase and uppercase are allowed):
    - Note that a group can not be made of numerical characters only;
    - Also note that names have to be unique in terms of the case-insensitive version;
- At least one group (of characters) is required;

A [name](#bot-name) can only be registered if it is available:
- a [name](#bot-name) is available if it was never registered;
- if the last 3Bot that registered that name is no longer active:
  - either because it is expired (because it did not pay any longer);
- the last 3Bot that owned it removed the [name](#bot-name) explicitly;
- the [name](#bot-name) is transferred to the 3Bot registering the [name](#bot-name);

[Network addresses](#network-address) only need to be unique within the context of a single 3Bot ([record](#record)). Meaning that a single 3Bot cannot define the same [(network) address](#network-address) more than once (as that wouldn't make any sense). But it is perfectly fine for multiple 3Bots to define the same [network address](#network-address) (each 3Bot once).

> Note that the consensus engine doesn't check the equality of IPv4 and IPv6 addresses.

## Types

3Bot transactions and [records](#record) consist of many different properties. This chapter will go over the types that require special attention, either because they're not intuitive or because they have been newly introduced only since and because of the 3Bot feature.

Please also ensure to read the [binary_encoding.md](binary_encoding.md) document, as all 3Bot transactions and [records](#bot-record) are binary-encoded using the tfchain binary encoding logic.

### Network Address

A network address defines an address on which a 3Bot can be reached using its public API.
Within the context of tfchain a network address can be a (network) hostnames (FDQN, see RFC 1178),
IPv4- (RFC 791) and IPv6 (RFC 2460) addresses.

The string format, also used for JSON encoding, of each type should be obvious.
Information about the binary encoding of network addresses can be found at [binary_encoding.md#Network-Address](binary_encoding.md#Network-Address).

### Bot Name

A Bot name, used to name/alias a 3Bot, allows a DNS-like service, such that a 3Bot
can be available (on its public API) using its name instead of directly using one of its network addresses.

Each name has the following rules:
- a name consists of groups of characters;
- at least one character group is required;
- character groups are separated by an ASCII dot (`.`);
- each character group consists of a minimum of 5 bytes;
- a character group can only start with an alphabetical ASCII character;
- all but the first character of a character group can also be ASCII numerical
  characters and the ASCII dash character (`-`);
- a maximum (total) size of 63 bytes.

Which in the reference Go implementation is represented by the following Regular Expression:

```regexp
^[A-Za-z]{1}[A-Za-z\-0-9]{3,61}[A-Za-z0-9]{1}(\.[A-Za-z]{1}[A-Za-z\-0-9]{3,55}[A-Za-z0-9]{1})*$
```

Note that even though both upper-case and lower-case ASCII alphabetical characters are accepted,
all ASCII alphabetical characters are normalized into lower-case characters. Meaning that for example
`"ThisIs.SomeBotName" == "thisis.somebotname"`.

The string format, also used for JSON encoding, are all ASCII characters encoded directly as an UTF-8 encoded string.
The binary encoding works analog to the string encoding, except that it is encoded into an UTF-8 character slice, instead of a string.

### Public Key

The string format, also used for JSON encoding, of a public (encryption) key, used to verify signatures,
works exactly the same as Rivine public keys. The binary encoding is a different story however,
as tfchain choose for a compactor format of it, within the context of 3Bot transactions and records.
Information about the binary encoding of public keys, within the context of the tfchain 3Bot feature
can be found at [binary_encoding.md#Public-Key](binary_encoding.md#Public-Key).

[tfchain]: http://github.com/threefoldfoundation/
[rivine]: http://github.com/threefoldtech/rivine
[ed25519]: https://en.wikipedia.org/wiki/EdDSA
