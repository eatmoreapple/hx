#!/bin/bash

# Complete HX Documentation Build & Deploy Script
set -e

echo "🚀 Building Complete HX Documentation (English + Chinese)"
echo "=========================================================="

# Function to print status
print_status() {
    echo "📋 $1"
}

# Function to print success
print_success() {
    echo "✅ $1"
}

# Function to print error
print_error() {
    echo "❌ $1"
}

# Clean previous builds
print_status "Cleaning previous builds..."
rm -rf _build/
print_success "Cleaned previous builds"

# Step 1: Build English documentation
print_status "Building English documentation..."
make html
if [ $? -eq 0 ]; then
    print_success "English documentation built successfully"
else
    print_error "Failed to build English documentation"
    exit 1
fi

# Create multi-language directory structure
print_status "Setting up multi-language structure..."
mkdir -p _build/html-complete/{en,zh_CN}
cp -r _build/html/* _build/html-complete/en/
print_success "English documentation copied to multi-language structure"

# Step 2: Extract messages for translation
print_status "Extracting translatable messages..."
make gettext
if [ $? -eq 0 ]; then
    print_success "Messages extracted successfully"
else
    print_error "Failed to extract messages"
    exit 1
fi

# Step 3: Apply complete translations
print_status "Applying complete translations..."
python3 complete_translate.py
if [ $? -eq 0 ]; then
    print_success "All translations applied successfully"
else
    print_error "Failed to apply translations"
    exit 1
fi

# Step 4: Build Chinese documentation
print_status "Building Chinese documentation..."
make -e SPHINXOPTS="-D language='zh_CN'" html
if [ $? -eq 0 ]; then
    print_success "Chinese documentation built successfully"
else
    print_error "Failed to build Chinese documentation"
    exit 1
fi

# Copy Chinese documentation
print_status "Copying Chinese documentation..."
cp -r _build/html/* _build/html-complete/zh_CN/
print_success "Chinese documentation copied to multi-language structure"

# Step 5: Create language selection homepage
print_status "Creating language selection homepage..."
cat > _build/html-complete/index.html << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>HX Documentation - Go HTTP Framework</title>
    <meta name="description" content="HX - A lightweight, flexible HTTP framework for Go that simplifies request handling and data extraction">
    <meta name="keywords" content="Go, HTTP, framework, web, API, REST, middleware, type-safe">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #333;
        }
        
        .container {
            background: white;
            border-radius: 20px;
            padding: 3rem;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            text-align: center;
            max-width: 600px;
            width: 90%;
        }
        
        .logo {
            font-size: 4rem;
            font-weight: bold;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 1rem;
        }
        
        .subtitle {
            color: #666;
            font-size: 1.2rem;
            margin-bottom: 2rem;
            line-height: 1.6;
        }
        
        .languages {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 1.5rem;
            margin-top: 2rem;
        }
        
        .language-card {
            display: block;
            padding: 2rem 1.5rem;
            background: #f8f9fa;
            border: 2px solid transparent;
            border-radius: 15px;
            text-decoration: none;
            color: #333;
            transition: all 0.3s ease;
            position: relative;
            overflow: hidden;
        }
        
        .language-card::before {
            content: '';
            position: absolute;
            top: 0;
            left: -100%;
            width: 100%;
            height: 100%;
            background: linear-gradient(90deg, transparent, rgba(255,255,255,0.6), transparent);
            transition: left 0.5s;
        }
        
        .language-card:hover::before {
            left: 100%;
        }
        
        .language-card:hover {
            transform: translateY(-5px);
            border-color: #667eea;
            box-shadow: 0 10px 25px rgba(102, 126, 234, 0.2);
        }
        
        .language-flag {
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
            display: block;
        }
        
        .language-name {
            font-size: 1.3rem;
            font-weight: 600;
            margin-bottom: 0.5rem;
        }
        
        .language-desc {
            color: #666;
            font-size: 0.9rem;
        }
        
        .features {
            margin-top: 3rem;
            padding-top: 2rem;
            border-top: 1px solid #eee;
        }
        
        .features-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
            gap: 1rem;
            margin-top: 1rem;
        }
        
        .feature {
            padding: 1rem;
            background: #f8f9fa;
            border-radius: 10px;
            font-size: 0.85rem;
            color: #666;
        }
        
        .feature-icon {
            font-size: 1.5rem;
            margin-bottom: 0.5rem;
            display: block;
        }
        
        @media (max-width: 768px) {
            .container {
                padding: 2rem;
                margin: 1rem;
            }
            
            .logo {
                font-size: 3rem;
            }
            
            .languages {
                grid-template-columns: 1fr;
            }
            
            .features-grid {
                grid-template-columns: 1fr 1fr;
            }
        }
        
        .github-link {
            position: absolute;
            top: 20px;
            right: 20px;
            color: white;
            font-size: 1.5rem;
            text-decoration: none;
            opacity: 0.8;
            transition: opacity 0.3s;
        }
        
        .github-link:hover {
            opacity: 1;
        }
    </style>
</head>
<body>
    <a href="https://github.com/eatmoreapple/hx" class="github-link" title="View on GitHub">⭐</a>
    
    <div class="container">
        <div class="logo">HX</div>
        <div class="subtitle">
            A lightweight, flexible HTTP framework for Go<br>
            Choose your preferred language to get started
        </div>
        
        <div class="languages">
            <a href="en/" class="language-card">
                <span class="language-flag">🇺🇸</span>
                <div class="language-name">English</div>
                <div class="language-desc">Complete documentation</div>
            </a>
            
            <a href="zh_CN/" class="language-card">
                <span class="language-flag">🇨🇳</span>
                <div class="language-name">中文</div>
                <div class="language-desc">完整文档</div>
            </a>
        </div>
        
        <div class="features">
            <div class="features-grid">
                <div class="feature">
                    <span class="feature-icon">🚀</span>
                    Lightweight & Fast
                </div>
                <div class="feature">
                    <span class="feature-icon">💪</span>
                    Type-safe
                </div>
                <div class="feature">
                    <span class="feature-icon">🔄</span>
                    Auto Binding
                </div>
                <div class="feature">
                    <span class="feature-icon">🛠</span>
                    Extensible
                </div>
            </div>
        </div>
    </div>
</body>
</html>
EOF

print_success "Language selection homepage created"

# Step 6: Create sitemap for SEO
print_status "Creating sitemap..."
cat > _build/html-complete/sitemap.xml << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
        xmlns:xhtml="http://www.w3.org/1999/xhtml">
    
    <url>
        <loc>https://hx.readthedocs.io/</loc>
        <changefreq>daily</changefreq>
        <priority>1.0</priority>
        <xhtml:link rel="alternate" hreflang="en" href="https://hx.readthedocs.io/en/"/>
        <xhtml:link rel="alternate" hreflang="zh-CN" href="https://hx.readthedocs.io/zh_CN/"/>
    </url>
    
    <url>
        <loc>https://hx.readthedocs.io/en/</loc>
        <changefreq>weekly</changefreq>
        <priority>0.9</priority>
        <xhtml:link rel="alternate" hreflang="zh-CN" href="https://hx.readthedocs.io/zh_CN/"/>
    </url>
    
    <url>
        <loc>https://hx.readthedocs.io/zh_CN/</loc>
        <changefreq>weekly</changefreq>
        <priority>0.9</priority>
        <xhtml:link rel="alternate" hreflang="en" href="https://hx.readthedocs.io/en/"/>
    </url>
    
</urlset>
EOF

print_success "Sitemap created"

# Step 7: Create robots.txt
print_status "Creating robots.txt..."
cat > _build/html-complete/robots.txt << 'EOF'
User-agent: *
Allow: /

Sitemap: https://hx.readthedocs.io/sitemap.xml
EOF

print_success "Robots.txt created"

# Step 8: Generate build statistics
print_status "Generating build statistics..."
EN_FILES=$(find _build/html-complete/en -name "*.html" | wc -l)
ZH_FILES=$(find _build/html-complete/zh_CN -name "*.html" | wc -l)
TOTAL_SIZE=$(du -sh _build/html-complete | cut -f1)

# Step 9: Create README for the build
cat > _build/html-complete/README.md << EOF
# HX Documentation Build

This directory contains the complete multi-language documentation for HX.

## Build Information

- **Build Date**: $(date)
- **English Pages**: $EN_FILES
- **Chinese Pages**: $ZH_FILES  
- **Total Size**: $TOTAL_SIZE
- **Languages**: English (en), Chinese (zh_CN)

## Directory Structure

\`\`\`
html-complete/
├── index.html          # Language selection homepage
├── sitemap.xml         # SEO sitemap
├── robots.txt          # Search engine instructions
├── en/                 # English documentation
│   ├── index.html
│   ├── installation.html
│   ├── quickstart.html
│   ├── api.html
│   ├── examples.html
│   └── advanced.html
└── zh_CN/              # Chinese documentation
    ├── index.html
    ├── installation.html
    ├── quickstart.html
    ├── api.html
    ├── examples.html
    └── advanced.html
\`\`\`

## Deployment

This build is ready for deployment to:
- ReadTheDocs
- GitHub Pages
- Netlify
- Any static web server

## Features

- 🌍 Multi-language support (English/Chinese)
- 📱 Responsive design
- 🔍 SEO optimized
- ⚡ Fast loading
- 🎨 Modern UI
- 🔄 Language switching

Generated with ❤️ by HX Documentation Build System
EOF

print_success "Build documentation created"

# Final summary
echo ""
echo "🎉 Complete Documentation Build Finished!"
echo "=========================================="
echo "📊 Build Summary:"
echo "   📄 English pages: $EN_FILES"
echo "   📄 Chinese pages: $ZH_FILES"
echo "   💾 Total size: $TOTAL_SIZE"
echo "   🌐 Languages: English, Chinese"
echo "   📱 Features: Responsive, SEO-optimized, Multi-language"
echo ""
echo "📂 Output Directory: _build/html-complete/"
echo "🌐 Language Selection: _build/html-complete/index.html"
echo "🇺🇸 English Docs: _build/html-complete/en/"
echo "🇨🇳 Chinese Docs: _build/html-complete/zh_CN/"
echo ""
echo "🚀 Ready for deployment to ReadTheDocs, GitHub Pages, or any web server!"

# Start local server if requested
if [ "$1" == "serve" ]; then
    echo ""
    echo "🌐 Starting local server on http://localhost:8000"
    echo "   Press Ctrl+C to stop"
    cd _build/html-complete
    python3 -m http.server 8000
fi