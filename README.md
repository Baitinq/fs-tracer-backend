# fs-tracer-backend

This repository contains the backend services for the fs-tracer project, which processes and stores file system modification data collected by the eBPF agent. 

## Overview

The backend is composed of several microservices that work together to handle incoming file system events, store them in a database, and provide an API for accessing and managing this data. The services are containerized and deployed using Kubernetes for scalability and reliability.

Related:
- https://github.com/baitinq/fs-tracer
- https://github.com/baitinq/fs-tracer-frontend
