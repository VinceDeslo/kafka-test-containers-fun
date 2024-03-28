# kafka-test-containers-fun

Tiny little program to play around with Kafka test containers I may want to use for some local integration testing.

```
go run cmd/kafka/main.go
```

### Roadmap
- [x] Spin up a test container
- [x] Fetch the broker URL
- [x] Build an admin client
- [x] Create a topic
- [ ] Set up a schema registry entry
- [ ] Produce a message
- [ ] Consume a message

### Observations

- Startup may be a bit slow if images aren't already fetched, not great for test durations and local DevEx.
- Config won't necessarily align with an existing staging/production cluster.
- Possibly fragile testing since this is only an abstraction of a live environment.
- A bit of a fun sandbox for experimentation with new event schemas, etc.
