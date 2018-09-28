# 3bot registration

> This is a technical specification, with all details worked out
> based on the functional spec at <https://github.com/rivine/home/blob/master/specs/3bot_registration.md>
> as well as discussions about it.

A _3bot_ —a virtual service— is a digital avatar for a person or
group of persons running in a ThreeFold Container.

The ThreeFold chain ([tfchain][tfchain]) will be used to store the information
for every registered 3bot. A registration happens by creating a 3bot Registration Tx
containing the initial information about a 3bot on the [tfchain][tfchain].
This has several benefits:

- All 3bot registration records are distributed and available at all times;
- It is ensured that all DNS names (used by 3bots) are unique and owned by a single public key;
- It is guaranteed that the required fees are paid without the need of a central service.

Fees are paid in TFT, as part of the initial registration Tx as well as any update Tx after that.
Records of a 3bot will exist forever in the `tfchain` registry, but it can be checked if a 3bot
is active by comparing its queried expiration epoch with the current epoch. It is up to to
3bot DNS services to keep its own DNS registry up to date, taking into account the expiration
of names as well as any manual updates applied to 3bot records over time.

## 3bot records

Explorers will keep a single record per 3bot —created for a new registration (Tx)
and kept up to date with every update Tx that applies an update to it— containing the following info:

- **Unique ID**: a unique incremental/sequential (4-byte integral) identifier, assigned to every registered 3bot. The order of the ID is based on the total 3bot Transaction Count. Meaning that if the `unique ID` counter is at 10, and a block is registered containing 3 transactions of which the first and third are 3bot registration Tx's, the one of the first Tx will be assigned unique ID `10`, while the latter Tx will be assigned `11`. If a few blocks later there is another 3bot registration Tx it will be assigned `12` and so on...;
- **List of Names**: inspired by DNS names, it are one or multiple optional names that can be assigned to a 3bot, such that you can reach a 3bot using one of its names, rather than having to directly use its IP address or hostname. The [tfchain][tfchain] registry defines no link between **the list of names** and **the list of addresses**, this is a detail that has to be worked out by the services (such as 3bot DNS services) that consume this data;
- **List of Network Addresses**: IPv4/6 addresses or (domain) hostnames that can be used to reach a 3bot on. It is optional and can be left empty (if and only if there is at least one name registered) as to be able to register a bot simply to reserve one or multiple names for it already, without the 3bot actually being active yet);
- **Public Key**: The unique Public Key (the [ed25519][ed25519] algorithm is the only supported one for the initial deployment of this feature) that is used by the 3bot to proof that it owns the rights for this 3bot, as to be able to make any future updates as well as the initial registration;
- **Expiration Epoch Time**: Expiration Epoch Time, defining until when the names for a given 3bot are active/claimed. Beyond this Epoch time the names will still be stored in the record, but should be seen as inactive by the consumer of this data (e.g. 3bot DNS services). This implies also that when a 3bot is expired, that any 3bot (including this 3bot) can (re)claim the expired names;

Note that a single 3bot will get a unique ID assigned only once, at the point of registration. It will never change ID because of an update, and it will never lose it (deleting a 3bot record is impossible, as the blockchain never forgets).

> For now a 3bot can only get to know its unique ID once its registration Tx is accepted by the consensus as part of a created block. Once that is the case, an up-to-date explorer node will be able to return the 3bot's record (including its unique ID) given the correct (string/text encoded) public key. See [the Rest API](#rest-api) chapter for more information.
>
> In the future we can probably handle this more elegantly using the [Rivine][rivine]-developed Electrum module. This way the 3bot doesn't have to poll an explorer node until it knows its unique ID, until than however there is no other option. See the [Rivine][rivine] issue at <https://github.com/rivine/rivine/issues/408> for more information about the upcoming Electrum module.

Extra information, which is not strictly required in order to consume the data, that could be stored by an explorer for a given 3bot record:

- A list of identifiers of all the transactions that have affected the 3bot record, which includes all update Tx's (should those exist) as well as the initial registration Tx;

## Fees

Registering a new 3bot as well updating an existing 3bot (record) requires
additional fees on top of the regular Tx fee (which is currently defined at the minimum of `0.1 TFT`).

Here are the base values used for the computation of the total additional fees:

- each (DNS) name costs `10 TFT` (with a maximum of 5 names and thus a maximum of `50 TFT`);
- the first 3 network addresses requires no additional fees;
- every extra network address costs `5 TFT` (a maximum amount of 10 network addresses, and thus a maximum of `35 TFT`);
- deleting data never requires additional fees;
- registration costs are `80 TFT`;
- per month is a `10 TFT` additional fee required for operational costs;
- modification of a properties costs `40 TFT` (with number of months being an exception);

In terms of a 3bot record:

- a month is 30 days;
- years do not exist in 3bot terminology, giving us `360` days for `12` months (**not** `365`);

Depending upon the number of months, an automatic discount is applied,
decreasing the amount of additional fees required:

- 15% discount is applied if at least 3 months, but less than 12 months is paid at once;
- 30% discount is applied if at least 12 months, but less than 24 months is paid at once;
- 50% discount is applied if at least 24 months is paid at once;

These discounts are given for 2 reasons:

- to reward the 3bot's tokens paid up front;
- to reward the 3bot of saving us precious bytes on the blockchain, given that the 3bot won't have to extend the Expiration time (using the Nr of months) as fast;

All this gives us the following formula to compute the total amount of required additional fees for a registration Tx Fee:

> F<sub>additional</sub> = `80 TFT` +
>   (
>       (C<sub>names</sub> * `10 TFT`) +
>       ((C<sub>addresses</sub> < 3 ? 0 : C<sub>addresses</sub>-3) * `5 TFT`) +
>       `10` TFT
>   ) * T<sub>months</sub> * R<sub>discount</sub>
>
> where:
>  - `R` is one of {`1.0`, `0.85`, `0.7`, `0.5`} (see discounts)

For update Tx's the formula gets a bit more hairy:

> F<sub>additional</sub> = `X TFT` +
>   ((
>       (C<sub>new names</sub> * `10 TFT`) +
>       ((C<sub>new addresses</sub> < 3 ? 0 : C<sub>new addresses</sub>-3) * `5 TFT`) +
>       `10` TFT
>   ) * T<sub>months</sub> * R<sub>discount(T)</sub>) +
>   ((
>       (C<sub>remaining names</sub> * `10 TFT`) +
>       ((C<sub>remaining addresses</sub> < 3 ? 0 : C<sub>remaining addresses</sub>-3) * `5 TFT`) +
>       `10 TFT`
>   ) * N<sub>months</sub> * R<sub>discount(N)</sub>)
>
> where:
>  - `X` equals:
>    - `0` if only Nr of months is defined
>    - `40` if any other property are (also) defined
>  - `R` is one of {`1.0`, `0.85`, `0.7`, `0.5`} (see discounts)
>  - T<sub>months</sub> equals the the total amount of months (remaining months the both is active + the given Nr of months)
>  - N<sub>months</sub> equals the the given Nr of months
>  - remaining meaning the names/addresses that were not removed and not added,
>    but already registered in a previous registration/update Tx
>  - no refunds are given, meaning that if you remove an address and/or name
>    which was already paid for (T<sub>months</sub> - N<sub>months</sub>) amount of months, those months are lost
>    - ⚠ Note that removing a (DNS) name makes it immediately available for any 3bot to claim it

### Examples

#### a minimal bot

In order to minimize the costs for a 3bot one can therefore choose
to register only the required data would give us:

- no (DNS) names;
- 1 to 3 network addresses;

Which would give us the following example additional fee table for the registration for the 3bot:

|number of months|additional fees in TFT|total discount in TFT|discount per month in TFT|
|-|-|-|-|
|1|`90`|`0`|`0`|
|3|`110`|`0`|`0`|
|12|`164`|`36`|`3`|
|24|`200`|`120`|`5`|

#### a typical bot

A more typical bot would have:

- (at least) one DNS name;
- 2 to 3 network addresses;

Which would give us the following example additional fee table for the registration for the 3bot:

|number of months|additional fees in TFT|total discount in TFT|discount per month in TFT|
|-|-|-|-|
|1|`100`|`0`|`0`|
|3|`140`|`0`|`0`|
|12|`248`|`72`|`6`|
|24|`320`|`240`|`10`|

> ⚠ It is no coincidence that the registration of a "typical" 3bot
> for 1 month costs exactly `100 TFT`. 

## Consensus Rules

Once you understand how the [Fees](#fees) work and what properties a 3bot record contains,
you'll notice that the consensus rules are straightforward.

Here is the complete list of rules applied on all 
3bot Registration/Update Tx's:

- The total sum of Miner Fees has to equal at least the minimum Tx fee (`0.1 TFT`);
- The additional fees have to at equal at least the amount of additional fees computed as described in [the fees chapter](#fees), anything extra is considered a donation towards the Threefold Foundation (fair CLI tools would warn the user for this though);
- All fees (meaning the combination of miner and additional fees) should be funded with given coin inputs;
- Each coin input has to be valid according to the standard rules;
- The refund coin output is optional and can be defined only to allow change given back to a wallet of choice);
- No extra coin inputs can be defined than needed (meaning that if you need to pay 100 TFT, and have already 2 coin inputs of `30 TFT` and `70 TFT`, than a third coin input would not be allowed, as it would result in a pure coin transfer, which is not allowed as part of a 3bot Tx);
- At any _resulting_ point no more than 5 (DNS) names can be registered for a single 3bot (_resulting_ meaning that if you update a 3bot that already has 4 DNS names you can add 2 DNS names ONLY if you also remove 1 in that same update Tx);
- At any _resulting_ point no more than 10 network addresses can be registered for a single 3bot (_resulting_ meaning that if you update a 3bot that already has 9 DNS names you can add 2 DNS names ONLY if you also remove 1 in that same update Tx);
- All DNS names have to be valid (more about this later);
- All network addresses have to be valid, a network address can be: IPv4, IPv6 or a (domain) hostname);
- At any resulting point the number of months has to be in the inclusive range of `[1, 24]` (meaning after the combination of the remaining months plus the newly added number of months);
- The signature has to be valid:
  - meaning the input data is as expected, and completely based on the given Tx data;
  - the signature is signed using the private key paired with the known/given public key (only at registration the public key has to be given);

Each (DNS) name has to be formatted according to the following rules:
- It can be maximum 63 bytes long;
- It can consists of a group of characters, separated by the `.` (dot) character;
- Each group has to have at least 5 characters and can contain only numerical and alphabetical ASCII characters (both lowercase and uppercase are allowed). Note that a group can not be made of numerical characters only;
- At least one group is required;

A (DNS) name can only be registered if it is available:
- a (DNS) name is available if it was never registered;
- if the last 3bot that registered that (DNS) name is no longer active:
  - either because it is expired (because it did not pay any longer);
  - or because it was deleted;
  - or because it deleted the (DNS) name;

Network addresses only need to be unique within the context of a single 3bot (record). Meaning that a single 3bot cannot define the same (network) address more than once (as that wouldn't make any sense). But it is perfectly fine for multiple 3bots to define the same IP (each 3bot once).

## Updates

Once a 3bot is registered, its record can be updated without having to register. This saves you in fees and saves the blockchain in the amount of bytes to store.

Please read [the fees chapter](#fees) to know the total amount of required additional fees for any given combination of updates. Also make sure to understand [the consensus rules](#consensus-rules) as it defines the limits of an update (as well as registration).

The following updates can happen:

- the number of months can be extended;
- one or multiple (DNS) name(s) can be added;
- one or multiple (DNS) name(s) can be removed;
- one or multiple network address(es) can be added;
- one or multiple network address(es) can be removed;

A 3bot (record) cannot be deleted (the blockchain never forgets). You can however deactivate it, by ensuring all network addresses are removed. No refunds are given. Should you want you can also remove all (DNS) names to free them up already (again no refunds are given), otherwise they'll expire once the record's Expiration Epoch time has been reached.

## REST API

In order for this feature to work, no extra endpoints are required. In order to create transactions the user (or program) can use the existing `POST /transactionpool/transactions` endpoint, also used to for example register a coin creation transaction.

While the explorer module is optional, it will be useful to index all 3bot records there as well, as to ensure any user can easily fetch a 3bot record.

### Explorer

The explorer module will provide the following endpoints for this feature:

#### `GET /explorer/3bot/:id`

`:id` can be the unique ID as well as the 3bot's public key (string/text encoded).

If found, an object in the structure of the following example is returned:

```javascript
{
    // unique ID (4 byte unsigned integer)
    "id": 42,
    // list of registered (DNS) names,
    // each name is a string
    "names": [],
    // list of registered network addresses,
    // each address is a string
    "addresses": [],
    // string encoded public key, defined in the registration Tx
    "publickey": "ed25519:28c1edd4c35f662cccfa7fc02194959d75855c02d342c1131b110c9e96764d9b",
    // Unix Epoch Time in seconds,
    // of when the 3bot is to expire,
    // meaning the 3bot (DNS) names will no longer
    // be pointing to the network addresses of this 3bot
    "expiration": 943916400,
}
```

#### `GET /explorer/3bot/:id/transactions`

`:id` can be the unique ID as well as the 3bot's public key (string/text encoded).

If found, an object in the structure of the following example is returned:

```javascript
{
    // all transactions which affected the 3bot record
    "transactions": [
        // the first Tx is always the registration Tx
        "84168dcf36f98a804dc52a0d285f3cb6a8b9ffa8ee69385a54b2d65d455a8060",
        // any other Tx is an update Tx,
        // in the order that it happened.
        "a6aca83fe8f51e939db0431e78f59686b5bd9d1b744308fe958fb9a9f7c17b9c",
    ]
}
```

## Explorer Web API

You will be able to search for a 3bot in the official web explorer. For this to work there will be on the home page a third search function:

```
Search a 3bot: [id.or.publickey.....................] [Go]
```

If the given id/publickey is valid, it will redirect you to a new `3bot` page which shows the most up-to-date record of a 3bot.

## Transactions

The following transaction types allow a 3bot to be registered and updated. It also allows two 3bots to transfer (DNS) names from one to the other.

The goal is to keep the footprint as low as possible for the 3bot transaction types in terms of space. Therefore the binary encoded should be done as compact as possible, breaking perhaps some conventions set in other (standard) Tx types).

### Registration Tx

The (3bot) Registration Tx is be used to register a 3bot. Some specifics for this Tx:

- it can only be used to register a new bot (identified by its public key);
- no network addresses or (DNS) names can be removed, given there is nothing yet registered for this Tx;
- it is the only 3bot transaction where the public key is registered on the chain;
- NrOfMonths has to be at least `1` (max `24`);
- Even though both the list of names and list of addresses are optional, at least one name or address has to be given, so only one of them is optional, but the 3bot can choose which;

Read up on [the Consensus Rules chapter](#consensus-rules) to learn about the other requirements/rules.

JSON-Encoded the Tx will look as follows:

```javascript
{
    "version": 144, // 0x90
    "data": {
         // optional network address, max 10
        "addresses": [
            // can be a domain name
            "network.address.a.com",
            "network.address.b.com",
            // can be an IPv4 address
            "83.200.201.201",
            // can also be an IPv6 address
            "2001:db8:85a3::8a2e:370:7334"
        ],
        // optional (dns) names, maximum 5
        "names": [ 
            // see consensus rules for name formatting requirements
            "char5",
            "char5.char5"
        ],
        // NOTE that even though both addresses and names
        // are optional, at least one address or name is required,
        // it is not valid for both to be empty

        // uint8, one of inclusive range: [0,24]
        "nrofmonths": 1,

        // the regular Tx fee has to be paid as well and defined explicitly.
        // Additional fees are assumed to equal to `sum(coinInputs)-txFee-coinOutput`,
        // in other words the additional fees are implicitly defined.
        "txfee": "100000000",

        // coin inputs used to fund the fees (same as regular Txs)
        // this allows also by the way for any party to pay,
        // but most logically is that the 3bot pays, it is however not enforced
        "coininputs": [],
        // optional, and is only allowed to be used for a single refund
        "coinoutput": {},

        // public key (ed25519 by default, and only algo supported for now)
        "publickey": "ed25519:28c1edd4c35f662cccfa7fc02194959d75855c02d342c1131b110c9e96764d9b",

        // signature (verifiable using the given public key),
        // (ed25519 by default, and only algo supported for now)
        "signature": "433f5283a82fac28dacddeb98ed57aabb649a4aad2e7813af8a0009e0d035f625724dc7ef9ef39e75aef10fc77c4ed43fa0ce09f80c77d81ffd0c04ee7ca8c00"
    }
}
```

For the most part the binary encoding of the Registration Tx is straightforward, but there are some specificities:

- An address is encoded in a new type, which has a single byte to denote the type.
  The last 5 bits can be used by the type for any other data it requires (such as the length for the domain hostname).
  The encoding of those bytes depends upon the type of network address (e.g. IPv4 requires 4 bytes);
- A name is encoded using a new type, where the string bytes are prefixed with one byte to indicate the length;
- The addresses and names are encoded directly after one another,
  with the entire group of bytes a single byte to define the amount of addresses and names
  (the first 4 bits defines amount of addresses, and other 4 bits the amount of names);

> NOTE: See [The Compact Binary Properties Chapter](#compact-binary-properties) to see how other (common) properties,
> are also highly optimized as to keep the binary-encoded Tx as small as possible.

> NOTE: for the address type we have 3 different possibilities. IPv4, IPv6 and domain host names.
> The first 2 have fixed-size values, and only the latter (domain host names) requires a length.
> For this one we make use of the remaining 5 bits:
>  * Is the first bit equal to `0`, than the length is defined by the remaining 4 bits;
>  * Otherwise the length is defined the next {1,2,3,4} byte(s) depending if the 2nd, 3rd, 4th or 5th bit is 1;

### Update Tx

The (3bot) Update Tx is be used to update the record of an existing 3bot. Some specifics for this Tx:

- it can only be used to update an existing bot;
- it does not contain the public key, only a signature, which have to be made using the private key paired to the already registered public key;
- the public key can be (and has to be) looked up using the given ID;
- network names and (DNS) names can be removed, as well as added;
- NrOfMonths can be `0` (but still only `24`);
- both the list of (DNS) names and network addresses are optional;

Read up on [the Consensus Rules chapter](#consensus-rules) to learn about the other requirements/rules.

JSON-Encoded the Tx will look as follows:

```javascript
{
    "version": 145, // 0x91
    "data": {
        // unique ID of 3bot to be updated
        "id": 1,

        // optional added/removed network addresses
        "addresses": {
            // addresses added
            // (cannot be registered yet for this 3bot)
            "add": [ // optional, max 10
                // can be a domain name
                "network.address.a.com",
                // can be an IPv4 address
                "83.200.201.201",
                // can also be an IPv6 address
                "2001:db8:85a3::8a2e:370:7334"
            ],
            // addresses removed
            // (have to be previously registered for this 3bot)
            "rem": [ // optional, max 10
                // can be a domain name
                "network.address.a.com",
                // can be an IPv4 address
                "83.200.201.201",
                // can also be an IPv6 address
                "2001:db8:85a3::8a2e:370:7334"
            ],
        ],
        // optional added/removed (DNS) names
        "names": {
            // (DNS) names added
            // (cannot be registered yet for this 3bot)
            "add": [ // optional, max 5
                "aaaaa.bbbbb",
            ],
            // (DNS) names removed
            // (have to be previously registered for this 3bot)
            "rem": [ // optional max 5
                "char5.char6",
            ],
        },

        // uint8, one of inclusive range: [0,24]
        "nrofmonths": 1,

        // the regular Tx fee has to be paid as well and defined explicitly.
        // Additional fees are assumed to equal to `sum(coinInputs)-txFee-coinOutput`,
        // in other words the additional fees are implicitly defined.
        "txfee": "100000000",

        // coin inputs used to fund the fees (same as regular Txs)
        // this allows also by the way for any party to pay,
        // but most logically is that the 3bot pays, it is however not enforced
        "coininputs": [],
        // optional, and is only allowed to be used for a single refund
        "coinoutput": {},

        // signature (verifiable using the previously registered public key),
        // (ed25519 by default, and only algo supported for now)
        "signature": "433f5283a82fac28dacddeb98ed57aabb649a4aad2e7813af8a0009e0d035f625724dc7ef9ef39e75aef10fc77c4ed43fa0ce09f80c77d81ffd0c04ee7ca8c00"
    }
}
```

For the most part the binary encoding of the Update Tx is straightforward, but there are some specificities:

- An address is encoded in a new type, which has a single byte to denote the type.
  The last 5 bits can be used by the type for any other data it requires (such as the length for the domain hostname).
  The encoding of those bytes depends upon the type of network address (e.g. IPv4 requires 4 bytes);
- A name is encoded using a new type, where the string bytes are prefixed with one byte to indicate the length;
- The added and removed network addresses are encoded directly after one another,
  with the entire group of bytes a single byte to define the amount of addresses added and removed
  (the first 4 bits defines amount of addresses added, and the other 4 bits the amount of addresses removed);
- The added and removed (DNS) names are encoded directly after one another,
  with the entire group of bytes a single byte to define the amount of names added and removed
  (the first 4 bits defines amount of names added, and the other 4 bits the amount of names removed);

> NOTE: See [The Compact Binary Properties Chapter](#compact-binary-properties) to see how other (common) properties,
> are also highly optimized as to keep the binary-encoded Tx as small as possible.

> NOTE: for the address type we have 3 different possibilities. IPv4, IPv6 and domain host names.
> The first 2 have fixed-size values, and only the latter (domain host names) requires a length.
> For this one we make use of the remaining 5 bits:
>  * Is the first bit equal to `0`, than the length is defined by the remaining 4 bits;
>  * Otherwise the length is defined the next {1,2,3,4} byte(s) depending if the 2nd, 3rd, 4th or 5th bit is 1;

To keep the Update Tx as compact as possible for small updates
(including the smallest possible update where only the months of values is updated),
the Tx (binary) encoding defines the following compressed (1 byte) value:

```
[ 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 ]
| NrOfMonths        | V | V | flag to indicate if a refund is given
|                   | V | flag to indicate if any name is added/removed
|                   | flag to indicate if any address is added/removed
```
> This byte is encoded instead of an uint8 value (that you might have expected as the type for NrOfMonths).

### Name Transfer Tx

The (3bot) Transfer Tx is be used to transfer one or multiple (DNS) names owned by one bot, to another bot. of an existing 3bot. Some specifics for this Tx:

- this Tx involves two different and already registered 3bots;
- the receiving both can be inactive up to that point, but the sending bot has to be active;
- the fees are most likely paid by the receiving bot, but this is not required;
- NOTE That the fees are similar to a regular Update Tx:
  - the update fee (`40 TFT`) will have to be paid;
  - in case the receiving 3bot is inactive, it will also have to pay all fees for the already registered info;
  - Nr of months can be extended, but note that it will to be paid for all other registered info as well;
- the (DNS) names transferred (and thus list) have to be owned by an active 3bot up to that point;
- at least one added (DNS) name is required which is transferred from the sending 3bot
- the unique ID of both sender and receiver (3bot) are  to be given, as well the both signatures;

Read up on [the Consensus Rules chapter](#consensus-rules) to learn about the other requirements/rules.

JSON-Encoded the Tx will look as follows:

```javascript
{
    "version": 146, // 0x92
    "data": {
        "sender": {
            // 4byte unique ID of existing 3bot
            "id": 1,
            // signature (verifiable using the given public key),
            // (ed25519 by default, and only algo supported for now)
            // (!) Signature includes the public key of both parties as input,
            //     but the signatures not, so when signing, the initial party,
            //     has to make sure to already define the public key of the other
            //     party as well
            "signature": "signature_sender"
        },
        "receiver": {
            // 4byte unique ID of existing 3bot
             "id": 1,
            // signature (verifiable using the given public key),
            // (ed25519 by default, and only algo supported for now)
            // (!) Signature includes the public key of both parties as input,
            //     but the signatures not, so when signing, the initial party,
            //     has to make sure to already define the public key of the other
            //     party as well
            "signature": "signature_receiver"
        },

        // optional added/removed network addresses
        "addresses": {
            // addresses added
            // (cannot be registered yet for this 3bot)
            "add": [ // optional, max 10
                // can be a domain name
                "network.address.a.com",
                // can be an IPv4 address
                "83.200.201.201",
                // can also be an IPv6 address
                "2001:db8:85a3::8a2e:370:7334"
            ],
            // addresses removed
            // (have to be previously registered for this 3bot)
            "rem": [ // optional, max 10
                // can be a domain name
                "network.address.a.com",
                // can be an IPv4 address
                "83.200.201.201",
                // can also be an IPv6 address
                "2001:db8:85a3::8a2e:370:7334"
            ],
        ],
        // required added/removed (DNS) names
        "names": {
            // at least one transferred DNS name is required!!!
            // (DNS) names added
            // (cannot be registered yet for this 3bot)
            "add": [ // optional, max 5
                // a name is either transferred from
                // the sending 3bot, or not claimed yet
                "aaaaa.bbbbb",
            ],
            // (DNS) names removed
            // (have to be previously registered for the receiving 3bot)
            // useful, as to make place for the new DNS
            "rem": [ // optional max 5
                "char5.char6",
            ],
        },

        // uint8, one of inclusive range: [0,24]
        "nrofmonths": 1,

        // the regular Tx fee has to be paid as well and defined explicitly.
        // Additional fees are assumed to equal to `sum(coinInputs)-txFee-coinOutput`,
        // in other words the additional fees are implicitly defined.
        "txfee": "100000000",

        // coin inputs used to fund the fees (same as regular Txs)
        // this allows also by the way for any party to pay,
        // but most logically is that the 3bot pays, it is however not enforced
        "coininputs": [],
        // optional, and is only allowed to be used for a single refund
        "coinoutput": {},
    }
}
```

For the most part the binary encoding of the Update Tx is straightforward, but there are some specificities:

- An address is encoded in a new type, which has a single byte to denote the type.
  The last 5 bits can be used by the type for any other data it requires (such as the length for the domain hostname).
  The encoding of those bytes depends upon the type of network address (e.g. IPv4 requires 4 bytes);
- A name is encoded using a new type, where the string bytes are prefixed with one byte to indicate the length;
- The added and removed network addresses are encoded directly after one another,
  with the entire group of bytes a single byte to define the amount of addresses added and removed
  (the first 4 bits defines amount of addresses added, and the other 4 bits the amount of addresses removed);
- The added and removed (DNS) names are encoded directly after one another,
  with the entire group of bytes a single byte to define the amount of names added and removed
  (the first 4 bits defines amount of names added, and the other 4 bits the amount of names removed);

> NOTE: See [The Compact Binary Properties Chapter](#compact-binary-properties) to see how other (common) properties,
> are also highly optimized as to keep the binary-encoded Tx as small as possible.

> NOTE: for the address type we have 3 different possibilities. IPv4, IPv6 and domain host names.
> The first 2 have fixed-size values, and only the latter (domain host names) requires a length.
> For this one we make use of the remaining 5 bits:
>  * Is the first bit equal to `0`, than the length is defined by the remaining 4 bits;
>  * Otherwise the length is defined the next {1,2,3,4} byte(s) depending if the 2nd, 3rd, 4th or 5th bit is 1;

To keep the Update Tx as compact as possible for small updates
(including the smallest possible update where only the months of values is updated),
the Tx (binary) encoding defines the following compressed (1 byte) value:

```
[ 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 ]
| NrOfMonths        | V | V | flag to indicate if a refund is given
|                   | V | ALWAYS 1, as a (DNS) name transfer HAS to include at least one added name
|                   | flag to indicate if any address is added/removed
```
> This byte is encoded instead of an uint8 value (that you might have expected as the type for NrOfMonths).

## Compact Binary Properties

The binary encoding of some properties of the new transactions proposed
for this feature have already been discussed in [The Transactions Chapter](#transactions).
These transactions also have properties which are the same or similar to those of standard/regular Tx's.
With an aim on keeping the memory footprint of these properties as small as possible,
such is already achieved with the new properties discussed in [The Transactions Chapter](#transactions),
it is needed to decode these more "classic" types in a new way as well, at least already when used as part of these Transaction Types.

### Dynamic Slices

In [Rivine][rivine] _slice_ types (which include the `string` type) are prefixed with 8 bytes,
containing an 64-bit unsigned integer to indicate the length of the _slice_ value.
Because the dynamic nature of _slices_ the encoding of the length together with its value is unavoidable.
It is however a waste to always use 8 bytes to indicate the slice. It is especially weird as
the [Rivine][rivine] encoding library defines a maximum limit of ~5MB for the size of a slice,
which, when expressed in bytes, gives a number that fits in 4 bytes. Meaning 4 bytes are ALWAYS wasted.
We can however do better than just save 4 bytes, for some instances, without giving up the ability to
encode slices that do require a length that cannot be encoded in less than 4 bytes.

```
         b   b   b   b     b
Given: [ 0 | 1 | 2 | ... | N - 1 ]

bit(0)=0 --> (A) length fits in 1 byte, with a max value (N) of (2^7)-1 (127 Bytes)
bit(0)=1 -+-> bit(1)=0 --> (B) length fits in 2 bytes, with a max value (N) of (2^14)-1 (16+ KB)
          |
          +-> bit(1)=1 -+-> bit(2)=0 --> (C) length fits in 3 bytes, with a max value (N) of (2^21)-1 (2+ MB)
                        |
                        +-> bit(2)=1 --> (D) length fits in 4 bytes, with a max value (N) of (2^29)-1 (536+ MB)
```

This compact scheme allows us for most slices (as most slices we use in Tx's do really fit in 127 bytes or less)
to encode the length as a single byte, while for the rare cases where we need more than 127 bytes,
we can still save 6 bytes by encoding the length as 2 bytes (including the 2bit prefix of case B).

It allows us to encode much compactor, without becoming more CPU intensive.

> It will require a overwrite of type, but the binary encoding of the `Currency` type
> can make use of this optimization as well.

### Dynamic Integers

In [Rivine][rivine] all integral types are little-endian binary encoded as unsigned 64-bit integers,
regardless of their actual type, resulting always in 8 bytes.

This is a waste, and given we assume the decode-callee always gives typed values to decode into,
we can very easily vary between the requirement of 1, 2, 4 or 8 bytes, without having to
use any kind of prefix byte.

+ 1 byte: `uint8`, `int8`
+ 2 bytes: `uint16`, `int16`
+ 4 bytes: `uint32`, `int32`
+ 8 bytes: `uint64`, `int64`, `uint`, `int`

While it is true, that it would still waste for example 3 bytes in the case of a `uint32` value of 255 or less,
optimizing at this level would start to be about trade-offs. If we continue with the `uint32` type example:
using 2 bits you could represent whether it requires 1, 2, 3 or 4 bytes. It does however mean you can now only have
a maximum value of `2^28` (268,435,456) instead of `2^32`.

In general this isn't really something you would want to do. For the rare cases where you do think you could use it,
it is probably best to implement a specialized solution just for it.

### Optimized Public Keys and Signatures

In [Rivine][rivine] public keys and signatures are binary-encoded in a very inefficient way.
Signatures are prefixed with 8 bytes, to indicate the length. Public keys have a overhead of 24 bytes,
8 bytes for the hash length, and 16 bytes for the algorithm specifier.

The length of both the public key and signature should however be static per algorithm type.
Therefore a pair of a public key and signature only requires in fact a 1 byte overhead (in total),
allowing 256 different algorithms, which seems more sufficient, given that today we have only 1 algorithm that we support.

For stand-alone signatures this would mean that we also still get this 1 byte overhead,
which would do as good as the smallest possible [Dynamic Slice](#dynamic-slices) proposed earlier.
Given that a signature is a byte slice, it makes this proposed optimization even for stand-alone signatures
an efficient approach, that wouldn't be beaten by giving up the type-info for the dynamic slice info.
On top of that. Encoding the a type byte, instead of a 1-byte length byte, gives us on top of that the advantage,
that we can make our signatures typed in-memory, providing for an optional cheap pre-check,
based on the public key's type.

### Optimized Coin Inputs and Outputs

Inspired by all optimizations proposed in this spec,
it becomes clear that we can also save a lot of bytes by optimizing
all coin inputs and coin outputs.

For Coin inputs that means optimizing the Fulfillments, as the ParentID is an array.
For Coin outputs that means optimizing the encoding of the Currency and the Conditions.

[tfchain]: http://github.com/threefoldfoundation/
[rivine]: http://github.com/rivine/rivine
[ed25519]: https://en.wikipedia.org/wiki/EdDSA
