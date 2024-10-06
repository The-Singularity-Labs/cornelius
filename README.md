# cornelius
Utility for syncing S3 compatible storage to Arweave. 


```sh
docker run  -v $(pwd)/test:/etc/cornelius -it --rm cornelius:latest -c /etc/cornelius/config.yaml -x ardrive ---debug=true -l text
```

### TODO

- [ ] switch cli to basic sdk
- [ ] Support IPFS bridge tags
- [ ] Graceful termination
- [ ] Custom gateway
- [ ] Handle redundant pipelines (avoid race condition on new files)
- [ ] Remove dependency on ardrive cli
