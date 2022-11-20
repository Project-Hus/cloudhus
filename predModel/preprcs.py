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

  # sex embedding
  df['sex'].replace([True, False], [1, 0], inplace=True)
  # limbs embedding
  df['arm_length'].replace(['short','medium','long'], [0.25,0.5,0.75], inplace=True)
  df['leg_length'].replace(['short','medium','long'], [0.25,0.5,0.75], inplace=True)

  df['program'] = df['program'][:-1].append(pd.Series(0)).to_numpy()

  return df

# MinMax scaler for [ age, height, weight, fat_rate, s/b/d ]
MIN = {
  'age':10,
  'height': 120,
  'weight': 40,
  'fat_rate': 0,
  'squat': 20,
  'benchpress': 20,
  'deadlift': 20
}
MAX = {
  'age':150,
  'height': 220,
  'weight': 250,
  'fat_rate': 50,
  'squat': 550,
  'benchpress': 500,
  'deadlift': 550
}
def recScaler(f_type, val):
  res = 0
  if f_type == 'age':
    return 

def recDescaler(f_type, race):