# Security policy

## Supported versions

Public v2.0.2 is the current supported stable line. Retained v1 reports must
identify their affected versions separately. Main prepares v3; it is not a
published security-fix claim for either retained line.

## Reporting

Report vulnerabilities through GitHub private vulnerability reporting. Do not
open a public issue containing exploit details, customer schedules, credentials,
or database contents.

Include the affected version, safe reproduction, impact, and suggested
mitigation. Severity, acknowledgement and remediation targets, embargo,
advisory ranges and coordinated releases follow the pinned
[ecosystem vulnerability process](https://github.com/faustbrian/go-library-tools/blob/77bfd78c12a853f0d490bb27a3fbcb5f34330772/docs/ecosystem/security/vulnerability-management.md).

The package has no network transport, exporter, background worker, global
registry, unsafe code, or cgo. See the detailed [security model](docs/security.md).
