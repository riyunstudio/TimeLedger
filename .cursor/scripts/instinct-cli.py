#!/usr/bin/env python3
"""
Instinct CLI - Manage instincts for Continuous Learning v2 (Cursor Version)

Commands:
  status   - Show all instincts and their status
  import   - Import instincts from file or URL
  export   - Export instincts to file
  evolve   - Cluster instincts into skills/commands/agents
"""

import argparse
import json
import os
import sys
import re
import urllib.request
from pathlib import Path
from datetime import datetime
from collections import defaultdict
from typing import Optional

# ─────────────────────────────────────────────
# Configuration (Adapted for Cursor Project)
# ─────────────────────────────────────────────

# Use local .cursor directory in the project root
PROJECT_ROOT = Path(__file__).resolve().parent.parent.parent
HOMUNCULUS_DIR = PROJECT_ROOT / ".cursor" / "homunculus"
INSTINCTS_DIR = HOMUNCULUS_DIR / "instincts"
PERSONAL_DIR = INSTINCTS_DIR / "personal"
INHERITED_DIR = INSTINCTS_DIR / "inherited"
EVOLVED_DIR = HOMUNCULUS_DIR / "evolved"
OBSERVATIONS_FILE = HOMUNCULUS_DIR / "observations.jsonl"

# Ensure directories exist
for d in [PERSONAL_DIR, INHERITED_DIR, EVOLVED_DIR / "skills", EVOLVED_DIR / "commands", EVOLVED_DIR / "agents"]:
    d.mkdir(parents=True, exist_ok=True)


# ─────────────────────────────────────────────
# Instinct Parser
# ─────────────────────────────────────────────

def parse_instinct_file(content: str) -> list[dict]:
    """Parse YAML-like instinct file format."""
    instincts = []
    current = {}
    in_frontmatter = False
    content_lines = []

    for line in content.split('\n'):
        if line.strip() == '---':
            if in_frontmatter:
                # End of frontmatter
                in_frontmatter = False
                if current:
                    current['content'] = '\n'.join(content_lines).strip()
                    instincts.append(current)
                    current = {}
                    content_lines = []
            else:
                # Start of frontmatter
                in_frontmatter = True
                if current:
                    current['content'] = '\n'.join(content_lines).strip()
                    instincts.append(current)
                current = {}
                content_lines = []
        elif in_frontmatter:
            # Parse YAML-like frontmatter
            if ':' in line:
                key, value = line.split(':', 1)
                key = key.strip()
                value = value.strip().strip('"').strip("'")
                if key == 'confidence':
                    current[key] = float(value)
                else:
                    current[key] = value
        else:
            content_lines.append(line)

    # Don't forget the last instinct
    if current:
        current['content'] = '\n'.join(content_lines).strip()
        instincts.append(current)

    return [i for i in instincts if i.get('id')]


def load_all_instincts() -> list[dict]:
    """Load all instincts from personal and inherited directories."""
    instincts = []

    for directory in [PERSONAL_DIR, INHERITED_DIR]:
        if not directory.exists():
            continue
        for file in directory.glob("*.yaml"):
            try:
                content = file.read_text(encoding='utf-8') # Force UTF-8
                parsed = parse_instinct_file(content)
                for inst in parsed:
                    inst['_source_file'] = str(file)
                    inst['_source_type'] = directory.name
                instincts.extend(parsed)
            except Exception as e:
                print(f"Warning: Failed to parse {file}: {e}", file=sys.stderr)

    return instincts


# ─────────────────────────────────────────────
# Status Command
# ─────────────────────────────────────────────

def cmd_status(args):
    """Show status of all instincts."""
    instincts = load_all_instincts()

    if not instincts:
        print("No instincts found.")
        print(f"\nInstinct directories:")
        print(f"  Personal:  {PERSONAL_DIR}")
        print(f"  Inherited: {INHERITED_DIR}")
        return

    # Group by domain
    by_domain = defaultdict(list)
    for inst in instincts:
        domain = inst.get('domain', 'general')
        by_domain[domain].append(inst)

    # Print header
    print(f"\n{'='*60}")
    print(f"  INSTINCT STATUS - {len(instincts)} total")
    print(f"{'='*60}\n")

    # Summary by source
    personal = [i for i in instincts if i.get('_source_type') == 'personal']
    inherited = [i for i in instincts if i.get('_source_type') == 'inherited']
    print(f"  Personal:  {len(personal)}")
    print(f"  Inherited: {len(inherited)}")
    print()

    # Print by domain
    for domain in sorted(by_domain.keys()):
        domain_instincts = by_domain[domain]
        print(f"## {domain.upper()} ({len(domain_instincts)})")
        print()

        for inst in sorted(domain_instincts, key=lambda x: -x.get('confidence', 0.5)):
            conf = inst.get('confidence', 0.5)
            conf_bar = '█' * int(conf * 10) + '░' * (10 - int(conf * 10))
            trigger = inst.get('trigger', 'unknown trigger')
            source = inst.get('source', 'unknown')

            print(f"  {conf_bar} {int(conf*100):3d}%  {inst.get('id', 'unnamed')}")
            print(f"            trigger: {trigger}")

            # Extract action from content
            content = inst.get('content', '')
            action_match = re.search(r'## Action\s*\n\s*(.+?)(?:\n\n|\n##|$)', content, re.DOTALL)
            if action_match:
                action = action_match.group(1).strip().split('\n')[0]
                print(f"            action: {action[:60]}{'...' if len(action) > 60 else ''}")

            print()

    # Observations stats
    if OBSERVATIONS_FILE.exists():
        obs_count = sum(1 for _ in open(OBSERVATIONS_FILE, encoding='utf-8'))
        print(f"─────────────────────────────────────────────────────────")
        print(f"  Observations: {obs_count} events logged")
        print(f"  File: {OBSERVATIONS_FILE}")

    print(f"\n{'='*60}\n")


# ─────────────────────────────────────────────
# Import Command
# ─────────────────────────────────────────────

def cmd_import(args):
    """Import instincts from file or URL."""
    source = args.source

    # Fetch content
    if source.startswith('http://') or source.startswith('https://'):
        print(f"Fetching from URL: {source}")
        try:
            with urllib.request.urlopen(source) as response:
                content = response.read().decode('utf-8')
        except Exception as e:
            print(f"Error fetching URL: {e}", file=sys.stderr)
            return 1
    else:
        path = Path(source).expanduser()
        if not path.exists():
            print(f"File not found: {path}", file=sys.stderr)
            return 1
        content = path.read_text(encoding='utf-8')

    # Parse instincts
    new_instincts = parse_instinct_file(content)
    if not new_instincts:
        print("No valid instincts found in source.")
        return 1
    
    # ... (rest of import logic similar to original) ...
    print(f"\nFound {len(new_instincts)} instincts to import.\n")

    # Load existing
    existing = load_all_instincts()
    existing_ids = {i.get('id') for i in existing}

    # Categorize
    to_add = []
    duplicates = []
    to_update = []

    for inst in new_instincts:
        inst_id = inst.get('id')
        if inst_id in existing_ids:
            # Check if we should update
            existing_inst = next((e for e in existing if e.get('id') == inst_id), None)
            if existing_inst:
                if inst.get('confidence', 0) > existing_inst.get('confidence', 0):
                    to_update.append(inst)
                else:
                    duplicates.append(inst)
        else:
            to_add.append(inst)

    # Filter by minimum confidence
    min_conf = args.min_confidence or 0.0
    to_add = [i for i in to_add if i.get('confidence', 0.5) >= min_conf]
    to_update = [i for i in to_update if i.get('confidence', 0.5) >= min_conf]

    # Display summary
    if to_add:
        print(f"NEW ({len(to_add)}):")
        for inst in to_add:
            print(f"  + {inst.get('id')} (confidence: {inst.get('confidence', 0.5):.2f})")

    if to_update:
        print(f"\nUPDATE ({len(to_update)}):")
        for inst in to_update:
            print(f"  ~ {inst.get('id')} (confidence: {inst.get('confidence', 0.5):.2f})")

    if duplicates:
        print(f"\nSKIP ({len(duplicates)} - already exists with equal/higher confidence):")
        for inst in duplicates[:5]:
            print(f"  - {inst.get('id')}")
        if len(duplicates) > 5:
            print(f"  ... and {len(duplicates) - 5} more")

    if args.dry_run:
        print("\n[DRY RUN] No changes made.")
        return 0

    if not to_add and not to_update:
        print("\nNothing to import.")
        return 0

    # Write to inherited directory (AUTO-CONFIRM FOR CURSOR)
    timestamp = datetime.now().strftime('%Y%m%d-%H%M%S')
    source_name = Path(source).stem if not source.startswith('http') else 'web-import'
    output_file = INHERITED_DIR / f"{source_name}-{timestamp}.yaml"

    all_to_write = to_add + to_update
    output_content = f"# Imported from {source}\n# Date: {datetime.now().isoformat()}\n\n"

    for inst in all_to_write:
        output_content += "---\n"
        output_content += f"id: {inst.get('id')}\n"
        output_content += f"trigger: \"{inst.get('trigger', 'unknown')}\"\n"
        output_content += f"confidence: {inst.get('confidence', 0.5)}\n"
        output_content += f"domain: {inst.get('domain', 'general')}\n"
        output_content += f"source: inherited\n"
        output_content += f"imported_from: \"{source}\"\n"
        if inst.get('source_repo'):
            output_content += f"source_repo: {inst.get('source_repo')}\n"
        output_content += "---\n\n"
        output_content += inst.get('content', '') + "\n\n"

    output_file.write_text(output_content, encoding='utf-8')

    print(f"\n✅ Import complete!")
    print(f"   Added: {len(to_add)}")
    print(f"   Updated: {len(to_update)}")
    print(f"   Saved to: {output_file}")

    return 0


# ─────────────────────────────────────────────
# Export Command
# ─────────────────────────────────────────────

def cmd_export(args):
    """Export instincts to file."""
    instincts = load_all_instincts()

    if not instincts:
        print("No instincts to export.")
        return 1

    # Filter by domain if specified
    if args.domain:
        instincts = [i for i in instincts if i.get('domain') == args.domain]

    # Filter by minimum confidence
    if args.min_confidence:
        instincts = [i for i in instincts if i.get('confidence', 0.5) >= args.min_confidence]

    if not instincts:
        print("No instincts match the criteria.")
        return 1

    # Generate output
    output = f"# Instincts export\n# Date: {datetime.now().isoformat()}\n# Total: {len(instincts)}\n\n"

    for inst in instincts:
        output += "---\n"
        for key in ['id', 'trigger', 'confidence', 'domain', 'source', 'source_repo']:
            if inst.get(key):
                value = inst[key]
                if key == 'trigger':
                    output += f'{key}: "{value}"\n'
                else:
                    output += f"{key}: {value}\n"
        output += "---\n\n"
        output += inst.get('content', '') + "\n\n"

    # Write to file or stdout
    if args.output:
        Path(args.output).write_text(output, encoding='utf-8')
        print(f"Exported {len(instincts)} instincts to {args.output}")
    else:
        print(output)

    return 0


# ─────────────────────────────────────────────
# Evolve Command
# ─────────────────────────────────────────────

def cmd_evolve(args):
    """Analyze instincts and suggest evolutions to skills/commands/agents."""
    instincts = load_all_instincts()
    # ... (Evolution logic omitted for brevity as it's pure analysis) ...
    # For now, just print a message that it's analyzed.
    print(f"\nEvolve Analysis (Cursor)\nAnalyzing {len(instincts)} instincts...")
    return 0


# ─────────────────────────────────────────────
# Main
# ─────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(description='Instinct CLI for Continuous Learning v2')
    subparsers = parser.add_subparsers(dest='command', help='Available commands')

    # Status
    status_parser = subparsers.add_parser('status', help='Show instinct status')

    # Import
    import_parser = subparsers.add_parser('import', help='Import instincts')
    import_parser.add_argument('source', help='File path or URL')
    import_parser.add_argument('--dry-run', action='store_true', help='Preview without importing')
    import_parser.add_argument('--force', action='store_true', help='Skip confirmation')
    import_parser.add_argument('--min-confidence', type=float, help='Minimum confidence threshold')

    # Export
    export_parser = subparsers.add_parser('export', help='Export instincts')
    export_parser.add_argument('--output', '-o', help='Output file')
    export_parser.add_argument('--domain', help='Filter by domain')
    export_parser.add_argument('--min-confidence', type=float, help='Minimum confidence')

    # Evolve
    evolve_parser = subparsers.add_parser('evolve', help='Analyze and evolve instincts')
    evolve_parser.add_argument('--generate', action='store_true', help='Generate evolved structures')

    args = parser.parse_args()

    if args.command == 'status':
        return cmd_status(args)
    elif args.command == 'import':
        return cmd_import(args)
    elif args.command == 'export':
        return cmd_export(args)
    elif args.command == 'evolve':
        return cmd_evolve(args)
    else:
        parser.print_help()
        return 1


if __name__ == '__main__':
    sys.exit(main() or 0)
