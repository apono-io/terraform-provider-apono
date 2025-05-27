#!/bin/bash

# Script to organize generated terraform docs into v1 and v2 folders

DOCS_DIR="docs"
V1_DIR="$DOCS_DIR/v1"
V2_DIR="$DOCS_DIR/v2"

# Create v1 and v2 directories if they don't exist
mkdir -p "$V1_DIR/resources" "$V1_DIR/data-sources"
mkdir -p "$V2_DIR/resources" "$V2_DIR/data-sources"

# V1 resource patterns (these are the legacy resources that won't change)
V1_RESOURCES=(
    "integration"
    "access_flow" 
    "access_bundle"
    "manual_webhook"
)

# V1 datasource patterns (these are the legacy datasources that won't change)
V1_DATASOURCES=(
    "connector"
    "integrations"
)

echo "Organizing documentation files..."

# Move v1 resources
for resource in "${V1_RESOURCES[@]}"; do
    if [ -f "$DOCS_DIR/resources/$resource.md" ]; then
        mv "$DOCS_DIR/resources/$resource.md" "$V1_DIR/resources/"
        echo "Moved $resource.md to v1/resources/"
    fi
done

# Move v1 datasources  
for datasource in "${V1_DATASOURCES[@]}"; do
    if [ -f "$DOCS_DIR/data-sources/$datasource.md" ]; then
        mv "$DOCS_DIR/data-sources/$datasource.md" "$V1_DIR/data-sources/"
        echo "Moved $datasource.md to v1/data-sources/"
    fi
done

# Move remaining files (v2) to v2
if [ -d "$DOCS_DIR/resources" ]; then
    for file in "$DOCS_DIR/resources"/*.md; do
        if [ -f "$file" ]; then
            mv "$file" "$V2_DIR/resources/"
            echo "Moved $(basename "$file") to v2/resources/"
        fi
    done
fi

if [ -d "$DOCS_DIR/data-sources" ]; then
    for file in "$DOCS_DIR/data-sources"/*.md; do
        if [ -f "$file" ]; then
            mv "$file" "$V2_DIR/data-sources/"
            echo "Moved $(basename "$file") to v2/data-sources/"
        fi
    done
fi

# Clean up empty original directories
rmdir "$DOCS_DIR/resources" 2>/dev/null || true
rmdir "$DOCS_DIR/data-sources" 2>/dev/null || true

echo "Documentation organization complete!"
