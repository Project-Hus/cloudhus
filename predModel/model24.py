import sys
import json

#import tensorflow as tf

import preprcs

with open('./predModel/model24Input.json', 'r') as f:
  data = json.load(f)

data = preprcs.preprcs(data)

with open('./predModel/model24Output.json', 'w') as f:
  json.dump(data, f)

# Prediction for every methods 
# ============================

# ============================

print('Amethod 183 105 205\nBmethod 200 100 302\nCmethod 200 100 302\nDmethod 200 100 302\nEmethod 200 100 302',end='')
sys.stdout.flush()