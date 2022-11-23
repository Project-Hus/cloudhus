import sys
import json

import preprcs
from LSTM24 import LSTM24

with open('./predModel/model24Input.json', 'r') as f:
  data = json.load(f)

data = preprcs.preprcs(data)

#with open('./predModel/model24Output.txt', 'w') as f:
#  np.savetxt(f, data)

print(LSTM24(data), end='')
sys.stdout.flush()