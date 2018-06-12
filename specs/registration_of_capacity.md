# Registration of capacity

Registration of capacity could be done within TFChain,
by means of supporting it as a new kind of Transaction.
This would give us following benefits:

+ It would give us freedom about what properties to put into the transaction date;
+ It would allow us to leave out any unnecessary properties (e.g. block stake inputs/outputs);
+ It would make it easy for software running on top of TFChain, to sort out the capacity registrations;

Additionally it would make it possible and easy for the daemon consensus
logic to validate the registration of capacity, probably on a more basic level,
with a full proof of capacity upon request.

When encoded in the JSON format, the transaction could look as follows:

```javascript
{
    "version": 128, // could be any number,
                 // but doing it in sequence and starting from the
                 // 2nd half of the available range makes sense
    "data": {
        // we could leave the farm property out,
        // should we link the wallet address to the farm (somehow) instead
        "farm": {}, // id of the farm, as to link this capacity to a farm?
        "capacity": {}, // defines the entire specification of the capacity to be provided,
                        // we can do it as one object, or an array of tuples in the form of {type,value}
        "coininputs": [], // regular coin inputs
        "coinoutputs": [], // regular coin outputs
        "minerfees": ["42"], // regular miner fees
        "arbitrarydata": "optional notes, base64 encoded",
    }
}
```

While the transaction would still allow to attach any kind of input/output transfers,
the motivation to allow for both coin inputs is more because we need the coin inputs in order
to fund the miner fees. Because of that, we need the coin outputs in order to allow for refunds.
Abuse, if you wish to call it that, of this feature is however possible,
but we should not have to care about that, as it is not of any importance.

The capacity property is the important one.
How it will look will depend upon the units of capacity which we wish to register,
and also on whether or not certain units can be duplicated and/or omitted.

The farm property, as noted before, could be omitted.
If however we do not wish to implicitly link wallet addresses to farms, than we will need
to link it using such a farm property (object) instead.
It is currently unclear/undefined how the registration of the farm will happen.

> Even though the registration of a farm is undefined.
> I am certain that we want to enforce that only authorized addresses
> can register capacity. If so, we might indeed already know which capacity belongs to which farm,
> based on the wallet address+key+signature used, and thus, no explicit farm identification property might be required.

Arbitrary data could be omitted, but I vote to keep it in.
It might be useful to add some additional free-format data,
or even it could be used to implement an extension to the capacity of registration protocol,
by the software built on top of TFChain. Whatever the case, the TFChain daemon will ignore it,
just as is the case with the arbitrary data of regular transactions.

Implementing this type of transaction is trivial, due to how Rivine is written we can simply define this
as a new Transaction Type, for which we have to define the encoding and decoding logic. As well as any other extension logic.
Actually validating this registered capacity is a a lot more difficult and comes bundled with a lot of problems to overcome.
