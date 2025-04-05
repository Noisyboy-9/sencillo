# Sencillo: a simple scheduling framework
Sencillo is a scheduling framework designed to streamline scheduler development in Kubernetes environments. With Sencillo, developers can focus on developing decision-making algorithms. The framework handles connecting to the API server, managing pod creation events, and managing cluster resources. 

## Features
- 🚀 **Support for Various Clusters Types**: The framework supports only-cloud, only-edge and cloud-assisted edge 
  clusters.
- 🔀 **Multi-Algorithm Support**: Supports multiple scheduling algorithms in a single codebase with a dynamic 
  algorithm selection mechanism configurable through Kubernetes configmaps.
- 📊 **Accurate Cluster State Tracking**: Continuously monitors and maintains an up-to-date view of CPU and memory availability across all nodes.
- <img src="https://raw.githubusercontent.com/kubernetes/kubernetes/master/logo/logo.png" alt="Kubernetes" width="20"/> **Kubernetes Integration**: Seamlessly integrates with Kubernetes clusters for smooth deployment and management.
- 🧩 **Modular Architecture**: Easily extend and customize scheduling logic thanks to a clean and modular codebase.

## Repository Structure

- `cmd/`: Contains the main entry points for the different scheduler implementations.
- `configs/`: Configuration files for the schedulers.
- `deployments/`: Kubernetes deployment manifests for deploying the schedulers.
- `internal/`: Internal packages and modules used across the project.
- `Dockerfile`: Dockerfile for building the scheduler container images.
- `Makefile`: Makefile for automating build and deployment tasks.
- `docker-compose.yaml`: Docker Compose configuration for local development and testing.
- `go.mod` and `go.sum`: Go module files for dependency management.
- `main.go`: The main entry point for the application.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.16 or higher
- [Docker](https://www.docker.com/get-started)
- [Kubernetes](https://kubernetes.io/docs/setup/) cluster

### Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/Noisyboy-9/sencillo.git
   cd sencillo
   ```

2. **Build the Docker Image**:

   ```bash
   make docker-build
   ```

3. **Deploy to Kubernetes**:

   Ensure your `kubectl` is configured to communicate with your Kubernetes cluster, then run:

   ```bash
   kubectl apply -f deployments/
   ```

## Usage

After deployment, the schedulers will be running within your Kubernetes cluster, ready to handle scheduling tasks based on their implemented algorithms. You can monitor their behavior and performance using Kubernetes tools and logs.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes. Ensure that your code adheres to the project's coding standards and includes appropriate tests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

