"""Base injector class."""

from abc import ABC, abstractmethod
from ..payloads import XSSPayload


class BaseInjector(ABC):
    """Abstract base class for all injectors."""
    
    @abstractmethod
    def inject(self, payload: XSSPayload):
        """
        Inject the XSS payload into the target vector.
        
        Args:
            payload: The XSSPayload object to inject
        """
        pass
