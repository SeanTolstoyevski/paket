# Problems arising when using the **paket**.

You can use issues for discussions on how to fix these issues.

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
