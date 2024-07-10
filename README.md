# cornelius
Utility for syncing S3 compatible storage to Arweave. 


```sh
docker run  -v $(pwd)/test:/etc/cornelius -it --rm cornelius:latest -c /etc/cornelius/config.yaml -x ardrive ---debug=true -l text
```

### TODO

- [ ] Compile metrics 
- [ ] Custom gateway
- [ ] IAM auth
- [ ] Graceful termination
- [ ] Handle redundant pipelines (avoid race condition on new files)
- [ ] Remove dependency on ardrive cli
- [ ] Bulk uploads
- [ ] Support IPFS bridge tags