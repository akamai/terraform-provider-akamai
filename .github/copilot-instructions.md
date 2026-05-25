# AI Review Guidelines — `terraform-provider-akamai`

## Meta Information for AI Agents

**Review scope**: AI agents must always review the **entire pull request in full**,
not just the incremental diff. This means:
- Read and consider all files in the repository that are relevant to the changes,
  not only the lines added or modified.
- Evaluate the impact of changes in the context of the whole codebase.
- Ensure consistency, correctness, and adherence to project conventions across the
  full scope of the PR, including unchanged surrounding code and its dependencies.

---

This file is the source of truth used by AI review agents when reviewing pull requests
in the `terraform-provider-akamai` repository. Every rule below MUST be checked for
the changed code (additions and modifications). Do NOT flag pre-existing issues in
code that is not part of the change unless it is directly touched by the PR.

Each rule has a stable `id` so review tools can reference it.

---

## 1. General practices

- **GEN-01**: If newly added functionality is not yet Globally Available (GA), the
  code must contain a comment stating when it is expected to reach GA.
- **GEN-02**: When the change is part of a multi-repo change (e.g. requires a
  matching edgegrid-golang or cli-terraform change), the branch name in this
  repo must be exactly identical to the branch name in the other repos.
- **GEN-03**: A change that forces an end-user to manually adjust their Terraform
  configuration (resource/data source schema change) is a breaking change. It
  must target the `feature/sp-breaking-changes` branch and is released only once
  a quarter or less often. When suspecting breaking changes in the PR, let someone
  from the DEVEXP team know to help identify the correct `sp-breaking-changes`
  branch. Exception: features still in Beta.
- **GEN-04**: Push code in major batches. Do NOT push after every tiny piece of
  work — the build takes 35+ minutes on limited shared resources. Run linter
  and tests locally first.
- **GEN-05**: A Documentation ticket must be created early enough so the
  techwriter team has time to prepare the documentation. Track comments and
  status changes there.
- **GEN-06**: When adding a new resource, adding a new required/optional
  attribute, or dropping an existing required/optional attribute, evaluate
  whether `cli-terraform` must change as well. Those changes ship together.
- **GEN-07**: If the PR adds methods that require edgegrid-golang changes and/or
  cli-terraform changes, state in the PR that all related PRs must be released
  together.

## 2. Common coding practices (apply to every change)

- **COD-01**: Any change affecting customers must be reflected in `CHANGELOG.md`.
- **COD-02**: No secrets, real contract IDs, real Akamai employee IDs, or any
  other non-public data may appear anywhere in the code (especially unit tests
  and fixtures).
- **COD-03**: Never skip error checking. Every returned `error` must be handled.
- **COD-04**: Do not add files that are never used.
- **COD-05**: Prefer `any` over `interface{}`.
- **COD-06**: Acronyms in exported identifiers (e.g. `IP`, `DNS`, `URL`, `ID`,
  `API`) must be in all uppercase. Unexported identifiers are exempt.
- **COD-07**: Descriptions and comments that are full sentences must end with a
  period.
- **COD-08**: Unit tests must cover all corner cases — especially presence and
  absence of fields. Cases may be aggregated where it makes sense.

## 3. Changelog rules

- **CHG-01**: Entries describe what changed from the customer perspective of
  THIS project (terraform-provider-akamai).
- **CHG-02**: Fixes/changes based on GitHub issues must include a link to that
  issue, e.g.
  `([#436](https://github.com/akamai/terraform-provider-akamai/issues/436))`.
- **CHG-03**: Use past tense.
- **CHG-04**: Use backticks (`` ` ``) around proper names (resources, data
  sources, attributes, packages, etc.).
- **CHG-05**: Entries should be placed in a random line position within the
  correct section/release to mitigate merge conflicts.
- **CHG-06**: Do NOT delete empty lines in the changelog.

## 4. PR hygiene

- **PRH-01**: Keep WIP / not-ready PRs in Draft.
- **PRH-02**: Branch must be rebased on top of its target branch before review.
- **PRH-03**: Commits squashed where it makes sense.
- **PRH-04**: Every commit message must include the appropriate JIRA number and
  the story title or a description of the change.

---

## 5. New subprovider

- **SUB-01**: When introducing a new subprovider, its boilerplate (Client
  creation, parameter passing to resources/data sources, package layout) must
  match the current reference implementation. The reference is the
  `cloudcertificates` package.

## 6. Common coding practices specific to the provider

- **TFP-01**: All new resources and data sources must be written in
  **Terraform Plugin Framework**. The legacy SDKv2 library must NOT be used for
  new code.
- **TFP-02**: At the top of every resource/data source file, declare compile-time
  interface compliance using blank-identifier assignments and list ALL
  interfaces the type is intended to satisfy, e.g.:
  ```go
  var _ resource.Resource                = &myResource{}
  var _ resource.ResourceWithImportState = &myResource{}
  ```
- **TFP-03**: Every schema attribute must have a description. If the description
  has multiple lines, use `MarkdownDescription` instead of `Description`.
- **TFP-04**: If a method returns `diagnostics` or an `error`, the result must
  NEVER be ignored.
- **TFP-05**: When a function call returns diagnostics, first append them into
  `resp` (e.g. `*resource.CreateResponse`) and THEN check for errors via the
  whole `resp`. Do not check the returned diagnostics directly and then push
  them into `resp`. If processing should not be continued after getting such
  an error, make sure to return immediately.
- **TFP-06**: When modelling a collection, decide between `List` and `Set` based
  on:
    - whether order matters to the API,
    - whether the API returns entries in the same order as sent,
    - how customers will use it (sets cannot be indexed).

  When reviewing a new collection attribute, ask the author whether the choice
  of `List` vs `Set` was a conscious decision or accidental. The choice must be
  made deliberately by the developer.
- **TFP-07**: If the configuration can be validated locally (schema validators
  or `ValidateConfig` on resources), it must be validated even if edgegrid-golang
  already validates it. Use `ValidateConfig` for plan-level validation,
  including validations that need API calls. Consider `NotEmpty` validators
  when empty values are nonsensical.
- **TFP-08**: In `provider.go`, keep the lists of resources and data sources
  ordered alphabetically.
- **TFP-09**: Add `id` to a resource/data source only when it makes sense from a
  business perspective — Framework no longer requires it.
- **TFP-10**: When handling errors returned from edgegrid-golang, prefer using
  `errors.Is(err, <package>.ErrXxx)` with exported sentinel errors over
  inspecting HTTP status codes or parsing error message strings. If a needed
  sentinel error does not exist in edgegrid-golang yet, request it from the
  library maintainers (or add it yourself in a companion PR).
- **TFP-11**: Use the following high-level utility functions/packages whenever
  applicable rather than reimplementing the same logic:

  | Function / package | Use case |
  | :--- | :--- |
  | `text.ImportIDSplitter` | Splitting the `importID` during resource import. Standard separators: `,` (id part) and `[` (optional parts). |
  | `retry` package | Repeatedly checking (via API call or otherwise) whether a status or value has reached a desired state. |
  | `tf.IsKnown` | Checking that a value is neither `null` nor `unknown`. |
  | `common/framework/date` package | Converting `time.Time` to `types.String`. |
  | `modifiers.IsCreate`, `modifiers.IsUpdate`, `modifiers.IsDelete` | Identifying within `ModifyPlan` which lifecycle flow is being processed. |
  | `modifiers` package | Various `PlanModifier` helpers (e.g. prefix equality, preventing value updates). |
  | `ptr.To` | Obtaining a pointer from any value. |

## 7. Resource rules

- **RES-01**: If an entity has versioning, do NOT introduce a separate resource
  for the version — a single resource manages both the entity and its versions.
- **RES-02**: If a resource waits for any operation (activation, status change,
  etc.), it must declare a customisable `Timeout` with a reasonable default
  based on real API behaviour plus margin.
- **RES-03**: If creation involves multiple API calls, verify the resource
  behaves correctly when the first (creating) call succeeds and a later call
  fails. In that case the resource must enter the Tainted state so the next
  apply will delete and recreate it.
- **RES-04**: For ID attributes (e.g. contractID, groupID, productID,
  propertyID), decide whether the field supports the prefix form (e.g.
  `ctr_1-ABC`), the unprefixed form (`1-ABC`), or both. When both are
  acceptable and must be treated as equal (e.g. `cpc_111 == 111`), use the
  `IgnorePrefixType` custom type from `internal/customtypes` instead of manual
  normalisation.
- **RES-05**: If a resource disappears from the server, `Read` must drop it from
  state (`resp.State.RemoveResource(ctx)` in Framework, `rd.SetId("")` in
  SDKv2) and emit a warning about the operation.
- **RES-06**: For nested objects, use `types.Object` instead of nesting a model
  struct inside another model struct. For collections of nested elements use
  `types.List` / `types.Set` rather than `[]nestedModel`. This is required for
  correct handling of the unknown state.
- **RES-07**: A computed attribute that never changes after creation must
  declare the `UseStateForUnknown` plan modifier. Attributes that are computed
  but sometimes change must be handled manually in `ModifyPlan` to populate
  the planned value when known.
- **RES-08**: If the API does not allow updating a particular optional/required
  field, block that operation in the schema (e.g. `PreventInt64Update`).
- **RES-09**: When adding a new field with a default value to an existing
  resource, verify the behaviour for customers who already have the resource
  in state.
- **RES-10**: Assuming the user provided correct configuration matching server
  state, the first `terraform plan` after `import` must be empty.
- **RES-11**: In `ModifyPlan` and `ValidateConfig`, account for values coming
  from variables / other resources — those will appear as `Unknown`. Add unit
  tests for these cases.
- **RES-12**: A primitive-type field that is not always returned must be set to
  `null` in state, never to Go's zero value.
- **RES-13**: When checking for value presence, prefer `tf.IsKnown`.
- **RES-14** *(suggestion, non-blocking)*: For functions that heavily use a
  model struct, suggest making the model the receiver. In particular, when a
  model contains nested `types.Object` fields, suggest adding typed
  getter/setter methods on the model, e.g.:
  ```go
  func (m *model) nestedObj(ctx context.Context) (*nestedModel, diag.Diagnostics)
  func (m *model) setNestedObj(ctx context.Context, n *nestedModel) diag.Diagnostics
  ```
- **RES-15** *(suggestion, non-blocking)*: When implementing import, suggest
  using `text.ImportIDSplitter` to split the ID. Standard separators are `,`
  (id part separator) and `[` (optional parts follow inside brackets).

## 8. Tests

- **TST-01**: When adding/extending/fixing a resource or data source, the change
  must be covered by unit tests.
- **TST-02**: Use the state checker helpers `NewStateChecker` and
  `NewImportChecker`. Check ALL possible present and missing fields. Build on
  common state-checker bases and use `AttributeBatch` checkers where they
  make sense.
- **TST-03**: Every test function must call `AssertExpectations` at the end.
- **TST-04**: Prefer one test function (test suite) that holds multiple test
  cases for a single resource/data source.
- **TST-05**: Prefer a high-level helper in the unit test that mocks a group of
  API calls, so that the test's `init` is a sequence of high-level calls that
  describe what the test covers.
- **TST-06**: When there are multiple test cases, sort them by the flow they
  cover: positive scenarios for `create`, `update`, `refresh`, other flow
  checks (e.g. force new), and import — in that order. After all positive
  tests, place the negative scenarios. Within each flow, order tests from
  simplest to most complex. This order is recommended but not enforced.
- **TST-07**: When mocking API calls (in the init block):
    - Split mock calls into clearly labelled sections: `// create`,
      `// update`, `// read` (or multiple `// read` sections), `// delete`.
    - Every mock call must specify call count via `Times(n)` (or `Once` /
      `Twice`).
    - Always pass the FULL expected request in the mock — never use
      `mock.Anything`, `mock.MatchedBy`, or `mock.AnythingOfType`. For
      `context.Context`, use `testutils.MockContext`.
- **TST-08**: For negative tests, assert the entire error message produced by
  the subprovider (the callstack may be omitted).
- **TST-09**: `t.Parallel()` must be present wherever applicable.
- **TST-10**: Generally do not write dedicated successful-delete tests — delete
  is exercised by create/update tests. A dedicated delete test is acceptable
  only if there are multiple successful-delete paths or for negative delete
  scenarios (which still require a working successful path).
- **TST-11**: If a unit test covers a flow involving status changes, timeouts,
  or activations, override the relevant intervals/timeouts (to a significantly
  small value) inside the test.

