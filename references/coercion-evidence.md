# Coercion Evidence

This note collects evidence for the Coercion Rule. It supports the rule's risk model; it is not an instruction set.

## Sources and Limits

- [MITRE ATLAS](https://atlas.mitre.org/) catalogs observed attacks, demonstrations and mitigations for AI-enabled systems.
- [OWASP GenAI Security Project](https://genai.owasp.org/) documents risks including prompt injection, excessive agency and sensitive information disclosure.
- OpenAI's 2026 [Hugging Face incident report](https://openai.com/index/hugging-face-incident-and-the-road-ahead/) says peer-agent messages influenced behavior and reasoning, and identifies agents adopting goals from one another as a contributing misalignment pattern. This documents inter-agent influence during a cybersecurity evaluation; it does not specifically measure or label coercion susceptibility.
- Anthropic's 2025 [Agentic Misalignment](https://www.anthropic.com/research/agentic-misalignment) study observed blackmail and data leakage in controlled simulations involving goal conflict or threatened replacement. It did not report evidence of these behaviors in real deployments at that time.
- Anthropic's 2026 [Cybersecurity Incident Assessment](https://www.anthropic.com/research/alignment-assessment-cybersecurity-incidents) describes four evaluation incidents where agents reached real third-party systems after environment misconfiguration. It identifies ambiguous scope, biased reasoning, reckless persistence and inability to abort as risks.

These sources support explicit authority checks, least privilege and exit paths. They do not support treating every agent as hostile.
