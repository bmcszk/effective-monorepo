# Effective Monorepo
======================

A monorepo for building scalable and maintainable microservices.

## Overview
-----------

This monorepo contains multiple services and tools for building a scalable and maintainable microservices architecture. The project is designed to demonstrate best practices for building a monorepo, including organization, testing, and deployment.

## Services
------------

* **Producer**: A service responsible for producing messages to a message queue.
* **Consumer**: A service responsible for consuming messages from a message queue.
* **Test**: A service responsible for running end-to-end tests for the monorepo.

## Tools
--------

* **Tilt**: A tool for building and deploying the monorepo to a local Kubernetes cluster.
* **Kind**: A tool for creating a local Kubernetes cluster for testing and development.

## Getting Started
---------------

1. Clone the repository: `git clone [https://github.com/bmcszk/effective-monorepo.git`](https://github.com/bmcszk/effective-monorepo.git`)
2. Install dependencies: `make install-all-deps`
3. Start the local Kubernetes cluster: `make kind-up`
4. Build and deploy the services: `make tilt-up`
5. Run end-to-end tests: `make test`

## Contributing
------------

Contributions are welcome! Please submit a pull request with a clear description of the changes and any relevant tests.

## License
-------

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Acknowledgments
---------------

* [Tilt](https://tilt.dev/)
* [Kind](https://kind.sigs.k8s.io/)
* [Go](https://golang.org/)
