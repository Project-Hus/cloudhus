import { Injectable } from '@nestjs/common';
import { SBD } from 'src/dto/sbd.dto';
import { spawnSync } from 'child_process';
import { RecordWeekly } from 'src/dto/RecordWeekly';
import { RecordFixed } from 'src/dto/RecordFixed';
import { RecordService } from 'utils/record/record.service';

@Injectable()
export class PredService {
  constructor(
    private readonly recordService : RecordService
  ) {}

  getPred(recordFixed: RecordFixed, records: RecordWeekly[]): SBD {
    try{
      this.recordService.processRecords(recordFixed, records);
      const pythonProcess = spawnSync('python',["./predModel/model24.py"]);
      const sbd = pythonProcess.stdout.toString().split(' ').map(Number); 
      return {
          squat: sbd[0],
          benchpress: sbd[1],
          deadlift: sbd[2],
      };
    } catch (error) {
      return {squat:0, benchpress:0, deadlift:0};
    }
  }
}
