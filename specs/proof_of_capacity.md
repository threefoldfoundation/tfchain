# Proof of capacity

## The problem

In the threefold grid, the farmers are rewarded with token for the amount of capacity their provide to the grid.  
This capacity commitment needs to be registered on the blockchain.
This involves multiple things:

- a node need to be link to a farmer
- the capacity that get registered on the blockchain needs to be provable

## Current status

As of today, the capacity of a node is computed automatically by some code and then publish to a centralized database.  
There are already some problem with this flow:

- There is no real proof of capacity. Capacity can easily be faked, which would lead to the evil famer to get paid for capcity he doesn't provide
- Anyone can today register capacity, even if they don't actually provide any capacity to the grid.

This shows that today, the validity of the capacity registered is only based on the trust that farmer will behave well and nobody will try to fake capacity.

We need to come up with some algorithm that provide some real proof of capacity. We also need to be able to limit the right to register capacity to only people that are verified and accepted as farmer.