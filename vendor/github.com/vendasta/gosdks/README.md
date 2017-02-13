# Vendasta Go Sdks
This repository includes all of our internal Golang sdks

The `pb` folder includes Golang code generated from our protos at https://github.com/vendasta/vendastaapis/
These can be rebuilt by running `inv build`, which will do the checkout all over again and regenerate all of the protos.

You can run tests by running `inv test`
