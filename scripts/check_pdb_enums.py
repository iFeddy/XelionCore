import sys

with open('pdo/Zone.pdb', 'rb') as f:
    data = f.read()

# Look for strings like NC_ or PROTO_NC_
import re
matches = set(re.findall(rb'NC_[A-Z0-9_]+', data))
print(f"Total NC_ strings in Zone.pdb: {len(matches)}")

sample = sorted([m.decode('latin1') for m in matches])[:30]
for s in sample:
    print(" ", s)
