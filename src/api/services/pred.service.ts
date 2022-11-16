import { Injectable } from '@nestjs/common';
import { SBD } from 'src/dto/sbd.dto';
import { spawnSync } from 'child_process';
import { RecordWeekly } from 'src/dto/RecordWeekly';

@Injectable()
export class PredService {
  getPred(records: RecordWeekly[]): SBD {
    // get 24weeks record and make prediction
    try{
        console.log(__dirname)
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
