import { Module } from '@nestjs/common';
import { RecordProcessService } from 'utils/recordProcess/recordProcess.service';
import { ApiController } from './api.controller';
import { PredService } from './services/pred.service';

@Module({
  imports: [],
  controllers: [ ApiController ],
  providers: [ PredService, RecordProcessService ],
})
export class ApiModule {}
