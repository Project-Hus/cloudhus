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
# MinMax scaler for [ age, height, weight, fat_rate, s/b/d ]
def recScaler(f_type, val):
  return (val-MIN[f_type])/(MAX[f_type]-MIN[f_type])

def recDescaler(f_type, val):
  return val * (MAX[f_type]-MIN[f_type]) + MIN[f_type]

def aF(v) :
  return recScaler('age', v)
def hF(v) :
  return recScaler('height', v)
def wF(v) :
  return recScaler('weight', v)
def fF(v) :
  return recScaler('fat_rate', v)
def sF(v) :
  return recScaler('squat', v)
def bF(v) :
  return recScaler('benchpress', v)
def dF(v) :
  return recScaler('deadlift', v)