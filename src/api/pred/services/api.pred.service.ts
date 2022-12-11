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
import { RecordPrismaService } from 'src/prisma_services/record/record.prisma.service';
import { Prisma, TrainingProgramRec } from '@prisma/client';
import { UserPrismaService } from 'src/prisma_services/user/user.prisma.service';

/**
 * API Service for training program suggestion and prediction.
 */
@Injectable()
export class PredApiService {
  constructor(
    private readonly recordProcessService : RecordProcessService,
    private readonly recordPrismaService : RecordPrismaService,
    private readonly userPrismaService : UserPrismaService
  ){}

  async getPred(recordFixed: RecordFixed, records: RecordWeekly[]): Promise<RecordOutput[]> {
    try{
      const recordsAttached: RecordInput[] = // Attaching constants and variables
      this.recordProcessService.processRecords(recordFixed, records);

      const latest_user = await this.userPrismaService.latestUsers()
      const latest_id = latest_user[0]['id'];
      
      await this.userPrismaService.createUser(
        {
          email_google: `no${latest_id+1}@no.no`,
          token_google: `no${latest_id+1}`,
          user_name: `n${latest_id+1}`,
          password: "no",
          age: recordFixed.age,
          sex: recordFixed.sex,
          height: recordFixed.height,
          arm_length: recordFixed.arm_length,
          leg_length: recordFixed.leg_length,
        }
      )

      const records_db: Prisma.TrainingProgramRecCreateManyInput[] =
        await records.map((e)=> {
          return {
            program_id: e.program+1,
            user_id: latest_id,
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
      await this.recordPrismaService.createProgramRecs( records_db )
      // transfer records to the model by file
      writeFileSync('./predModel/model24Input.json', JSON.stringify(recordsAttached))
      return [
        {
          method: 'General 5x5 Program',
          squat: records[records.length-1].squat*1.01615,
          benchpress: records[records.length-1].benchpress*1.00751,
          deadlift: records[records.length-1].deadlift*1.0172
        },
        {
          method: 'Beginner Powerlifting Programs',
          squat: records[records.length-1].squat*1.0152,
          benchpress: records[records.length-1].benchpress*1.0081,
          deadlift: records[records.length-1].deadlift*1.0162
        },
        {
          method: 'Intermediate Powerlifting Programs',
          squat: records[records.length-1].squat*1.016,
          benchpress: records[records.length-1].benchpress*1.008,
          deadlift: records[records.length-1].deadlift*1.0101
        },
      ]
      /* ==========================
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
      =================== */
    } catch (error) {
      console.log(error)
      return [{method: 'failed', squat:0, benchpress:0, deadlift:0}];
    }
  }
}
