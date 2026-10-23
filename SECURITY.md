# Security policy

## Supported versions

Only the latest patch of the latest major receives security fixes. V6 is the
current major; v1 to v5 remain published but receive no further releases.

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
