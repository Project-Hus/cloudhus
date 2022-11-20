'''
{"sex":true,
"age":20,
"height":183,
"arm_length":"medium"
,"leg_length":"long",
"weight":83,
"fat_rate":20,
"program":11,
"squat":120,
"benchpress":60,
"deadlift":140}
'''

import pandas as pd
import numpy as np

def preprcs(raw):
  df = pd.DataFrame(raw)

  df['sex'].replace([True, False], [1, 0], inplace=True)

  df['arm_length'].replace(['short','medium','long'], [0.25,0.5,0.75], inplace=True)
  df['leg_length'].replace(['short','medium','long'], [0.25,0.5,0.75], inplace=True)

  df['program'] = df['program'][:-1].append(pd.Series(0)).to_numpy()
  
  return df