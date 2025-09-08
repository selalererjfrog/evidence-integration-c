#!/usr/bin/env python3
"""
Dependency Analysis Script for AI Translate Service
Analyzes package sizes and suggests optimizations to reduce disk space usage.
"""

import subprocess
import sys
import json
from pathlib import Path

def get_package_info(package_name):
    """Get package information using pip show"""
    try:
        result = subprocess.run(
            [sys.executable, "-m", "pip", "show", package_name],
            capture_output=True,
            text=True,
            check=True
        )
        info = {}
        for line in result.stdout.split('\n'):
            if ':' in line:
                key, value = line.split(':', 1)
                info[key.strip()] = value.strip()
        return info
    except subprocess.CalledProcessError:
        return None

def analyze_requirements(requirements_file):
    """Analyze requirements file and get package sizes"""
    packages = []
    
    with open(requirements_file, 'r') as f:
        for line in f:
            line = line.strip()
            if line and not line.startswith('#'):
                # Extract package name (remove version specifiers)
                package_name = line.split('==')[0].split('>=')[0].split('<=')[0].split('[')[0]
                packages.append(package_name)
    
    print(f"📦 Analyzing {len(packages)} packages from {requirements_file}")
    print("=" * 80)
    
    total_size = 0
    package_details = []
    
    for package in packages:
        info = get_package_info(package)
        if info:
            location = info.get('Location', 'Unknown')
            size_mb = 0
            
            # Calculate package size
            if location != 'Unknown':
                try:
                    import os
                    size_bytes = sum(
                        os.path.getsize(os.path.join(dirpath, filename))
                        for dirpath, dirnames, filenames in os.walk(location)
                        for filename in filenames
                    )
                    size_mb = size_bytes / (1024 * 1024)
                except:
                    size_mb = 0
            
            package_details.append({
                'name': package,
                'version': info.get('Version', 'Unknown'),
                'size_mb': size_mb,
                'location': location
            })
            total_size += size_mb
    
    # Sort by size (largest first)
    package_details.sort(key=lambda x: x['size_mb'], reverse=True)
    
    print(f"{'Package':<25} {'Version':<15} {'Size (MB)':<12} {'Location'}")
    print("-" * 80)
    
    for pkg in package_details:
        print(f"{pkg['name']:<25} {pkg['version']:<15} {pkg['size_mb']:<12.1f} {pkg['location']}")
    
    print("-" * 80)
    print(f"{'TOTAL':<25} {'':<15} {total_size:<12.1f}")
    
    return package_details, total_size

def suggest_optimizations(package_details):
    """Suggest optimizations based on package analysis"""
    print("\n🔍 OPTIMIZATION SUGGESTIONS")
    print("=" * 80)
    
    # Large packages that could be optimized
    large_packages = [pkg for pkg in package_details if pkg['size_mb'] > 100]
    
    if large_packages:
        print("📦 Large packages (>100MB) that could be optimized:")
        for pkg in large_packages:
            print(f"  - {pkg['name']}: {pkg['size_mb']:.1f}MB")
    
    # Specific suggestions
    suggestions = []
    
    # Check for torch
    torch_pkg = next((pkg for pkg in package_details if 'torch' in pkg['name'].lower()), None)
    if torch_pkg and torch_pkg['size_mb'] > 500:
        suggestions.append({
            'package': 'torch',
            'suggestion': 'Consider using torch-cpu only or torch-lite for smaller footprint',
            'potential_savings': '300-500MB'
        })
    
    # Check for transformers
    transformers_pkg = next((pkg for pkg in package_details if 'transformers' in pkg['name'].lower()), None)
    if transformers_pkg and transformers_pkg['size_mb'] > 200:
        suggestions.append({
            'package': 'transformers',
            'suggestion': 'Consider using specific model components only',
            'potential_savings': '50-100MB'
        })
    
    # Check for development dependencies in production
    dev_packages = ['pytest', 'black', 'flake8', 'mypy', 'pytest-cov']
    for dev_pkg in dev_packages:
        pkg = next((pkg for pkg in package_details if dev_pkg in pkg['name'].lower()), None)
        if pkg:
            suggestions.append({
                'package': dev_pkg,
                'suggestion': 'Move to requirements-dev.txt - not needed in production',
                'potential_savings': f"{pkg['size_mb']:.1f}MB"
            })
    
    if suggestions:
        print("\n💡 Specific optimization suggestions:")
        for suggestion in suggestions:
            print(f"  - {suggestion['package']}: {suggestion['suggestion']}")
            print(f"    Potential savings: {suggestion['potential_savings']}")
    else:
        print("✅ No major optimization opportunities found")
    
    return suggestions

def main():
    """Main analysis function"""
    print("🚀 AI Translate Service - Dependency Analysis")
    print("=" * 80)
    
    # Analyze current requirements
    current_packages, current_size = analyze_requirements('requirements.txt')
    
    # Analyze optimized requirements if it exists
    if Path('requirements-optimized.txt').exists():
        print("\n" + "=" * 80)
        optimized_packages, optimized_size = analyze_requirements('requirements-optimized.txt')
        
        print(f"\n📊 COMPARISON")
        print(f"Current requirements: {current_size:.1f}MB")
        print(f"Optimized requirements: {optimized_size:.1f}MB")
        print(f"Potential savings: {current_size - optimized_size:.1f}MB ({(current_size - optimized_size)/current_size*100:.1f}%)")
    
    # Generate suggestions
    suggestions = suggest_optimizations(current_packages)
    
    print(f"\n📋 SUMMARY")
    print(f"Total packages analyzed: {len(current_packages)}")
    print(f"Total size: {current_size:.1f}MB")
    print(f"Optimization suggestions: {len(suggestions)}")
    
    if suggestions:
        total_savings = sum(
            float(s['potential_savings'].replace('MB', '')) 
            for s in suggestions 
            if 'MB' in s['potential_savings']
        )
        print(f"Potential total savings: {total_savings:.1f}MB")

if __name__ == "__main__":
    main()
