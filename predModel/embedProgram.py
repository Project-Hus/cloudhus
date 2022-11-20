PROGRAM_ID = [
    'Rest',
]

def space_generator_2d(length = 5, base=[0,-1,1]):
  space = []
  for i in base:
    space.append([i])
  new_space = []
  for i in range(length) :
    for j in space:
      for k in base:
        new_space.append(j+[k])
    if i!= length-1:
      space = new_space[:]
      new_space = []
  return space