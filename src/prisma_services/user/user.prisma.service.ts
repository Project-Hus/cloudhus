import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { User, Prisma } from '@prisma/client';

@Injectable()
export class UserPrismaService {
    constructor(private prisma: PrismaService) {}

    async createUser(data: Prisma.UserCreateInput): Promise<User> {
        return this.prisma.user.create({
            data,
        });
    }

    // get latest [take] users. default 1 person.
    async latestUsers(params: {
        //skip?: number;
        take?: number;
        //cursor?: Prisma.UserWhereUniqueInput;
        //where?: Prisma.UserWhereUniqueInput;
        //orderBy?: Prisma.UserOrderByWithRelationInput;
    } = { take:1}): Promise</*User[]*/object[]> {
        const { /* skip, */ take, /* cursor, where, orderBy */ } = params;
        return this.prisma.user.findMany({
            //skip,
            take,
            //cursor,
            //where,
            select: { id: true },
            orderBy : { id: 'desc'},
        });
    }

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


