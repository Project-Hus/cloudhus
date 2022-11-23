import { Module } from '@nestjs/common';
import { PrismaService } from 'src/services/prisma/prisma.service';
import { RecordService } from 'src/services/record/record.service';
import { UserService } from 'src/services/user/user.service';
import { RecordProcessService } from 'utils/recordProcess/recordProcess.service';
import { ApiController } from './api.controller';
import { PredService } from './services/pred.service';

@Module({
  imports: [],
  controllers: [ ApiController ],
  providers: [ PredService, RecordProcessService, RecordService, PrismaService,
  UserService ],
})
export class ApiModule {}
