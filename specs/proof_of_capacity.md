# Proof of capacity

## The problem

In the threefold grid, the farmers are rewarded with token for the amount of capacity they provide to the grid.  
This capacity commitment needs to be registered on the blockchain.
This involves multiple things:

- a node need to be linked to a farmer
- the capacity that gets registered on the blockchain needs to be provable

## Current status

As of today, the capacity of a node is computed automatically by some code and than published to a centralized database.  
There are already some problem with this flow:

- There is no real proof of capacity. Capacity can easily be faked, which would lead to the evil farmer to get paid for capacity he doesn't provide
- Anyone can today register capacity, even if they don't actually provide any capacity to the grid.

This shows that today, the validity of the capacity registered is only based on the trust that farmer will behave well and nobody will try to fake capacity.

We need to come up with some algorithm that provide some real proof of capacity. We also need to be able to limit the right to register capacity to only people that are verified and accepted as farmer.

We also need to come up with a way to register a farm and authorizing certain address to manage this farm (as in add themselves to that farm, or be linked to it since creation, even though I'm fairly sure we also might need to add more addresses after creation).


## Possible solutions

We need to find a solution to measure the capacity of 4 units:

- CPU
- Memory
- Hard drive disk storage
- Solid state drive storage

### CPU and Memory
For CPU and memory, we could use some kind of computational puzzle that are CPU or memory bound to ensure that the capacity a node reports can indeed be verified by solving one or more of these puzzles.  
For this to work the difficulty of a puzzle need to be tunable and link with a fix minimum amount of CPU power or amount of memory.

On problem with this approach is that once the node will be used to run workloads, both the CPU and memory will be used and not available to solve puzzle anymore. So we don't have a solution that allow us to re-check the capacity of a node over time while it is being used.

### Disk storage
For disk storage it also exists some kind of proof of capacity algorithm. Something like [burst](https://www.burst-coin.org/proof-of-capacity) is using:

> Burst mining style solves all these problems (and is also ASIC-resistant) by allowing HDD mining – miners secure the network with their disk space. It can be seen as a “condensed Proof-of-Work”: you compute once (a process called plotting) and cache the results of your work on hard disk space. Then mining only requires to read through your cache – your HDD is idle most of the time and reads through the plot files only for a few seconds for each block.

This also doesn't allow to do a capacity check while the node is being used cause it required to store the data on disk for a long period, which we can't do in our case.