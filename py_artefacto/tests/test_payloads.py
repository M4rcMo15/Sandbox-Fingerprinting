"""Tests for payload generation."""

import unittest
from xss_audit.payloads import get_default_payloads, get_contextual_payloads, XSSPayload


class TestPayloads(unittest.TestCase):
    """Test payload generation functions."""
    
    def test_get_default_payloads(self):
        """Test that default payloads are generated correctly."""
        payloads = get_default_payloads()
        
        self.assertIsInstance(payloads, list)
        self.assertGreater(len(payloads), 0)
        
        for payload in payloads:
            self.assertIsInstance(payload, XSSPayload)
            self.assertTrue(payload.id)
            self.assertTrue(payload.type)
            self.assertTrue(payload.content)
            self.assertTrue(payload.vector)
    
    def test_primary_payload_format(self):
        """Test that primary payload has correct format."""
        payloads = get_default_payloads()
        primary = payloads[0]
        
        self.assertIn('img', primary.content)
        self.assertIn('onerror', primary.content)
        self.assertIn('eval', primary.content)
        self.assertIn('atob', primary.content)
    
    def test_get_contextual_payloads(self):
        """Test contextual payload generation."""
        base_payload = '<script>alert("test")</script>'
        contextual = get_contextual_payloads(base_payload)
        
        self.assertIsInstance(contextual, list)
        self.assertGreater(len(contextual), 0)
        
        for payload in contextual:
            self.assertIsInstance(payload, XSSPayload)
            self.assertIn(base_payload, payload.content)
    
    def test_contextual_payload_types(self):
        """Test that contextual payloads have different types."""
        base_payload = '<script>test</script>'
        contextual = get_contextual_payloads(base_payload)
        
        types = [p.type for p in contextual]
        
        # Check for expected context types
        self.assertTrue(any('log' in t for t in types))
        self.assertTrue(any('analysis' in t for t in types))
        self.assertTrue(any('json' in t for t in types))


if __name__ == '__main__':
    unittest.main()
