# Developing Roadmap / planning for Paket

This file it contains explanations for making better use, adding new features, and major problems.

Any ideas, pull requests, suggestions or something like that are appreciated.

As a side note, I don't have any education in cryptography. And this package was not created for everything to be perfect. I had a similar need for my various licensed projects and I code for it.

I will complete these items when I have time and can figure out how to design it.


* [ ] Adding complex, interdependent hash generation method for PBDFK2 (**See footnote 1** for more information).
* [x] Support for GCM and other algorithms should be added.
	- CFB - OK
	- CTR - OK
	- GCM - OK
	- OFB - OK
	- **CBC - not completed**. I'm still thinking about how to implement it.

## Completeds

* [x] replace the hash values in the table with []byte
 - This saved us from **stringify jobs**. But it can complicate the cmd tool.
* [x] Panic occurs when several file requests are made at the same time. **With goroutines**.

### Footnote 1

The basic idea here is to combine several different hash functions with different iteration counts.
Instead of using a single hash function, creating a complex but slow hashing system that is interdependent.  
Example (pseudo):

1. MD5, iteration number: 55
2. SHA256, iteration number: 2,000
3. MD5, iteration number: 5,000
4. SHA512, iteration number: 12,000

The value returned from each hash function is passed to the other hash function when the iteration ends.

The problem is that this will also require creating a ready-made interface.  
According to my plan, these declarations can be specified in a JSON file. The Paket reads the JSON file, creates the interface.

```cmd

paket other_flags_and_keys itertemplate=my_iter.json

```
