# Jarvis [![CD](https://github.com/LucasGois1/jarvis/actions/workflows/cd.yaml/badge.svg?branch=master)](https://github.com/LucasGois1/jarvis/actions/workflows/cd.yaml)

A ChatGPT based chat system

### Features
 * Sessions
 * Historic
 * Token control
 * Rest and gRPC server with stream

### Deployment

Countinuous deployment with GithubActions and ArgoCD cluster. Ever a new feature is merged on branch master, an Action commits the new version of the service and ArgoCD is notified to sync the service and re-deploy pods
