#!/bin/bash

# Multi-language documentation build script
set -e

echo "Building HX Documentation in multiple languages..."

# Clean previous builds
echo "Cleaning previous builds..."
rm -rf _build/

# Build English (default)
echo "Building English documentation..."
make html
mkdir -p _build/html-multi/en
cp -r _build/html/* _build/html-multi/en/

# Extract translatable messages
echo "Extracting translatable messages..."
make gettext

# Update translation files
echo "Updating translation files..."
sphinx-intl update -p _build/gettext -l zh_CN

# Build Chinese documentation
echo "Building Chinese documentation..."
make -e SPHINXOPTS="-D language='zh_CN'" html
mkdir -p _build/html-multi/zh_CN
cp -r _build/html/* _build/html-multi/zh_CN/

# Create index page for language selection
echo "Creating language selection index..."
cat > _build/html-multi/index.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>HX Documentation</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            max-width: 600px;
            margin: 100px auto;
            padding: 20px;
            text-align: center;
        }
        .language-card {
            display: inline-block;
            margin: 20px;
            padding: 30px;
            border: 2px solid #e1e4e8;
            border-radius: 8px;
            text-decoration: none;
            color: #24292e;
            transition: all 0.2s;
        }
        .language-card:hover {
            border-color: #2980B9;
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(41, 128, 185, 0.15);
        }
        .language-name {
            font-size: 1.5em;
            font-weight: bold;
            margin-bottom: 10px;
        }
        .language-desc {
            color: #586069;
        }
        h1 {
            color: #2980B9;
            margin-bottom: 40px;
        }
    </style>
</head>
<body>
    <h1>HX Documentation</h1>
    <p>Choose your preferred language:</p>
    
    <a href="en/" class="language-card">
        <div class="language-name">English</div>
        <div class="language-desc">English Documentation</div>
    </a>
    
    <a href="zh_CN/" class="language-card">
        <div class="language-name">中文</div>
        <div class="language-desc">中文文档</div>
    </a>
</body>
</html>
EOF

echo "Multi-language documentation built successfully!"
echo "English: _build/html-multi/en/"
echo "Chinese: _build/html-multi/zh_CN/"
echo "Language selector: _build/html-multi/index.html"

# Start local server if requested
if [ "$1" == "serve" ]; then
    echo "Starting local server on http://localhost:8000"
    cd _build/html-multi
    python3 -m http.server 8000
fi