# PRD: Kubernetes Cleanup Sidecar

## Problem Statement

During development and testing cycles, accumulated data in RabbitMQ queues and etcd database can cause:
- Stale chess game states affecting new deployments
- Queue message buildup leading to processing delays
- Development environment inconsistencies
- Manual cleanup overhead for developers

## Objective

Implement an optional Kubernetes init container (sidecar) that automatically cleans queue topics and database data during deployment, providing a fresh state for services.

## Requirements

### Functional Requirements

**FR1: RabbitMQ Queue Cleanup**
- Delete messages from specified topics (`monorepo`, `monorepo-dlq`)
- Purge dead letter queues
- Handle connection failures gracefully
- Log cleanup operations for debugging

**FR2: etcd Database Cleanup** 
- Delete all keys with configurable prefix (`monorepo`)
- Preserve system/infrastructure keys
- Handle etcd connection failures
- Support etcd cluster configurations

**FR3: Configuration Management**
- Enable/disable cleanup via Helm values
- Configure which resources to clean (queue, db, or both)
- Environment-specific configuration (dev/staging/production)
- Backward compatibility with existing deployments

**FR4: Deployment Integration**
- Run as init container before main application containers
- Block main containers until cleanup completes
- Fail deployment if cleanup encounters critical errors
- Support both producer and consumer deployments

### Non-Functional Requirements

**NFR1: Performance**
- Complete cleanup within 30 seconds for typical dev environments
- Minimal resource usage (128Mi memory, 100m CPU)
- Parallel cleanup operations where possible

**NFR2: Reliability**
- Retry failed operations with exponential backoff
- Graceful handling of missing resources
- Clear error reporting and exit codes

**NFR3: Security**
- Use service account authentication for K8s resources
- Minimal required permissions for cleanup operations
- No sensitive data logging

**NFR4: Observability**
- Structured logging for cleanup operations
- Integration with existing monitoring
- Clear success/failure indicators

## Technical Design

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│ Kubernetes Pod                                                  │
│ ┌─────────────────────┐ ┌─────────────────────────────────────┐ │
│ │ Init Container      │ │ Main Container                      │ │
│ │ (cleanup-sidecar)   │ │ (producer/consumer)                 │ │
│ │                     │ │                                     │ │
│ │ 1. Connect to       │ │ Waits for init container            │ │
│ │    RabbitMQ         │ │ completion before starting          │ │
│ │ 2. Purge queues     │ │                                     │ │
│ │ 3. Connect to etcd  │ │                                     │ │
│ │ 4. Delete keys      │ │                                     │ │
│ │ 5. Exit success     │ │                                     │ │
│ └─────────────────────┘ └─────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

### Implementation Components

1. **Cleanup Container**: Go binary with RabbitMQ and etcd clients
2. **Helm Configuration**: Values and templates for optional deployment
3. **Docker Image**: Lightweight container with cleanup binary
4. **Documentation**: Configuration and usage guidelines

### Configuration Schema

```yaml
cleanup:
  enabled: false           # Enable cleanup sidecar
  targets:
    queue: true           # Clean RabbitMQ queues
    database: true        # Clean etcd keys
  resources:
    requests:
      memory: "64Mi"
      cpu: "50m"
    limits:
      memory: "128Mi" 
      cpu: "100m"
  retry:
    attempts: 3
    backoff: "5s"
```

## Success Criteria

1. **Functional Success**
   - Cleanup sidecar successfully purges RabbitMQ queues
   - etcd keys with specified prefix are deleted
   - Main containers start only after cleanup completion
   - Configuration changes take effect without deployment issues

2. **Operational Success**
   - Development environment refresh time reduced by 80%
   - Zero manual cleanup interventions required
   - No production incidents due to accidental cleanup
   - Clear documentation and usage patterns established

3. **Technical Success**
   - Cleanup completes within 30-second timeout
   - Memory usage stays under 128Mi limit
   - Error conditions are properly handled and logged
   - Integration tests pass in CI/CD pipeline

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Accidental production cleanup | Critical | Default disabled, environment-specific configs |
| Cleanup failure blocks deployment | High | Timeout and retry mechanisms, bypass option |
| Resource connection issues | Medium | Retry logic, graceful degradation |
| Breaking existing deployments | Medium | Backward compatibility, feature flags |

## Timeline

- **Week 1**: Implementation and unit tests
- **Week 2**: Integration with Helm charts and local testing  
- **Week 3**: Documentation and CI/CD integration
- **Week 4**: Production readiness review and release

## Acceptance Criteria

- [ ] Cleanup sidecar container implemented and tested
- [ ] Helm chart updated with cleanup configuration options
- [ ] Local development environment supports cleanup feature
- [ ] Documentation covers configuration and troubleshooting
- [ ] CI/CD pipeline includes cleanup sidecar validation
- [ ] Feature can be safely enabled/disabled per environment