import { Module } from '@nestjs/common';
import { RecordProcessService } from './recordProcess.service';

@Module({
  imports: [ ],
  controllers: [ ],
  providers: [ RecordProcessService ],
  exports: [ RecordProcessService ]
})
export class RecordModule {}
