import { Module } from '@nestjs/common';

import { PrismaService } from 'src/prisma_services/prisma/prisma.service';
import { UserPrismaService } from 'src/prisma_services/user/user.prisma.service';

import { AuthApiController } from './auth.api.controller';
import { AuthApiService } from './services/auth.api.service';

@Module({
  imports: [],
  controllers: [ AuthApiController ],
  providers: [ 
    PrismaService,
    UserPrismaService,
    AuthApiService
  ],
}) export class AuthApiModule {}
