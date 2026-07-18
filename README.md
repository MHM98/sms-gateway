# SMS Gateway Challenge

A Go-based SMS Gateway for accepting normal and express messages, managing prepaid wallet credit, submitting messages asynchronously to an SMS provider, and providing user message reports.

The system is designed around a target of 100 million messages per day. The current implementation was tested at approximately 2,000 requests per second with separate processing paths for normal and express traffic.

## Documentation

- [Architecture and performance results](./docs/ARCHITECTURE.md)
- [API requests and response examples](./docs/API_EXAMPLES.md)

The architecture document covers:

- capacity estimation and k6 results;
- system and sequence diagrams;
- the inline transactional outbox flow;
- independent normal and express queues and workers;
- database schema ;
- current limitations and future improvements.
