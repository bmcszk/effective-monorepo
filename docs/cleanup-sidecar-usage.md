# Cleanup Sidecar Usage Guide

## Overview

The cleanup sidecar is a Kubernetes init container that optionally cleans RabbitMQ queues and etcd database before starting main application containers. This ensures fresh state for development and testing environments.

## Configuration

### Enabling Cleanup

Enable the cleanup sidecar in your Helm values:

```yaml
cleanup:
  enabled: true  # Enable cleanup sidecar (default: false)
  targets:
    queue: true     # Clean RabbitMQ queues (default: true)
    database: true  # Clean etcd keys (default: true)
```

### Advanced Configuration

```yaml
cleanup:
  enabled: true
  targets:
    queue: true
    database: true
  resources:
    requests:
      memory: "64Mi"
      cpu: "50m"
    limits:
      memory: "128Mi"
      cpu: "100m"
  retry:
    attempts: 3      # Number of retry attempts (default: 3)
    backoff: "5s"    # Backoff duration between retries (default: 5s)
  timeout: "30s"     # Overall cleanup timeout (default: 30s)
  image:
    repository: effective-monorepo/cleanup-sidecar
    tag: build
```

## Environment-Specific Usage

### Development Environment

```yaml
# values-dev.yaml
cleanup:
  enabled: true
  targets:
    queue: true
    database: true
```

### Staging Environment

```yaml
# values-staging.yaml
cleanup:
  enabled: true
  targets:
    queue: true
    database: false  # Preserve staging data
```

### Production Environment

```yaml
# values-prod.yaml
cleanup:
  enabled: false  # Never enable in production
```

## Deployment Commands

### Deploy with Cleanup Enabled

```bash
# Deploy with cleanup for development
helm upgrade --install monorepo ./helm --values values-dev.yaml

# Deploy specific service with cleanup
helm upgrade --install monorepo ./helm --set cleanup.enabled=true
```

### Deploy without Cleanup

```bash
# Standard deployment (cleanup disabled by default)
helm upgrade --install monorepo ./helm

# Explicitly disable cleanup
helm upgrade --install monorepo ./helm --set cleanup.enabled=false
```

## What Gets Cleaned

### RabbitMQ Cleanup (`targets.queue: true`)
- Purges all messages from `monorepo` topic queue
- Purges all messages from `monorepo-dlq` dead letter queue  
- Maintains queue declarations for service connectivity

### etcd Cleanup (`targets.database: true`)
- Deletes all keys with prefix `monorepo`
- Preserves system and infrastructure keys
- Maintains etcd cluster health

## Behavior

### Init Container Sequence
1. **Cleanup runs first** - Before any main containers start
2. **Main containers wait** - Until cleanup completes successfully
3. **Deployment fails** - If cleanup encounters critical errors

### Logging
- Structured JSON logs for monitoring integration
- Clear success/failure indicators
- Resource-specific cleanup status

### Error Handling
- **Retry Logic**: Configurable attempts with exponential backoff
- **Timeout Protection**: Prevents hanging deployments
- **Graceful Degradation**: Missing resources don't cause failures

## Troubleshooting

### Cleanup Taking Too Long
```bash
# Check init container logs
kubectl logs deployment/consumer -c cleanup-sidecar

# Increase timeout
helm upgrade monorepo ./helm --set cleanup.timeout=60s
```

### Connection Failures
```bash
# Check network connectivity
kubectl exec deployment/consumer -c cleanup-sidecar -- ping rabbitmq.infra
kubectl exec deployment/consumer -c cleanup-sidecar -- ping etcd.infra

# Check service endpoints
kubectl get endpoints rabbitmq etcd -n infra
```

### Resource Not Found Errors
These are typically safe to ignore - the cleanup sidecar handles missing resources gracefully:
```
INFO queue cleaned queue=monorepo messages_purged=0 total_messages=0
INFO no keys found to delete prefix=monorepo
```

### Disabling Specific Targets
```bash
# Only clean queues, skip database
helm upgrade monorepo ./helm --set cleanup.enabled=true --set cleanup.targets.database=false

# Only clean database, skip queues  
helm upgrade monorepo ./helm --set cleanup.enabled=true --set cleanup.targets.queue=false
```

## Best Practices

1. **Environment Safety**: Never enable cleanup in production
2. **Selective Cleanup**: Use different targets for different environments
3. **Resource Limits**: Monitor memory/CPU usage and adjust limits
4. **Timeout Tuning**: Set appropriate timeouts for your environment size
5. **Testing**: Verify cleanup behavior in development before staging

## Integration with CI/CD

```yaml
# Example GitLab CI pipeline
deploy-dev:
  script:
    - helm upgrade --install monorepo ./helm 
        --values values-dev.yaml 
        --set cleanup.enabled=true
  environment:
    name: development

deploy-prod:
  script:
    - helm upgrade --install monorepo ./helm 
        --values values-prod.yaml 
        --set cleanup.enabled=false
  environment:
    name: production
```