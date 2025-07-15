# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#project-information

project = 'HX'
copyright = '2025, eatmoreapple'
author = 'eatmoreapple'
release = 'v1.0.0'

# -- General configuration ---------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#general-configuration

extensions = [
    'sphinx.ext.autodoc',
    'sphinx.ext.viewcode',
    'sphinx.ext.napoleon',
    'sphinx.ext.githubpages',
    'sphinx.ext.intersphinx',
]

templates_path = ['_templates']
exclude_patterns = ['_build', 'Thumbs.db', '.DS_Store']

# Internationalization
language = 'en'
locale_dirs = ['locale/']
gettext_compact = False

# Supported languages
languages = ['en', 'zh_CN']

# -- Options for HTML output -------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#options-for-html-output

html_theme = 'sphinx_rtd_theme'
html_static_path = ['_static']

# Theme options for language switching
html_theme_options = {
    'display_version': True,
    'prev_next_buttons_location': 'bottom',
    'style_external_links': False,
    'vcs_pageview_mode': '',
    'style_nav_header_background': '#2980B9',
    # Collapse navigation (False makes it tree-like)
    'collapse_navigation': False,
    # Breadcrumbs
    'navigation_depth': 4,
    # TOC options
    'includehidden': True,
    'titles_only': False
}

# Custom CSS for language switcher
html_css_files = [
    'css/custom.css',
]

# Language selector
html_context = {
    'current_language': 'en',
    'languages': [
        ('en', 'English'),
        ('zh_CN', '中文'),
    ],
    'language_links': True,
    # ReadTheDocs specific
    'current_version': 'latest',
    'versions': [('latest', 'latest')],
}

# -- Options for LaTeX output ------------------------------------------------

latex_elements = {
    'papersize': 'letterpaper',
    'pointsize': '10pt',
    'preamble': '',
    'fncychap': '\\usepackage[Bjornstrup]{fncychap}',
    'printindex': '\\footnotesize\\raggedright\\printindex',
}

# Grouping the document tree into LaTeX files. List of tuples
# (source start file, target name, title,
#  author, documentclass [howto, manual, or own class]).
latex_documents = [
    ('index', 'hx.tex', 'HX Documentation',
     'eatmoreapple', 'manual'),
]

# -- Extension configuration -------------------------------------------------

# Napoleon settings
napoleon_google_docstring = True
napoleon_numpy_docstring = True
napoleon_include_init_with_doc = False
napoleon_include_private_with_doc = False
napoleon_include_special_with_doc = True
napoleon_use_admonition_for_examples = False
napoleon_use_admonition_for_notes = False
napoleon_use_admonition_for_references = False
napoleon_use_ivar = False
napoleon_use_param = True
napoleon_use_rtype = True