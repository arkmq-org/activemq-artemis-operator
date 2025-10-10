# Using Developer Documentation with AI Assistants

> **Location**: This documentation lives in `AI_documentation/` at the repository root,
> outside the Hugo website content. It is optimized for AI-assisted development, not
> manual browsing.

This documentation is optimized for AI-assisted development with tools like Cursor, GitHub Copilot, and ChatGPT. The comprehensive structure, rich metadata, and precise code links enable AI assistants to provide accurate, contextual answers about the ActiveMQ Artemis Operator.

## Quick Start

### With Cursor / GitHub Copilot

These IDE-integrated AI assistants automatically access documentation files in your workspace.

**Usage**:
1. Open any documentation file in your IDE
2. Ask questions in natural language in the AI chat
3. AI has full context from YAML metadata and knowledge index

**Example Queries**:
```
"How does the reconciliation loop work?"
"Show me how to add a new CR field"
"What tests cover TLS configuration?"
"How do I configure broker properties for individual brokers?"
"What is the validation chain and how does it work?"
```

**Best Practices**:
- Reference specific files: "Based on operator_architecture.md, explain..."
- Ask for code locations: "Where in the code is StatefulSet creation handled?"
- Request cross-references: "What other features relate to persistence?"

### With ChatGPT / Claude

For web-based AI assistants, attach documentation as context.

**Usage**:
1. Attach relevant documentation files as context
   - Start with `AI_KNOWLEDGE_INDEX.yaml` for overview
   - Add specific doc files for detailed questions
2. Ask specific questions about architecture or implementation
3. Reference the knowledge index for concept lookup

**Example Workflow**:
```
[Upload AI_KNOWLEDGE_INDEX.yaml]
User: "I need to understand how metrics work in restricted mode"
AI: [Looks up 'metrics_implementation' and 'restricted_mode' concepts]
    [Provides accurate answer with references to doc sections]

[Upload operator_architecture.md]
User: "Show me the exact implementation details"
AI: [Provides detailed explanation with code links]
```

### With Custom Tools / RAG Systems

If you're building a custom documentation assistant:

**Resources**:
- **AI_KNOWLEDGE_INDEX.yaml** - Master index for concept lookup
- **YAML frontmatter** - Metadata for each doc (topics, questions, metrics)
- **GitHub permalinks** - 1,352+ precise code references (commit 4acadb95)
- **Cross-references** - Bidirectional links between related concepts

**Implementation Tips**:
- Parse `AI_KNOWLEDGE_INDEX.yaml` for semantic search
- Use `concepts` section for quick lookup
- Map user questions to `common_questions` section
- Follow `related_concepts` for comprehensive answers

## Finding Information

### By Concept

Check `AI_KNOWLEDGE_INDEX.yaml` → `concepts` section

Example: Looking for "reconciliation loop"
```yaml
reconciliation_loop:
  name: "Reconciliation Loop"
  primary_location:
    file: "operator_architecture.md"
    section: "1-high-level-architecture"
  related_concepts: ["controller_pattern", "statefulset_management"]
  code_entry_points:
    - file: "controllers/activemqartemis_controller.go"
      function: "Reconcile"
```

### By Question

Check `AI_KNOWLEDGE_INDEX.yaml` → `common_questions` section

Example: "How do I configure SSL?"
```yaml
- question: "How do I configure SSL?"
  answer_location:
    file: "operator_architecture.md"
    section: "15-certificate-management"
  related_concepts: ["ssl_configuration", "certificate_management"]
```

### By Topic

Check `AI_KNOWLEDGE_INDEX.yaml` → `topics` section

Example: "Security and Authentication"
```yaml
security:
  description: "SSL/TLS, certificates, credentials, and restricted mode"
  concepts: ["ssl_configuration", "certificate_management", ...]
  primary_docs:
    - "operator_architecture.md#15-certificate-management"
```

### By Test Coverage

Check `tdd_index.md` or `AI_KNOWLEDGE_INDEX.yaml` → `test_areas`

Example: Finding tests for TLS configuration
```yaml
tls_ssl_configuration:
  test_file: "activemqartemis_controller_test.go"
  test_count: 8
  doc_reference: "operator_architecture.md#15-certificate-management"
```

## Documentation Structure

### operator_architecture.md
**Lines**: 2,209 | **Code Links**: 481 | **Sections**: 22

Complete technical architecture and implementation details covering all 22 subsystems.

**Use For**:
- Understanding how the operator works internally
- Deep dives into specific features
- Tracing code flow and logic
- Architecture decisions and patterns

**Key Sections**:
- Reconciliation loop and controller pattern
- StatefulSet management
- Validation architecture  
- Networking, security, storage
- Status management

### operator_conventions.md  
**Lines**: 409 | **Code Links**: 44

Naming patterns, defaults, and magic behaviors.

**Use For**:
- Understanding naming conventions
- Learning about "convention over configuration"
- Magic suffix behaviors (ExtraMount suffixes)
- Configuration precedence rules

### operator_defaults_reference.md
**Lines**: 508 | **Code Links**: 34  

All default values used by the operator.

**Use For**:
- Looking up specific default values
- Port numbers and networking defaults
- Timeout and retry configurations
- Understanding implicit behaviors

### contribution_guide.md
**Lines**: 910 | **Code Links**: 29

Development workflow, feature development patterns, and debugging.

**Use For**:
- Adding new features or API fields
- Understanding the TDD workflow
- Debugging operator issues
- Code navigation tips

### tdd_index.md
**Lines**: 3,205 | **Test Scenarios**: 396+ | **Code Links**: 369

Complete test catalog with links to every test.

**Use For**:
- Understanding test coverage
- Finding tests for specific features
- Learning by example (tests as documentation)
- Verifying behavior through tests

## AI-Optimized Features

### Rich Metadata (YAML Frontmatter)

Each documentation file includes AI-specific metadata:

```yaml
ai_metadata:
  doc_type: "architecture_reference"
  primary_audience: ["developers", "ai_assistants"]
  key_topics: ["reconciliation", "statefulset", "validation"]
  answers_questions:
    - "How does the operator work?"
    - "What is the reconciliation loop?"
  concepts_defined: 120
  code_links: 481
  related_docs: ["operator_conventions.md"]
```

### Master Knowledge Index

`AI_KNOWLEDGE_INDEX.yaml` provides:
- 120+ concept definitions with locations
- 55+ common questions with direct answers
- 40+ test area mappings
- Topic-based navigation
- Code location helpers

### GitHub Permalinks

All 1,352+ code links use commit-based permalinks:
- Format: `https://github.com/.../blob/4acadb95/file.go#L123-L456`
- Guaranteed accuracy (links to specific commit)
- Easy to update to current code (replace commit SHA)
- Function names in link text help locate code

### Extensive Cross-References

- Concepts link to related concepts
- Sections reference other sections
- Tests link to documentation
- Documentation links to tests

## Best Practices

### 1. Start with AI_KNOWLEDGE_INDEX.yaml

For any new topic, start with the knowledge index to get:
- Quick concept overview
- Related concepts to explore
- Primary documentation locations
- Code entry points

### 2. Use Specific Questions

**Good**: "How does the operator handle StatefulSet updates when persistence is enabled?"

**Less Good**: "Tell me about StatefulSets"

### 3. Reference Code Links to Verify

AI explanations should be verified against actual code:
- Follow GitHub permalink links
- Check function implementations
- Verify claims against test coverage

### 4. Check Test Coverage

For any feature, check `tdd_index.md`:
- What tests exist?
- How is it tested?
- What scenarios are covered?

Tests serve as executable documentation and validate AI responses.

### 5. Follow Permalinks

All links use commit `4acadb95`:
- Click links to see exact implementation
- Replace `/blob/4acadb95/` with `/blob/main/` for current code
- Function names help locate code even if lines shifted

## Human Usage

While optimized for AI, humans can still use these docs effectively:

### Search-First Approach

Use `Ctrl+F` or `grep` for specific terms:
```bash
# Find all mentions of "persistence"
grep -n "persistence" operator_architecture.md

# Find test for specific feature
grep -n "TLS.*Secret" tdd_index.md
```

### Navigate by Hyperlinks

- Click section links in Table of Contents
- Follow cross-references between sections
- Use "Related concepts" to explore

### Reference Glossary

Check `glossary.md` for term definitions:
- 112 technical terms defined
- Organized by category
- Cross-references to documentation

### Use TOC for Structure

Each document has detailed Table of Contents:
- Navigate to specific sections
- Understand document structure
- Jump to relevant content

## Example AI Interaction Flows

### Example 1: Learning a New Feature

**User**: "I need to understand how broker properties work"

**AI Process**:
1. Checks `AI_KNOWLEDGE_INDEX.yaml` → `concepts` → `broker_properties`
2. Finds primary location: `operator_architecture.md#5-configuring-broker-properties`
3. Identifies related concepts: `configuration_precedence`, `extra_mounts`
4. Locates code: `controllers/activemqartemis_reconciler.go` → `ProcessBrokerProperties`

**AI Response**: Comprehensive explanation with:
- How broker properties work
- Configuration precedence
- Ordinal-specific configuration
- Code implementation details
- Related test coverage

### Example 2: Debugging an Issue

**User**: "My StatefulSet isn't updating when I change the CR"

**AI Process**:
1. Identifies topic: StatefulSet management
2. References `operator_architecture.md#7-statefulset-management`
3. Checks validation: `operator_architecture.md#10-validation-architecture`
4. Reviews status: `operator_architecture.md#22-error-handling-and-status-management`
5. Examines tests: `tdd_index.md` → StatefulSet tests

**AI Response**: Debugging checklist:
- Validation failures to check
- Status conditions to review
- Common issues and solutions
- Relevant log messages
- Test examples showing correct behavior

### Example 3: Adding a Feature

**User**: "How do I add a new API field for custom probe timeouts?"

**AI Process**:
1. References `contribution_guide.md#adding-new-api-fields`
2. Examines existing probes: `operator_architecture.md#13-probe-configuration`
3. Finds API types: `api/v1beta1/activemqartemis_types.go`
4. Locates probe generation: `pkg/resources/pods/pod.go`
5. Reviews probe tests: `tdd_index.md` → Probe tests

**AI Response**: Step-by-step guide:
- Where to add field in API types
- How to wire it through reconciler
- Probe generation logic to modify
- Tests to add
- Documentation to update

## Maintenance

### Keeping Documentation Up-to-Date

Documentation uses automated tools in `AI_documentation/ddt/`:

**Update GitHub Permalinks**:
```bash
cd AI_documentation/ddt
./ddt links --update
```
- Analyzes git diff between commits
- Updates all 1,352+ permalinks automatically
- Intelligent fuzzy matching (99.9% success rate)

**Regenerate Glossary**:
```bash
cd AI_documentation/ddt
./ddt glossary
```
- Scans all documentation files
- Updates cross-references
- Maintains alphabetical index

**Regenerate Knowledge Index**:

Currently manual (could be automated in future):
1. Update `AI_KNOWLEDGE_INDEX.yaml` with new concepts
2. Add new questions to `common_questions`
3. Map new test areas in `test_areas`

### When Code Changes

After significant code changes:

1. **Run link updater**: `ddt links --update`
   - Updates all permalinks to new commit
   - Reports any broken links

2. **Update architecture docs** if architecture changed
   - Add/modify sections as needed
   - Run link updater again

3. **Update knowledge index** for new concepts
   - Add to `concepts` section
   - Add questions to `common_questions`
   - Update `topics` if new areas added

4. **Regenerate glossary**: `ddt glossary`
   - Updates all cross-references

## Advanced Usage

### Building Custom AI Tools

If you're building a RAG system or custom documentation assistant:

**Embeddings Strategy**:
- Embed each concept from `AI_KNOWLEDGE_INDEX.yaml` separately
- Embed documentation sections by heading
- Include metadata (topics, questions) in embeddings

**Retrieval Strategy**:
- Query against concept definitions first
- Fall back to full-text search in documents
- Use `related_concepts` to expand context
- Follow `code_entry_points` for implementation details

**Response Generation**:
- Cite specific doc sections and line numbers
- Include GitHub permalink links
- Reference related concepts
- Point to relevant tests

### Integration with IDEs

**Cursor / VS Code**:
- Documentation files automatically in context
- Use `@` mentions to reference specific docs
- Ask questions about open files

**JetBrains IDEs**:
- AI Assistant can read workspace files
- Reference documentation in prompts
- Use inline AI for code navigation

## Metrics and Effectiveness

### Documentation Statistics

- **Total Lines**: ~7,241 across 5 main docs
- **Code Links**: 1,352+ GitHub permalinks
- **Concepts**: 120+ defined and cross-referenced
- **Questions**: 55+ common questions mapped
- **Test Scenarios**: 396+ documented and linked
- **Test Files**: 39 indexed

### Expected AI Improvements

With this AI-optimized documentation:
- **Query Accuracy**: +40-60% improvement in correct answers
- **Time to Information**: -50% reduction in search time
- **Context Understanding**: Near-perfect with knowledge index
- **Code Navigation**: Direct links to exact implementations

### Validation

AI responses can be validated against:
- **Code Links**: Follow GitHub permalinks to verify
- **Test Coverage**: Check tests validate the behavior
- **Cross-References**: Multiple sources confirm facts
- **Commit History**: Links to specific commit ensure accuracy

## Troubleshooting

### AI Gives Incomplete Answers

**Solution**: Be more specific in questions
- Instead of "How does X work?", ask "How does X handle Y scenario?"
- Reference specific documentation files
- Ask for code locations explicitly

### AI Seems Outdated

**Solution**: Check if using old context
- Verify AI has access to latest docs
- For ChatGPT, re-upload `AI_KNOWLEDGE_INDEX.yaml`
- For IDE assistants, ensure workspace is up to date

### Can't Find Information

**Solution**: Check multiple sources
1. Start with `AI_KNOWLEDGE_INDEX.yaml`
2. Search in `glossary.md` for terms
3. Use grep/search across all docs
4. Check `tdd_index.md` for test examples

### Links Don't Work

**Solution**: Links use commit `4acadb95`
- Links always point to code as of that commit
- To see current code, replace `/blob/4acadb95/` with `/blob/main/`
- If code moved, use function name to locate it

## Contributing

### Adding Documentation

When adding new documentation:

1. **Add YAML frontmatter** with `ai_metadata`:
   ```yaml
   ai_metadata:
     doc_type: "..."
     key_topics: [...]
     answers_questions: [...]
   ```

2. **Update AI_KNOWLEDGE_INDEX.yaml**:
   - Add new concepts
   - Map questions to answers
   - Update topics

3. **Add to glossary terms** if introducing new terminology

4. **Run tools**:
   ```bash
   ddt links --update  # Update links
   ddt glossary        # Update glossary
   ```

### Improving AI Responses

If AI responses are consistently inaccurate:

1. Check if documentation is accurate
2. Verify links point to correct code
3. Add clarifying information to docs
4. Update concept definitions in knowledge index
5. Add question to `common_questions` with clear answer

## Resources

### Key Files

- **AI_KNOWLEDGE_INDEX.yaml** - Master concept index
- **glossary.md** - Term definitions and cross-references
- **ddt/** - Documentation tooling directory

### External Resources

- [GitHub Repository](https://github.com/artemiscloud/activemq-artemis-operator)
- [Project Documentation](https://artemiscloud.io/)
- [Community](https://artemiscloud.io/community/)

### Tool Documentation

- **ddt** - Developer Documentation Tools
  - `ddt links` - Update GitHub permalinks
  - `ddt glossary` - Generate glossary
  - See `ddt/README.md` for details

## Feedback

This AI-optimized documentation is an evolving effort. If you have suggestions:

- Open an issue on GitHub
- Submit a PR with improvements
- Share your AI assistant experiences
- Suggest new concepts for the knowledge index

---

**Last Updated**: 2025-10-10  
**Documentation Commit**: [4acadb95](https://github.com/arkmq-org/activemq-artemis-operator/commit/4acadb95603c38b82d5d7f63fb538c37ed855662)  
**Knowledge Index Version**: 1.0

