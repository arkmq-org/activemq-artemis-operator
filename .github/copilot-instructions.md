# GitHub Copilot - ActiveMQ Artemis Operator Context

## Project Overview

This is a Kubernetes operator for Apache ActiveMQ Artemis, written in Go using the controller-runtime framework.

## Documentation

Comprehensive AI-optimized documentation is in `AI_documentation/`:

### Key Files

1. **AI_KNOWLEDGE_INDEX.yaml** - Start here for any question
   - 120+ concepts with definitions and code locations
   - 55+ common questions mapped to exact answers
   - Test area mappings
   - Code entry points

2. **operator_architecture.md** (2,223 lines)
   - Complete technical architecture
   - 22 major subsystems
   - 481 GitHub permalinks to code

3. **tdd_index.md** (3,230 lines)
   - 396+ test scenarios
   - Complete test coverage catalog
   - TDD feature inventory

4. **operator_conventions.md**
   - Naming conventions
   - Default values
   - Magic behaviors (suffix-based)

## When Writing Code

**Check Conventions**:
- Naming patterns: `operator_conventions.md`
- Default values: `operator_defaults_reference.md`
- Architecture patterns: `operator_architecture.md`

**Follow Patterns**:
- Controller pattern with reconciliation loop
- StatefulSet-based broker management
- Validation chain for CRs
- Status condition management

**Test Coverage**:
- Check `tdd_index.md` for similar test patterns
- Follow TDD approach
- Test scenarios are well-documented

## Code Guidelines

1. **Naming**: Follow conventions in `operator_conventions.md`
   - Resources: `{cr-name}-{resource-type}-{ordinal}`
   - Functions: Clear, descriptive names
   - Constants: ALL_CAPS with descriptive names

2. **Structure**: Follow controller pattern
   - Reconcile functions in `controllers/`
   - Resource generation in `pkg/resources/`
   - Utilities in `pkg/utils/`

3. **Validation**: Use validation chain pattern
   - See `operator_architecture.md#10-validation-architecture`
   - Chain validators together
   - Return early on errors

4. **Status**: Update status conditions
   - Valid, Deployed, Ready, ConfigApplied
   - See `operator_architecture.md#22-error-handling-and-status-management`

## Key Concepts (Quick Reference)

From AI_KNOWLEDGE_INDEX.yaml:

- **reconciliation_loop**: Core control loop (controllers/activemqartemis_controller.go)
- **statefulset_management**: Pod lifecycle (pkg/resources/statefulsets/)
- **validation_chain**: CR validation (controllers/activemqartemis_reconciler.go)
- **broker_properties**: Configuration mechanism (see architecture docs)
- **persistence**: PVC management and storage
- **metrics_implementation**: Prometheus integration (Jolokia/JMX Exporter)

## Code Link Reference

All documentation uses GitHub permalinks to commit `4acadb95`:
- Links point to exact code as of that commit
- To see current: replace `/blob/4acadb95/` with `/blob/main/`
- Function names help locate code if moved

## Getting Context

For detailed context, reference these files in your prompts:
- `@AI_documentation/AI_KNOWLEDGE_INDEX.yaml` - Concept lookup
- `@AI_documentation/operator_architecture.md` - Architecture details
- `@AI_documentation/tdd_index.md` - Test examples

## Documentation Structure

The knowledge index provides semantic navigation:
- `concepts`: 120+ concept definitions
- `common_questions`: 55+ Q&A mappings  
- `topics`: 10 functional areas
- `test_areas`: 40+ test mappings
- `code_locations`: Quick code navigation

## Link Formatting for IDE Compatibility

When citing documentation or code:
- Use FILE PATH ONLY for clickable links: `AI_documentation/operator_architecture.md`
- Mention sections separately: "See section 1: High-Level Architecture"
- Use line numbers for code: `controllers/activemqartemis_controller.go:156-180`
- DON'T use anchor syntax: `file.md#anchor` (won't open in IDE)

