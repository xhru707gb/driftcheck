# driftcheck

> Compares live cloud infrastructure state against Terraform definitions and reports configuration drift.

---

## Installation

```bash
go install github.com/yourusername/driftcheck@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftcheck.git
cd driftcheck
go build -o driftcheck ./cmd/driftcheck
```

---

## Usage

Point `driftcheck` at your Terraform state file and target cloud provider to detect drift:

```bash
driftcheck --state terraform.tfstate --provider aws --region us-east-1
```

**Example output:**

```
[DRIFT] aws_instance.web_server
  expected: instance_type = "t3.micro"
  actual:   instance_type = "t3.small"

[DRIFT] aws_security_group.app_sg
  expected: ingress.0.cidr_blocks = ["10.0.0.0/8"]
  actual:   ingress.0.cidr_blocks = ["0.0.0.0/0"]

2 drift(s) detected across 14 resources checked.
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--state` | Path to Terraform state file | `terraform.tfstate` |
| `--provider` | Cloud provider (`aws`, `gcp`, `azure`) | `aws` |
| `--region` | Cloud region to query | `us-east-1` |
| `--output` | Output format (`text`, `json`) | `text` |
| `--fail-on-drift` | Exit with code 1 if drift is found | `false` |

---

## Requirements

- Go 1.21+
- Valid cloud provider credentials (e.g., AWS credentials via environment or `~/.aws/credentials`)

---

## License

MIT © 2024 yourusername