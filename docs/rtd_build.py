#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
ReadTheDocs build script for HX documentation
Ensures both English and Chinese versions are built
"""

import os
import subprocess
import sys

def run_command(cmd, cwd=None):
    """Run a command and return success status"""
    print(f"Running: {cmd}")
    try:
        result = subprocess.run(cmd, shell=True, check=True, cwd=cwd, 
                              capture_output=True, text=True)
        print(f"Success: {result.stdout}")
        return True
    except subprocess.CalledProcessError as e:
        print(f"Error: {e}")
        print(f"Stderr: {e.stderr}")
        return False

def main():
    """Main build function"""
    print("🚀 Building HX Documentation for ReadTheDocs...")
    
    # Get the docs directory
    docs_dir = os.path.join(os.path.dirname(__file__))
    
    # Step 1: Apply translations
    print("📝 Applying translations...")
    if not run_command("python3 complete_translate.py", cwd=docs_dir):
        print("❌ Translation failed")
        sys.exit(1)
    
    # Step 2: Build Chinese version if we're in post_build
    if os.environ.get('READTHEDOCS_OUTPUT'):
        print("🇨🇳 Building Chinese documentation...")
        if not run_command("make -e SPHINXOPTS=\"-D language='zh_CN'\" html", cwd=docs_dir):
            print("❌ Chinese build failed")
            sys.exit(1)
        
        # Copy Chinese docs to ReadTheDocs output
        output_dir = os.environ['READTHEDOCS_OUTPUT']
        zh_dir = os.path.join(output_dir, 'html', 'zh_CN')
        
        print(f"📁 Creating Chinese directory: {zh_dir}")
        os.makedirs(zh_dir, exist_ok=True)
        
        print("📋 Copying Chinese documentation...")
        if not run_command(f"cp -r _build/html/* {zh_dir}/", cwd=docs_dir):
            print("❌ Failed to copy Chinese docs")
            sys.exit(1)
        
        print("✅ Chinese documentation built successfully!")
    
    print("🎉 Build completed successfully!")

if __name__ == "__main__":
    main()