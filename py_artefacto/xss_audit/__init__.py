"""XSS Sandbox Audit - Security testing tool for sandbox environments."""

from .auditor import XSSAuditor
from .payloads import XSSPayload, get_default_payloads

__version__ = "1.0.0"
__all__ = ["XSSAuditor", "XSSPayload", "get_default_payloads"]
