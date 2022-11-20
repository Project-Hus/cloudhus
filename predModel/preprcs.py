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

from recScaler import recScaler, recDescaler
from recScaler import aF, hF, wF, fF, sF, bF, dF

from embedProgram import space_generator_2d

def preprcs(raw):
  df = pd.DataFrame(raw)

  # sex embedding
  df['sex'].replace([True, False], [1, 0], inplace=True)
  # limbs embedding
  df['arm_length'].replace(['short','medium','long'], [0.25,0.5,0.75], inplace=True)
  df['leg_length'].replace(['short','medium','long'], [0.25,0.5,0.75], inplace=True)
  # program embedding
  df['program'] = df['program'][:-1].append(pd.Series(0)).to_numpy()

  # MinMax scaler for [ age, height, weight, fat_rate, s/b/d ]
  df['age'] = df['age'].map(aF)
  df['height'] = df['height'].map(hF)
  df['weight'] = df['weight'].map(wF)
  df['fat_rate'] = df['fat_rate'].map(fF)
  df['squat'] = df['squat'].map(sF)
  df['benchpress'] = df['benchpress'].map(bF)
  df['deadlift'] = df['deadlift'].map(dF)

  # 243 whole vector integer space
  v_space = np.array(space_generator_2d())
  program_vec = np.zeros((24,5))

  for i in df['program']:
    program_vec[i] = v_space[i]

  for i in range(5):
    df[f'program_{i}'] = program_vec[:,i]

  df.drop('program', inplace=True, axis=1)

  return df.to_numpy()