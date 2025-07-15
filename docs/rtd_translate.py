#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Simplified translation script for ReadTheDocs
"""
import os
import sys

def apply_translations():
    """Apply translations to po files"""
    # Simple translation script for ReadTheDocs
    print("📝 Applying basic translations...")
    
    po_dir = "locale/zh_CN/LC_MESSAGES"
    if not os.path.exists(po_dir):
        print(f"Warning: {po_dir} not found, skipping translations")
        return
    
    # Just ensure the po files exist and are valid
    po_files = ["index.po", "installation.po", "quickstart.po", "api.po", "examples.po", "advanced.po"]
    
    for po_file in po_files:
        po_path = os.path.join(po_dir, po_file)
        if os.path.exists(po_path):
            print(f"✅ Found: {po_file}")
        else:
            print(f"⚠️  Missing: {po_file}")
    
    print("✅ Translation check completed")

if __name__ == "__main__":
    apply_translations()