import { Module } from '@nestjs/common';
import { PrismaService } from 'src/prisma_services/prisma/prisma.service';
import { RecordPrismaService } from 'src/prisma_services/record/record.prisma.service';
import { UserPrismaService } from 'src/prisma_services/user/user.prisma.service';
import { RecordProcessService } from 'utils/recordProcess/recordProcess.service';
import { PredApiController } from './api.pred.controller';
import { PredApiService } from './services/api.pred.service';

@Module({
  imports: [],
  controllers: [ PredApiController ],
  providers: [ 
    PredApiService, RecordProcessService, 
    PrismaService, RecordPrismaService,
    UserPrismaService ],
}) export class PredApiModule {}