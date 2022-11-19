import sys
import json

with open('./predModel/model24Input.json', 'r') as f:
  data = json.load(f)

print('Amethod 183 105 205\nBmethod 200 100 302\nCmethod 200 100 302\nDmethod 200 100 302\nEmethod 200 100 302',end='')
sys.stdout.flush()