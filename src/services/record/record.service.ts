import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { Prisma, TrainingProgramRec, User } from '@prisma/client';

@Injectable()
export class RecordService {
    constructor(private prisma: PrismaService) {}
    /*
    async user(
        userWhereUniqueInput: Prisma.UserWhereUniqueInput,
    ): Promise<User | null> {
        return this.prisma.user.findUnique({
            where: userWhereUniqueInput,
        });
    }

    async users(params: {
        skip?: number;
        take?: number;
        cursor?: Prisma.UserWhereUniqueInput;
        where?: Prisma.UserWhereUniqueInput;
        orderBy?: Prisma.UserOrderByWithRelationInput;
    }): Promise<User[]> {
        const { skip, take, cursor, where, orderBy } = params;
        return this.prisma.user.findMany({
            skip,
            take,
            cursor,
            where,
            orderBy,
        });
    }
    */
/*
id Int @default(autoincrement()) @id
  program Program @relation(fields: [program_id], references: [id])
  program_id Int
  start DateTime
  end DateTime
  sq Float
  bp Float
  dl Float
*/
    async createProgramRecs(data:Prisma.TrainingProgramRecCreateManyInput[]): Promise<Prisma.BatchPayload> {
        return this.prisma.trainingProgramRec.createMany({
            data: data,
            //skipDuplicates: true,
        });
    }
    
    async updateUser(params: {
        where: Prisma.UserWhereUniqueInput;
        data: Prisma.UserUpdateInput;
    }): Promise<User> {
        const {where, data} = params;
        return this.prisma.user.update({
            data,
            where,
        });
    }

    async deleteUser(where: Prisma.UserWhereUniqueInput): Promise<User> {
        return this.prisma.user.delete({
            where,
        });
    }
}
