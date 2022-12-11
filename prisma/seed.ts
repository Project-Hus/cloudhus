// npx prisma db seed

import { Prisma, PrismaClient } from '@prisma/client'
const prisma = new PrismaClient()

function space_generator_2d(length = 1, base=[0,1]) : Promise<any[]> {
    return new Promise((resolve, reject) => {
        let space = [];
    for (let i of base)
      { space = [...space, [i]]; }
    let new_space = [];
    for (let i = 0; i < length -1 ; i++ ) {
      for ( let j of space ) {
        for ( let k of base )
          new_space = [...new_space, [...j, k]]
      }
      if ( i!== length-2) {
        space = [...new_space]
        new_space = []
      }
    }
    resolve(length > 1 ? new_space : space); 
      })
  }

async function main() {
    
    await prisma.user.create({
        data: {
           email_google: 'a@a.a',
           token_google: 'abcd',
           user_name: 'lifthus',
           password: '5486',
           age: 25,
           sex: 'male',
           height: 183,
           arm_length: 'medium',
           leg_length: 'medium'
        }
    });
    await prisma.trainingProgramType.create({
        data : { type: 'powerlifting' }
    });

    const vector_space = await space_generator_2d(5, [0,-1,1]);
    for ( const i in vector_space ) {
        await prisma.programVector.create({data:{
            c0: vector_space[i][0],
            c1: vector_space[i][1],
            c2: vector_space[i][2],
            c3: vector_space[i][3],
            c4: vector_space[i][4],
        }});
    }
    
    const initial_programs = [
        '휴식', // 0 (DB에서는 인덱스 +1 )
        'Kizen Powerlifting Peaking Program',
        'nSuns Programs',
        'Jim Wendler 5/3/1 Programs',
        'Calgary Barbell Programs',
        'Sheiko Programs',
        'Candito Program',
        'Juggernaut Method Base Template',
        'Greg Nuckols 28 Programs',
        'Beginner Powerlifting Programs',
        'Intermediate Powerlifting Programs',
        'Madcow 5x5 Program',
        'General 5x5 Program',
        '기타(etc)' // index 242
    ] 
    for (const i in initial_programs) {
        await prisma.trainingProgram.create({
                data : {
                    type_id: 1,
                    author: 1,
                    name: initial_programs[i],
                    description: 'no',
                    vector: Number(i) !== (initial_programs.length)-1 ? Number(i)+1 : 243
                }
            })
    }
}
main()
  .then(async () => {
    await prisma.$disconnect()
  })
  .catch(async (e) => {
    console.error(e)
    await prisma.$disconnect()
    process.exit(1)
  })