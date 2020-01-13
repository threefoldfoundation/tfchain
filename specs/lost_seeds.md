# Lost seeds

We have more and more early token holders who loose their seed and get their tft back from the foundation or the person/company they bought them from.

Besides the fact that this costs the refunder TFT, this also diminishes the total amount of TFT available and it might be a fraud.

Since these people are well known and as such have an implicit KYC, there is a technical solution we can implement to solve this:

- Block the old address
- Mint new tokens on their new address.

Since the tokens on the old addresses are locked, burnt or wathever you call it, the total amount of available tokens backed by capacity does not increase.

In order to avoid the impression that the foundation does this at random, a document might be needed in case legal consequences follow for abuse.

A hash of this document can be included in the locking transaction and the hash of the locking transaction can be added to the minting transaction so can proof we did not mint tokens that are not burnt.
