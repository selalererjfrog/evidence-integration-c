#!/bin/bash

# Convert Surefire XML test reports to JSON format
# Usage: ./convert-test-reports.sh [input_dir] [output_file]

set -e

# Default values
INPUT_DIR="${1:-target/surefire-reports}"
OUTPUT_FILE="${2:-test-evidence.json}"
REPOSITORY="${REPOSITORY:-unknown}"
COMMIT_SHA="${COMMIT_SHA:-unknown}"
BRANCH="${BRANCH:-unknown}"
TRIGGERED_BY="${TRIGGERED_BY:-unknown}"

echo "🔍 Converting test reports from XML to JSON..."
echo "📁 Input directory: $INPUT_DIR"
echo "📄 Output file: $OUTPUT_FILE"

# Function to escape JSON strings
escape_json() {
    echo "$1" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g' | sed 's/\n/\\n/g' | sed 's/\r/\\r/g' | sed 's/\t/\\t/g'
}

# Function to extract XML attribute value
extract_xml_attr() {
    local file="$1"
    local attr="$2"
    local default="$3"
    
    if command -v xmllint >/dev/null 2>&1; then
        # Use xmllint and take only the first line to avoid multi-line output
        xmllint --xpath "string(/testsuite/@$attr)" "$file" 2>/dev/null | head -n1 | tr -d '\n' || echo "$default"
    else
        # Fallback to grep - take only the first match
        grep -o "$attr=\"[^\"]*\"" "$file" | head -n1 | cut -d'"' -f2 || echo "$default"
    fi
}

# Function to extract test case details
extract_test_cases() {
    local file="$1"
    local class_name="$2"
    
    if command -v xmllint >/dev/null 2>&1; then
        # Extract test cases using xmllint
        xmllint --xpath "//testcase" "$file" 2>/dev/null | while read -r line; do
            if [[ $line =~ name=\"([^\"]+)\" ]]; then
                test_name="${BASH_REMATCH[1]}"
                # Extract status (passed/failed/skipped)
                if [[ $line =~ \<failure ]]; then
                    status="failed"
                elif [[ $line =~ \<skipped ]]; then
                    status="skipped"
                else
                    status="passed"
                fi
                echo "{\"class\":\"$class_name\",\"name\":\"$test_name\",\"status\":\"$status\"}"
            fi
        done
    else
        # Fallback to grep for test cases
        grep -o 'name="[^"]*"' "$file" | cut -d'"' -f2 | while read -r test_name; do
            echo "{\"class\":\"$class_name\",\"name\":\"$test_name\",\"status\":\"passed\"}"
        done
    fi
}

# Function to safely add floating point numbers
add_floating_point() {
    local a="$1"
    local b="$2"
    
    # Extract integer parts only to avoid complex arithmetic
    local a_int=${a%.*}
    local b_int=${b%.*}
    
    # Handle empty strings
    if [ -z "$a_int" ]; then a_int=0; fi
    if [ -z "$b_int" ]; then b_int=0; fi
    
    # Simple integer addition
    local result=$((a_int + b_int))
    
    # Return as decimal for consistency
    echo "$result.000"
}

# Function to safely calculate percentage
calculate_percentage() {
    local numerator="$1"
    local denominator="$2"
    
    if [ "$denominator" -eq 0 ]; then
        echo "0"
    elif command -v awk >/dev/null 2>&1; then
        echo "$numerator $denominator" | awk '{printf "%.2f", ($1 * 100) / $2}'
    elif command -v bc >/dev/null 2>&1; then
        echo "scale=2; $numerator * 100 / $denominator" | bc -l 2>/dev/null || echo "0"
    else
        # Fallback: simple integer calculation
        echo $((numerator * 100 / denominator))
    fi
}

# Initialize JSON structure
cat > "$OUTPUT_FILE" << EOF
{
  "test_summary": {
    "total_tests": 0,
    "passed_tests": 0,
    "failed_tests": 0,
    "skipped_tests": 0,
    "test_duration": 0,
    "success_rate": 0
  },
  "test_details": [],
  "test_cases": []
}
EOF

# Check if input directory exists
if [ ! -d "$INPUT_DIR" ]; then
    echo "⚠️ Input directory '$INPUT_DIR' not found"
    echo "📄 Created empty test evidence JSON: $OUTPUT_FILE"
    exit 0
fi

# Find all XML files
xml_files=$(find "$INPUT_DIR" -name "*.xml" -type f)

if [ -z "$xml_files" ]; then
    echo "⚠️ No XML test reports found in '$INPUT_DIR'"
    echo "📄 Created empty test evidence JSON: $OUTPUT_FILE"
    exit 0
fi

# Initialize counters
total_tests=0
passed_tests=0
failed_tests=0
skipped_tests=0
total_duration=0
test_details_json=""
test_cases_json=""

# Process each XML file
echo "📊 Processing XML test reports..."
for xml_file in $xml_files; do
    echo "  📄 Processing: $(basename "$xml_file")"
    
    # Extract test class name
    class_name=$(basename "$xml_file" .xml)
    
    # Extract test statistics
    tests=$(extract_xml_attr "$xml_file" "tests" "0")
    failures=$(extract_xml_attr "$xml_file" "failures" "0")
    skipped=$(extract_xml_attr "$xml_file" "skipped" "0")
    time_raw=$(extract_xml_attr "$xml_file" "time" "0")
    
    echo "    Raw values - tests: '$tests', failures: '$failures', skipped: '$skipped', time: '$time_raw'"
    
    # Ensure all values are integers
    tests=${tests%.*}
    failures=${failures%.*}
    skipped=${skipped%.*}
    time=${time_raw%.*}
    
    # Handle empty strings
    if [ -z "$tests" ]; then tests=0; fi
    if [ -z "$failures" ]; then failures=0; fi
    if [ -z "$skipped" ]; then skipped=0; fi
    if [ -z "$time" ]; then time=0; fi
    
    echo "    Processed values - tests: $tests, failures: $failures, skipped: $skipped, time: $time"
    
    # Calculate passed tests
    passed=$((tests - failures - skipped))
    
    # Update totals
    total_tests=$((total_tests + tests))
    passed_tests=$((passed_tests + passed))
    failed_tests=$((failed_tests + failures))
    skipped_tests=$((skipped_tests + skipped))
    
    # Add duration (extract integer part only)
    time_int=${time%.*}
    if [ -z "$time_int" ]; then time_int=0; fi
    total_duration=$((total_duration + time_int))
    
    # Add test class details to JSON
    if [ -n "$test_details_json" ]; then
        test_details_json="$test_details_json,"
    fi
    test_details_json="$test_details_json
    {
      \"class_name\": \"$(escape_json "$class_name")\",
      \"tests\": $tests,
      \"passed\": $passed,
      \"failed\": $failures,
      \"skipped\": $skipped,
      \"duration\": $time
    }"
    
    # Extract individual test cases
    test_cases=$(extract_test_cases "$xml_file" "$class_name")
    if [ -n "$test_cases" ]; then
        while IFS= read -r test_case; do
            if [ -n "$test_cases_json" ]; then
                test_cases_json="$test_cases_json,"
            fi
            test_cases_json="$test_cases_json
      $test_case"
        done <<< "$test_cases"
    fi
done

# Calculate success rate
success_rate=$(calculate_percentage "$passed_tests" "$total_tests")

# Create final JSON file
cat > "$OUTPUT_FILE" << EOF
{
  "test_summary": {
    "total_tests": $total_tests,
    "passed_tests": $passed_tests,
    "failed_tests": $failed_tests,
    "skipped_tests": $skipped_tests,
    "test_duration": $total_duration,
    "success_rate": $success_rate
  },
  "test_details": [$test_details_json
  ],
  "test_cases": [$test_cases_json
  ]
}
EOF

echo "✅ Test evidence JSON created successfully: $OUTPUT_FILE"
echo "📊 Test Summary:"
echo "   Total Tests: $total_tests"
echo "   Passed: $passed_tests"
echo "   Failed: $failed_tests"
echo "   Skipped: $skipped_tests"
echo "   Success Rate: ${success_rate}%"
echo "   Duration: ${total_duration}s"

# Display file size
file_size=$(wc -c < "$OUTPUT_FILE")
echo "📄 File size: ${file_size} bytes"

# Validate JSON format
if command -v jq >/dev/null 2>&1; then
    if jq empty "$OUTPUT_FILE" 2>/dev/null; then
        echo "✅ JSON validation passed"
    else
        echo "❌ JSON validation failed"
        exit 1
    fi
elif command -v python3 >/dev/null 2>&1; then
    if python3 -m json.tool "$OUTPUT_FILE" >/dev/null 2>&1; then
        echo "✅ JSON validation passed"
    else
        echo "❌ JSON validation failed"
        exit 1
    fi
else
    echo "⚠️ No JSON validator available, skipping validation"
fi

# Generate Markdown report
echo "📝 Generating Markdown test report..."
MARKDOWN_FILE="${OUTPUT_FILE%.json}.md"

# Function to get status emoji
get_status_emoji() {
    local status="$1"
    case "$status" in
        "passed") echo "✅" ;;
        "failed") echo "❌" ;;
        "skipped") echo "⏭️" ;;
        *) echo "❓" ;;
    esac
}

# Function to get summary emoji
get_summary_emoji() {
    local success_rate="$1"
    if [ "$(echo "$success_rate >= 90" | bc -l 2>/dev/null || echo "0")" = "1" ]; then
        echo "🎉"
    elif [ "$(echo "$success_rate >= 70" | bc -l 2>/dev/null || echo "0")" = "1" ]; then
        echo "⚠️"
    else
        echo "🚨"
    fi
}

# Create markdown content
cat > "$MARKDOWN_FILE" << EOF
# Test Results Report

## 📊 Test Summary

$(get_summary_emoji "$success_rate") **Overall Success Rate: ${success_rate}%**

| Metric | Value |
|--------|-------|
| **Total Tests** | $total_tests |
| **Passed** | ✅ $passed_tests |
| **Failed** | ❌ $failed_tests |
| **Skipped** | ⏭️ $skipped_tests |
| **Duration** | ⏱️ ${total_duration}s |

## 📋 Test Details by Class

EOF

# Add test details table
if [ -n "$test_details_json" ]; then
    echo "| Class | Tests | Passed | Failed | Skipped | Duration |" >> "$MARKDOWN_FILE"
    echo "|-------|-------|--------|--------|---------|----------|" >> "$MARKDOWN_FILE"
    
    # Parse test_details_json to create table rows
    echo "$test_details_json" | while IFS= read -r line; do
        if [[ $line =~ \"class_name\":\ \"([^\"]+)\" ]]; then
            class_name="${BASH_REMATCH[1]}"
        elif [[ $line =~ \"tests\":\ ([0-9]+) ]]; then
            tests="${BASH_REMATCH[1]}"
        elif [[ $line =~ \"passed\":\ ([0-9]+) ]]; then
            passed="${BASH_REMATCH[1]}"
        elif [[ $line =~ \"failed\":\ ([0-9]+) ]]; then
            failed="${BASH_REMATCH[1]}"
        elif [[ $line =~ \"skipped\":\ ([0-9]+) ]]; then
            skipped="${BASH_REMATCH[1]}"
        elif [[ $line =~ \"duration\":\ ([0-9]+) ]]; then
            duration="${BASH_REMATCH[1]}"
            # Output the complete row
            echo "| \`$class_name\` | $tests | ✅ $passed | ❌ $failed | ⏭️ $skipped | ${duration}s |" >> "$MARKDOWN_FILE"
        fi
    done
else
    echo "*No test details available*" >> "$MARKDOWN_FILE"
fi

# Add footer
cat >> "$MARKDOWN_FILE" << EOF

---

**Report generated on:** $(date)  

EOF

echo "✅ Markdown report created: $MARKDOWN_FILE"
