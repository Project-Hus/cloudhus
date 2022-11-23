import { 
    Controller,
    Get,
    Param,
    Post,
    Body,
    Put,
    Delete
   } from '@nestjs/common';
import { RecordFixed } from 'src/dto/RecordFixed';
import { RecordOutput } from 'src/dto/RecordOutput';
import { RecordWeekly } from 'src/dto/RecordWeekly';
import { PredService } from './services/pred.service';

  @Controller('api')
  export class ApiController {
    constructor(
      private readonly predService: PredService,
      ) {}
  
    /**
     * API which suggests best 3 methods based on training records.
     * @param records Training records for 24 weeks
     * @returns Best 3 methods and the results
     */
    @Post('pred') // POST /api/pred
    async getPred(
      @Body('recordFixed') recordFixed: RecordFixed,
      @Body('recordWeekly') records: RecordWeekly[]): Promise<RecordOutput[]> {
      
      return this.predService.getPred(recordFixed, records);
    } 
  }
  