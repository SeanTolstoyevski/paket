# Problems arising when using the **paket**.

You can use issues for discussions on how to fix these issues.


* [ ] Support for GCM and other algorithms should be added.

* [ ] replace the hash values in the table with []byte
 - This saved us from **stringify jobs**. But it can complicate the cmd tool.

## Completeds

* [x] Panic occurs when several file requests are made at the same time. **With goroutines**.
