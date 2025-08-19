#!/bin/bash

# Change to the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Change to project root directory (one level up from proto/)
cd ..

echo "=== Proto Generation Started ==="
echo "Working directory: $(pwd)"
echo

# Counter for tracking generation
total_files=0
generated_files=0

# Find all proto files in subdirectories of proto/
for proto_dir in proto/*/; do
    if [ -d "$proto_dir" ]; then
        dir_name=$(basename "$proto_dir")
        echo "📁 Scanning directory: $dir_name"
        
        # Check if directory contains .proto files
        proto_files=($(find "$proto_dir" -name "*.proto" -type f))
        
        if [ ${#proto_files[@]} -eq 0 ]; then
            echo "   ⚠️  No .proto files found in $dir_name"
            continue
        fi
        
        echo "   📄 Found ${#proto_files[@]} proto file(s) in $dir_name"
        
        # Generate each proto file
        for proto_file in "${proto_files[@]}"; do
            total_files=$((total_files + 1))
            echo "   🔧 Generating: $proto_file"
            
            if protoc --go_out=. --go-grpc_out=. "$proto_file" 2>/dev/null; then
                generated_files=$((generated_files + 1))
                echo "   ✅ Successfully generated: $proto_file"
            else
                echo "   ❌ Failed to generate: $proto_file"
                echo "      Running with verbose output:"
                protoc --go_out=. --go-grpc_out=. "$proto_file"
            fi
        done
        echo
    fi
done

# Summary
echo "=== Generation Summary ==="
echo "📊 Total proto files found: $total_files"
echo "✅ Successfully generated: $generated_files"
echo "❌ Failed: $((total_files - generated_files))"

if [ $generated_files -eq $total_files ] && [ $total_files -gt 0 ]; then
    echo "🎉 All proto files generated successfully!"
    exit 0
elif [ $total_files -eq 0 ]; then
    echo "⚠️  No proto files found in any subdirectory of proto/"
    exit 1
else
    echo "⚠️  Some files failed to generate"
    exit 1
fi