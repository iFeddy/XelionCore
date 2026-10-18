import re

with open('docs/TablasLegacy.md', 'r') as f:
    content = f.read()

matches = re.findall(r'\|\s*`World00_Character\.bak`\s*\|\s*\*\*`?([^`*]+)`?\*\*\s*\|\s*([^|]+)\|', content)
print(f"Total tables found: {len(matches)}")
for tname, cols in sorted(matches, key=lambda x: x[0]):
    print(f"=== {tname} ===")
    print(cols.strip())
