import { Module } from '@nestjs/common';
import { RecordService } from './record.service';

@Module({
  imports: [ ],
  controllers: [ ],
  providers: [ RecordService ],
  exports: [ RecordService ]
})
export class RecordModule {}
