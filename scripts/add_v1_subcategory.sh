#!/bin/bash

# Script to add subcategory: "v1" to specific v1 documentation files

set -e

# Define the docs directory
DOCS_DIR="docs"

# Define specific V1 files that need subcategory: "v1"
V1_RESOURCES=(
    "integration"
    "access_bundle"
    "manual_webhook"
)

V1_DATA_SOURCES=(
    "connector"
    "integrations"
    "integration"
)

# Function to add subcategory to frontmatter
add_subcategory() {
    local file="$1"
    
    if [ -f "$file" ]; then
        echo "Processing: $file"
        
        # Check if subcategory already exists with empty value and replace it
        if grep -q '^subcategory: ""' "$file"; then
            sed -i '' 's/^subcategory: ""/subcategory: "v1"/' "$file"
            echo "  Replaced empty subcategory with v1"
        fi
    else
        echo "File not found: $file"
    fi
}

# Main execution
echo "Adding subcategory 'v1' to specific v1 documentation files..."
echo "============================================================"

# Process V1 resources
echo ""
echo "Processing V1 resources..."
for resource in "${V1_RESOURCES[@]}"; do
    resource_file="$DOCS_DIR/resources/$resource.md"
    add_subcategory "$resource_file"
done

# Process V1 data sources
echo ""
echo "Processing V1 data sources..."
for datasource in "${V1_DATA_SOURCES[@]}"; do
    datasource_file="$DOCS_DIR/data-sources/$datasource.md"
    add_subcategory "$datasource_file"
done

echo ""
echo "Done! Specified v1 documentation files have been updated."