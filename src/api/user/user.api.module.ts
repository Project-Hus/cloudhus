import { Module } from '@nestjs/common';
import { PrismaService } from 'src/prisma_services/prisma/prisma.service';
import { UserPrismaService } from 'src/prisma_services/user/user.prisma.service';
import { UserApiController } from './user.api.controller';

@Module({
  imports: [],
  controllers: [ UserApiController ],
  providers: [ 
    PrismaService,
    UserPrismaService 
  ],
}) export class UserApiModule {}
