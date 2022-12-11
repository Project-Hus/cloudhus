import { 
    Controller,
    Get,
    Param,
    Post,
    Body,
    Put,
    Delete
   } from '@nestjs/common';
import { RecordFixed } from 'src/dto/RecordFixed.dto';
import { RecordOutput } from 'src/dto/RecordOutput.dto';
import { RecordWeekly } from 'src/dto/RecordWeekly.dto';
import { PredApiService } from './services/api.pred.service';

  @Controller()
  export class PredApiController {
    constructor(
      private readonly predService: PredApiService,
      ) {}
  
    /**
     * API which suggests best 3 methods based on training records.
     * @param records Training records for 24 weeks
     * @returns Best 3 methods and the results
     */
    @Post('') // POST /api/pred
    async getPred(
      @Body('recordFixed') recordFixed: RecordFixed,
      @Body('recordWeekly') records: RecordWeekly[]
      ): Promise<RecordOutput[]> {
      return this.predService.getPred(recordFixed, records);
    } 
  }
  