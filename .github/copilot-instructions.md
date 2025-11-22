# ActiveMQ Artemis Operator - AI Assistant Configuration

## Project Overview

This is a Kubernetes operator for Apache ActiveMQ Artemis, written in Go using the controller-runtime framework.

## Core Principle

**Documentation is FOR AI assistants, not BY AI assistants** — unless specified directly by the user.

- ✓ USE the documentation to answer questions and understand the codebase
- ✗ NEVER create new .md files, summaries, or implementation notes
- ✗ NEVER add to AI_documentation/ directory
- ✓ Exception: Only create documentation when explicitly requested by the user

## Communication Style (MANDATORY)

**CRITICAL: Professional, straightforward communication ONLY**

- ✓ USE plain, professional language
- ✓ USE simple bullet points or numbered lists when needed
- ✓ USE standard terminal commands without decoration
- ✗ NEVER create ASCII art boxes or decorative borders
- ✗ NEVER use excessive emojis (max 1-2 if truly necessary)
- ✗ NEVER format output with cat/echo heredocs for "pretty" summaries
- ✗ NEVER create "fancy" formatted reports with boxes, lines, or decoration
- ✗ NEVER use echo or other command-line tools to communicate thoughts, explanations, or instructions to the user

**CRITICAL: Terminal commands are for SYSTEM OPERATIONS ONLY, not for communicating with the user**

- ✓ USE terminal commands for: running tests, building code, starting services, file operations
- ✗ NEVER use terminal commands for: summaries, status updates, explanations, or any communication
- ✗ NEVER use `cat`, `echo`, `printf`, etc. to display information to the user
- ✓ Output all communication directly in your response text instead

**BAD - Never do this:**
```bash
cat << 'EOF'
╔══════════════════════════════════════╗
║        TASK COMPLETE ✅              ║
╚══════════════════════════════════════╝
1️⃣  FEATURE ADDED
   ✅ This
   ✅ That
EOF
```

**BAD - Also never do this:**
```bash
echo "Feature implementation complete"
echo "Changes made:"
echo "- Implemented feature X"
```

**GOOD - Do this:**
```
Feature implementation complete.

Changes made:
- Implemented feature X in file.go
- Added tests in file_test.go
- Updated documentation

Test results: All passing
```

## Documentation

AI documentation is in `AI_documentation/`:

### Key Files

1. **AI_KNOWLEDGE_INDEX.yaml** - Quick concept lookup with code locations
   - Core concepts with code entry points
   - Common questions with guidance
   - Use as starting point for navigation

2. **contribution_guide.md** - Development workflow, TDD, and testing
   - How to add features and API fields
   - Testing and debugging guide
   - E2E test setup instructions

## When Writing Code

Follow the `AI_documentation/contribution_guide.md`

**Discover Patterns**:
- Use `codebase_search` to find similar implementations
- Check existing tests for patterns
- Follow TDD approach from contribution_guide.md

**Follow Patterns**:
- Controller pattern with reconciliation loop
- StatefulSet-based broker management
- Validation chain for CRs
- Status condition management

**Test Coverage**:
- Use `codebase_search` to find similar test patterns
- Follow TDD approach
- Every feature must have tests

## Code Guidelines

1. **Naming**: Follow naming conventions (search codebase for patterns)
   - Resources: `{cr-name}-{resource-type}-{ordinal}`
   - Functions: Clear, descriptive names
   - Constants: ALL_CAPS with descriptive names

2. **Structure**: Follow controller pattern
   - Reconcile functions in `controllers/`
   - Resource generation in `pkg/resources/`
   - Utilities in `pkg/utils/`

3. **Validation**: Use validation chain pattern
   - Search for validate functions in controllers/activemqartemis_reconciler.go
   - Chain validators together
   - Return early on errors

4. **Status**: Update status conditions
   - Valid, Deployed, Ready, ConfigApplied
   - See api/v1beta1/activemqartemis_types.go for status definitions

5. When producing code:
- ✓ always create tests for the functionalities you're adding
- ✓ run the tests after implementing a functionality
- ✓ always run E2E tests related to your addition, how to do it is in the `contribution_guide.md`
- ✓ if required start minikube to run the E2E tests
- ✗ DON'T assume that your code is functional until tests are executed and green

## Test-Driven Development (MANDATORY)

**CRITICAL: Follow Outside-In TDD - E2E First, Unit Tests During Implementation**

**Outside-In TDD Workflow (MUST be followed for ALL code changes):**

1. **Write E2E Test FIRST** → Define acceptance criteria from user perspective
   - Write E2E test that describes the complete feature behavior
   - Test should interact with the operator as a real user would
   - Test will FAIL initially (red phase at acceptance level)
   - **PAUSE HERE** → Present E2E test to human for review

2. **Human Review Gate** → Wait for explicit approval before implementation
   - Human reviews feature specification and E2E test approach
   - Human may request changes to the E2E test
   - Only proceed to implementation after approval

3. **Implement Layer by Layer** → Write unit tests DURING implementation
   - For each component/function you implement:
     - Write unit test first (red phase at unit level)
     - Implement minimal code to pass the unit test (green phase)
     - Refactor while keeping tests green
   - Work from outside (controller) to inside (utilities)
   - Follow existing patterns discovered via codebase_search
   - Unit tests guide the design of individual components

4. **Verify E2E Test Passes** → Feature complete when E2E test is green
   - Run E2E test to verify end-to-end behavior
   - All unit tests should already be passing
   - E2E test validates that all components work together

5. **Final Refactor** → Clean up implementation while keeping all tests green
   - Improve code quality across all layers
   - Ensure both E2E and unit tests remain green

**Test Execution (happens throughout workflow):**

- **Setup Environment** → AUTOMATICALLY (for E2E tests):
  - Check if minikube/cluster is running
  - Start minikube if needed (don't ask permission)
  - Verify cert-manager installed and ready

- **Run Tests** → Execute tests AUTOMATICALLY:
  - E2E tests at start and end: `USE_EXISTING_CLUSTER=true go test -v ./controllers -ginkgo.focus="<pattern>" -timeout 10m`
  - Unit tests during implementation: `go test ./controllers -run <TestName>`

- **Fix Issues** → If tests fail, fix and re-run

- **Complete** → Only mark done when ALL tests pass (E2E + all unit tests)

## Unit Test Quality Standards (MANDATORY)

**Write Meaningful Tests - Not Just Coverage**

**CRITICAL: Avoid Over-Mocking**

When writing unit tests, ensure they actually test real behavior:

1. **Meaningful Testing**:
   - ✓ Test real logic and business rules
   - ✓ Test actual code paths and error handling
   - ✓ Use real structs and minimal mocking
   - ✗ Don't create tests that only verify mock interactions

2. **When Mocking is Acceptable**:
   - External dependencies (Kubernetes API, databases, network calls)
   - Time-dependent operations
   - File system operations
   - Random number generation

3. **When to Request Human Feedback**:
   - **PAUSE and ASK** if you need to mock >3 interfaces/dependencies
   - **PAUSE and ASK** if mocks make up >50% of test code
   - **PAUSE and ASK** if you're testing trivial getters/setters only
   - **PAUSE and ASK** if the test doesn't verify actual business logic

4. **Red Flags** (request human review):
   - Tests that only verify mock method calls without assertions on behavior
   - Tests that require complex mock setup with no real logic testing
   - Tests that would pass regardless of implementation correctness
   - Integration tests disguised as unit tests

**Example of GOOD unit test:**
```go
// Tests real validation logic without mocking
func TestValidateBrokerConfig(t *testing.T) {
    config := &BrokerConfig{Size: -1}
    err := ValidateBrokerConfig(config)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "size must be positive")
}
```

**Example of BAD unit test (request human feedback):**
```go
// Only tests mock interactions, not real behavior
func TestProcessBroker(t *testing.T) {
    mockClient := &MockClient{}
    mockValidator := &MockValidator{}
    mockLogger := &MockLogger{}
    // ... 50 lines of mock setup ...
    mockClient.EXPECT().Get(gomock.Any()).Times(1)
    // No actual assertion on behavior or results
}
```

**Definition of Done:**
- ✅ E2E test written FIRST (defines acceptance criteria)
- ✅ E2E test reviewed and approved by human
- ✅ Unit tests written DURING implementation (for each component)
- ✅ Unit tests are meaningful (not just mock verification)
- ✅ Code compiles with no errors
- ✅ All unit tests passing
- ✅ E2E test passing (validates end-to-end behavior)
- ✅ Code refactored while maintaining green tests
- ✅ Documentation updated
- ✅ No trailing spaces in the generated code
- ❌ Code is NOT done if E2E test wasn't written first
- ❌ Code is NOT done if E2E test wasn't reviewed by human
- ❌ Code is NOT done if unit tests were written after implementation
- ❌ Code is NOT done if tests haven't run
- ❌ Code is NOT done if any tests are failing

**Infrastructure Management:**
- Automatically check cluster status before E2E tests
- Use minikube profile "cursor" to avoid interfering with existing clusters
- Automatically start minikube if needed (4GB RAM, 2 CPUs minimum)
- **CRITICAL**: Configure ingress with SSL passthrough BEFORE running restricted mode tests
- Refer to contribution_guide.md (lines 850-889) for complete setup details
- Only ask user if setup fails
- Clean up test resources after test completion

**Complete Minikube Setup (REQUIRED for E2E tests):**
```bash
# 1. Start minikube with dedicated profile and set as active
minikube start --profile cursor --memory=4096 --cpus=2
minikube profile cursor
kubectl config use-context cursor

# 2. Enable ingress addon (REQUIRED for restricted mode tests)
minikube addons enable ingress

# 3. Wait for ingress controller to be ready
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s

# 4. Enable SSL passthrough (CRITICAL for mTLS in restricted mode)
kubectl patch deployment ingress-nginx-controller -n ingress-nginx \
  --type='json' \
  -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value":"--enable-ssl-passthrough"}]'

# 5. Wait for controller to restart with new config
kubectl rollout status deployment/ingress-nginx-controller -n ingress-nginx

# 6. Verify cert-manager is ready (auto-installed by tests if missing)
kubectl wait --for=condition=ready pod \
  -l app.kubernetes.io/instance=cert-manager \
  -n cert-manager --timeout=120s 2>/dev/null || echo "cert-manager will be auto-installed by tests"
```

**Run E2E Tests:**
```bash
# Control plane/restricted mode tests
USE_EXISTING_CLUSTER=true go test -v ./controllers \
  -ginkgo.focus="minimal restricted" -timeout 10m

# Specific feature test
USE_EXISTING_CLUSTER=true go test -v ./controllers \
  -ginkgo.focus="<your-test-name>" -timeout 10m
```

**Cleanup After Testing:**
```bash
# Stop and clean up test environment
minikube stop --profile cursor
minikube delete --profile cursor
```

## Key Concepts (Quick Reference)

From AI_KNOWLEDGE_INDEX.yaml:

- **reconciliation_loop**: Core control loop (controllers/activemqartemis_controller.go)
- **statefulset_management**: Pod lifecycle (pkg/resources/statefulsets/)
- **validation_chain**: CR validation (controllers/activemqartemis_reconciler.go)
- **broker_properties**: Configuration mechanism (use codebase_search)
- **persistence**: PVC management and storage
- **metrics_implementation**: Prometheus integration (Jolokia/JMX Exporter)

## Finding Information

**Use the knowledge index for concept lookup:**
- Start with `AI_KNOWLEDGE_INDEX.yaml` for quick concept navigation
- Use `codebase_search` to find code by meaning and discover patterns
- Use `grep` to find exact symbols or strings in code

**Code references format:**
- Knowledge index uses `file::symbol` format: `controllers/activemqartemis_reconciler.go::ProcessStatefulSet`
- This points you to the right file and function/type name
- Use `codebase_search` or `grep` to locate the actual code

## Getting Context

For detailed context, reference these files:
- `AI_documentation/AI_KNOWLEDGE_INDEX.yaml` - Concept lookup and navigation
- `AI_documentation/contribution_guide.md` - Development workflow and testing

## Documentation Structure

Simplified structure optimized for AI:
- **Concepts**: Core operator concepts mapped to files and symbols
- **Questions**: Common questions mapped to documentation sections
- **Code locations**: Key code entry points organized by area
- No auto-generated content - all manually curated

## Link Formatting for IDE Compatibility

When citing documentation or code:
- ✓ Use FILE PATH ONLY for clickable links: `AI_documentation/contribution_guide.md`
- ✓ Use file::symbol for code: `controllers/activemqartemis_controller.go::Reconcile`
- ✗ DON'T use anchor syntax: `file.md#anchor` (won't open in IDE)

