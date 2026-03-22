"""Tests for XSS auditor."""

import unittest
from xss_audit import XSSAuditor, XSSPayload


class TestXSSAuditor(unittest.TestCase):
    """Test XSSAuditor class."""
    
    def test_auditor_initialization(self):
        """Test auditor initializes correctly."""
        auditor = XSSAuditor()
        
        self.assertIsNotNone(auditor.payloads)
        self.assertGreater(len(auditor.payloads), 0)
        
        # Check all injectors are initialized
        self.assertIsNotNone(auditor.file_injector)
        self.assertIsNotNone(auditor.registry_injector)
        self.assertIsNotNone(auditor.process_injector)
        self.assertIsNotNone(auditor.powershell_injector)
        self.assertIsNotNone(auditor.cmd_injector)
    
    def test_auditor_with_custom_payloads(self):
        """Test auditor with custom payloads."""
        custom_payload = XSSPayload(
            id="test_001",
            type="test",
            content="<script>test</script>",
            vector="all",
            description="Test payload"
        )
        
        auditor = XSSAuditor(payloads=[custom_payload])
        
        self.assertEqual(len(auditor.payloads), 1)
        self.assertEqual(auditor.payloads[0].id, "test_001")
    
    def test_injection_methods_exist(self):
        """Test that all injection methods exist."""
        auditor = XSSAuditor()
        
        methods = [
            'inject_files',
            'inject_registry',
            'inject_processes',
            'inject_powershell',
            'inject_cmd',
            'inject_eventlogs',
            'inject_environment',
            'inject_network',
            'inject_shortcuts',
            'inject_tasks',
            'inject_services',
            'inject_clipboard',
            'inject_debug',
            'inject_wmi',
            'inject_dns',
        ]
        
        for method in methods:
            self.assertTrue(hasattr(auditor, method))
            self.assertTrue(callable(getattr(auditor, method)))


if __name__ == '__main__':
    unittest.main()
