import { Prisma, PrismaClient } from '@prisma/client'
const prisma = new PrismaClient()
async function main() {
    Prisma
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
    const initial_programs = [
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
    ]
    for (const e of initial_programs) {
    await prisma.trainingProgram.create({
            data : {

            }
        })
    }

    await prisma.
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