import { Injectable } from '@nestjs/common';
import { spawnSync } from 'child_process';
import { writeFileSync } from 'fs';

/* DTOs */
import { RecordWeekly } from 'src/dto/RecordWeekly';
import { RecordFixed } from 'src/dto/RecordFixed';
import { RecordInput } from 'src/dto/RecordInput';
import { RecordOutput } from 'src/dto/RecordOutput';

/* Service */
import { RecordProcessService } from 'utils/recordProcess/recordProcess.service';
import { RecordService } from 'src/services/record/record.service';
import { Prisma, TrainingProgramRec } from '@prisma/client';
import { UserService } from 'src/services/user/user.service';

/**
 * API Service for training program suggestion and prediction.
 */
@Injectable()
export class PredService {
  constructor(
    private readonly recordProcessService : RecordProcessService,
    private readonly recordService : RecordService,
    private readonly userService : UserService
  ) {}

  async getPred(recordFixed: RecordFixed, records: RecordWeekly[]): Promise<RecordOutput[]> {
    try{
      const recordsAttached: RecordInput[] = // Attaching constants and variables
      this.recordProcessService.processRecords(recordFixed, records);
      
      const new_user : Prisma.UserCreateInput = 

      const records_db: Prisma.TrainingProgramRecCreateManyInput[] =
        await records.map((e)=> {
          return {
            program_id: 1,
            user_id: 1,
            start: new Date(),
            end: new Date(),
            comment: '',
            weight: e.weight,
            fat_rate: e.fat_rate,
            squat: e.squat,
            benchpress: e.benchpress,
            deadlift: e.deadlift,
          }
        })
      await this.recordService.createProgramRecs( records_db )
      // transfer records to the model by file
      writeFileSync('./predModel/model24Input.json', JSON.stringify(recordsAttached))
      // spawn a prediction model and get the result
      const pythonProcess = spawnSync('python',["./predModel/model24.py"]);
      // get the result and return
      return pythonProcess.stdout.toString().trim()
      .split('\n')
      .map((preds): RecordOutput =>{
        const pred = preds.split(' ')
        return {
          method: pred[0],
          squat: Number(pred[1]),
          benchpress: Number(pred[2]),
          deadlift: Number(pred[3]),
        }
      })
    } catch (error) {
      return [{method: 'failed', squat:0, benchpress:0, deadlift:0}];
    }
  }
}
