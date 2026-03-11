import re
from typing import Optional
from datetime import datetime
from pathlib import WindowsPath

def validate_username(username: str) -> bool:
    if not username or len(username) < 3:
        return False
    pattern = re.compile(r'^[a-zA-Z0-9_]+$')
    return pattern.match(username) is not None

def validate_message_content(content: str) -> bool:
    if not content or len(content) > 1000:
        return False
    return True

def sanitize_input(text: str) -> str:
    dangerous_chars = ['<', '>', '&', '"', "'"]
    for char in dangerous_chars:
        text = text.replace(char, '')
    return text.strip()

def validate_file_path(path: str) -> Optional[str]:
    path_obj = WindowsPath(path)
    if path_obj.exists():
        return str(path_obj.absolute())
    return None

class MessageValidator:
    def __init__(self):
        self.max_length = 1000
        self.min_length = 1

    def validate(self, message: dict) -> bool:
        if 'content' not in message:
            return False

        content = message['content']
        if not isinstance(content, str):
            return False

        if len(content) < self.min_length or len(content) > self.max_length:
            return False

        return True

def format_timestamp(ts: datetime) -> str:
    return ts.strftime('%Y-%m-%d %H:%M:%S')
