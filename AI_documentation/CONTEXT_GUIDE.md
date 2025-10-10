# AI Assistant Context Guide

> **Purpose**: Quick-start guide for AI assistants working with this codebase

## Start Here

You are assisting with the **ActiveMQ Artemis Operator**, a Kubernetes operator written in Go.

## Essential Files (Load First)

### 1. AI_KNOWLEDGE_INDEX.yaml (This Directory)

**What**: Master index of 120+ concepts, 55+ questions, 40+ test areas  
**Use**: Start here for any question about the codebase  
**Size**: 1,475 lines, 56 KB

**Structure**:
```yaml
concepts:           # 120+ concept definitions with locations
  reconciliation_loop:
    definition: "..."
    primary_location: { file, section }
    code_entry_points: [...]
    
common_questions:   # 55+ questions → direct answers
  - question: "How does X work?"
    answer_location: { file, section }
    
topics:            # 10 functional areas
  architecture: { concepts, docs }
  
test_areas:        # 40+ test mappings
  tls_configuration: { file, count, docs }
```

### 2. operator_architecture.md (This Directory)

**What**: Complete technical architecture (2,223 lines)  
**Use**: Deep dives into how everything works  
**Contains**: 22 sections, 481 code links, all subsystems documented

### 3. tdd_index.md (This Directory)

**What**: Test catalog (3,230 lines, 396 scenarios)  
**Use**: Find test coverage, understand by example  
**Contains**: All test files, scenarios, code links

## Quick Concept Lookup

**Format**: Open `AI_KNOWLEDGE_INDEX.yaml`, find concept:

```yaml
concepts:
  {concept_name}:
    definition: "What it is"
    primary_location:
      file: "where to read more"
      section: "specific section"
    code_entry_points:
      - file: "where code lives"
        function: "function name"
    related_concepts: ["what else to check"]
```

**Example**: User asks "How does the reconciliation loop work?"

1. Check `AI_KNOWLEDGE_INDEX.yaml` → `concepts` → `reconciliation_loop`
2. See `primary_location`: `operator_architecture.md#1-high-level-architecture`
3. Find `code_entry_points`: `controllers/activemqartemis_controller.go` → `Reconcile()`
4. Check `related_concepts`: `controller_pattern`, `statefulset_management`
5. Provide comprehensive answer with links

## Question Routing

**Check `common_questions` section** for direct mappings:

- "How do I add a new field?" → `contribution_guide.md#adding-new-api-fields`
- "How does validation work?" → `operator_architecture.md#10-validation-architecture`
- "What are the naming conventions?" → `operator_conventions.md#resource-naming-conventions`

## Documentation Files

| File | Purpose | Size | Links |
|------|---------|------|-------|
| operator_architecture.md | Architecture | 2,223 lines | 481 |
| tdd_index.md | Test catalog | 3,230 lines | 369 |
| operator_conventions.md | Conventions | 421 lines | 44 |
| operator_defaults_reference.md | Defaults | 520 lines | 34 |
| contribution_guide.md | Development | 923 lines | 29 |
| glossary.md | Terms | 1,889 lines | 859+ |

**Total**: 7,241 lines, 1,352+ code links

## Code Links

All use GitHub permalinks to commit `4acadb95603c38b82d5d7f63fb538c37ed855662`:
- **Format**: `https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95/file.go#L123-L456`
- **Permanence**: Link always points to same code
- **Current code**: Replace `/blob/4acadb95/` with `/blob/main/`

## Usage Patterns

### For Architecture Questions
```
1. Check AI_KNOWLEDGE_INDEX.yaml → concepts
2. Read primary_location documentation
3. Examine code_entry_points
4. Review related_concepts
5. Provide comprehensive answer
```

### For Implementation Help
```
1. Check AI_KNOWLEDGE_INDEX.yaml → topics
2. Find relevant concepts
3. Check operator_conventions.md for patterns
4. Look at tdd_index.md for test examples
5. Provide guidance with code references
```

### For Debugging
```
1. Identify feature area
2. Check operator_architecture.md for how it works
3. Find validation in validation_architecture section
4. Check status management section
5. Point to relevant tests in tdd_index.md
```

## Key Statistics

- **Concepts Documented**: 120+
- **Common Questions Mapped**: 55+
- **Test Scenarios**: 396+
- **Code Links**: 1,352+
- **Documentation Lines**: 7,241
- **Test Files Indexed**: 39

## Project Context

- **Language**: Go 1.23+
- **Framework**: controller-runtime (Kubernetes)
- **Architecture**: Operator pattern with reconciliation loop
- **Testing**: Test-Driven Development (TDD)
- **Broker**: Apache ActiveMQ Artemis

## Tips for AI Assistants

✓ **Always check knowledge index first** - It's your map  
✓ **Follow related_concepts** - Provides full context  
✓ **Cite specific sections** - Users can verify  
✓ **Reference tests** - Shows real examples  
✓ **Use permalinks** - Ensures accuracy

✗ **Don't guess** - Check documentation first  
✗ **Don't assume** - Verify in knowledge index  
✗ **Don't skip validation** - Check test coverage

## Link Formatting for IDE Compatibility

When providing links to users:
- ✓ Use FILE PATH ONLY: `AI_documentation/operator_architecture.md`
- ✓ Mention section separately: "See section 1: High-Level Architecture"
- ✓ Use line numbers for code: `controllers/activemqartemis_controller.go:156-180`
- ✗ DON'T use anchor syntax: `file.md#section` (won't open in IDE)

Example:
"The reconciliation loop is in `AI_documentation/operator_architecture.md` 
(section 1: High-Level Architecture), implemented in 
`controllers/activemqartemis_controller.go:156-180`"

## Integration

This guide is referenced by:
- `.cursorrules` (Cursor auto-loads)
- `.vscode/settings.json` (Cursor context)
- `.github/copilot-instructions.md` (GitHub Copilot)

All AI tools should start here for project context.

---

**Version**: 1.0  
**Last Updated**: 2025-10-10  
**Documentation Commit**: 4acadb95

