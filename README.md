# ActiveMQ Artemis Operator

This project is a [Kubernetes](https://kubernetes.io/) [operator](https://coreos.com/blog/introducing-operators.html)
to manage the [Apache ActiveMQ Artemis](https://activemq.apache.org/artemis/) message broker.

> [!NOTE]
> **🤖 For AI Assistants**: This project has comprehensive documentation
> optimized for AI in [`AI_documentation/`](AI_documentation/). Start with
> [`AI_KNOWLEDGE_INDEX.yaml`](AI_documentation/AI_KNOWLEDGE_INDEX.yaml) for
> semantic navigation of 120+ concepts, 55+ common questions, and 1,352+ code
> links. See [CONTEXT_GUIDE.md](AI_documentation/CONTEXT_GUIDE.md) for usage.

## Status ##

The current api version of all main CRDs managed by the operator is **v1beta1**.

## Quickstart

The [quickstart.md](docs/getting-started/quick-start.md) provides simple steps to quickly get the operator up and running
as well as deploy/managing broker deployments.

## Building

The [building.md](docs/help/building.md) describes how to build operator and how to test your changes

## OLM integration

The [bundle.md](docs/help/bundle.md) contains instructions for how to build operator bundle images and integrate it into [Operator Liftcycle Manager](https://olm.operatorframework.io/) framework.

## Debugging operator inside a container

Install delve in the `builder` container, i.e. `RUN go install github.com/go-delve/delve/cmd/dlv@latest`
Disable build optimization, i.e. `go build -gcflags="all=-N -l"`
Copy delve to the `base-env` container, i.e. `COPY --from=builder /go/bin/dlv /bin`
Execute operator using delve, i.e. `/bin/dlv exec --listen=0.0.0.0:40000 --headless=true --api-version=2 --accept-multiclient ${OPERATOR} $@`

## AI-Optimized Documentation

Comprehensive developer documentation optimized for AI assistants is available in [`AI_documentation/`](AI_documentation/):

- **[AI Assistant Guide](AI_documentation/AI_ASSISTANT_GUIDE.md)** - How to use with Cursor, Copilot, ChatGPT
- **[Knowledge Index](AI_documentation/AI_KNOWLEDGE_INDEX.yaml)** - Master concept index (120+ concepts)
- **[Architecture](AI_documentation/operator_architecture.md)** - Complete technical architecture (2,200+ lines)
- **[Test Catalog](AI_documentation/tdd_index.md)** - 396+ test scenarios with code links
- **[Conventions](AI_documentation/operator_conventions.md)** - Naming patterns and defaults
- **[Contribution Guide](AI_documentation/contribution_guide.md)** - Development workflow and debugging

These docs are comprehensive and detailed (7,000+ lines, 1,352+ code links), designed for AI-assisted development rather than manual reading.

**For AI assistants**: Start with `AI_KNOWLEDGE_INDEX.yaml` for semantic navigation.

**For humans**: Use AI tools (Cursor, Copilot, ChatGPT) to query the documentation.

