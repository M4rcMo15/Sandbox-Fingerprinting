"""Injector modules for different system vectors."""

from .file_injector import FileInjector
from .registry_injector import RegistryInjector
from .process_injector import ProcessInjector
from .powershell_injector import PowerShellInjector
from .cmd_injector import CMDInjector
from .eventlog_injector import EventLogInjector
from .environment_injector import EnvironmentInjector
from .network_injector import NetworkInjector
from .shortcut_injector import ShortcutInjector
from .task_injector import ScheduledTaskInjector
from .service_injector import ServiceInjector
from .clipboard_injector import ClipboardInjector
from .debug_injector import DebugInjector
from .wmi_injector import WMIInjector
from .dns_injector import DNSInjector

__all__ = [
    "FileInjector",
    "RegistryInjector",
    "ProcessInjector",
    "PowerShellInjector",
    "CMDInjector",
    "EventLogInjector",
    "EnvironmentInjector",
    "NetworkInjector",
    "ShortcutInjector",
    "ScheduledTaskInjector",
    "ServiceInjector",
    "ClipboardInjector",
    "DebugInjector",
    "WMIInjector",
    "DNSInjector",
]
